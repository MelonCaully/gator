package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MelonCaully/gator/internal/config"
)

type state struct {
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
	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	fmt.Println("User has been set")
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Reading config: %+v\n", cfg)

	s := &state{
		Cfg: &cfg,
	}
	c := commands{
		Command: make(map[string]func(*state, command) error),
	}

	c.register("login", handlerLogin)
	if len(os.Args) < 2 {
		log.Fatal("No command present")
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := c.run(s, cmd); err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Reading config again: %+v\n", cfg)
}
