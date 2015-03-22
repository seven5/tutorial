// +build js

package main

import (
	"fmt"
	"time"

	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client" //client side

	"tutorial/client/uicommon"
	"tutorial/shared"
)

var (
	loginParent = s5.NewHtmlId("div", "login-parent")
	bottomItem  = s5.NewHtmlId("div", "bottom-item")
)

const (
	SLUG_LEN = 480
)

//per page state, this lives as long as the page is displayed and there is
//only one instance per page
type indexPage struct {
	*uicommon.StandardPage
	posts        *s5.Collection
	isNormalUser s5.BooleanAttribute
	myId         s5.StringAttribute
}

//create new index page
func newIndexPage() *indexPage {
	result := &indexPage{}
	result.posts = s5.NewList(result)
	result.isNormalUser = s5.NewBooleanSimple(true)
	result.myId = s5.NewStringSimple("")

	result.StandardPage = uicommon.NewStandardPage(
		shared.URLGen.IndexPage(),
		true,
		func(u *shared.UserRecord) {
			if u != nil {
				if u.Admin {
					result.isNormalUser.Set(false)
				}
				result.myId.Set(u.UserUdid)
				result.StandardPage.AddNewPostNav(false)

			}
		})

	return result
}

//called to add a post to the screen
func (self *indexPage) addPost(p *shared.Post) {
	model := self.newPostModel(p)
	self.posts.Add(model)
}

//called to remove a post from the *server* and if that succeeds, *then*
//remove it from the client display
func (self *indexPage) deletePost(model *postModel) {
	chDelete, chDeleteErr := s5.AjaxDelete(model.orig, shared.URLGen.Post(model.orig.Id, false))
	go func() {
		select {
		case <-chDelete:
			self.posts.Remove(model)
		case err := <-chDeleteErr:
			self.StandardPage.DisplayErrorText(err.Message, true)
		}
	}()
}

// constraint function that calculates if a given post could be edited or
// deleted by the current user
func hideEditDelete(raw []s5.Equaler) s5.Equaler {
	isNormal := raw[0].(s5.BoolEqualer).B
	myudid := raw[1].(s5.StringEqualer).S
	postersUdid := raw[2].(s5.StringEqualer).S

	if !isNormal {
		return s5.BoolEqualer{B: false}
	}
	return s5.BoolEqualer{B: myudid != postersUdid}
}

//called after the dom is ready
func (self *indexPage) Start() {

	self.StandardPage.GetLoggedInUser()
	self.StandardPage.AddBlogNav(true)

	postsURL := shared.URLGen.Posts(shared.MOST_RECENT, 10)
	var posts []*shared.Post
	chPosts, chPostErr := s5.AjaxIndex(&posts, postsURL)

	go func() {
		select {
		case raw := <-chPosts:
			posts := raw.(*[]*shared.Post)
			for i := 0; i < len(*posts); i++ {
				self.addPost((*posts)[i])
			}
		case perr := <-chPostErr:
			self.StandardPage.DisplayErrorText("Unable to retreive posts: "+perr.Message, true)
		}
	}()
}

func main() {
	s5.Main(newIndexPage())
}

//
// POST MODEL
//
//
type postModel struct {
	s5.ModelName
	date           s5.StringAttribute
	title          s5.StringAttribute
	authorName     s5.StringAttribute
	authorUdid     s5.StringAttribute
	hideEditDelete s5.BooleanAttribute
	orig           *shared.Post
}

//
// Needed for Collections
//
func (self *postModel) Equal(e s5.Equaler) bool {
	if e == nil {
		return false
	}
	other := e.(*postModel)
	return self.Id() == other.Id()
}

//
// Create a post model with attributes based on the shared.Post
//
func (self *indexPage) newPostModel(p *shared.Post) *postModel {
	result := new(postModel)
	result.ModelName = s5.NewModelName(result)

	result.date = s5.NewStringSimple(p.Created.Format(time.RFC822))
	if p.Updated.After(p.Created) {
		result.date = s5.NewStringSimple(p.Updated.Format(time.RFC822))
	}
	result.title = s5.NewStringSimple(p.Title)
	result.authorName = s5.NewStringSimple(fmt.Sprintf("%s %s", p.Author.FirstName, p.Author.LastName))
	result.authorUdid = s5.NewStringSimple(p.AuthorUdid)
	result.hideEditDelete = s5.NewBooleanSimple(true)

	result.hideEditDelete.Attach(s5.NewSimpleConstraint(
		hideEditDelete, self.isNormalUser, self.myId, result.authorUdid))

	result.orig = p
	return result
}

//
// JOINER API
//

func (self *indexPage) Add(newLen int, m s5.Model) {
	model := m.(*postModel)

	tree := s5.DIV(
		s5.Class(uicommon.Row),
		s5.DIV(
			s5.Class(uicommon.BlogEntry),
			s5.Class(uicommon.Col12),
			s5.ModelId(model),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.Class(uicommon.Shaded),
				s5.A(
					s5.HtmlAttrConstant(s5.HREF, shared.URLGen.PostView(model.orig.Id)),
					s5.Class(uicommon.Col10),
					s5.Class(uicommon.H3),
					s5.TextEqual(model.title),
				),
				s5.SPAN(
					s5.Class(uicommon.Col2),
					s5.Class(uicommon.TopSpace),
					s5.BUTTON(
						s5.Class(uicommon.Btn),
						s5.Class(uicommon.BtnTiny),
						s5.Class(uicommon.BtnPrimary),
						s5.Text("Edit"),
						s5.CssExistence(uicommon.Hide, model.hideEditDelete),
						s5.Event(s5.CLICK, func(evt jquery.Event) {
							evt.PreventDefault()
							uicommon.SetCurrentPage(shared.URLGen.PostEdit(model.orig.Id))
						}),
					),
					s5.BUTTON(
						s5.Class(uicommon.Btn),
						s5.Class(uicommon.BtnTiny),
						s5.Class(uicommon.BtnDanger),
						s5.Class(uicommon.LeftSpace),
						s5.Text("Delete"),
						s5.CssExistence(uicommon.Hide, model.hideEditDelete),
						s5.Event(s5.CLICK, func(evt jquery.Event) {
							evt.PreventDefault()
							self.deletePost(model)
						}),
					),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.Class(uicommon.Shaded),
				s5.SPAN(
					s5.Class(uicommon.Col12),
					s5.Class(uicommon.H5),
					s5.SPAN(
						s5.TextEqual(model.authorName),
					),
					s5.SPAN(
						s5.Text("@"),
					),
					s5.SPAN(
						s5.TextEqual(model.date),
					),
				),
			),
			s5.DIV(
				s5.Class(uicommon.Row),
				s5.SPAN(
					s5.Class(uicommon.BlogText),
					s5.Class(uicommon.ColOffset1),
					s5.Class(uicommon.Col10),
					s5.HtmlConstant(model.orig.TextShort),
				),
			),
		),
	).Build()
	bottomItem.Dom().Before(tree)
}

func (self *indexPage) Remove(i int, m s5.Model) {
	s5.HtmlIdFromModel("div", m).Dom().Remove()
}
