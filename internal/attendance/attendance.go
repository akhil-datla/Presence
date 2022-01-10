package attendance

import (
	"encoding/csv"
	"fmt"
	"main/internal/participants"
	"main/internal/platform/dbmanager"
	"main/internal/platform/uuid"
	"main/internal/sessions"
	"os"
	"time"
)

//Attendance defines the fields that Attendance has
type Attendance struct {
	ID              string `json:"attendanceID" storm:"id"`
	SessionID       string `json:"sessionID" storm:"index"`
	ParticipantName string `json:"participantName"`
	ParticipantID   string `json:"particpantID" storm:"index"`
	TimeIn          string `json:"timeIn" storm:"index"`
	TimeOut         string `json:"timeOut" storm:"index"`
}

//New creates a new bucket for Attendance and creates a pointer to Attendance
func New() {
	dbmanager.CreateBucket(&Attendance{})
}

//CheckIn checks in a Participant given the SessionID and ParticipantID
func CheckIn(sessionID string, participantID string) error {
	attendance := &Attendance{}
	attendance.SessionID = sessionID
	attendance.ParticipantID = participantID
	var participant participants.Participant
	err := dbmanager.Query("ID", participantID, &participant)

	if err != nil {
		return err
	}

	attendance.ParticipantName = participant.FirstName + " " + participant.LastName

	attendance.ID = uuid.New()
	attendance.TimeIn = time.Now().Format(time.RFC3339)
	err = dbmanager.Save(attendance)
	return err

}

//CheckOut checks out a Participant given the SessionID and ParticipantID
func CheckOut(sessionID, participantID string) error {
	var attendanceRecords []*Attendance
	err := dbmanager.GroupQuery("SessionID", sessionID, &attendanceRecords)
	if err != nil {
		return err
	}
	for _, attendance := range attendanceRecords {
		if attendance.ParticipantID == participantID {
			attendance.TimeOut = time.Now().Format(time.RFC3339)
			err = dbmanager.Save(attendance)
			return err
		}
	}
	return nil
}

func ClearAttendance(sessionID string) error {
	var attendanceList []*Attendance
	err := dbmanager.GroupQuery("SessionID", sessionID, &attendanceList)

	for _, att := range attendanceList {
		dbmanager.Delete(&att)
	}

	return err
}

//GetAttendance gets the Attendance records given the SessionID
func GetAttendance(sessionID string) ([]*Attendance, error) {
	var attendanceList []*Attendance
	err := dbmanager.GroupQuery("SessionID", sessionID, &attendanceList)
	return attendanceList, err
}

// FilterAttendance gets the Attendance Records for
// a given sessionID, time in IEEE 8601 format as a string, and the mode either as "before" or "after"
func FilterAttendance(sessionID, timeIn, mode string) ([]*Attendance, error) {
	var attendanceList []*Attendance
	var filteredAttendanceList []*Attendance
	dbmanager.GroupQuery("SessionID", sessionID, &attendanceList)

	givenTime, err := time.Parse(time.RFC3339, timeIn)
	if err != nil {
		return nil, err
	}
	switch {
	case mode == "before":
		for _, attendance := range attendanceList {
			recordedTime, err := time.Parse(time.RFC3339, attendance.TimeIn)
			if err != nil {
				return nil, err
			}
			if recordedTime.Before(givenTime) {
				filteredAttendanceList = append(filteredAttendanceList, attendance)
			}
		}
	case mode == "after":
		for _, attendance := range attendanceList {
			recordedTime, err := time.Parse(time.RFC3339, attendance.TimeIn)
			if err != nil {
				return nil, err
			}
			if recordedTime.After(givenTime) {
				filteredAttendanceList = append(filteredAttendanceList, attendance)
			}
		}
	default:
		return nil, err
	}
	return filteredAttendanceList, nil
}

//GenerateCSV generates a CSV file of the Attendance records given the SessionID and returns the name of the file
func GenerateCSV(sessionID string) (string, error) {

	var session sessions.Session
	err := dbmanager.Query("ID", sessionID, &session)

	if err != nil {
		return "", err
	}

	var attendanceList []*Attendance
	dbmanager.GroupQuery("SessionID", sessionID, &attendanceList)

	csvName := fmt.Sprintf("%s.csv", time.Now().Format("2006-01-02:15:04:05"))

	csvFile, err := os.Create(csvName)

	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	categoryNames := []string{"ParticpantID", "ParticipantName", "Time In", "Time Out"}
	writer.Write(categoryNames)

	for _, attendance := range attendanceList {
		var row []string
		row = append(row, attendance.ParticipantID)
		row = append(row, attendance.ParticipantName)
		row = append(row, attendance.TimeIn)
		row = append(row, attendance.TimeOut)
		writer.Write(row)
	}

	writer.Flush()

	return csvName, nil
}
