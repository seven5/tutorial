package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/coocood/qbs"
	_ "github.com/lib/pq"
	s5 "github.com/seven5/seven5"

	"github.com/seven5/tutorial/shared"
)

func setupForTestsNoServer(port int) {
	if testConfig == nil {
		os.Setenv("PORT", fmt.Sprint(port))
		testConfig = setup()
		if port != testConfig.heroku.Port() {
			panic(fmt.Sprintf("Unable to set PORT to %d: %d", port, testConfig.heroku.Port()))
		}
	}
}

func checkComponentResult(t *testing.T, msg string, cr s5.ComponentResult,
	expectedStatus int, pathOrRedir string) {

	if cr.Status != expectedStatus {
		t.Errorf("%s: expected status %d but got %d", msg, expectedStatus, cr.Status)
		return
	}
	if expectedStatus == http.StatusOK || expectedStatus == http.StatusMovedPermanently {

		if expectedStatus == http.StatusOK && pathOrRedir != cr.Path {
			t.Errorf("%s: bad path, expected '%s' but got '%s'", msg, pathOrRedir, cr.Path)
		} else if expectedStatus == http.StatusMovedPermanently && pathOrRedir != cr.Redir {
			t.Errorf("%s: bad redir, expected '%s' but got '%s'", msg, pathOrRedir, cr.Redir)
		}
	}
}

func TestRouting(t *testing.T) {
	setupForTestsNoServer(5000)

	//
	// UGH DB SETUP AN EXTRA TIME FOR THIS TEST
	//
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		fmt.Fprintf(os.Stderr, "failed to get DATABASE_URL from environment")
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("unable to open %s: %v", os.Getenv("DATABASE_URL"), err)
	}
	qbs.RegisterWithDb("postgres", db, qbs.NewPostgres())

	//
	// MATCHER UNDER TEST, PULLED FROM OUR CONFIG
	//
	matcher := testConfig.matcher

	result := matcher.Match(nil, "/")

	checkComponentResult(t, "route /", result, http.StatusMovedPermanently, "/posts/index.html")

	result = matcher.Match(nil, "/en/web/post/")
	checkComponentResult(t, "route /en/web/post/", result, http.StatusOK, "/en/web/post/index.html")

	result = matcher.Match(nil, "/en/web/posts/")
	checkComponentResult(t, "route /en/web/posts/", result, http.StatusOK, "/en/web/post/index.html")

	result = matcher.Match(nil, "/posts")
	checkComponentResult(t, "route /posts", result, http.StatusOK, "/post/index.html") //SINGULAR

	result = matcher.Match(nil, "/posts/index.html")
	checkComponentResult(t, "route /posts/index.html", result, http.StatusOK, "/post/index.html")

	result = matcher.Match(nil, "/post/1")
	checkComponentResult(t, "route /post/1", result, http.StatusOK, "/post/view.html")

	result = matcher.Match(nil, "/post/1/")
	checkComponentResult(t, "route /post/1/", result, http.StatusOK, "/post/view.html")

	result = matcher.Match(nil, "/post/1/view")
	checkComponentResult(t, "route /post/1/view", result, http.StatusOK, "/post/view.html")

	result = matcher.Match(nil, "/en/web/post/1/view.html")
	checkComponentResult(t, "route /en/web/post/1/view.html", result, http.StatusOK, "/en/web/post/view.html")

	result = matcher.Match(nil, "/post/1/view.css")
	checkComponentResult(t, "route /post/1/view.css", result, http.StatusOK, "/post/view.css")

	result = matcher.Match(nil, "/en/web/post/1/view.css")
	checkComponentResult(t, "route /en/web/post/1/view.css", result, http.StatusOK, "/en/web/post/view.css")

	result = matcher.Match(nil, "/post/1/edit")
	checkComponentResult(t, "route /post/1/edit", result, http.StatusUnauthorized, "")

	result = matcher.Match(nil, "/post/new")
	checkComponentResult(t, "route /post/new", result, http.StatusUnauthorized, "")

	result = matcher.Match(nil, "/post/new.html")
	checkComponentResult(t, "route /post/new.html", result, http.StatusUnauthorized, "")

	result = matcher.Match(nil, "/posts/index.js")
	checkComponentResult(t, "route /posts/index.js", result, http.StatusOK, "/post/index.js")

	result = matcher.Match(nil, "/post/new.js") //would end up with a 404 when served up
	checkComponentResult(t, "route /post/new.js", result, http.StatusOK, "/post/new.js")

	//
	// SETUP FOR A TEST AS ADMIN
	//
	q, err := qbs.GetQbs()
	if err != nil {
		t.Fatalf("unable to get qbs: %v", err)
	}
	defer q.Close()

	var user shared.UserRecord
	user.UserUdid = "515f7619-8ea2-427f-8cf3-7a9201c747dd" //mary
	if err := q.Find(&user); err != nil {
		t.Fatalf("unable to find mary: %v", err)
	}

	session, err := testConfig.sm.Assign(user.EmailAddr, &user, time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("Unable create session for mary: %v", err)
	}
	pbundle := s5.NewTestPBundle(map[string]string{}, map[string]string{}, session,
		testConfig.sm, map[string]string{}, map[reflect.Type]interface{}{})

	result = matcher.Match(pbundle, "/post/1/edit")
	checkComponentResult(t, "route /post/1/edit", result, http.StatusOK, "/post/edit.html")

	result = matcher.Match(pbundle, "/post/1/")
	checkComponentResult(t, "route /post/1/", result, http.StatusOK, "/post/view.html")

	result = matcher.Match(pbundle, "/post/new")
	checkComponentResult(t, "route /post/new/", result, http.StatusOK, "/post/new.html")

	result = matcher.Match(pbundle, "/post/new.html")
	checkComponentResult(t, "route /post/new.html", result, http.StatusOK, "/post/new.html")

}
