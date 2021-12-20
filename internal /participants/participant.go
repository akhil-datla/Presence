package participants

import (
	"main/internal/platform/dbmanager"
	"main/internal/platform/uuid"

	"golang.org/x/crypto/bcrypt"
)

//Participant defines the fields participant has
type Participant struct {
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

//New creates a new bucket for Participant
func New() {
	dbmanager.CreateBucket(&Participant{})
}

//AddParticipant adds an Participant to the database given participant information
func AddParticipant(firstname string, lastname string, em string, pwd string) (string, error) {
	participant := &Participant{}
	participant.FirstName = firstname
	participant.LastName = lastname
	participant.Email = em
	var err error
	participant.Password, err = hash(pwd)
	if err != nil {
		return "", err
	}
	participant.ID = uuid.New()
	err = dbmanager.Save(participant)
	if err != nil {
		return "", err
	}
	return participant.ID, err
}

//GetParticipant returns participant information
func GetParticipant(id string) (Participant, error) {
	var participant Participant
	err := dbmanager.Query("ID", id, &participant)
	if err != nil {
		return Participant{}, err
	}
	participant.Password = ""
	return participant, err
}

//RemoveParticipant removes an Participant
func RemoveParticipant(id string) error {
	var participant Participant
	err := dbmanager.Query("ID", id, &participant)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&participant)
	return err
}

//UpdateParticipant updates an Participant's information
func UpdateParticipant(id, email, password, firstName, lastName string) error {
	hashedPwd, err := hash(password)
	if err != nil {
		return err
	}
	err = dbmanager.Update(&Participant{ID: id, Email: email, Password: hashedPwd})
	return err
}

//AuthenticateParticipant checks the login credentials of an Participant
func AuthenticateParticipant(email, pwd string) (string, error) {
	var par Participant
	err := dbmanager.Query("Email", email, &par)
	if comparePassword(par.Password, pwd) {
		return par.ID, err
	}
	return "", err
}
