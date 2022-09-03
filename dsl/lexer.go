package dsl

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type item struct {
	Typ  itemType
	Val  string
	Pos  int
	Line int
}

type itemType int

const (
	ItemError itemType = iota
	ItemEOF
	ItemDot
	ItemWord
	ItemNumber
	ItemEqual
	ItemNotEqual
	ItemGreater
	ItemLesser
	ItemRule
	ItemActivate
	ItemDeactivate
	ItemSet
	ItemWhen
	ItemThen
	ItemComment
	ItemOpenSquareBracket
	ItemCloseSquareBracket
	ItemComma
	ItemOpenBracket
	ItemCloseBracket
	ItemRandomAction
	ItemOrderedAction
	ItemDoAction
	ItemAs

	ItemCCode
)

var keywords = map[string]itemType{
	"ccode": ItemCCode,

	"rule":          ItemRule,
	"activate":      ItemActivate,
	"deactivate":    ItemDeactivate,
	"set":           ItemSet,
	"when":          ItemWhen,
	"then":          ItemThen,
	"RandomAction":  ItemRandomAction,
	"OrderedAction": ItemOrderedAction,
	"Do":            ItemDoAction,
	"as":            ItemAs,
}

const eof = -1

type lexer struct {
	input string //input text
	start int
	pos   int
	width int       //with of last rune
	items chan item //channel of scanned items
	line  int
}

func (i item) String() string {
	switch i.Typ {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return i.Val
	}
	return fmt.Sprintf("%q", i.Val)
}

func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for state := lexItem; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func Lex(input string) chan item {
	l := &lexer{
		input: input,
		items: make(chan item),
		line:  1,
	}
	go l.run() //Run states async
	return l.items
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos], l.pos, l.line}
	l.start = l.pos
}

func (l *lexer) next() rune {
	lx := len(l.input)
	if l.pos >= lx {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return r
}

// skipp to here, ignore buffer
func (l *lexer) ignore() {
	l.start = l.pos
}

// move back one rune (WARNING: can be done only once per next())
func (l *lexer) backup() {
	l.pos -= l.width
}

// look at the next rune without consuming
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// errorf emits an error token with the formatted arguments and returns the terminal state.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{Pos: l.pos, Typ: ItemError, Val: fmt.Sprintf(format, args...)}
	return nil
}

// isValidIdent reports whether r is either a letter or a digit
func isValidIdent(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_'
}

///Lexer stat functions below

func lexItem(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(ItemEOF)
			return nil
		case r == '\n':
			l.line++
			l.ignore()
		case r == ' ' || r == '\t':
			l.ignore()
		case r == '.':
			l.emit(ItemDot)
		case r == '=':
			l.emit(ItemEqual)
		case r == '!':
			l.emit(ItemNotEqual)
		case r == '>':
			l.emit(ItemGreater)
		case r == '<':
			l.emit(ItemLesser)
		case r == '[':
			l.emit(ItemOpenSquareBracket)
		case r == ']':
			l.emit(ItemCloseSquareBracket)
		case r == '{':
			l.emit(ItemOpenBracket)
		case r == '}':
			l.emit(ItemCloseBracket)
		case r == ',':
			l.emit(ItemComma)
		case r == '/':
			return lexCommentBegin
		case unicode.IsDigit(r):
			return lexNumber
		case unicode.IsLetter(r):
			return lexText

		default:
			return l.errorf("unexpected token %q", r)
		}
	}
}

func lexNumber(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			//keep the number
		default:
			l.backup()
			l.emit(ItemNumber)
			return lexItem
		}
	}
}

func lexText(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isValidIdent(r):
			//keep the char
		default:
			l.backup()
			if typ, ok := keywords[l.current()]; ok {
				l.emit(typ)
				return lexItem
			}

			l.emit(ItemWord)
			return lexItem
		}
	}
}

func lexCommentBegin(l *lexer) stateFn {
	lineComment := false
	blockComment := false
	for {
		switch r := l.next(); {
		case r == '/':
			lineComment = true
			break
		case r == '*':
			blockComment = true
			break
		case lineComment && r == '\n':
			l.emit(ItemComment)
			l.line++
			return lexItem
		case blockComment && r == '*':
			r2 := l.next()
			if r2 == '\n' {
				l.line++
			}
			if r2 == '/' {
				l.emit(ItemComment)
				return lexItem
			}
			l.backup()
		default:
			if lineComment || blockComment {
				//swallow
				if r == '\n' {
					l.line++
				}
			} else {
				return l.errorf("unexpected token after '/'", r)
			}
		}

	}
}
