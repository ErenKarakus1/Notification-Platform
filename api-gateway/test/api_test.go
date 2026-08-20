package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8084"

func TestRegister(t *testing.T) {
	body := map[string]string{
		"email":    "gatewaytest@example.com",
		"password": "password123",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		baseURL+"/auth/register",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 201 or 409, got %d", resp.StatusCode)
	}
}

func TestLogin(t *testing.T) {
	body := map[string]string{
		"email":    "gatewaytest@example.com",
		"password": "password123",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestNotificationsWithoutToken(t *testing.T) {
	resp, err := http.Post(
		baseURL+"/notifications",
		"application/json",
		bytes.NewBuffer([]byte(`{}`)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestNotificationsWithInvalidToken(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/00000000-0000-0000-0000-000000000000",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLoginWithWrongPassword(t *testing.T) {
	body := map[string]string{
		"email":    "gatewaytest@example.com",
		"password": "wrong-password",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLoginWithUnknownEmail(t *testing.T) {
	body := map[string]string{
		"email":    "does-not-exist@example.com",
		"password": "password123",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLoginWithInvalidBody(t *testing.T) {
	resp, err := http.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer([]byte(`invalid json`)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRegisterWithInvalidBody(t *testing.T) {
	resp, err := http.Post(
		baseURL+"/auth/register",
		"application/json",
		bytes.NewBuffer([]byte(`invalid json`)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestNotificationsWithMalformedAuthorizationHeader(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/00000000-0000-0000-0000-000000000000",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Basic something")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestNotificationsWithEmptyBearerToken(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/00000000-0000-0000-0000-000000000000",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer ")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestNotificationsWithGarbageToken(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/00000000-0000-0000-0000-000000000000",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer abc.def.ghi")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestNotificationsInvalidID(t *testing.T) {
	token := getTestToken(t)

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/not-a-uuid",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func getTestToken(t *testing.T) string {
	t.Helper()

	body := map[string]string{
		"email":    "gatewaytest@example.com",
		"password": "password123",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected login 200, got %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if result.Token == "" {
		t.Fatal("expected token")
	}

	return result.Token
}

func TestCreateNotification(t *testing.T) {
	token := getTestToken(t)

	body := map[string]string{
		"recipient": "gatewaytest@example.com",
		"channel":   "email",
		"subject":   "Gateway test",
		"body":      "Testing notification creation",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/notifications",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestNotificationHappyPath(t *testing.T) {
	token := getTestToken(t)

	// Create notification
	body := map[string]string{
		"recipient": "gatewaytest@example.com",
		"channel":   "email",
		"subject":   "Gateway test",
		"body":      "Testing notification creation",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/notifications",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	if created.ID == "" {
		t.Fatal("expected notification ID")
	}

	// Get notification
	getReq, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/"+created.ID,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	getReq.Header.Set("Authorization", "Bearer "+token)

	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	t.Logf("notification created and retrieved successfully: %s", created.ID)
}

func TestFullNotificationFlow(t *testing.T) {
	token := getTestToken(t)

	body := map[string]string{
		"recipient": "gatewaytest@example.com",
		"channel":   "email",
		"subject":   "Full flow test",
		"body":      "Testing the complete gateway flow",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/notifications",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var notification struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&notification); err != nil {
		t.Fatal(err)
	}

	if notification.ID == "" {
		t.Fatal("expected notification ID")
	}

	getReq, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/notifications/"+notification.ID,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	getReq.Header.Set("Authorization", "Bearer "+token)

	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}
}

func TestRateLimiter(t *testing.T) {
	token := getNewTestToken(t)

	for i := 1; i <= 101; i++ {
		req, err := http.NewRequest(
			http.MethodGet,
			baseURL+"/notifications/00000000-0000-0000-0000-000000000000",
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if i <= 100 {
			if resp.StatusCode == http.StatusTooManyRequests {
				t.Fatalf("request %d was rate limited too early", i)
			}
		}

		if i == 101 {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Fatalf(
					"expected request 101 to return 429, got %d",
					resp.StatusCode,
				)
			}
		}
	}

	t.Log("rate limiter correctly returned 429 after 100 requests")
}

func getNewTestToken(t *testing.T) string {
	t.Helper()

	email := fmt.Sprintf(
		"ratelimit-%d@example.com",
		time.Now().UnixNano(),
	)

	password := "password123"

	registerBody := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonBody, err := json.Marshal(registerBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		baseURL+"/auth/register",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected register 201, got %d", resp.StatusCode)
	}

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonBody, err = json.Marshal(loginBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected login 200, got %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	return result.Token
}
