package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/gopherjs/gopherjs" //force godep to save this package
	_ "github.com/gopherjs/jquery"   //force godep to save this package
	_ "github.com/tools/godep"       //force godep to save this package

	_ "github.com/lib/pq"
	//it is customary to use the named import version for seven5 as "s5"
	"github.com/coocood/qbs"
	s5 "github.com/seven5/seven5"

	"tutorial/resource"
	"tutorial/shared"
)

const (
	NAME             = "fresno"
	PWD_RESET_PAGE   = "pwdreset.html"
	PWD_RESET_PREFIX = "oops"
)

type fresnoConfig struct {
	serveMux   *s5.ServeMux
	pwdHandler *s5.SimplePasswordHandler
	base       *s5.BaseDispatcher

	heroku  s5.DeploymentEnvironment
	cm      s5.CookieMapper
	sm      s5.ValidatingSessionManager
	matcher s5.ComponentMatcher
}

func idToPost(id int64) (*shared.Post, error) {
	//check to see if that ID is valid
	q, err := qbs.GetQbs()
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var postId shared.Post
	postId.Id = id
	err = q.Find(&postId)
	return &postId, err
}

func existCheck(pb s5.PBundle, id int64) (bool, error) {
	_, err := idToPost(id)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

//only logged in users can post
func newCheck(pb s5.PBundle) (bool, error) {
	if pb == nil || pb.Session == nil {
		return false, nil
	}
	return true, nil
}

//anybody can view a post but to edit you have to be an admin or
//you have to be the owner of the post
func viewEditCheck(pb s5.PBundle, id int64, isView bool) (bool, error) {
	if isView {
		return true, nil
	}
	if pb == nil || pb.Session() == nil {
		return false, nil
	}
	//
	// We must look at the actual post to see if you are the person
	// who wrote it.
	//
	user := pb.Session().UserData().(*shared.UserRecord)
	if user.Admin {
		return true, nil
	}
	p, err := idToPost(id)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil //should never happen because of exist check
	}

	return p.AuthorUdid == user.UserUdid, nil
}

//
// setup does the bulk of the work of configuring the application.
//
func setup() *fresnoConfig {

	result := &fresnoConfig{}

	herokuName := os.Getenv("HEROKU_NAME")
	if herokuName == "" {
		log.Fatalf("unable to get HEROKU_NAME from environment!")
	}
	result.heroku = s5.NewHerokuDeploy(herokuName, NAME)

	//
	// We supply our own session manager since we want to do some user-specific
	// things.
	//
	result.cm = s5.NewSimpleCookieMapper(NAME)
	result.sm = newFresnoValidatingSessionManager()

	//utility to handle passwords
	result.pwdHandler = s5.NewSimplePasswordHandler(result.sm, result.cm)

	//matcher for dealing with the typed in URLs
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		log.Fatalf("unable to get STATIC_DIR from environment!")
	}
	//what do we do if given empty URL, note we use the SINGULAR here
	homepage := s5.ComponentResult{
		Status: http.StatusMovedPermanently,
		Redir:  "/post/index.html",
	}

	postComponent := s5.NewSimpleIdComponent("post", existCheck, newCheck, viewEditCheck)
	indexComponent := s5.NewIndexOnlyComponent("posts", "post/index.html") //SINGULAR for index

	result.matcher = s5.NewSimpleComponentMatcher(result.cm, result.sm, staticDir,
		homepage, result.heroku.IsTest(), postComponent, indexComponent)

	//base dispatcher is the "root" of the dispatching for REST resources
	result.base = s5.NewBaseDispatcher(result.sm, result.cm)

	//we are going to be using QBS as a storage engine
	store := result.heroku.GetQbsStore()

	//the serve mux (works like http.ServeMux)
	result.serveMux = s5.NewServeMux()

	//plus allows disptachers to be configured
	result.serveMux.Dispatch("/rest/", result.base)

	//add static files
	result.serveMux.Handle("/", result.matcher)
	//
	//user record
	//
	result.base.ResourceSeparateUdid("userrecord",
		&shared.UserRecord{},
		nil, //index
		s5.QbsWrapFindUdid(&resource.UserRecordResource{}, store),
		s5.QbsWrapPost(&resource.UserRecordResource{}, store), //post
		nil, //put
		nil) //delete

	result.base.ResourceSeparate("post",
		&shared.Post{},
		s5.QbsWrapIndex(&resource.PostResource{}, store),
		s5.QbsWrapFind(&resource.PostResource{}, store),
		s5.QbsWrapPost(&resource.PostResource{}, store),
		s5.QbsWrapPut(&resource.PostResource{}, store),
		s5.QbsWrapDelete(&resource.PostResource{}, store),
	)

	result.serveMux.HandleFunc(shared.URLGen.Me(), result.pwdHandler.MeHandler)
	result.serveMux.HandleFunc(shared.URLGen.Auth(), result.pwdHandler.AuthHandler)

	return result

}

func main() {
	config := setup()
	log.Printf("[SERVE] (IsTest=%v) waiting on :%d", config.heroku.IsTest(), config.heroku.Port())
	log.Fatalf("%s", http.ListenAndServe(fmt.Sprintf(":%d", config.heroku.Port()), config.serveMux))
}
