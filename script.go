package main

import (
	"arkham-script/eval"
	"log"
	"math/rand"
	"time"
)

var events = []string{
	"card_reveal", "card_played", "round_end", "pre_play_card",
}
var engine *eval.Engine

func main() {
	log.Printf("Start")

	engine = eval.NewEngine()
	engine.InitCards()

	go sendRandomEvents()

	for true {
		time.Sleep(time.Second)
	}

	log.Printf("End")
}

func sendRandomEvents() {
	for true {
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
		ev := events[rand.Intn(len(events))]
		log.Printf("Sendin Event: %s", ev)
		engine.IncomingEvent(ev)
	}
}
