// +build js

package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	s5 "github.com/seven5/seven5/client"
	"strconv"
	"strings"
	"time"

	"tutorial/client/uicommon"
	"tutorial/shared"
)

var (
	errRegion  s5.HtmlId
	errText    s5.HtmlId
	postParent = s5.NewHtmlId("div", "post-parent")
)

type viewPostPage struct {
	*uicommon.StandardPage
	id       int64
	myid     string
	authorid string
	admin    bool
	addedNav bool
}

func newViewPostPage(id int64) *viewPostPage {
	result := &viewPostPage{
		id: id,
	}
	result.StandardPage = uicommon.NewStandardPage(
		shared.URLGen.PostView(id),
		true,
		func(u *shared.UserRecord) {
			if u != nil {
				result.myid = u.UserUdid
				result.admin = u.Admin
				result.checkForSameAuthor()
			}
		})
	return result
}

//we need to do this check in two places, because don't know what order we will
//get myid and and author id
func (self *viewPostPage) checkForSameAuthor() {
	if !self.addedNav {
		if self.admin || self.myid == self.authorid {
			self.StandardPage.AddLeftNav(shared.URLGen.PostEdit(self.id),
				"Edit #"+fmt.Sprint(self.id), false)
			self.addedNav = true
		}
	}
}

func (self *viewPostPage) displayPost(p *shared.Post) {
	tree :=
		s5.DIV(
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.Class(uicommon.TextCenter),
				s5.H3(
					s5.Class(uicommon.Col11),
					s5.Text(p.Title),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.SPAN(
					s5.Class(uicommon.Col3),
					s5.Text("Written By"),
				),
				s5.SPAN(
					s5.Class(uicommon.Col9),
					s5.Text(p.Author.FirstName+" "+p.Author.LastName),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.SPAN(
					s5.Class(uicommon.Col3),
					s5.Text("Last Modification At"),
				),
				s5.SPAN(
					s5.Class(uicommon.Col9),
					s5.Text(p.Updated.Format(time.RFC850)),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.SPAN(
					s5.Class(uicommon.Col3),
					s5.Text("Created At"),
				),
				s5.SPAN(
					s5.Class(uicommon.Col9),
					s5.Text(p.Updated.Format(time.RFC850)),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.SPAN(
					s5.Class(uicommon.Col12),
					s5.HR(),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.SPAN(
					s5.Class(uicommon.Col12),
					s5.HtmlConstant(p.Text),
				),
			),
		).Build()
	postParent.Dom().Append(tree)
}

func (self *viewPostPage) Start() {
	self.StandardPage.GetLoggedInUser()
	self.StandardPage.AddBlogNav(false)
	self.StandardPage.AddLeftNav("#", "View #"+fmt.Sprint(self.id), true)
	content, errorChan := s5.AjaxGet(&shared.Post{}, shared.URLGen.Post(self.id)+"?markdown=true")
	go func() {
		select {
		case raw := <-content:
			p := raw.(*shared.Post)
			self.authorid = p.AuthorUdid
			self.checkForSameAuthor()
			self.displayPost(p)
		case err := <-errorChan:
			self.StandardPage.DisplayErrorText("failed to read post:"+err.Message, true)
		}
	}()
}
func main() {
	path := js.Global.Get("location").Get("pathname").String()
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		print("can't understand current url!", js.Global.Get("location"))
		return
	}
	idRaw := parts[2]
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		print("can't understand current url (id part?)!", js.Global.Get("location"))
		return
	}
	s5.Main(newViewPostPage(id))
}
