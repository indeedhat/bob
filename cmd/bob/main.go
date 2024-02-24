package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/indeedhat/bob/internal/templates"
	"github.com/indeedhat/bob/internal/ui"
)

func main() {
	t, err := templates.Load("templates")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = tea.NewProgram(ui.New(t)).Run(); err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(0)
	}
}
