package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required")
	}
	s.cfg.CurrentUserName = cmd.arguments[0]
	err := s.cfg.SetUser(s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	fmt.Println("Current user has been changed to: " + s.cfg.CurrentUserName)
	return nil
}
