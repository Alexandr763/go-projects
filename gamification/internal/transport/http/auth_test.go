package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"gamification/internal/transport/http/apigen"
	"github.com/golang-jwt/jwt/v5"
)

func TestRegisterLoginProfile(t *testing.T) {
	handler := testHandler(t)
	session := register(t, handler, "Owner@Example.com", "longenough")
	if session.Email != "owner@example.com" || session.Coins != 0 || session.Access == "" || session.Refresh == "" {
		t.Fatalf("session = %+v", session)
	}
	if strings.Contains(mustRaw(t, handler, http.MethodPost, "/api/v1/auth/login", `{"email":"missing@example.com","password":"longenough"}`, ""), "password_hash") {
		t.Fatal("hash leaked")
	}

	unknown := call(handler, http.MethodPost, "/api/v1/auth/login", `{"email":"missing@example.com","password":"longenough"}`, "", "")
	wrong := call(handler, http.MethodPost, "/api/v1/auth/login", `{"email":"owner@example.com","password":"wrongpass"}`, "", "")
	if unknown.Code != http.StatusUnauthorized || unknown.Body.String() != wrong.Body.String() {
		t.Fatalf("unknown=%d %s wrong=%d %s", unknown.Code, unknown.Body.String(), wrong.Code, wrong.Body.String())
	}

	again := register(t, handler, "other@example.com", "longenough")
	profile := call(handler, http.MethodGet, "/api/v1/profile", "", "Bearer "+session.Access, "")
	if profile.Code != http.StatusOK {
		t.Fatalf("profile %d %s", profile.Code, profile.Body.String())
	}
	var got apigen.Profile
	if err := json.Unmarshal(profile.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Id != session.Id || got.Email != "owner@example.com" || strings.Contains(profile.Body.String(), "password") {
		t.Fatalf("profile body = %s", profile.Body.String())
	}
	if got.Id == again.Id {
		t.Fatal("profile returned the other user")
	}

	dup := call(handler, http.MethodPost, "/api/v1/auth/register", `{"email":"owner@example.com","password":"longenough"}`, "", "")
	if dup.Code != http.StatusConflict || strings.Contains(dup.Body.String(), "owner@example.com") {
		t.Fatalf("duplicate %d %s", dup.Code, dup.Body.String())
	}
}

func TestRefreshAndLogout(t *testing.T) {
	handler := testHandler(t)
	session := register(t, handler, "fresh@example.com", "longenough")

	rotated := call(handler, http.MethodPost, "/api/v1/auth/refresh", `{"refresh":"`+session.Refresh+`"}`, "", "")
	if rotated.Code != http.StatusOK {
		t.Fatalf("refresh %d %s", rotated.Code, rotated.Body.String())
	}
	var next apigen.Session
	if err := json.Unmarshal(rotated.Body.Bytes(), &next); err != nil {
		t.Fatal(err)
	}
	old := call(handler, http.MethodPost, "/api/v1/auth/refresh", `{"refresh":"`+session.Refresh+`"}`, "", "")
	if old.Code != http.StatusUnauthorized {
		t.Fatalf("old refresh %d", old.Code)
	}

	logout := call(handler, http.MethodPost, "/api/v1/auth/logout", `{"refresh":"`+next.Refresh+`"}`, "Bearer "+next.Access, "")
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout %d %s", logout.Code, logout.Body.String())
	}
	after := call(handler, http.MethodPost, "/api/v1/auth/refresh", `{"refresh":"`+next.Refresh+`"}`, "", "")
	if after.Code != http.StatusUnauthorized {
		t.Fatalf("refresh after logout %d", after.Code)
	}
}

func TestRejectedTokens(t *testing.T) {
	handler := testHandler(t)
	session := register(t, handler, "token@example.com", "longenough")

	other, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   session.Id.String(),
		Issuer:    "gamification",
		Audience:  jwt.ClaimStrings{"gamification"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte("other-secret-other-secret-other"))
	if err != nil {
		t.Fatal(err)
	}
	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   session.Id.String(),
		Issuer:    "gamification",
		Audience:  jwt.ClaimStrings{"gamification"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	}).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(session.Access, ".")
	parts[2] = parts[2][:len(parts[2])-1] + "x"
	tampered := strings.Join(parts, ".")
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"` + session.Id.String() + `","iss":"gamification","aud":["gamification"],"exp":9999999999}`))
	none := header + "." + payload + "."

	for _, token := range []string{other, expired, tampered, none} {
		rec := call(handler, http.MethodGet, "/api/v1/profile", "", "Bearer "+token, "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("token status = %d body = %s", rec.Code, rec.Body.String())
		}
	}
}

func TestLoginRateLimit(t *testing.T) {
	handler := testHandler(t)
	body := `{"email":"limit@example.com","password":"longenough"}`
	for i := 0; i < 10; i++ {
		rec := call(handler, http.MethodPost, "/api/v1/auth/login", body, "", "")
		if rec.Code == http.StatusTooManyRequests {
			t.Fatalf("limited on attempt %d", i+1)
		}
	}
	rec := call(handler, http.MethodPost, "/api/v1/auth/login", body, "", "")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("11th status = %d", rec.Code)
	}
	bad := call(handler, http.MethodPost, "/api/v1/auth/login", `{`, "", "")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("bad json after limit = %d", bad.Code)
	}
}

func mustRaw(t *testing.T, handler http.Handler, method, path, body, auth string) string {
	t.Helper()
	return call(handler, method, path, body, auth, "").Body.String()
}
