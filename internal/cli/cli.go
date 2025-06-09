package cli

import (
	"errors"
	"fmt"

	"github.com/JadedPigeon/Gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if cmd.name == "" || len(cmd.args) == 0 {
		return errors.New("login command requires a username")
	}
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}
	fmt.Println("Setting current user to", s.cfg.CurrentUser)
	return nil
}

// func (c *commands) register(name string, f func(*state, command) error) {
// 	c.handlers[name] = f
// }

// func (c *commands) run(s *state, cmd command) error {
// 	handler, ok := c.handlers[cmd.name]
// 	if !ok {
// 		return fmt.Errorf("unknown command: %s", cmd.name)
// 	}
// 	return handler(s, cmd)
// }
