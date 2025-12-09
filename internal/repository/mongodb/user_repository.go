package mongodb

import (
	"context"
	"go-auth-backend/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type userRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{coll: db.Collection("users")}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	user.ID = primitive.NewObjectID()
	_, err := r.coll.InsertOne(ctx, user)
	return err
}
