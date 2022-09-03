package main

import (
	"arkham-script/dsl"
	"arkham-script/eval"
	"log"
	"os"
)

func main() {

	ctx := eval.New()

	script, err := os.ReadFile("scripts/0815.arkham")
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

	log.Printf("execute script")
	ctx.EvalRuledefinition(ast)

}
