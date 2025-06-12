package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JadedPigeon/Gator/internal/config"
	"github.com/JadedPigeon/Gator/internal/database"
	"github.com/JadedPigeon/Gator/internal/rss"
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

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("reset command does not take any arguments")
	}
	err := s.DB.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all users: %v", err)
	}
	fmt.Println("All users deleted successfully.")
	if err := s.Cfg.SetUser(""); err != nil {
		return fmt.Errorf("error resetting current user in config: %v", err)
	}
	fmt.Println("All users deleted")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("users command does not take any arguments")
	}
	users, err := s.DB.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %v", err)
	}
	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}
	fmt.Println("Users:")
	for _, user := range users {
		output := "* " + user.Name
		if user.Name == s.Cfg.CurrentUser {
			output += " (current)"
		}
		fmt.Println(output)
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func HandlerAddFeeds(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("feeds command requires a name and URL")
	} else if len(cmd.Args) > 2 {
		return errors.New("feeds command takes only a name and URL")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("current user %s does not exist", s.Cfg.CurrentUser)
		}
		return fmt.Errorf("error checking current user: %v", err)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}
	fmt.Printf("Feed added:\n- ID: %s\n- Name: %s\n- URL: %s\n", feed.ID, feed.Name, feed.Url)
	if _, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}); err != nil {
		return fmt.Errorf("error following feed after creation: %v", err)
	}
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("feeds command does not take any arguments")
	}

	feeds, err := s.DB.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %v", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	for _, feed := range feeds {
		output := fmt.Sprintf("* %s (%s), - Added by %s", feed.Name, feed.Url, feed.UserName)
		fmt.Println(output)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("follow command requires a feed URL")
	}
	feedURL := cmd.Args[0]
	feed, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed with URL %s does not exist", feedURL)
		}
		return fmt.Errorf("error checking feed: %v", err)
	}
	user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("current user %s does not exist", s.Cfg.CurrentUser)
		}
		return fmt.Errorf("error checking current user: %v", err)
	}
	follow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %v", err)
	}
	fmt.Printf("Successfully followed feed %s (%s) for user %s\n", follow.FeedName, feed.Url, follow.UserName)
	return nil
}

func HandlerFollowing(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("following command does not take any arguments")
	}
	user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("current user %s does not exist", s.Cfg.CurrentUser)
		}
		return fmt.Errorf("error checking current user: %v", err)
	}

	follows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving followed feeds: %v", err)
	}
	if len(follows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}
	fmt.Println("Feeds you are following:")
	for _, follow := range follows {
		fmt.Printf("* %s (%s)\n", follow.FeedName, follow.FeedUrl)
	}
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
