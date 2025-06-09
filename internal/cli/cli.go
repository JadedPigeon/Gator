package cli

import (
	"errors"
	"fmt"

	"github.com/JadedPigeon/Gator/internal/config"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if cmd.Name == "" || len(cmd.Args) == 0 {
		return errors.New("login command requires a username")
	}
	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}
	fmt.Println("Setting current user to", s.Cfg.CurrentUser)
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}
