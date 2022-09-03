package engine

import "log"

type Engine struct {
}

type ListenerCallback func()

func (e Engine) RegisterEventListener(event string, callback ListenerCallback) {
	log.Printf("Listener for %s registered", event)
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
