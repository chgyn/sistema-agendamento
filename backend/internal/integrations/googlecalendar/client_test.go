package googlecalendar_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/integrations/googlecalendar"
)

func TestBuildAuthURL(t *testing.T) {
	client := googlecalendar.NewClient()
	raw := client.BuildAuthURL("client-123", "https://app.test.com/api/v1/integrations/google-calendar/callback", "state-xyz")

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("URL inválida gerada: %v", err)
	}
	q := u.Query()
	if q.Get("client_id") != "client-123" {
		t.Errorf("client_id inesperado: %s", q.Get("client_id"))
	}
	if q.Get("redirect_uri") != "https://app.test.com/api/v1/integrations/google-calendar/callback" {
		t.Errorf("redirect_uri inesperada: %s", q.Get("redirect_uri"))
	}
	if q.Get("access_type") != "offline" || q.Get("prompt") != "consent" || q.Get("include_granted_scopes") != "true" {
		t.Errorf("parâmetros OAuth esperados ausentes: %v", q)
	}
	if q.Get("state") != "state-xyz" {
		t.Errorf("state inesperado: %s", q.Get("state"))
	}
	scope := q.Get("scope")
	for _, reqScope := range []string{
		"https://www.googleapis.com/auth/calendar.events",
		"https://www.googleapis.com/auth/calendar.calendarlist.readonly",
		"https://www.googleapis.com/auth/userinfo.email",
	} {
		if !strings.Contains(scope, reqScope) {
			t.Errorf("escopo mínimo ausente: %s", reqScope)
		}
	}
}

func TestExchangeCodeAndRefreshToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		vals, _ := url.ParseQuery(string(bodyBytes))

		switch vals.Get("grant_type") {
		case "authorization_code":
			if vals.Get("code") != "auth-code-1" || vals.Get("client_id") != "cid" || vals.Get("client_secret") != "csec" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(domain.OAuthTokenResponse{
				AccessToken:  "access-1",
				RefreshToken: "refresh-1",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
			})
		case "refresh_token":
			if vals.Get("refresh_token") == "revoked-token" {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error":             "invalid_grant",
					"error_description": "Token has been expired or revoked.",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(domain.OAuthTokenResponse{
				AccessToken: "access-refreshed",
				ExpiresIn:   3600,
			})
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	client := googlecalendar.NewClientWithEndpoints(googlecalendar.Endpoints{
		OAuthTokenURL: srv.URL,
	})
	ctx := context.Background()

	tok, err := client.ExchangeCode(ctx, "cid", "csec", "https://cb.test", "auth-code-1")
	if err != nil {
		t.Fatalf("ExchangeCode falhou: %v", err)
	}
	if tok.AccessToken != "access-1" || tok.RefreshToken != "refresh-1" {
		t.Fatalf("tokens inesperados: %+v", tok)
	}

	refreshed, err := client.RefreshToken(ctx, "cid", "csec", "refresh-1")
	if err != nil {
		t.Fatalf("RefreshToken falhou: %v", err)
	}
	if refreshed.AccessToken != "access-refreshed" {
		t.Fatalf("access token renovado inesperado: %s", refreshed.AccessToken)
	}

	_, err = client.RefreshToken(ctx, "cid", "csec", "revoked-token")
	if !errors.Is(err, domain.ErrGoogleCalendarAuthExpired) {
		t.Fatalf("esperava ErrGoogleCalendarAuthExpired para invalid_grant, obteve: %v", err)
	}
}

func TestRevokeUserInfoAndListCalendars(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/revoke":
			bodyBytes, _ := io.ReadAll(r.Body)
			vals, _ := url.ParseQuery(string(bodyBytes))
			if vals.Get("token") != "tok-to-revoke" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "/userinfo":
			if r.Header.Get("Authorization") != "Bearer acc-123" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"email": "agenda@barbearia.com"})
		case "/users/me/calendarList":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{"id": "primary", "summary": "Principal", "primary": true, "accessRole": "owner"},
					{"id": "team-cal", "summary": "Equipe", "primary": false, "accessRole": "writer"},
					{"id": "holidays", "summary": "Feriados", "primary": false, "accessRole": "reader"},
					{"id": "freebusy", "summary": "Bloqueios", "primary": false, "accessRole": "freeBusyReader"},
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := googlecalendar.NewClientWithEndpoints(googlecalendar.Endpoints{
		OAuthRevokeURL:     srv.URL + "/revoke",
		UserInfoURL:        srv.URL + "/userinfo",
		CalendarAPIBaseURL: srv.URL,
	})
	ctx := context.Background()

	if err := client.RevokeToken(ctx, "tok-to-revoke"); err != nil {
		t.Fatalf("RevokeToken falhou: %v", err)
	}

	email, err := client.GetUserEmail(ctx, "acc-123")
	if err != nil {
		t.Fatalf("GetUserEmail falhou: %v", err)
	}
	if email != "agenda@barbearia.com" {
		t.Fatalf("email inesperado: %s", email)
	}

	cals, err := client.ListCalendars(ctx, "acc-123")
	if err != nil {
		t.Fatalf("ListCalendars falhou: %v", err)
	}
	if len(cals) != 2 {
		t.Fatalf("esperava 2 calendários graváveis (owner/writer), obteve %d: %+v", len(cals), cals)
	}
}

func TestEventsCRUD_AndIdempotentStatusCodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("sendUpdates") != "none" {
			t.Errorf("sendUpdates=none obrigatório, obteve %q", r.URL.Query().Get("sendUpdates"))
		}

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/calendars/primary/events":
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["id"] == "conflict-id" {
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(`{"error":{"code":409,"message":"already exists"}}`))
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"id": body["id"].(string)})

		case r.Method == http.MethodPatch && r.URL.Path == "/calendars/primary/events/evt-1":
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "evt-1"})

		case r.Method == http.MethodPatch && r.URL.Path == "/calendars/primary/events/missing-404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"code":404,"message":"not found"}}`))

		case r.Method == http.MethodPatch && r.URL.Path == "/calendars/primary/events/gone-410":
			w.WriteHeader(http.StatusGone)
			_, _ = w.Write([]byte(`{"error":{"code":410,"message":"Resource has been deleted"}}`))

		case r.Method == http.MethodDelete && r.URL.Path == "/calendars/primary/events/evt-1":
			w.WriteHeader(http.StatusNoContent)

		case r.Method == http.MethodDelete && r.URL.Path == "/calendars/primary/events/missing-404":
			w.WriteHeader(http.StatusNotFound)

		case r.Method == http.MethodDelete && r.URL.Path == "/calendars/primary/events/gone-410":
			w.WriteHeader(http.StatusGone)

		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	client := googlecalendar.NewClientWithEndpoints(googlecalendar.Endpoints{
		CalendarAPIBaseURL: srv.URL,
	})
	ctx := context.Background()
	input := domain.CalendarEventInput{
		EventID: "af123456",
		Summary: "Corte Masculino - João",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(30 * time.Minute),
	}

	id, err := client.CreateEvent(ctx, "tok", "primary", input)
	if err != nil || id != "af123456" {
		t.Fatalf("CreateEvent falhou: id=%s err=%v", id, err)
	}

	conflictInput := input
	conflictInput.EventID = "conflict-id"
	_, err = client.CreateEvent(ctx, "tok", "primary", conflictInput)
	if !errors.Is(err, googlecalendar.ErrEventConflict) {
		t.Fatalf("esperava ErrEventConflict para 409, obteve %v", err)
	}

	updatedID, err := client.UpdateEvent(ctx, "tok", "primary", "evt-1", input)
	if err != nil || updatedID != "evt-1" {
		t.Fatalf("UpdateEvent falhou: id=%s err=%v", updatedID, err)
	}

	if _, err := client.UpdateEvent(ctx, "tok", "primary", "missing-404", input); !errors.Is(err, googlecalendar.ErrEventNotFound) {
		t.Fatalf("esperava ErrEventNotFound para 404, obteve %v", err)
	}
	if _, err := client.UpdateEvent(ctx, "tok", "primary", "gone-410", input); !errors.Is(err, googlecalendar.ErrEventNotFound) {
		t.Fatalf("esperava ErrEventNotFound para 410, obteve %v", err)
	}

	if err := client.DeleteEvent(ctx, "tok", "primary", "evt-1"); err != nil {
		t.Fatalf("DeleteEvent 204 falhou: %v", err)
	}
	if err := client.DeleteEvent(ctx, "tok", "primary", "missing-404"); err != nil {
		t.Fatalf("DeleteEvent 404 deveria ser idempotente, obteve: %v", err)
	}
	if err := client.DeleteEvent(ctx, "tok", "primary", "gone-410"); err != nil {
		t.Fatalf("DeleteEvent 410 deveria ser idempotente, obteve: %v", err)
	}
}
