package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/augustfrih/gator/internal/config"
	"github.com/augustfrih/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	var sta state
	sta.cfg = &cfg

	db, err := sql.Open("postgres", cfg.DbURL)
	dbQueries := database.New(db)
	sta.db = dbQueries

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments were provided")
		os.Exit(1)
	}
	var cmd command

	cmd.name, cmd.arguments = args[1], args[2:]

	err = cmds.run(&sta, cmd)
	if err != nil {
		fmt.Printf("Couldnt run command %s. Error: %s", cmd.name, err)
		os.Exit(1)
	}
}
