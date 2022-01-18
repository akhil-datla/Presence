package apiserver

import (
	"main/internal/attendance"
	"main/internal/users"
	"main/internal/sessions"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	ID, err := users.AddUser(c.FormValue("firstName"), c.FormValue("lastName"), c.FormValue("email"), c.FormValue("password"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, ID)
}

func AuthenticateUser(c echo.Context) error {
	ID, err := users.AuthenticateUser(c.FormValue("email"), c.FormValue("password"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, ID)
}

func GetUser(c echo.Context) error {
	user, err := users.GetUser(c.FormValue("ID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	err := users.UpdateUser(c.FormValue("ID"), c.FormValue("email"), c.FormValue("password"), c.FormValue("firstName"), c.FormValue("lastName"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Updated User")
}

func DeleteUser(c echo.Context) error {
	err := users.RemoveUser(c.FormValue("ID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Deleted User")
}

func CreateSession(c echo.Context) error {
	ID, err := sessions.AddSession(c.FormValue("ID"), c.FormValue("name"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, ID)
}

func GetSessions(c echo.Context) error {
	ses, err := sessions.GetSessions(c.FormValue("ID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, ses)
}

func UpdateSession(c echo.Context) error {
	err := sessions.UpdateSession(c.FormValue("sesID"), c.FormValue("ID"), c.FormValue("name"))
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
	err := attendance.CheckIn(c.FormValue("sesID"), c.FormValue("ID"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Checked In")
}

func CheckOut(c echo.Context) error {
	err := attendance.CheckOut(c.FormValue("sesID"), c.FormValue("ID"))
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
