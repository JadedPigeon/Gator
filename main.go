package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/JadedPigeon/Gator/internal/cli"
	"github.com/JadedPigeon/Gator/internal/config"
	"github.com/JadedPigeon/Gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Load the config from disk
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}
	fmt.Println("Current DB URL:", cfg.DBURL)
	fmt.Println("Current User:", cfg.CurrentUser)

	// Open the DB connection using the config's DB URL
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize the SQLC-generated query wrapper
	dbQueries := database.New(db)

	// Store both the DB and config in the state
	s := &cli.State{
		Cfg: &cfg,
		DB:  dbQueries,
	}
	// Create a new commands structure and initialize it
	cmds := cli.Commands{
		Handlers: make(map[string]func(*cli.State, cli.Command) error),
	}

	// Register a handler function for the login command
	cmds.Register("login", cli.HandlerLogin)
	cmds.Register("register", cli.HandlerRegister)
	cmds.Register("reset", cli.HandlerReset)
	cmds.Register("users", cli.HandlerUsers)
	cmds.Register("agg", cli.HandlerAgg)
	cmds.Register("addfeed", cli.HandlerFeeds)

	// Use os.Args to get the command-line arguments passed in by the user
	if len(os.Args) < 2 {
		log.Fatal("no command provided, please specify a command")
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

	// // Database setup
	// sql.Open("postgres", cfg.DBURL)
	// db, err := sql.Open("postgres", dbURL)
	// dbQueries := database.New(db)

	// type state struct {
	// db  *database.Queries
	// cfg *config.Config
	// }
}
