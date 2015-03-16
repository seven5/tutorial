package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	s5 "github.com/seven5/seven5"

	"tutorial/shared"
)

const (
	TEST_PORT = 8230
)

var (
	testConfig *fresnoConfig
)

func setupForTests(port int) {
	if testConfig == nil {
		os.Setenv("PORT", fmt.Sprint(port))
		testConfig = setup()
		if port != testConfig.heroku.Port() {
			panic(fmt.Sprintf("Unable to set PORT to %d: %d", port, testConfig.heroku.Port()))
		}
		go http.ListenAndServe(fmt.Sprintf(":%d", port), testConfig.serveMux)
	}
}

func TestValidateCredentials(t *testing.T) {

	setupForTests(TEST_PORT)

	checkVerifyCredentials(t, TEST_PORT, "foo@bar.com", "baz", 401)
	checkVerifyCredentials(t, TEST_PORT, "joe@example.com", "bad2bone", 401)

	//these are created by  migration #1
	checkVerifyCredentials(t, TEST_PORT, "joe@example.com", "seekret", 200)
	checkVerifyCredentials(t, TEST_PORT, "mary@example.com", "bigseekret", 200)
}

func checkVerifyCredentials(t *testing.T, port int, username, password string, expected int) {
	var sas s5.PasswordAuthParameters
	sas.Username = username
	sas.Password = password
	sas.Op = s5.AUTH_OP_LOGIN

	rd := strings.NewReader(encodeAuth(t, &sas))
	resp, err := http.Post(testUrl(port, shared.URLGen.Auth()), "application/json", rd)
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	if resp.StatusCode != expected {
		t.Errorf("unexpected response from post (expected %d): %d, %s",
			expected, resp.StatusCode, resp.Status)
	}
}
func decodeCookieToSessionUser(t *testing.T, cookie *http.Cookie) *shared.UserRecord {
	if cookie.Name != "fresno-seven5-session" {
		t.Errorf("expected cookie to be named %s but got %s", "fresno-seven5-session", cookie.Name)
	}
	value := cookie.Value

	//go through the side door to test this, a real client can't do this because
	//this implies SERVER_SESSION_KEY
	sr, err := testConfig.sm.Find(value)
	if err != nil {
		t.Fatalf("Unable to do find on cookie value: %v", err)
	}
	if sr == nil || sr.Session == nil || sr.UniqueId != "" {
		t.Fatalf("should have received a session since we were logged in but didn't: %+v", sr)
	}
	session := sr.Session
	ud := session.UserData()
	return ud.(*shared.UserRecord)
}

func checkGetStatus(t *testing.T, client *http.Client, path string, expected int) []byte {
	resp, err := client.Get(testUrl(TEST_PORT, path))
	if err != nil {
		t.Fatalf("Unable to get to rest user record endpoint: %v", err)
	}
	if resp.StatusCode != expected {
		t.Fatalf("Unexpected status fetching %s: %s", path, resp.Status)
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body of GET repsonse: %v", err)
	}
	resp.Body.Close()
	return all
}

func checkLogoutCurrentUser(t *testing.T, client *http.Client, expected int) {
	var sas s5.PasswordAuthParameters
	sas.Op = s5.AUTH_OP_LOGOUT
	rd := strings.NewReader(encodeAuth(t, &sas))

	resp, err := client.Post(testUrl(TEST_PORT, shared.URLGen.Auth()), "application/json", rd)
	if err != nil {
		t.Fatalf("Unable to post to auth endpoint: %v", err)
	}
	if resp.StatusCode != expected {
		t.Fatalf("Did not receive expected response to logging out (expected %d): %s",
			expected, resp.Status)
	}
}

func checkLoginAsJoe(t *testing.T, client *http.Client) {
	var sas s5.PasswordAuthParameters
	sas.Username = "joe@example.com"
	sas.Password = "seekret"
	sas.Op = s5.AUTH_OP_LOGIN
	rd := strings.NewReader(encodeAuth(t, &sas))

	resp, err := client.Post(testUrl(TEST_PORT, shared.URLGen.Auth()), "application/json", rd)
	if err != nil {
		t.Fatalf("Unable to post to auth endpoint: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Did not succeed in authenticating: %s", resp.Status)
	}
}

func TestCreateNewUser(t *testing.T) {
	setupForTests(TEST_PORT)
	var ur shared.UserRecord
	ur.EmailAddr = "foo" + fmt.Sprint(time.Now().Nanosecond()) + "@bar.com"
	ur.FirstName = "foo"
	ur.LastName = "bar"
	ur.Password = "sumpin"
	ur.Admin = true //hee hee hee
	rd := strings.NewReader(encodeUserRecord(t, &ur))

	client := &http.Client{}
	resp, err := client.Post(testUrl(TEST_PORT, shared.URLGen.UserRecordResource()), "application/json", rd)
	if err != nil {
		t.Fatalf("Unable to post to user record endpoint: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("Did not succeed in creating user: %s", resp.Status)
	}
}

func TestMe(t *testing.T) {
	setupForTests(TEST_PORT)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("unable to create cookie jar: %v", err)
	}
	client := &http.Client{
		Jar: jar,
	}
	checkGetStatus(t, client, shared.URLGen.Me(), 401)
	checkLoginAsJoe(t, client)

	b := checkGetStatus(t, client, shared.URLGen.Me(), 200)
	//decode body as json
	buf := bytes.NewBuffer(b)
	dec := json.NewDecoder(buf)
	var ur shared.UserRecord
	if err := dec.Decode(&ur); err != nil {
		t.Fatalf("Unable to decode response from endpoint me: %v", err)
	}
	if ur.UserUdid != "df12ba96-71c7-436d-b8f6-2d157d5f8ff1" ||
		ur.FirstName != "Joe" ||
		ur.LastName != "Smith" ||
		ur.EmailAddr != "joe@example.com" {
		t.Errorf("Unexpected value found in body: %+v", ur)
	}
	if ur.Password != "" {
		t.Errorf("Should not have received password back from server!")
	}
}

func TestCookieOnLoginLogout(t *testing.T) {
	setupForTests(TEST_PORT)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("unable to create cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	//
	// LOGIN AS JOE
	//
	checkLoginAsJoe(t, client)

	//
	// RETREIVE COOKIE STORED ON CLIENT
	//
	u, err := url.Parse(testUrl(TEST_PORT, ""))
	if err != nil {
		t.Fatalf("Unable to understand url: %s", testUrl(TEST_PORT, ""))
	}
	cookies := jar.Cookies(u)
	if len(cookies) != 1 {
		t.Errorf("Unexpected number of cookies: %d (not 1):", len(cookies))
	}
	cookieval := cookies[0]

	user := decodeCookieToSessionUser(t, cookies[0])
	if user.FirstName != "Joe" || user.EmailAddr != "joe@example.com" || user.UserUdid != "df12ba96-71c7-436d-b8f6-2d157d5f8ff1" {
		t.Errorf("wrong user found: %+v", user)
	}
	//
	// TRY TO READ JOE'S RECORD (as joe)
	//
	joePath := shared.URLGen.UserRecord(user.UserUdid)
	checkGetStatus(t, client, joePath, 200)
	//
	// TRY TO READ MARY'S RECORD (as joe)
	//
	maryPath := shared.URLGen.UserRecord("515f7619-8ea2-427f-8cf3-7a9201c747dd")
	checkGetStatus(t, client, maryPath, 401)

	//
	// LOGOUT CURRENT USER (joe), check cookie dropped
	//
	cookies = client.Jar.Cookies(u)
	if len(cookies) != 1 {
		t.Errorf("Unexpected number of cookies at start: %d (not 1):", len(cookies))
	}
	checkLogoutCurrentUser(t, client, 200)
	cookies = jar.Cookies(u)
	if len(cookies) != 0 {
		t.Errorf("Unexpected number of cookies: %d (not 0):", len(cookies))
	}

	//
	// TRY TO LOGOUT AGAIN
	//
	checkLogoutCurrentUser(t, client, 400)

	//
	// TRY TO READ JOE'S AND MARY'S RECORD, NOT LOGGED IN
	//
	checkGetStatus(t, client, maryPath, 401)
	checkGetStatus(t, client, joePath, 401)

	//
	// USE PREVIOUSLY CREATED SESSION TO READ JOES RECORD
	//
	jar.SetCookies(u, []*http.Cookie{cookieval})
	checkGetStatus(t, client, joePath, 200)

}

//
// HELPERS
//
func encodeAuth(t *testing.T, sas *s5.PasswordAuthParameters) string {
	return encodeAnything(t, sas)
}

func encodeAnything(t *testing.T, i interface{}) string {
	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	if err := enc.Encode(i); err != nil {
		t.Fatalf("failed to encode: %v", err)
	}
	return body.String()
}

func encodeUserRecord(t *testing.T, ur *shared.UserRecord) string {
	return encodeAnything(t, ur)
}

func testUrl(port int, path string) string {
	return fmt.Sprintf("http://localhost:%d%s", port, path)
}
