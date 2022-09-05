package eval

import (
	"arkham-script/dsl"
	"log"
	"reflect"
)

type EvaluationContext struct {
	engine Engine
	ccode  string
}

func New() EvaluationContext {
	return EvaluationContext{}
}

func (e *EvaluationContext) SetData(name string, dp interface{}) {
}

func (e *EvaluationContext) EvalCardScript(ast dsl.Ast) {
	//log.Printf("process: %s\n", reflect.TypeOf(ast))
	switch n := ast.(type) {

	case *dsl.ProgramNode:
		e.evalNodeList(n.Statements)
	//case *dsl.ActivateStatementNode:
	//e.activate(n)
	//case *dsl.DeactivateStatementNode:
	//e.deactivate(n)
	//case *dsl.DoActionStatement:
	//e.defineDo(n)
	//case *dsl.RandomActionStatementNode:
	//e.defineRandomAction(n)
	//case *dsl.OrderedActionStatementNode:
	//e.defineOrderedAction(n)
	//case *dsl.RuleStatement:
	//e.defineRule(n)

	case *dsl.CCode:
		e.execRegisterCardCode(n)
	case *dsl.On:
		e.execOn(n)
	case *dsl.Emit:
		e.execEmit(n)
	case *dsl.Print:
		e.execPrint(n)
	case *dsl.Test:
		e.execTest(n)
	case *dsl.Damage:
		e.execDamage(n)
	case *dsl.Intercept:
		e.execIntercept(n)
	default:
		log.Printf("Unexpexted: %s\n", reflect.TypeOf(n))
		panic("Ouch")
	}
}
func (e *EvaluationContext) evalNodeList(nodes []dsl.Node) {
	for _, s := range nodes {
		e.EvalCardScript(s)
	}
}

func (e *EvaluationContext) execRegisterCardCode(n *dsl.CCode) {
	e.ccode = n.CCode
}

func (e *EvaluationContext) execOn(n *dsl.On) {
	e.engine.RegisterEventListener(n.Event, func() {
		e.evalNodeList(n.Programm.Statements)
	})
}

func (e *EvaluationContext) execEmit(n *dsl.Emit) {
	log.Printf("Emit event with name: %s with arg: %v", n.Event, n.Arguments)
}

func (e *EvaluationContext) execPrint(n *dsl.Print) {
	log.Printf("PRINT: %s", n.Text)
}

func (e *EvaluationContext) execTest(n *dsl.Test) {
	log.Printf("Player has do do a Test: %s for %d", n.What, n.Level)
	log.Printf("If he succeeds:")
	e.evalNodeList(n.Success.Statements)
	log.Printf("If he fails:")
	e.evalNodeList(n.Failure.Statements)

}

func (e *EvaluationContext) execDamage(n *dsl.Damage) {
	log.Printf("%s takes %d Damage to %s", n.Who, n.Amount, n.Where)
}

func (e *EvaluationContext) execIntercept(n *dsl.Intercept) {
	log.Printf("Register Intercept on %s event", n.When)
	e.evalNodeList(n.Program.Statements)
}

/*
func (e *EvaluationContext) activate(n *dsl.ActivateStatementNode) {
	action := e.namedActions[n.NamedRuleOrAction]
	rule := e.rules[n.NamedRuleOrAction]
	if action != nil {
		action.SetActivation(true)
	} else if rule != nil {
		rule.SetActivation(true)
	} else {
		err := fmt.Errorf("Error Item to activate not found: \"%s\" Line: %v", n.NamedRuleOrAction, n.Position)
		panic(err)
	}
}

func (e *EvaluationContext) deactivate(n *dsl.DeactivateStatementNode) {
	action := e.namedActions[n.NamedRuleOrAction]
	rule := e.rules[n.NamedRuleOrAction]
	if action != nil {
		action.SetActivation(false)
	} else if rule != nil {
		rule.SetActivation(false)
	} else {
		err := fmt.Errorf("Error Item to deactivate not found: \"%s\" Line: %v", n.NamedRuleOrAction, n.Position)
		panic(err)
	}
}



func (e *EvaluationContext) defineRandomAction(s *dsl.RandomActionStatementNode) {
	e.namedActions[s.ActionName] = &NamedAction{
		activatable: activatable{false},
		name: s.ActionName,
		actions: s.Actions,
		executeFunc:randomActionExecute,
	}
	e.actionNameToIndex = append(e.actionNameToIndex,s.ActionName)

}

func (e *EvaluationContext) defineOrderedAction(s *dsl.OrderedActionStatementNode) {
	e.namedActions[s.ActionName] = &NamedAction{
		activatable: activatable{false},
		name: s.ActionName,
		orderedActionIndex: 0,
		executeFunc:orderedActionExecute,
		actions: s.Actions,
	}
	e.actionNameToIndex = append(e.actionNameToIndex,s.ActionName)
}

func (e *EvaluationContext) defineRule(rule *dsl.RuleStatement) {
	e.rules[rule.RuleName] = &NamedRule{
		name:rule.RuleName,
		activatable: activatable{false},
		statement: rule,
	}
	e.ruleNameToIndex = append(e.ruleNameToIndex, rule.RuleName)
}


func (e *EvaluationContext) defineDo(s *dsl.DoActionStatement) {
	e.namedActions[s.ActionName] = &NamedAction{
		activatable: activatable{false},
		name: s.ActionName,
		executeFunc: doActionExecute,
	}
	e.actionNameToIndex = append(e.actionNameToIndex,s.ActionName)

}
func (e *EvaluationContext) EvalForNextAction() string {
	var res string
	for _, k := range e.ruleNameToIndex{
		v := e.rules[k]
		res = evalRule(v,e)
		if res != "" {
			break
		}
	}
	if res == ""{
		for _,k := range e.actionNameToIndex {
			v := e.namedActions[k]
			res = v.executeFunc(v,e)
			if res != "" {
				break
			}
		}
	}

	return "Next action is: " + res
}



//fixme alle rules die aktiv sind nicht unur die letzte

func (e *EvaluationContext) getProp(node *dsl.PropertieNode) Prop {

	n, err := strconv.Atoi(node.Object)
	if err == nil{
		return Prop{ n }
	}

	obj := e.data[node.Object]
	if obj == nil {
		panic(fmt.Errorf("Unknown Object: %s", node.Object))
	}
	if node.Attribute == ""{
		return Prop{ obj }
	}else {
		s := structs.New(obj)
		v := s.Field(node.Attribute).Value()
		return Prop{
			val: v,
		}
	}
}



func evalRule(rule *NamedRule,ctx *EvaluationContext) string {
	if !rule.Status(){
		return ""
	}
	color.Green("Evaluating: %s", rule.name)

	obj := ctx.getProp(rule.statement.WhenStatement.Object)
	comp := ctx.getProp(rule.statement.WhenStatement.Comparator)
	compRes := obj.compareTo(comp, rule.statement.WhenStatement.Operator)
	if !compRes{
		return ""
	}
	return ""
}

type Prop struct{
	val interface{}
}

func (p *Prop) compareTo(obj Prop, s string) bool {
	fmt.Printf("compare %v with %v and operator %v\n", p.val, obj.val, s)

	return false
}

type NamedRule struct {
	name string
	activatable
	statement *dsl.RuleStatement
}

type NamedAction struct {
	name string
	activatable
	executeFunc executeFunc
	orderedActionIndex int
	actions *dsl.WordListNode
}

type activatable struct {
	status bool
}

func (d *activatable) SetActivation(doActivate bool){
	d.status = doActivate
}

func (d *activatable) Status() bool{
	return d.status
}
*/
/* Functions to execute the actions when needed */

/**
Define the function layout for the namedaction execute
*/

/*
type executeFunc func(action *NamedAction, ctx *EvaluationContext) string


func doActionExecute(d *NamedAction, ctx *EvaluationContext) string {
	if !d.Status(){
		return ""
	}
	return d.name
}

func randomActionExecute(d *NamedAction, ctx *EvaluationContext) string {
	if !d.Status(){
		return ""
	}
	idx := ctx.d100.RollSides(len(d.actions.Actions))-1
	return d.actions.Actions[idx]
}

func orderedActionExecute(d *NamedAction, ctx *EvaluationContext) string{
	if !d.Status(){
		return ""
	}
	if d.orderedActionIndex == len(d.actions.Actions){
		d.orderedActionIndex = 0
	}
	r :=  d.actions.Actions[d.orderedActionIndex]
	d.orderedActionIndex++
	return r

}
*/
