package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	//it is customary to use the named import version for seven5 as "s5"
	s5 "github.com/seven5/seven5"

	"tutorial/resource"
	"tutorial/shared"
)

const (
	NAME = "fresno"
)

func main() {

	herokuName := os.Getenv("HEROKU_NAME")
	if herokuName == "" {
		log.Fatalf("unable to get HEROKU_NAME from environment!")
	}
	heroku := s5.NewHerokuDeploy(herokuName, NAME)
	//get the port from the environment, PORT variable
	port := heroku.Port()

	//
	// For now, we will use the default implementation of sessions.
	//
	sm := s5.NewSimpleSessionManager()

	//base dispatcher is the "root" of the dispatching for REST resources
	base := s5.NewBaseDispatcher(NAME, sm)

	//we are going to be using QBS as a storage engine
	store := heroku.GetQbsStore()

	//the serve mux (works like http.ServeMux)
	serveMux := s5.NewServeMux()
	//plus allows disptachers to be configured
	serveMux.Dispatch("/rest/", base)
	//add static files
	serveMux.Handle("/", s5.NewSimpleStaticFilesServer("/", heroku.IsTest()))
	//
	//just implementing "find" method right now
	//
	base.ResourceSeparateUdid("userrecord",
		&shared.UserRecord{},
		nil, //index
		s5.QbsWrapFindUdid(&resource.UserRecordResource{}, store),
		s5.QbsWrapPost(&resource.UserRecordResource{}, store), //post
		nil, //put
		nil) //delete

	base.ResourceSeparate("post",
		&shared.Post{},
		nil, //index
		s5.QbsWrapFind(&resource.PostResource{}, store),
		nil, //post
		nil, //put
		nil) //delete

	log.Printf("[SERVE] (IsTest=%v) waiting on :%d", heroku.IsTest(), port)
	log.Fatalf("%s", http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux))

}
