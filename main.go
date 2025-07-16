package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/MelonCaully/gator/internal/config"
	"github.com/MelonCaully/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	Cfg *config.Config
}

func main() {
	// Read data from the config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Reading config: %+v\n", cfg)

	// Intialize state and commmands struct
	s := &state{
		Cfg: &cfg,
	}
	c := commands{
		Commands: make(map[string]func(*state, command) error),
	}

	// Establish a connection with the database
	db, err := sql.Open("postgres", s.Cfg.DbURL)
	if err != nil {
		log.Fatalf("database not found: %v", err)
	}

	// Create a new Queries instance from the database connection and assign it to the state
	dbQueries := database.New(db)
	s.db = dbQueries

	// Checks if command was entered
	if len(os.Args) < 2 {
		log.Fatal("No command present")
	}

	// register commands
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", handlerAddFeed)
	c.register("feeds", handlerFeeds)

	// Initialize the command struct
	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	// Run the command
	if err := c.run(s, cmd); err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	// Read data from the config... again
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Reading config again: %+v\n", cfg)
}
