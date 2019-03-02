package store

import (
	"errors"
	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"mine-stats/models"
)

func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.orm.One("Username", username, &user)

	return &user, err
}

func (s *Store) GetUserBySessionID(sessionID string) (*models.User, error) {
	if sessionID == "" {
		return nil, errors.New("empty cookie")
	}
	var user models.User
	err := s.orm.One("SessionID", sessionID, &user)

	return &user, err
}

func (s *Store) GetUserByID(ID int) (*models.User, error) {
	var user models.User
	err := s.orm.One("ID", ID, &user)

	return &user, err
}

func CheckPasswordWithHash(hash, plainPwd string) (bool, error) {
	byteHash := []byte(hash)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Store) AddUser(username, password string) (user *models.User, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user = &models.User{
		Username: username,
		Hash: string(hash),
	}

	err = s.orm.Save(user)

	return
}

func (s *Store) VerifyLogin(username, plainPwd string) (user *models.User, err error) {
	user, err = s.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	result, err := CheckPasswordWithHash(user.Hash, plainPwd)
	if !result {
		return nil, err
	}

	return
}

func (s *Store) UpdateUserWithSessionID(user *models.User) (*models.User, error) {
	id, err := gonanoid.Nanoid(32)
	if err != nil {
		return nil, err
	}

	user.SessionID = id

	err = s.orm.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}