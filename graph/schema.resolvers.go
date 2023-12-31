package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/graph/model"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	userName := input.Name

	if userName == "" {
		return nil, fmt.Errorf("invalid input provided")
	}

	databaseUser, databaseUserError := r.Database.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userName,
	})

	if databaseUserError != nil {
		return nil, fmt.Errorf("error creating user: %v", databaseUserError)
	}

	user := model.User{
		ID:   databaseUser.ID.String(),
		Name: databaseUser.Name,
	}

	return &user, nil
}

// CreateFeed is the resolver for the createFeed field.
func (r *mutationResolver) CreateFeed(ctx context.Context, input model.CreateFeedInput) (*model.Feed, error) {
	name, url, userId := input.Name, input.URL, input.UserID
	if name == "" || url == "" {
		return nil, fmt.Errorf("invalid input provided")
	}

	dbFeed, err := r.Database.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    uuid.MustParse(userId),
	})

	if err != nil {
		return nil, fmt.Errorf("error creating feed: %v", err)
	}

	feed := model.Feed{
		ID:   dbFeed.ID.String(),
		Name: dbFeed.Name,
		URL:  dbFeed.Url,
	}

	return &feed, nil
}

// CreateFeedFollow is the resolver for the createFeedFollow field.
func (r *mutationResolver) CreateFeedFollow(ctx context.Context, input *model.CreateFeedFollowInput) (*model.FeedFollow, error) {
	userId, feedId := input.UserID, input.FeedID
	if feedId == "" || userId == "" {
		return nil, fmt.Errorf("invalid input provided")
	}

	dbUser, err := r.Database.GetUserById(ctx, uuid.MustParse(userId))
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	dbFeed, err := r.Database.GetFeed(ctx, uuid.MustParse(feedId))
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %v", err)
	}

	dbFeedFollow, err := r.Database.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    dbUser.ID,
		FeedID:    dbFeed.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating feed follow: %v", err)
	}

	feedFollow := model.FeedFollow{
		ID:   dbFeedFollow.ID.String(),
		Feed: &model.Feed{ID: dbFeed.ID.String(), Name: dbFeed.Name, URL: dbFeed.Url},
		User: &model.User{ID: dbUser.ID.String(), Name: dbUser.Name},
	}

	return &feedFollow, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	parsedStringUUID, parsingErr := uuid.Parse(id)

	if parsingErr != nil {
		return nil, fmt.Errorf("user not found")
	}

	databaseUser, databaseUserError := r.Resolver.Database.GetUserById(ctx, parsedStringUUID)

	if databaseUserError != nil {
		return nil, fmt.Errorf("error fetching users: %v", databaseUserError)
	}

	user := model.User{
		ID:   databaseUser.ID.String(),
		Name: databaseUser.Name,
	}
	return &user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	databaseUsers, databaseUserError := r.Resolver.Database.GetUsers(ctx)

	if databaseUserError != nil {
		return nil, fmt.Errorf("error fetching users: %v", databaseUserError)
	}

	users := make([]*model.User, 0)
	for _, dbUser := range databaseUsers {
		user := model.User{ID: dbUser.ID.String(), Name: dbUser.Name}
		users = append(users, &user)
	}

	return users, nil
}

// Feed is the resolver for the feed field.
func (r *queryResolver) Feed(ctx context.Context, id string) (*model.Feed, error) {
	parsedFeedId := uuid.MustParse(id)

	dbFeed, err := r.Resolver.Database.GetFeed(ctx, parsedFeedId)

	if err != nil {
		return nil, fmt.Errorf("error fetching feeds: %v", err)
	}

	feed := model.Feed{
		ID:   dbFeed.ID.String(),
		Name: dbFeed.Name,
		URL:  dbFeed.Url,
	}

	return &feed, nil
}

// Feeds is the resolver for the feeds field.
func (r *queryResolver) Feeds(ctx context.Context) ([]*model.Feed, error) {
	dbFeeds, err := r.Resolver.Database.GetFeeds(ctx)

	if err != nil {
		return nil, fmt.Errorf("error fetching feeds: %v", err)
	}

	feeds := make([]*model.Feed, 0)
	for _, dbFeed := range dbFeeds {
		feed := model.Feed{ID: dbFeed.ID.String(), Name: dbFeed.Name, URL: dbFeed.Url}
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}

// FeedFollows is the resolver for the feedFollows field.
func (r *queryResolver) FeedFollows(ctx context.Context, userID string) ([]*model.FeedFollow, error) {
	dbFeedFollows, err := r.Resolver.Database.GetFeedFollows(ctx, uuid.MustParse(userID))

	if err != nil {
		return nil, fmt.Errorf("error fetching feed follows: %v", err)
	}

	feedFollows := make([]*model.FeedFollow, 0)
	for _, dbFeedFollow := range dbFeedFollows {
		dbUser, err := r.Resolver.Database.GetUserById(ctx, dbFeedFollow.UserID)
		if err != nil {
			return nil, fmt.Errorf("error fetching feed follow user: %v", err)
		}
		dbFeed, err := r.Resolver.Database.GetFeed(ctx, dbFeedFollow.FeedID)
		if err != nil {
			return nil, fmt.Errorf("error fetching feed follow feed: %v", err)
		}
		user := model.User{ID: dbUser.ID.String(), Name: dbUser.Name}
		feed := model.Feed{ID: dbFeed.ID.String(), Name: dbFeed.Name, URL: dbFeed.Url}
		feedFollow := model.FeedFollow{ID: dbFeedFollow.ID.String(), User: &user, Feed: &feed}
		feedFollows = append(feedFollows, &feedFollow)
	}

	return feedFollows, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	dbPost, err := r.Resolver.Database.GetPost(ctx, uuid.MustParse(id))

	if err != nil {
		return nil, fmt.Errorf("error fetching post: %v", err)
	}

	parsedPubDateString := dbPost.PublishedAt.String()
	post := model.Post{
		ID:          dbPost.ID.String(),
		Title:       dbPost.Title,
		Description: &dbPost.Description.String,
		PublishedAt: &parsedPubDateString,
		URL:         dbPost.Url,
	}

	return &post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	dbPosts, err := r.Resolver.Database.GetPosts(ctx)

	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %v", err)
	}

	posts := make([]*model.Post, 0)
	for _, dbPost := range dbPosts {
		parsedPubDateString := dbPost.PublishedAt.String()
		post := model.Post{ID: dbPost.ID.String(), Title: dbPost.Title, Description: &dbPost.Description.String, PublishedAt: &parsedPubDateString, URL: dbPost.Url}
		posts = append(posts, &post)
	}

	return posts, nil
}

// PostsByUser is the resolver for the postsByUser field.
func (r *queryResolver) PostsByUser(ctx context.Context, userID string) ([]*model.Post, error) {
	dbPosts, err := r.Resolver.Database.GetPostsByUser(ctx, database.GetPostsByUserParams{UserID: uuid.MustParse(userID), Limit: 10})

	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %v", err)
	}

	posts := make([]*model.Post, 0)
	for _, dbPost := range dbPosts {
		parsedPubDateString := dbPost.PublishedAt.String()
		post := model.Post{ID: dbPost.ID.String(), Title: dbPost.Title, Description: &dbPost.Description.String, PublishedAt: &parsedPubDateString, URL: dbPost.Url}
		posts = append(posts, &post)
	}

	return posts, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
