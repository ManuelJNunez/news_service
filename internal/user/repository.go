package user

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ErrNotFound is used when it is not possible to find the requested User.
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExists is used when trying to create a user that already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

type Repository interface {
	FindOne(ctx context.Context, input LoginInput) (*UserOutput, error)
	Create(ctx context.Context, input LoginInput) (*UserOutput, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &mongoRepository{collection: collection}
}

func (r *mongoRepository) FindOne(ctx context.Context, input LoginInput) (*UserOutput, error) {
	var result User

	// Query MongoDB with provided credentials
	err := r.collection.FindOne(ctx, input).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		slog.Error("error finding user", slog.Any("error", err))
		return nil, err
	}

	slog.Info("user found successfully", slog.String("username", result.Username))
	userOutput := result.ToOutput()
	return &userOutput, nil
}

func (r *mongoRepository) Create(ctx context.Context, input LoginInput) (*UserOutput, error) {
	slog.Debug("creating user", slog.String("username", input.Username))

	// Check if user already exists
	var existing User
	err := r.collection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&existing)
	if err == nil {
		slog.Warn("user already exists", slog.String("username", input.Username))
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		slog.Error("error checking user existence", slog.String("username", input.Username), slog.Any("error", err))
		return nil, err
	}

	user := User{
		Username: input.Username,
		Password: input.Password,
	}

	// Insert user in collection
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		slog.Error("error creating user", slog.String("username", input.Username), slog.Any("error", err))
		return nil, err
	}

	// Prepare return value
	userOutput := &UserOutput{
		ID:       result.InsertedID.(bson.ObjectID).Hex(),
		Username: input.Username,
	}

	slog.Info("user created successfully", slog.String("username", userOutput.Username), slog.String("id", userOutput.ID))
	return userOutput, nil
}
