package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MelonCaully/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no url present in argument")
	}

	// Look up feed by url
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	followsFeed, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(followsFeed.UserName, followsFeed.FeedName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no url present in argument")
	}

	err := s.db.DeleteFeedFollowByUserAndFeedURL(context.Background(), database.DeleteFeedFollowByUserAndFeedURLParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to delete feed %s: %w", cmd.Args[0], err)
	}

	fmt.Printf("Succesfully unfollowed feed: %s\n", cmd.Args[0])
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch following feeds for %s: %w", user.Name, err)
	}

	if len(following) == 0 {
		fmt.Println("no feed followings")
		return nil
	}

	fmt.Printf("Feed follows for user: %s\n", user.Name)
	for _, feed := range following {
		fmt.Printf("* %s\n", feed.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
