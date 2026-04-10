package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/augustfrih/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("username is required")
	}

	userName := cmd.arguments[0]

	if _, err := s.db.GetUser(context.Background(), userName); err != nil {
		fmt.Printf("Cant login. %s user does not exists\n", userName)
		os.Exit(1)
	}


	s.cfg.CurrentUserName = userName
	err := s.cfg.SetUser(s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	fmt.Println("Current user has been changed to: " + s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("user is required")
	}

	name := cmd.arguments[0]
	if _, err := s.db.GetUser(context.Background(), name); err == nil {
		fmt.Printf("Cant register %s, user already exists", name)
		os.Exit(1)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		fmt.Printf("Couldnt create user %s. Error: %v", name, err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(name)

	fmt.Printf("User %s was created", name)
	fmt.Println(user)

	return nil
}

func reset(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("usage: 'go run . reset'")
	}
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database was reset succesfully")
	return nil
}

func users(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("usage: 'go run . users'")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println(user.Name + " (current)")
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil
}
