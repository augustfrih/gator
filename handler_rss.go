package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/augustfrih/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: go run . agg <time between reqs>")
	}

	time_between_reqs, err := time.ParseDuration(cmd.arguments[0])
	fmt.Printf("Collecting feeds every %s", time_between_reqs.String())
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <- ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}


}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("usage: go run . addfeed <name> <url>")
	}

	params := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}

	feed, err := s.db.AddFeed(context.Background(), params)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed created succesfully:")
	printFeedFollow(feed, user)
	fmt.Println()
	fmt.Println("=============================================")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) > 1 {
		return fmt.Errorf("Usage: go run . browse  <(optional)limit>")
	}
	limit := 2
	if len(cmd.arguments) == 1 {
		var err error
		limit, err = strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
	}
	params := database.GetPostsForUserParams {
		Limit: int32(limit),
		Name: user.Name,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	for _, post := range posts {
		printPost(post)
	}
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("usage: go run . feeds")
	}

	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println("=============================================")
		printFeedFollow(feed, user)
	}
	return nil
}

func printFeedFollow(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
	fmt.Printf("* UserName:      %s\n", user.Name)
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: go run . follow <url>")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	feed_follows, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	for _, feed_follow := range feed_follows {
		fmt.Println("=============================")
		fmt.Println(feed_follow.UserName)
		fmt.Println(feed_follow.FeedName)
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("usage: go run . following")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range feeds {
		fmt.Println(follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: go run . unfollow <url>")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Succesfully unfollowed %s\n", feed.Name)

	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Println("======================")
	fmt.Printf("Feed name:        %s\n", post.FeedName)
	fmt.Printf("Post title:       %s\n", post.Title)
	fmt.Printf("Published at:     %s\n", post.PublishedAt.Format(time.RFC3339))
	fmt.Printf("URL:              %s\n", post.Url)
}
