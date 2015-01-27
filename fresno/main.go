package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//it is customary to use the named import version for seven5 as "s5"
	s5 "github.com/seven5/seven5"
)

func main() {

	herokuName := os.Getenv("HEROKU_NAME")
	if herokuName == "" {
		log.Fatalf("unable to get HEROKU_NAME from environment!")
	}
	heroku := s5.NewHerokuDeploy(herokuName, "fresno")
	//get the port from the environment, PORT variable
	port := heroku.Port()

	//the serve mux
	serveMux := s5.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "it's alive! bwah, haha!")
	})

	log.Printf("[SERVE] (IsTest=%v) waiting on :%d", heroku.IsTest(), port)
	log.Fatalf("%s", http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux))

}
