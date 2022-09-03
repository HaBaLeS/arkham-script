package dsl

type Node interface {
	//Pos() Position
}

type Position struct {
	Line int
	Char int
}

func (p Position) Pos() Position {
	return p
}

type Ast Node

type ProgramNode struct {
	Position   int
	Statements []Node
}

type ActivateStatementNode struct {
	Position          int
	NamedRuleOrAction string
}

type DeactivateStatementNode struct {
	Position          int
	NamedRuleOrAction string
}

type DoActionStatement struct {
	Position        int
	ActionToExecute string
	ActionName      string
}

type RandomActionStatementNode struct {
	Position   int
	Actions    *WordListNode
	ActionName string
}

type OrderedActionStatementNode struct {
	Position   int
	Actions    *WordListNode
	ActionName string
}

type WordListNode struct {
	Position int
	Actions  []string
}

type RuleStatement struct {
	Position      int
	WhenStatement *WhenStatement
	ThenStatement *ThenStatement
	RuleName      string
}

type WhenStatement struct {
	Position   int
	Object     *PropertieNode
	Operator   string
	Comparator *PropertieNode
}

type ThenStatement struct {
	Position  int
	Statement Node
}

type PropertieNode struct {
	Position  int
	Object    string
	Attribute string
}

type CCode struct {
	Position int
	CCode    string
}
