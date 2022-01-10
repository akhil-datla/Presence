package apiserver

import (
	"main/internal/attendance"
	"main/internal/organizers"
	"main/internal/participants"
	"main/internal/sessions"
	"net/http"

	"github.com/labstack/echo/v4"
)

//AddOrganizer adds an organizer to the database
func AddOrganizer(c echo.Context) error {

	orgID, err2 := organizers.AddOrganizer(c.FormValue("firstName"), c.FormValue("lastName"), c.FormValue("email"), c.FormValue("password"))
	if err2 != nil {
		return err2
	}

	return c.String(http.StatusOK, orgID)
}

//ViewOrganizer returns the Organizer information
func ViewOrganizer(c echo.Context) error {

	org, err := organizers.GetOrganizer(c.FormValue("id"))

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, org)
}

//AuthenticateOrganizer logs in an Organizer
func AuthenticateOrganizer(c echo.Context) error {
	id, err := organizers.AuthenticateOrganizer(c.FormValue("email"), c.FormValue("password"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Error getting organizer information")
	}
	if id != "" {
		return c.String(http.StatusOK, id)
	}
	return c.String(http.StatusOK, "Invalid username and/or password")
}

//UpdateOrganizer updates an Organizer's information
func UpdateOrganizer(c echo.Context) error {

	err := organizers.UpdateOrganizer(c.FormValue("id"), c.FormValue("email"), c.FormValue("password"), c.FormValue("firstName"), c.FormValue("lastName"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Updated Organizer")
}

//DeleteOrganizer deletes an Organizer
func DeleteOrganizer(c echo.Context) error {

	err := organizers.RemoveOrganizer(c.FormValue("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "Deleted Organizer")
}

//AddParticipant adds an participant to the database
func AddParticipant(c echo.Context) error {

	parID, err2 := participants.AddParticipant(c.FormValue("firstName"), c.FormValue("lastName"), c.FormValue("email"), c.FormValue("password"))
	if err2 != nil {
		return err2
	}
	return c.String(http.StatusOK, parID)
}

//ViewParticipant returns the Participant information
func ViewParticipant(c echo.Context) error {

	par, err := participants.GetParticipant(c.FormValue("id"))

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, par)
}

//AuthenticateParticipant logs in a Participant
func AuthenticateParticipant(c echo.Context) error {
	id, err := participants.AuthenticateParticipant(c.FormValue("email"), c.FormValue("password"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Error getting participant information")
	}
	if id != "" {
		return c.String(http.StatusOK, id)
	}
	return c.String(http.StatusOK, "Invalid username and/or password")
}

//UpdateParticipant updates a Participant's information
func UpdateParticipant(c echo.Context) error {

	err := participants.UpdateParticipant(c.FormValue("id"), c.FormValue("email"), c.FormValue("password"), c.FormValue("firstName"), c.FormValue("lastName"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Updated Participant")
}

//DeleteParticipant deletes a Participant
func DeleteParticipant(c echo.Context) error {

	err := participants.RemoveParticipant(c.FormValue("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "Deleted Participant")
}

func CreateSession(c echo.Context) error {
	id, err := sessions.AddSession(c.FormValue("orgID"), c.FormValue("name"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, id)
}

func GetSessions(c echo.Context) error {
	ses, err := sessions.GetSessions(c.FormValue("orgID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, ses)
}

func UpdateSession(c echo.Context) error {
	err := sessions.UpdateSession(c.FormValue("id"), c.FormValue("orgID"), c.FormValue("name"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Updated Session")
}

func DeleteSession(c echo.Context) error {
	err := sessions.RemoveSession(c.FormValue("sesID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Deleted Session")
}

func GetAttendance(c echo.Context) error {
	attendanceList, err := attendance.GetAttendance(c.FormValue("sesID"))

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, attendanceList)

}

func FilterAttendance(c echo.Context) error {
	attendanceList, err := attendance.FilterAttendance(c.FormValue("sesID"), c.FormValue("time"), c.FormValue("mode"))

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, attendanceList)

}

func CheckIn(c echo.Context) error {
	err := attendance.CheckIn(c.FormValue("sesID"), c.FormValue("parID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Checked In")
}

func CheckOut(c echo.Context) error {
	err := attendance.CheckOut(c.FormValue("sesID"), c.FormValue("parID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Checked Out")
}

func ClearAttendance(c echo.Context) error {
	err := attendance.ClearAttendance(c.FormValue("sesID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Cleared Attendance Records")
}

func GenerateCSV(c echo.Context) error {
	fileName, err := attendance.GenerateCSV(c.FormValue("sesID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Attachment(fileName, "attendance.csv")
}
