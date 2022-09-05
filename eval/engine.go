package eval

import (
	"arkham-script/dsl"
	"log"
	"os"
)

type Engine struct {
	evalCtx *EvaluationContext
	scripts map[string]*cardEvents
}

type cardEvents struct {
	ccode       string
	canPlayCard *ListenerCallback
	onReveal    *ListenerCallback
	onRoundEnd  *ListenerCallback
}

func NewEngine() *Engine {
	return &Engine{
		evalCtx: &EvaluationContext{},
		scripts: make(map[string]*cardEvents),
	}
}

type ListenerCallback func()

func (e Engine) RegisterEventListener(event string, callback ListenerCallback) {
	log.Printf("Listener for %s registered", event)
	switch event {
		
	}
}

/*
 * - load card script
 * - execute cardScript
 */
func (e Engine) InitCard(ccode string) {
	//run card script to register all on ...
	script, err := os.ReadFile("scripts/" + ccode)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("parse scrtipt")
	ast, err := dsl.Parse(string(script))
	if err != nil {
		log.Fatalln(err)
	}

	if ast == nil {
		log.Fatalln("Nooooo its nulllll")
	}

	e.scripts[ccode] = &cardEvents{
		ccode: ccode,
	}
	e.evalCtx.EvalCardScript(ast)

}

/*
 * can be executed any time to ched if a card can be played.
 * MUST BE WITHOUT SIDE EFFECT
 */
func (e Engine) CanPlayCard(ccode string) bool {
	var ast any //getFromAstCache
	e.evalCtx.EvalCardScript(ast)

	return false
}

/**
Events:

-- With Cards --
pre_play_card
play_card
reveal_card
round_end

-- During Fight --

-- Player Actions --
pre_player_action
player_action


*/

/** play card

- Check can play
- play
- post play



*/
