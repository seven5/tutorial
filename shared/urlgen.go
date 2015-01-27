package shared

import ()

//visible implementation of URLGenerator
var URLGen URLGenerator

//A UrlGenerator is a utility type to insure that all your code
//can easily "agree" on the URLs inside your application.
type URLGenerator interface {
	IndexPage() string
}

//implementation
type urlgen struct {
	//probably shouldn't have state here because the client and
	//server share this code and typically don't share this state
}

//IndexPage is the "home" page of the application as an absolute path.
func (u *urlgen) IndexPage() string {
	return "/index.html"
}
