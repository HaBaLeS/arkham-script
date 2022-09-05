package main

import (
	"arkham-script/eval"
	"log"
)

func main() {
	log.Printf("Start")

	engine := eval.NewEngine()
	engine.InitCard("playground.arkham")

	log.Printf("End")
}
