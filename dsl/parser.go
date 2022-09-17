package dsl

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
)

func Parse(input string) (Ast, error) {
	parser := newParser(Lex(input))
	return parser.Parse()
}

type parser struct {
	tokens <-chan item

	lookahead [2]item
	peekCount int
}

func newParser(tokens <-chan item) *parser {
	return &parser{
		tokens: tokens,
	}
}

func (p *parser) Parse() (ast Ast, err error) {
	// Parsing uses panics to bubble up errors
	defer p.recover(&err)
	ast = p.program()
	return
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (p *parser) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		*errp = e.(error)
	}
	return
}

// peek returns but does not consume the next token.
func (p *parser) peek() item {
	if p.peekCount > 0 {
		return p.lookahead[p.peekCount-1]
	}
	p.peekCount = 1
	p.lookahead[1] = p.lookahead[0]
	p.lookahead[0] = <-p.tokens
	return p.lookahead[0]
}

// next returns the next token.
func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.lookahead[0] = <-p.tokens
	}
	return p.lookahead[p.peekCount]
}

// errorf formats the error and terminates processing.
func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("parser: %s", format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (p *parser) error(err error) {
	p.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type.
func (p *parser) expect(expected itemType) item {
	t := p.next()
	if t.Typ != expected {
		p.unexpected(t, expected)
	}
	return t
}

// unexpected complains about the token and terminates processing.
func (p *parser) unexpected(tok item, expected ...itemType) {
	expectedStrs := make([]string, len(expected))
	for i := range expected {
		expectedStrs[i] = fmt.Sprintf("%v", expected[i])
	}
	expectedStr := strings.Join(expectedStrs, ",")
	p.errorf("unexpected token %v with value %v at pos %d, expected: %v (Line: %d)", tok.Typ, tok.Val, tok.Pos, expectedStr, tok.Line)
}

/*
----------------- Parser Funcs here ----------------
*/

func (p *parser) program() *ProgramNode {
	prog := &ProgramNode{
		Position: 0,
	}
	for {
		peek := p.peek()

		log.Printf("Typ: %v Val: %q", peek.Typ, peek.Val)
		switch peek.Typ {
		case ItemEOF:
			return prog

		//case ItemRule:
		//	s := p.ruleStatement()
		//	prog.Statements = append(prog.Statements, s)

		case ItemCloseBracket:
			return prog

		default:
			s := p.blockStatement()
			if s != nil {
				prog.Statements = append(prog.Statements, s)
			}

		}
	}
}

func (p *parser) ccodeStatement() Node {
	r := p.expect(ItemCCode)
	n := p.expect(ItemString)
	return &CCode{
		r.Pos,
		n.Val,
	}
}

func (p *parser) onStatement() Node {
	r := p.expect(ItemOn)
	n := p.expect(ItemWord)
	p.expect(ItemOpenBracket)
	prog := p.program()
	p.expect(ItemCloseBracket)
	return &On{
		r.Pos,
		n.Val,
		prog,
	}
}

func (p *parser) testStatement() Node {
	r := p.expect(ItemTest)
	w := p.expect(ItemWord)
	p.expect(ItemAgainst)
	a := p.expect(ItemNumber)
	am, _ := strconv.Atoi(a.Val)

	p.expect(ItemOpenBracket)
	s := p.successStatement()
	f := p.failureStatement()
	p.expect(ItemCloseBracket)

	return &Test{
		r.Pos,
		w.Val,
		am,
		s,
		f,
	}
}

func (p *parser) blockStatement() Node {
	n := p.peek()
	switch n.Typ {
	case ItemError:
		panic(p.peek().Val)
		break
	case ItemComment:
		p.next() //consume comment
		fmt.Println("Ignoring: Comment")
		return nil
	case ItemCCode:
		return p.ccodeStatement()
	case ItemOn:
		return p.onStatement()
	case ItemEmit:
		return p.emitStatement()
	case ItemPrint:
		return p.printStatement()
	case ItemTest:
		return p.testStatement()
	case ItemDamage:
		return p.damageStatement()
	case ItemIntercept:
		return p.InterceptStatement()
	/*case ItemRandomAction:
		return p.randomActionStatement()
	case ItemOrderedAction:
		return p.orderedActionStatement()
	case ItemActivate:
		return p.activateStatement()
	case ItemDeactivate:
		return p.deactivateStatement()
	case ItemDoAction:
		return p.doActionStatement()
	*/

	default:
		p.errorf("Don't know what to do with token: \"%s\", (line: %d)", n.Val, n.Line)
	}
	return nil
}

func (p *parser) getPropertie() *PropertieNode {
	var a item
	v := p.next()
	if v.Typ == ItemNumber {
		return &PropertieNode{
			Position: v.Pos,
			Object:   v.Val,
		}
	} else if v.Typ == ItemWord {
		n := p.peek()
		if n.Typ == ItemDot {
			p.next()
			a = p.expect(ItemWord)
		}
		return &PropertieNode{
			v.Pos,
			v.Val,
			a.Val,
		}
	}
	return nil

}

func (p *parser) getOperator() string {
	switch n := p.peek(); {
	case n.Typ == ItemEqual:
		p.next()
		return n.Val
	case n.Typ == ItemGreater:
		p.next()
		return n.Val
	case n.Typ == ItemLesser:
		p.next()
		return n.Val
	case n.Typ == ItemNotEqual:
		p.next()
		return n.Val
	default:
		p.errorf("Excpected: <, >, =, ! got %v", n.Val)
	}
	return ""
}

func (p *parser) getWordList() *WordListNode {
	words := make([]string, 0, 10)
	p.expect(ItemOpenSquareBracket)
	for {
		switch t := p.next(); {
		case t.Typ == ItemCloseSquareBracket:
			return &WordListNode{
				t.Pos,
				words,
			}
		case t.Typ == ItemComma:
			//Ignore
		case t.Typ == ItemWord:
			words = append(words, t.Val)
		default:
			p.error(fmt.Errorf("Illegal expression in wordlist: %s of type: %v", t.Val, t.Typ))
		}
	}
}

func (p *parser) emitStatement() Node {
	r := p.expect(ItemEmit)
	e := p.expect(ItemWord)
	p.expect(ItemColon)
	a := p.getWordList()

	return &Emit{
		r.Pos,
		e.Val,
		a.Words,
	}
}

func (p *parser) printStatement() Node {
	r := p.expect(ItemPrint)
	s := p.expect(ItemString)
	return &Print{
		r.Pos,
		s.Val,
	}
}

func (p *parser) successStatement() *ProgramNode {
	p.expect(ItemSuccess)
	p.expect(ItemOpenBracket)
	prog := p.program()
	p.expect(ItemCloseBracket)
	return prog
}

func (p *parser) failureStatement() *ProgramNode {
	p.expect(ItemFailure)
	p.expect(ItemOpenBracket)
	prog := p.program()
	p.expect(ItemCloseBracket)
	return prog
}

func (p *parser) damageStatement() Node {
	var key item
	r := p.expect(ItemDamage)
	w := p.expectItemValue()
	p.expect(ItemColon)
	a := p.expectItemPropertie()

	t := p.expect(ItemWord)
	return &Damage{
		r.Pos,
		w,
		a,
		t.Val,
		key.Val,
	}
}

func (p *parser) InterceptStatement() Node {
	r := p.expect(ItemIntercept)
	w := p.expect(ItemWord)
	p.expect(ItemOpenBracket)
	prog := p.program()
	p.expect(ItemCloseBracket)
	return &Intercept{
		r.Pos,
		w.Val,
		prog,
	}
}

func (p *parser) expectItemValue() Node {
	var key item
	obj := p.expect(ItemWord)
	if p.peek().Typ == ItemOpenSquareBracket {
		p.next()
		key = p.expect(ItemWord)
		p.expect(ItemCloseSquareBracket)
	}
	return &Value{
		Position: obj.Pos,
		Map:      obj.Val,
		Key:      key.Val,
	}
}

func (p *parser) expectItemPropertie() Node {
	var key item
	var obj item

	if p.peek().Typ == ItemNumber {
		obj = p.expect(ItemNumber)
	} else {
		obj = p.expect(ItemWord)
		if p.peek().Typ == ItemDot {
			p.next()
			key = p.expect(ItemWord)
		}
	}

	return &Value{
		Position: obj.Pos,
		Map:      obj.Val,
		Key:      key.Val,
	}
}

/*
func (p *parser) ruleStatement() Node {
	r := p.expect(ItemRule)
	p.expect(ItemOpenBracket)
	w := p.whenStatement()
	t := p.thenStatement()
	p.expect(ItemCloseBracket)
	p.expect(ItemAs)
	n := p.expect(ItemWord)
	return &RuleStatement{
		r.Pos,
		w,
		t,
		n.Val,
	}
}


func (p *parser) whenStatement() *WhenStatement {
	n := p.expect(ItemWhen)
	o := p.getPropertie()
	op := p.getOperator()
	a := p.getPropertie()
	return &WhenStatement{
		n.Pos,
		o,
		op,
		a,
	}
}

func (p *parser) thenStatement() *ThenStatement {
	n := p.expect(ItemThen)
	b := p.blockStatement()
	return &ThenStatement{
		n.Pos,
		b,
	}
}*/

/*
func (p *parser) activateStatement() *ActivateStatementNode {
	a := p.expect(ItemActivate)
	w := p.expect(ItemWord)
	return &ActivateStatementNode{
		a.Pos,
		w.Val,
	}
}

func (p *parser) deactivateStatement() *DeactivateStatementNode {
	a := p.expect(ItemDeactivate)
	w := p.expect(ItemWord)
	return &DeactivateStatementNode{
		a.Pos,
		w.Val,
	}
}

func (p *parser) doActionStatement() *DoActionStatement {
	d := p.expect(ItemDoAction)
	w := p.expect(ItemWord)
	if p.peek().Typ == ItemAs {
		p.expect(ItemAs)
		n := p.expect(ItemWord)
		return &DoActionStatement{
			d.Pos,
			w.Val,
			n.Val,
		}
	}
	return &DoActionStatement{
		Position:        d.Pos,
		ActionToExecute: w.Val,
	}

}

func (p *parser) randomActionStatement() *RandomActionStatementNode {
	r := p.expect(ItemRandomAction)
	p.expect(ItemOpenSquareBracket)
	wl := p.getWordList()
	p.expect(ItemAs)
	n := p.expect(ItemWord)

	return &RandomActionStatementNode{
		r.Pos,
		wl,
		n.Val,
	}
}

func (p *parser) orderedActionStatement() *OrderedActionStatementNode {
	r := p.expect(ItemOrderedAction)
	p.expect(ItemOpenSquareBracket)
	wl := p.getWordList()
	p.expect(ItemAs)
	n := p.expect(ItemWord)

	return &OrderedActionStatementNode{
		r.Pos,
		wl,
		n.Val,
	}
}*/
