package user

import (
	"context"
	"errors"
	"encoding/json"
	"strconv"
	"fmt"
	"os"
	"time"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"   
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/dgrijalva/jwt-go"
)

var RepoErr = errors.New("Unable to handle Repo Request")

type repo struct {
	db     *mongo.Collection
	logger log.Logger
}

func NewRepo(db *mongo.Collection, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	if user.Email == "" || user.Password == "" {
		return RepoErr
	}

	_, err := repo.db.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repo) ViewUser(ctx context.Context, id string) (User, error) {
	var user User
	err := repo.db.FindOne(ctx, bson.M{"id" : id}).Decode(&user)
	
	if err != nil {
		return User{},err
	}

	return user, nil
}

func (repo *repo) UpdateUser(ctx context.Context, id string, name string, phone string, email string) error {

	_, err := repo.db.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{
			"$set": bson.M{
			  "name": name,
			  "phone": phone,
			  "email": email,
			},
		  },
	)
	
	if err != nil {
		return err
	}

	return nil
}

func (repo *repo) DeleteUser(ctx context.Context, id string) error {

	_, err := repo.db.DeleteOne(
		ctx,
		bson.M{"id": id},
	)
	
	if err != nil {
		return err
	}

	return nil
}

func (repo *repo) ListUser(ctx context.Context, limit string, offset string) (string, error) {
	var users []*User

	fmt.Println(limit)
	fmt.Println(offset)
	
	var iLimit int64
	iLimit = 10
	if (limit != "") {
		cLimit, err := strconv.ParseInt(limit, 10, 64)
		iLimit = cLimit
		if (err != nil) {
			return "", err
		}
	}

	var iOffset int64
	iOffset = 0
	if (offset != "") { 
		cOffset, err := strconv.ParseInt(offset, 10, 64)
		iOffset = cOffset
		if (err != nil) {
			return "", err
		}
	}

	options := options.Find()
	options.SetSkip(iOffset)
	options.SetLimit(iLimit)
	cur, err := repo.db.Find(ctx, bson.M{}, options)
	
	if err != nil {
		return "Unable to find data", err
	}

	for cur.Next(ctx) {
        var user User
        err = cur.Decode(&user)
        if err != nil {
            return "Error fetching data", err
        }
        users = append(users, &user)
    }

	usersJson, err := json.Marshal(users)
    if err != nil {
        return "", err
    }

	return string(usersJson), nil
}

func (repo *repo) AuthenticateUser(ctx context.Context, phone string, password string) (string, error) {
	var user User
	err := repo.db.FindOne(ctx, bson.M{"phone" : phone}).Decode(&user)
	
	if err != nil {
		return "",err
	}

	fmt.Println("Auth User Repo")

	if (user.Password == password) {
		os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
		atClaims := jwt.MapClaims{}
		atClaims["authorized"] = true
		atClaims["user_id"] = user.ID
		atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
		if err != nil {
			return "", err
		}
		return token, nil
	}

	return "", nil
}