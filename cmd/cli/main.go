package main

import (
	"flag"
	"log"
	"os"

	"github.com/3WDeveloper-GM/kanban_board_new/cmd/config"
	tea "github.com/charmbracelet/bubbletea"

	_ "github.com/lib/pq"
)

func main() {

	config.ClearTerminal()

	b := config.NewModel()

	flag.StringVar(&b.DSN, "db", os.Getenv("KANBAN_DSN"), "The kanbban DB DSN string.")
	flag.Parse()

	db, err := openDB(b.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	b.InitDB(db)
	b.InitBoard(0, 0)

	p := tea.NewProgram(b)
	if _, err := p.Run(); err != nil {
		log.Panic(err)
	}

}
