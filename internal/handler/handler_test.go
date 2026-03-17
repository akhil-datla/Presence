package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akhil-datla/Presence/internal/auth"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
)

const testSecret = "test-jwt-secret-for-handler-tests"

// testEnv holds the shared dependencies for handler tests.
type testEnv struct {
	echo *echo.Echo
	store *store.Store
	jwt  *auth.JWTService
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	jwtSvc := auth.NewJWTService(testSecret)
	e := echo.New()

	authH := NewAuthHandler(s, jwtSvc)
	userH := NewUserHandler(s)
	sessH := NewSessionHandler(s)
	attH := NewAttendanceHandler(s)

	api := e.Group("/api/v1")

	// Public routes
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// Simple JWT middleware for tests
	jwtMW := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return echo.NewHTTPError(401, map[string]string{"error": "missing authorization header"})
			}
			token := strings.TrimPrefix(header, "Bearer ")
			if token == header {
				return echo.NewHTTPError(401, map[string]string{"error": "invalid authorization format"})
			}
			claims, err := jwtSvc.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(401, map[string]string{"error": "invalid or expired token"})
			}
			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}

	protected := api.Group("", jwtMW)

	protected.GET("/users/me", userH.GetProfile)
	protected.PUT("/users/me", userH.UpdateProfile)
	protected.DELETE("/users/me", userH.DeleteProfile)

	protected.POST("/sessions", sessH.Create)
	protected.GET("/sessions", sessH.List)
	protected.GET("/sessions/:id", sessH.Get)
	protected.PUT("/sessions/:id", sessH.Update)
	protected.DELETE("/sessions/:id", sessH.Delete)

	protected.POST("/sessions/:id/checkin", attH.CheckIn)
	protected.POST("/sessions/:id/checkout", attH.CheckOut)
	protected.GET("/sessions/:id/attendance", attH.List)
	protected.GET("/sessions/:id/attendance/filter", attH.Filter)
	protected.DELETE("/sessions/:id/attendance", attH.Clear)
	protected.GET("/sessions/:id/export/csv", attH.ExportCSV)

	return &testEnv{echo: e, store: s, jwt: jwtSvc}
}

func (te *testEnv) request(method, path, body string, token ...string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if len(token) > 0 && token[0] != "" {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	rec := httptest.NewRecorder()
	te.echo.ServeHTTP(rec, req)
	return rec
}

// registerUser is a test helper that registers a user and returns the token.
func (te *testEnv) registerUser(t *testing.T, firstName, lastName, email, password string) string {
	t.Helper()
	body := `{"first_name":"` + firstName + `","last_name":"` + lastName + `","email":"` + email + `","password":"` + password + `"}`
	rec := te.request("POST", "/api/v1/auth/register", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register failed: status=%d body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}
	return resp.Data.Token
}

// ---------- Auth Tests ----------

func TestRegister(t *testing.T) {
	env := newTestEnv(t)

	t.Run("success", func(t *testing.T) {
		body := `{"first_name":"John","last_name":"Doe","email":"john@test.com","password":"password123"}`
		rec := env.request("POST", "/api/v1/auth/register", body)
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		data := resp["data"].(map[string]interface{})
		if data["token"] == nil || data["token"] == "" {
			t.Fatal("expected token in response")
		}
		user := data["user"].(map[string]interface{})
		if user["email"] != "john@test.com" {
			t.Fatalf("expected email=john@test.com, got %v", user["email"])
		}
		// Password should not be exposed
		if user["password"] != nil && user["password"] != "" {
			t.Fatal("password should not be in response")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		body := `{"first_name":"John","last_name":"Doe","email":"john@test.com","password":"password123"}`
		rec := env.request("POST", "/api/v1/auth/register", body)
		if rec.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", rec.Code)
		}
	})

	t.Run("invalid request missing fields", func(t *testing.T) {
		body := `{"email":"a@b.com"}`
		rec := env.request("POST", "/api/v1/auth/register", body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("short password", func(t *testing.T) {
		body := `{"first_name":"A","last_name":"B","email":"ab@test.com","password":"short"}`
		rec := env.request("POST", "/api/v1/auth/register", body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		rec := env.request("POST", "/api/v1/auth/register", `{invalid}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	env := newTestEnv(t)
	env.registerUser(t, "John", "Doe", "john@test.com", "password123")

	t.Run("success", func(t *testing.T) {
		body := `{"email":"john@test.com","password":"password123"}`
		rec := env.request("POST", "/api/v1/auth/login", body)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		data := resp["data"].(map[string]interface{})
		if data["token"] == nil || data["token"] == "" {
			t.Fatal("expected token in response")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		body := `{"email":"john@test.com","password":"wrongpassword"}`
		rec := env.request("POST", "/api/v1/auth/login", body)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		body := `{"email":"nobody@test.com","password":"password123"}`
		rec := env.request("POST", "/api/v1/auth/login", body)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		body := `{"email":""}`
		rec := env.request("POST", "/api/v1/auth/login", body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}

// ---------- User Profile Tests ----------

func TestGetProfile(t *testing.T) {
	env := newTestEnv(t)
	token := env.registerUser(t, "Jane", "Smith", "jane@test.com", "password123")

	t.Run("authenticated", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/users/me", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		data := resp["data"].(map[string]interface{})
		if data["email"] != "jane@test.com" {
			t.Fatalf("expected email=jane@test.com, got %v", data["email"])
		}
	})

	t.Run("unauthenticated", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/users/me", "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})
}

func TestUpdateProfile(t *testing.T) {
	env := newTestEnv(t)
	token := env.registerUser(t, "Jane", "Smith", "jane2@test.com", "password123")

	t.Run("update first name", func(t *testing.T) {
		body := `{"first_name":"Janet"}`
		rec := env.request("PUT", "/api/v1/users/me", body, token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["first_name"] != "Janet" {
			t.Fatalf("expected first_name=Janet, got %v", data["first_name"])
		}
	})

	t.Run("invalid email format", func(t *testing.T) {
		body := `{"email":"bademail"}`
		rec := env.request("PUT", "/api/v1/users/me", body, token)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}

func TestDeleteProfile(t *testing.T) {
	env := newTestEnv(t)
	token := env.registerUser(t, "Del", "User", "del@test.com", "password123")

	t.Run("success", func(t *testing.T) {
		rec := env.request("DELETE", "/api/v1/users/me", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		// Profile should no longer be accessible
		rec = env.request("GET", "/api/v1/users/me", "", token)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 after deletion, got %d", rec.Code)
		}
	})
}

// ---------- Session Tests ----------

func TestSessionCRUD(t *testing.T) {
	env := newTestEnv(t)
	token := env.registerUser(t, "Org", "Anizer", "org@test.com", "password123")

	var sessionID string

	t.Run("create session", func(t *testing.T) {
		body := `{"name":"Morning Standup"}`
		rec := env.request("POST", "/api/v1/sessions", body, token)
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		sessionID = data["id"].(string)
		if data["name"] != "Morning Standup" {
			t.Fatalf("expected name=Morning Standup, got %v", data["name"])
		}
	})

	t.Run("create session invalid input", func(t *testing.T) {
		body := `{"name":""}`
		rec := env.request("POST", "/api/v1/sessions", body, token)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("list sessions", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].([]interface{})
		if len(data) != 1 {
			t.Fatalf("expected 1 session, got %d", len(data))
		}
	})

	t.Run("get session", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID, "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["name"] != "Morning Standup" {
			t.Fatalf("expected name=Morning Standup, got %v", data["name"])
		}
	})

	t.Run("get non-existent session", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/nonexistent", "", token)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rec.Code)
		}
	})

	t.Run("update session", func(t *testing.T) {
		body := `{"name":"Afternoon Standup"}`
		rec := env.request("PUT", "/api/v1/sessions/"+sessionID, body, token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["name"] != "Afternoon Standup" {
			t.Fatalf("expected name=Afternoon Standup, got %v", data["name"])
		}
	})

	t.Run("delete session", func(t *testing.T) {
		rec := env.request("DELETE", "/api/v1/sessions/"+sessionID, "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify it's gone
		rec = env.request("GET", "/api/v1/sessions/"+sessionID, "", token)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 after deletion, got %d", rec.Code)
		}
	})
}

func TestSessionProtectedRoutes(t *testing.T) {
	env := newTestEnv(t)

	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/sessions"},
		{"GET", "/api/v1/sessions"},
		{"GET", "/api/v1/sessions/someid"},
		{"PUT", "/api/v1/sessions/someid"},
		{"DELETE", "/api/v1/sessions/someid"},
	}

	for _, r := range routes {
		t.Run(r.method+" "+r.path+" without auth", func(t *testing.T) {
			rec := env.request(r.method, r.path, "")
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", rec.Code)
			}
		})
	}
}

// ---------- Attendance Tests ----------

func TestAttendanceFlow(t *testing.T) {
	env := newTestEnv(t)
	token := env.registerUser(t, "Alice", "Wonder", "alice@test.com", "password123")

	// Create a session first
	body := `{"name":"Attendance Test"}`
	rec := env.request("POST", "/api/v1/sessions", body, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to create session: %d %s", rec.Code, rec.Body.String())
	}
	var sessResp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &sessResp)
	sessionID := sessResp["data"].(map[string]interface{})["id"].(string)

	t.Run("check in", func(t *testing.T) {
		rec := env.request("POST", "/api/v1/sessions/"+sessionID+"/checkin", "", token)
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["participant_name"] != "Alice Wonder" {
			t.Fatalf("expected participant_name=Alice Wonder, got %v", data["participant_name"])
		}
		if data["time_out"] != nil {
			t.Fatal("expected time_out to be nil on check-in")
		}
	})

	t.Run("check out", func(t *testing.T) {
		rec := env.request("POST", "/api/v1/sessions/"+sessionID+"/checkout", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["time_out"] == nil || data["time_out"] == "" {
			t.Fatal("expected time_out to be set after check-out")
		}
	})

	t.Run("get attendance", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID+"/attendance", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].([]interface{})
		if len(data) != 1 {
			t.Fatalf("expected 1 attendance record, got %d", len(data))
		}
	})

	t.Run("export csv", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID+"/export/csv", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/csv") {
			t.Fatalf("expected Content-Type text/csv, got %s", contentType)
		}

		csvBody := rec.Body.String()
		if !strings.Contains(csvBody, "ParticipantID") {
			t.Fatal("expected CSV header in response")
		}
		if !strings.Contains(csvBody, "Alice Wonder") {
			t.Fatal("expected participant name in CSV")
		}
	})

	t.Run("clear attendance", func(t *testing.T) {
		rec := env.request("DELETE", "/api/v1/sessions/"+sessionID+"/attendance", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify cleared
		rec = env.request("GET", "/api/v1/sessions/"+sessionID+"/attendance", "", token)
		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].([]interface{})
		if len(data) != 0 {
			t.Fatalf("expected 0 records after clear, got %d", len(data))
		}
	})
}

func TestAttendanceProtectedRoutes(t *testing.T) {
	env := newTestEnv(t)

	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/sessions/someid/checkin"},
		{"POST", "/api/v1/sessions/someid/checkout"},
		{"GET", "/api/v1/sessions/someid/attendance"},
		{"DELETE", "/api/v1/sessions/someid/attendance"},
		{"GET", "/api/v1/sessions/someid/export/csv"},
	}

	for _, r := range routes {
		t.Run(r.method+" "+r.path+" without auth", func(t *testing.T) {
			rec := env.request(r.method, r.path, "")
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", rec.Code)
			}
		})
	}
}

func TestAttendanceFilterEndpoint(t *testing.T) {
	env := newTestEnv(t)
	token := env.registerUser(t, "Bob", "Filter", "bob@test.com", "password123")

	// Create session and check in
	rec := env.request("POST", "/api/v1/sessions", `{"name":"Filter Test"}`, token)
	var sessResp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &sessResp)
	sessionID := sessResp["data"].(map[string]interface{})["id"].(string)

	env.request("POST", "/api/v1/sessions/"+sessionID+"/checkin", "", token)

	t.Run("missing query params", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID+"/attendance/filter", "", token)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid time format", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID+"/attendance/filter?mode=after&time=badtime", "", token)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("filter after past time", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID+"/attendance/filter?mode=after&time=2000-01-01T00:00:00Z", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		data := resp["data"].([]interface{})
		if len(data) != 1 {
			t.Fatalf("expected 1 record, got %d", len(data))
		}
	})

	t.Run("filter before past time returns nothing", func(t *testing.T) {
		rec := env.request("GET", "/api/v1/sessions/"+sessionID+"/attendance/filter?mode=before&time=2000-01-01T00:00:00Z", "", token)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		// data could be null (nil slice) or empty array
		if resp["data"] != nil {
			data, ok := resp["data"].([]interface{})
			if ok && len(data) != 0 {
				t.Fatalf("expected 0 records, got %d", len(data))
			}
		}
	})
}
