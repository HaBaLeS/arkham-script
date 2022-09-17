package eval

import (
	"arkham-script/dsl"
	"log"
	"os"
	"strings"
)

type Engine struct {
	evalCtx  *EvaluationContext
	scripts  map[string]*cardEvents
	listener map[string][]ListenerCallback
	evq      chan string
}

type cardEvents struct {
	ccode       string
	canPlayCard *ListenerCallback
	onReveal    *ListenerCallback
	onRoundEnd  *ListenerCallback
}

func NewEngine() *Engine {
	e := &Engine{
		scripts:  make(map[string]*cardEvents),
		evq:      make(chan string, 100),
		listener: make(map[string][]ListenerCallback, 0),
	}
	e.evalCtx = New(e)

	return e
}

type ListenerCallback func()

func (e *Engine) RegisterEventListener(event string, callback ListenerCallback) {
	log.Printf("Listener for %s registered", event)

	lq := e.listener[event]
	if lq == nil {
		lq = make([]ListenerCallback, 0)
	}
	e.listener[event] = append(lq, callback)
}

/*
 * - load card script
 * - execute cardScript
 */
func (e *Engine) InitCard(file string) {
	//run card script to register all on ...
	script, err := os.ReadFile("scripts/" + file)
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

	e.scripts[file] = &cardEvents{
		ccode: file,
	}
	e.evalCtx.EvalCardScript(ast)

}

/*
 * can be executed any time to ched if a card can be played.
 * MUST BE WITHOUT SIDE EFFECT
 */
func (e *Engine) CanPlayCard(ccode string) bool {
	var ast any //getFromAstCache
	e.evalCtx.EvalCardScript(ast)

	return false
}

func (e *Engine) InitCards() {
	dirs, _ := os.ReadDir("scripts/")
	for _, d := range dirs {
		if strings.HasSuffix(d.Name(), ".arkham") {
			e.InitCard(d.Name())
		}

	}
}

func (e *Engine) IncomingEvent(ev string) {
	e.evq <- ev
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
