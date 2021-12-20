package sessions

import (
	"encoding/json"
	"main/internal/platform/dbmanager"
	"main/internal/platform/uuid"
)

//Session defines the fields a session has
type Session struct {
	ID string `json:"sessionID" storm:"id"`
	OrganizerID string `json:"organizerID" storm:"index"` 
	Name string `json:"sessionName"` 
}

//New creates a new bucket for Session and creates a pointer to Sessions
func New() {
	dbmanager.CreateBucket(&Session{})
}

//AddSession adds a Session to the database given session information
func AddSession(orgID, name string) error {
	session := &Session{}
	session.OrganizerID = orgID
	session.Name = name
	var err error
	session.ID = uuid.New()
	err = dbmanager.Save(session)
	return err
}

//RemoveSession removes a Session
func RemoveSession(id string) error {
	var session Session
	err := dbmanager.Query("ID", id, &session)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&session)
	return err
}

//UpdateSession updates a Session's information
func UpdateSession(id string, orgID, name string) error {
	err := dbmanager.Update(&Session{ID: id, OrganizerID: orgID, Name: name})
	return err
}

func GetSessions(orgID string) (string, error) {
	var ses []Session
	err := dbmanager.GroupQuery("OrganizerID", orgID, &ses)
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(&ses)
	return string(bytes), err
}