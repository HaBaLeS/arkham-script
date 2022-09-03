package main

import (
	"arkham-script/dsl"
	"arkham-script/events"
	"log"
)

/*
*

	type Gen struct {
		CCode        string `yaml:"ccode"`
		OnRevealFunc struct {
			ScriptFunc struct {
				Name     string `yaml:"name"`
				ParamMap struct {
					Key string `yaml:"key"`
				} `yaml:"paramMap"`
			} `yaml:"scriptFunc"`
		} `yaml:"on_revealFunc"`
		OnEnterLocation []struct {
		} `yaml:"on_enter_location"`
		OnReveal []struct {
			PlayEnemyFormDiscard interface{} `yaml:"play_enemy_form_discard,omitempty"`
			EnemyType            string      `yaml:"enemy_type,omitempty"`
			PlayTo               string      `yaml:"play_to,omitempty"`

			PlaceEnemy []struct {
				Ccode    int    `yaml:"ccode,omitempty"`
				Name     string `yaml:"name,omitempty"`
				Location string `yaml:"location,omitempty"`
				OnDie    struct {
					ScriptFunc struct {
						Name     string `yaml:"name"`
						ParamMap struct {
							Key string `yaml:"key"`
						} `yaml:"paramMap"`
					} `yaml:"scriptFunc"`
				} `yaml:"on_die,omitempty"`
			} `yaml:"place_enemy,omitempty"`
			PlaceCard []struct {
				Ccode    int    `yaml:"ccode,omitempty"`
				Name     string `yaml:"name,omitempty"`
				Location string `yaml:"location,omitempty"`
			} `yaml:"place_card,omitempty"`
			PlaceLocation []struct {
				Ccode string `yaml:"ccode,omitempty"`
				Name  string `yaml:"name,omitempty"`
			} `yaml:"place_location,omitempty"`
			ForceTest     []interface{} `yaml:"force_test,omitempty"`
			PlayDoomToken []struct {
				Amount    int    `yaml:"amount,omitempty"`
				EmitEvent string `yaml:"emit_event,omitempty"`
			} `yaml:"play_doom_token,omitempty"`
		} `yaml:"on_reveal"`
		[]struct {
			MoveEnemys struct {
				To string `yaml:"to"`
			} `yaml:"move_enemys,omitempty"`
			ScriptFunc struct {
				Name     string `yaml:"name"`
				ParamMap struct {
					Key string `yaml:"key"`
				} `yaml:"paramMap"`
			} `yaml:"scriptFunc,omitempty"`
			ForceTest []interface{} `yaml:"force_test,omitempty"`
		} `yaml:"on_round_end"`
	}
*/
type CardScript struct {
	Events []events.ScriptEvent `yaml:"events"`
	CCode  string               `yaml:"ccode"`
}

func main() {

	cardScript := `

ccode C73119

//on blaevent 
`

	ast, err := dsl.Parse(cardScript)
	if err != nil {
		log.Fatalln(err)
	}

	if ast == nil {
		log.Fatalln("Nooooo its nulllll")
	}
}
