package main

import (
	"log"

	"github.com/indeedhat/bob/internal/templates"
)

func main() {
	t, err := templates.Load("templates")
	if err != nil {
		log.Fatal(err)
	}

	for k, g := range t {
		vars, err := g.Vars()
		log.Printf("%s: %v\n%s", k, vars, err)
	}
}
