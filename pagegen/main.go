package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	s5 "github.com/seven5/seven5"

	"tutorial/shared"
)

const (
	JSON_HELPER = "support/json_helper.tmpl"
)

//utility function for generating current year for footers
func year() int {
	return time.Now().Year()
}

//useful in page generation for generating links
func urlgen() shared.URLGenerator {
	return shared.URLGen
}

//this table adds functionality to the "pipelines" you can use in
//go templates.
var funcs = map[string]interface{}{
	"year":   year,
	"urlgen": urlgen,
}

var (
	debug    = flag.Bool("debug", false, "Enable debug mode (more verbose output)")
	jsonFile = flag.String("json", "", "Json file to use as data for template")
	dir      = flag.String("dir", "", "Directory to read templates and json files from")
	start    = flag.String("start", "", "template file to start processing")
	support  = flag.String("support", "", "Support directory (inside dir) that should be included with all templates (only .tmpl files read)")
)

//our main() is just a wrapper to collect some options and then call
//s5.PagegenOpts.Main()
func main() {
	flag.Parse()
	if *start == "" || *dir == "" {
		flag.Usage()
		return
	}
	if filepath.Ext(*start) != ".html" && filepath.Ext(*start) != ".css" {
		log.Printf("probably a bad idea to use on a file that isn't html or css (%s)", *start)
	}

	po := s5.PagegenOpts{
		Funcs:           funcs,
		BaseDir:         *dir,
		SupportDir:      *support,
		JsonSupportFile: JSON_HELPER,
		JsonFile:        *jsonFile,
		TemplateFile:    *start,
		Debug:           *debug,
		TemplateSuffix:  filepath.Ext(*start),
	}
	po.Main()
}
