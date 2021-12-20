package attendance

import (
	"main/internal/platform/dbmanager"
	"main/internal/platform/uuid"
	"time"
	"main/internal/participants"
)

//Attendance defines the fields that Attendance has
type Attendance struct {
	ID            string `json:"attendanceID" storm:"id"`
	SessionID     string `json:"sessionID" storm:"index"`
	ParticipantName string `json:"participantName"`
	ParticipantID string `json:"particpantID" storm:"index"`
	TimeIn      string `json:"timeIn" storm:"index"`
	TimeOut     string `json:"timeOut" storm:"index"`
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
