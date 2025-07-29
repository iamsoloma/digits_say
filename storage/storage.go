package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type User struct {
	ID           models.RecordID        `json:"id,omitempty"`
	State        map[string]interface{} `json:"State"`
	UserName     string                 `json:"UserName"`
	Name         string                 `json:"Name"`
	Surname      string                 `json:"Surname"`
	FullName     string                 `json:"FullName"`
	LanguageCode string                 `json:"LanguageCode"`
	Email        string                 `json:"Email"`
	Birthdate    string                 `json:"Birthdate"`
}

type DBConfig struct {
	ConnectionURL string
	Namespace     string
	Database      string
	Username      string
	Password      string
}

type Storage struct {
	DBConfig DBConfig
	db       *surrealdb.DB
	ctx      context.Context
	reauth   sync.RWMutex
}

func (s *Storage) ConnectToSurreal() (err error) {
	s.reauth.Lock()
	s.ctx = context.Background()
	s.db, err = surrealdb.Connect(context.Background(), s.DBConfig.ConnectionURL)
	if err != nil {
		return err
	}

	err = s.db.Use(s.ctx, s.DBConfig.Namespace, s.DBConfig.Database)
	if err != nil {
		return err
	}

	token, err := s.db.SignIn(s.ctx, surrealdb.Auth{
		Username: s.DBConfig.Username,
		Password: s.DBConfig.Password,
	})
	if err != nil {
		return err
	}
	if err := s.db.Authenticate(s.ctx, token); err != nil {
		return err
	}
	s.reauth.Unlock()

	return nil
}

func (s *Storage) Close() error {
	s.reauth.Lock()
	defer s.reauth.Unlock()
	if s.db != nil {
		if err := s.db.Close(s.ctx); err != nil {
			return err
		}
		s.db = nil
	}
	return nil
}

func (s *Storage) GetTokenExpirationTime() (time.Time, error) {
	type ExpResult struct {
		Exp int64 `json:"exp"`
	}

	res, err := surrealdb.Query[[]ExpResult](s.ctx, s.db, "SELECT exp FROM $token;", nil)
	if err != nil {
		return time.Now(), err
	}
	exp := (*res)[0].Result[0].Exp
	if time.Unix(exp, 0).IsZero() {
		return time.Now(), fmt.Errorf("token is not set or has no expiration time")
	} /* else {
		fmt.Println("Token expiration time:", time.Unix(exp, 0).UTC().String())
	}*/
	t := time.Unix(exp, 0).UTC()

	return t, nil
}

func (s *Storage) CheckToken() error {
	s.reauth.Lock()
	exp, err := s.GetTokenExpirationTime()
	if err != nil {
		s.Close()
		conerr := s.ConnectToSurreal()
		if conerr != nil {
			return fmt.Errorf("failed to reconnect to SurrealDB: %w", conerr)
		} else {
			s.reauth.Unlock()
			return nil
		}
	}
	if time.Now().After(exp) {
		err = s.ConnectToSurreal()
		if err != nil {
			return fmt.Errorf("failed to reauthenticate: %w", err)
		}
		s.reauth.Unlock()
		return nil
	} else {
		s.reauth.Unlock()
		return nil
	}
}

func (s *Storage) GetUserByID(id string) (user *User, exist bool, err error) {
	err = s.CheckToken()
	if err != nil {
		return nil, false, err
	}
	user, err = surrealdb.Select[User](s.ctx, s.db, models.RecordID{Table: "Users", ID: id})
	if err != nil {
		return nil, false, err
	}

	if user == nil || user.ID == (models.RecordID{}) {
		return nil, false, nil
	}

	return user, true, nil
}

func (s *Storage) UpdateUser(NewUser User) (updatedUser *User, err error) {
	err = s.CheckToken()
	if err != nil {
		return nil, err
	}
	updatedUser = &User{}
	updatedUser, err = surrealdb.Update[User](s.ctx, s.db, models.RecordID{Table: "Users", ID: NewUser.ID.ID.(string)}, NewUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *Storage) RegisterNewUser(user User) (userInDb *User, err error) {
	err = s.CheckToken()
	if err != nil {
		return nil, err
	}
	userInDb = &User{}

	userInDb, err = surrealdb.Create[User](s.ctx, s.db, models.Table("Users"), user)
	if err != nil {
		return userInDb, err
	}

	return userInDb, nil
}
