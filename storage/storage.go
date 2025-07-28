package storage

import (
	"context"
	"fmt"

	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type User struct {
	ID           models.RecordID `json:"id,omitempty"`
	State        string          `json:"State"`
	UserName     string          `json:"UserName"`
	Name         string          `json:"Name"`
	Surname      string          `json:"Surname"`
	FullName     string          `json:"FullName"`
	LanguageCode string          `json:"LanguageCode"`
	Email        string          `json:"Email"`
	Birthdate    string          `json:"Birthdate"`
}

type DBConfig struct {
	ConnectionURL string
	Namespace     string
	Database      string
	Username      string
	Password      string
}

func ConnectToSurreal(config DBConfig) (db *surrealdb.DB, err error) {
	ctx := context.Background()
	db, err = surrealdb.Connect(context.Background(), config.ConnectionURL)
	if err != nil {
		return db, err
	}

	err = db.Use(ctx, config.Namespace, config.Database)
	if err != nil {
		return db, err
	}

	token, err := db.SignIn(ctx, surrealdb.Auth{
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return db, err
	}
	if err := db.Authenticate(ctx, token); err != nil {
		return db, err
	}

	return db, nil
}

func GetUserByID(id string, db *surrealdb.DB) (user *User, exist bool, err error) {
	user, err = surrealdb.Select[User](context.Background(), db, models.RecordID{Table: "Users", ID: id})
	if err != nil {
		return nil, false, err
	}

	if user == nil || user.ID == (models.RecordID{}) {
		return nil, false, nil
	}

	return user, true, nil
}

func UpdateUser(NewUser User, db *surrealdb.DB) (updatedUser *User, err error) {
	updatedUser = &User{}
	updatedUser, err = surrealdb.Update[User](context.Background(), db, models.RecordID{Table: "Users", ID: NewUser.ID.ID.(string)}, NewUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func RegisterNewUser(user User, db *surrealdb.DB) (userInDb *User, err error) {
	userInDb = &User{}

	userInDb, err = surrealdb.Create[User](context.Background(), db, models.Table("Users"), user)
	if err != nil {
		return userInDb, err
	}
	fmt.Printf("Registered a new user with a map %+v\n", userInDb)

	return userInDb, nil
}