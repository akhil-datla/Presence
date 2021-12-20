package organizers

import (
	"main/internal/platform/dbmanager"
	"main/internal/platform/uuid"

	"golang.org/x/crypto/bcrypt"
)

//Organizer defines the fields organizer has
type Organizer struct {
	ID        string `json:"participantID" storm:"id"`
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

//New creates a new bucket for Organizer and creates a pointer to Organizers
func New() {
	dbmanager.CreateBucket(&Organizer{})
}

//AddOrganizer adds an Organizer to the database given organizer information
func AddOrganizer(firstname string, lastname string, em string, pwd string) (string, error) {
	organizer := &Organizer{}
	organizer.FirstName = firstname
	organizer.LastName = lastname
	organizer.Email = em
	var err error
	organizer.Password, err = hash(pwd)
	if err != nil {
		return "", err
	}
	organizer.ID = uuid.New()
	err = dbmanager.Save(organizer)
	if err != nil {
		return "", err
	}
	return organizer.ID, err
}

//RemoveOrganizer removes an Organizer
func RemoveOrganizer(id string) error {
	var organizer Organizer
	err := dbmanager.Query("ID", id, &organizer)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&organizer)
	return err
}

//GetOrganizer returns organizer information
func GetOrganizer(id string) (Organizer, error) {
	var organizer Organizer
	err := dbmanager.Query("ID", id, &organizer)
	if err != nil {
		return Organizer{}, err
	}
	organizer.Password = ""

	return organizer, err
}

//UpdateOrganizer updates an Organizer's information
func UpdateOrganizer(id, email, password, firstName, lastName string) error {
	hashedPwd, err := hash(password)
	if err != nil {
		return err
	}
	err = dbmanager.Update(&Organizer{ID: id, Email: email, Password: hashedPwd})
	return err
}

//AuthenticateOrganizer checks the login credentials of an Organizer
func AuthenticateOrganizer(email, pwd string) (string, error) {
	var org Organizer
	err := dbmanager.Query("Email", email, &org)
	if comparePassword(org.Password, pwd) {
		return org.ID, err
	}
	return "", err
}
