package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MelonCaully/gator/internal/config"
	"github.com/MelonCaully/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	Cfg *config.Config
}

type commands struct {
	Command map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.Command[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.Command == nil {
		c.Command = make(map[string]func(*state, command) error)
	}
	c.Command[name] = f
}

type command struct {
	Name string
	Args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("no command present")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.Args[0]); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(os.Stderr, "User '%s' does not exist.\n", cmd.Args[0])
			os.Exit(1)
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no username provided")
	}

	id := uuid.New()
	name := cmd.Args[0]
	now := time.Now()

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Fprintf(os.Stderr, "User with name '%s' already exists.\n", name)
			os.Exit(1)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("failed to set user in config: %w", err)
	}

	fmt.Println("User created successfully.")
	log.Printf("Created user: ID=%s, Name=%s, CreatedAt=%s", id.String(), name, now.Format(time.RFC3339))
	return nil
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
		Command: make(map[string]func(*state, command) error),
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
