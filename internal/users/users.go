package users

import (
	"errors"
	"main/internal/platform/dbmanager"
	"main/internal/platform/uuid"

	"golang.org/x/crypto/bcrypt"
)

//User defines the fields user has
type User struct {
	ID        string `json:"id" storm:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" storm:"index,unique"`
	Password  string `json:"password"`
}

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

//New creates a new bucket for User and creates a pointer to Users
func New() {
	dbmanager.CreateBucket(&User{})
}

//AddUser adds an User to the database given user information
func AddUser(firstname string, lastname string, em string, pwd string) (string, error) {
	user := &User{}
	user.FirstName = firstname
	user.LastName = lastname
	user.Email = em
	var err error
	user.Password, err = hash(pwd)
	if err != nil {
		return "", err
	}
	user.ID = uuid.New()
	err = dbmanager.Save(user)
	if err != nil {
		return "", err
	}
	return user.ID, err
}

//RemoveUser removes an User
func RemoveUser(id string) error {
	var user User
	err := dbmanager.Query("ID", id, &user)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&user)
	return err
}

//GetUser returns user information
func GetUser(id string) (User, error) {
	var user User
	err := dbmanager.Query("ID", id, &user)
	if err != nil {
		return User{}, err
	}
	user.Password = ""

	return user, err
}

//UpdateUser updates an User's information
func UpdateUser(id, email, password, firstName, lastName string) error {
	hashedPwd, err := hash(password)
	if err != nil {
		return err
	}
	err = dbmanager.Update(&User{ID: id, Email: email, Password: hashedPwd})
	return err
}

//AuthenticateUser checks the login credentials of an User
func AuthenticateUser(email, pwd string) (string, error) {
	var user User
	err := dbmanager.Query("Email", email, &user)
	if err != nil {
		return "", err
	}
	if comparePassword(user.Password, pwd) {
		return user.ID, nil
	} else {
		return "", errors.New("invalid credentials")
	}

}
