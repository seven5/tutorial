package shared

import (
	"fmt"
	"net/url"
)

const (
	MOST_RECENT  = 0
	OFFSET_PARAM = "offset"
	LIMIT_PARAM  = "limit"

	MARKDOWN_PARAM = "markdown"
)

func init() {
	URLGen = &urlgen{}
}

//URLGen is an implementation of URLGenerator for this application.  It is
//visible to both the client and the server and should be used for creating
//paths within the application.
var URLGen URLGenerator

//A UrlGenerator is a utility type to insure that all your code
//can easily "agree" on the URLs inside your application.
type URLGenerator interface {
	IndexPage() string
	LoginPage() string
	SignupPage() string
	Auth() string
	Me() string
	UserRecord(udid string) string
	UserRecordResource() string
	PostResource() string
	Posts(offset int, limit int) string
	Post(id int64, wantMarkdown bool) string
	PostView(id int64) string
	PostEdit(id int64) string
	NewPost() string
}

//implementation
type urlgen struct {
	//probably shouldn't have state here because the client and
	//server share this code and typically don't share this state
}

//IndexPage is the "home" page of the application as an absolute path.
func (u *urlgen) IndexPage() string {
	return "/posts/index.html"
}

//LoginPage is the page of the application where the user can log in.
func (u *urlgen) LoginPage() string {
	return "/login.html"
}

//SignupPage is the page of the application where the user can sign up.
func (u *urlgen) SignupPage() string {
	return "/signup.html"
}

//where the auth requests go, notable login and logout
func (u *urlgen) Auth() string {
	return "/auth"
}

//how to get who you are
func (u *urlgen) Me() string {
	return "/me"
}

//compute the rest resource user record for a given udid
func (u *urlgen) UserRecord(udid string) string {
	return u.UserRecordResource() + "/" + udid
}

//compute the rest resource part for the user record
func (u *urlgen) UserRecordResource() string {
	return "/rest/userrecord"
}

//compute the rest resource part for the post
func (u *urlgen) PostResource() string {
	return "/rest/post"
}

//compute the rest resource url for a given post
func (u *urlgen) Post(id int64, wantMarkdown bool) string {
	base := u.PostResource() + "/" + fmt.Sprint(id)
	if wantMarkdown {
		base += "?" + MARKDOWN_PARAM + "=true"
	}
	return base
}

//page for retreiving the new post form
func (u *urlgen) NewPost() string {
	return "/post/new.html"
}

//compute the html viewer  url for a given post
func (u *urlgen) PostView(id int64) string {
	return "/post/" + fmt.Sprint(id) + "/view"
}

//compute the html editor  url for a given post
func (u *urlgen) PostEdit(id int64) string {
	return "/post/" + fmt.Sprint(id) + "/edit"
}

//compute the url for a max number of posts
func (u *urlgen) Posts(start, limit int) string {
	values := url.Values{}
	values.Add(OFFSET_PARAM, fmt.Sprint(start))
	values.Add(LIMIT_PARAM, fmt.Sprint(limit))
	return "/rest/post?" + values.Encode()
}
