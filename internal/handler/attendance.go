package handler

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"time"

	"github.com/akhil-datla/Presence/internal/model"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
)

// AttendanceHandler handles attendance endpoints.
type AttendanceHandler struct {
	store *store.Store
}

// NewAttendanceHandler creates a new attendance handler.
func NewAttendanceHandler(s *store.Store) *AttendanceHandler {
	return &AttendanceHandler{store: s}
}

// CheckIn records a participant checking in.
func (h *AttendanceHandler) CheckIn(c echo.Context) error {
	sessionID := c.Param("id")
	uid := userID(c)

	// Look up user for their name
	user, err := h.store.GetUserByID(uid)
	if err != nil {
		return notFound(c, "user not found")
	}

	name := user.FirstName + " " + user.LastName
	att, err := h.store.CheckIn(sessionID, uid, name)
	if err != nil {
		return badRequest(c, "failed to check in: "+err.Error())
	}
	return created(c, att)
}

// CheckOut records a participant checking out.
func (h *AttendanceHandler) CheckOut(c echo.Context) error {
	att, err := h.store.CheckOut(c.Param("id"), userID(c))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(c, "no open check-in found")
		}
		return badRequest(c, "failed to check out: "+err.Error())
	}
	return ok(c, att)
}

// List returns attendance records for a session.
func (h *AttendanceHandler) List(c echo.Context) error {
	records, err := h.store.GetAttendance(c.Param("id"))
	if err != nil {
		return serverError(c)
	}
	if records == nil {
		records = []model.Attendance{}
	}
	return ok(c, records)
}

// Filter returns filtered attendance records.
func (h *AttendanceHandler) Filter(c echo.Context) error {
	sessionID := c.Param("id")
	mode := c.QueryParam("mode")
	timeStr := c.QueryParam("time")

	if mode == "" || timeStr == "" {
		return badRequest(c, "query parameters 'mode' and 'time' are required")
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return badRequest(c, "invalid time format, use RFC3339 (e.g. 2024-01-15T09:00:00Z)")
	}

	records, err := h.store.FilterAttendance(sessionID, t, mode)
	if err != nil {
		return badRequest(c, err.Error())
	}
	return ok(c, records)
}

// Clear deletes all attendance records for a session.
func (h *AttendanceHandler) Clear(c echo.Context) error {
	err := h.store.ClearAttendance(c.Param("id"), userID(c))
	if err != nil {
		return notFound(c, "session not found or not owned by you")
	}
	return ok(c, map[string]string{"message": "attendance records cleared"})
}

// ExportCSV streams attendance records as a CSV download.
func (h *AttendanceHandler) ExportCSV(c echo.Context) error {
	records, err := h.store.GetAttendance(c.Param("id"))
	if err != nil {
		return serverError(c)
	}

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="attendance_%s.csv"`, time.Now().Format("2006-01-02")))
	c.Response().WriteHeader(200)

	w := csv.NewWriter(c.Response())
	_ = w.Write([]string{"ParticipantID", "ParticipantName", "TimeIn", "TimeOut"})

	for _, r := range records {
		timeOut := ""
		if r.TimeOut != nil {
			timeOut = r.TimeOut.Format(time.RFC3339)
		}
		_ = w.Write([]string{r.ParticipantID, r.ParticipantName, r.TimeIn.Format(time.RFC3339), timeOut})
	}

	w.Flush()
	return nil
}
