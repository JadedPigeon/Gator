package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JadedPigeon/Gator/internal/cli"
	"github.com/JadedPigeon/Gator/internal/config"
)

func main() {
	// Load the config from disk
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}
	fmt.Println("Current DB URL:", cfg.DBURL)
	fmt.Println("Current User:", cfg.CurrentUser)

	// Store the config pointer in the state
	s := &cli.State{
		Cfg: &cfg,
	}

	// Create a new commands structure and initialize it
	cmds := cli.Commands{
		Handlers: make(map[string]func(*cli.State, cli.Command) error),
	}

	// Register a handler function for the login command
	cmds.Register("login", cli.HandlerLogin)

	// Use os.Args to get the command-line arguments passed in by the user
	if len(os.Args) < 2 {
		fmt.Errorf("no command provided, please specify a command")
	}
	args := os.Args[1:]
	cmd := args[0]
	cmdArgs := args[1:]

	// Create a Command struct with the command name and arguments
	command := cli.Command{
		Name: cmd,
		Args: cmdArgs,
	}

	// Run the command using the registered handlers
	commandErr := cmds.Run(s, command)
	if commandErr != nil {
		log.Fatalf("error running command '%s': %v", command.Name, commandErr)
	}
	fmt.Println("Command executed successfully.")
}
