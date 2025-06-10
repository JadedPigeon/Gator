package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JadedPigeon/Gator/internal/config"
	"github.com/JadedPigeon/Gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Cfg *config.Config
	DB  *database.Queries
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
	username := cmd.Args[0]
	_, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %s does not exist", username)
		}
		return fmt.Errorf("error checking user: %v", err)
	}
	if err := s.Cfg.SetUser(username); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}
	fmt.Println("Setting current user to", s.Cfg.CurrentUser)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if cmd.Name == "" || len(cmd.Args) == 0 {
		return errors.New("register command requires a username")
	}
	username := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user %s already exists", username)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("unexpected error checking user: %v", err)
	}

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	fmt.Printf("User %s created with ID %s\n", user.Name, user.ID)
	if err := s.Cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("error setting user in config: %v", err)
	}
	fmt.Printf("User %q registered successfully!\n", user.Name)
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
