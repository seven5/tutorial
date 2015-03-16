// +build js

package uicommon

import (
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client" //client side

	"tutorial/shared"
)

//
// StandardPage is the type shared by all the pages in the application. It's
// primarily responsible for displaying errors in a consistent way and for
// managing the contents of the navigation bar at the top of the page.
//
type StandardPage struct {
	errRegion       s5.HtmlId
	errText         s5.HtmlId
	navRight        s5.HtmlId
	navLeft         s5.HtmlId
	includeRightNav bool
	user            *shared.UserRecord
	logoutTarget    string
	arrivalFunc     func(*shared.UserRecord)
}

// NewStandardPage creates a new standard page in the app.  This should not be
// called "standalone" but should be called as part of the initialization of
// of another type.  The logout target provides the url to navigate to, should
// the user click logout.  The includeRightNav option, if set, creates extra
// navigation elements in the far right for sign up, login, or, if the user
// is logged in, the current user's name and a log out button.   The userArrival
// function, if supplied, is called when the results of attempting to log in
// the user are known.  See #GetLoggedInUser.
func NewStandardPage(logoutTarget string, includeRightNav bool, userArrival func(*shared.UserRecord)) *StandardPage {
	self := &StandardPage{}
	self.errRegion = s5.NewHtmlId("div", "err-region")
	self.errText = s5.NewHtmlId("h5", "err-text")
	self.navLeft = s5.NewHtmlId("ul", "nav-left")
	self.navRight = s5.NewHtmlId("ul", "nav-right")
	self.logoutTarget = logoutTarget
	self.includeRightNav = includeRightNav
	self.arrivalFunc = userArrival
	return self
}

func (self *StandardPage) CurrentUser() *shared.UserRecord {
	return self.user
}

func (self *StandardPage) AddCurrentUserNav() {
	tree :=
		s5.LI(
			s5.A(
				s5.Text(self.user.FirstName+" "+self.user.LastName),
				s5.Event(s5.CLICK, func(evt jquery.Event) {
					evt.PreventDefault() //avoid the page going to top
				}),
			),
		).Build()
	self.navRight.Dom().Append(tree)
}

func (self *StandardPage) AddLogInNav() {
	self.AddRightNav(shared.URLGen.LoginPage(), "Log In")
}
func (self *StandardPage) AddSignUpNav() {
	self.AddRightNav(shared.URLGen.SignupPage(), "Sign Up")
}

func (self *StandardPage) AddRightNav(target string, text string) {
	tree :=
		s5.LI(
			s5.A(
				s5.HtmlAttrConstant(s5.HREF, target),
				s5.Text(text),
			),
		).Build()
	self.navRight.Dom().Append(tree)
}

//PerformLogOut is called when the user clicks the logout button.  Note that
//this may fail!
func (self *StandardPage) PerformLogOut() {
	var pap s5.PasswordAuthParameters

	pap.Username = self.user.EmailAddr
	pap.Op = s5.AUTH_OP_LOGOUT

	contentCh, errCh := s5.AjaxPost(&pap, shared.URLGen.Auth())

	go func() {
		select {
		case <-contentCh:
			SetCurrentPage(self.logoutTarget)
		case err := <-errCh:
			print("err is", err)
			self.DisplayErrorText("Log Out Trouble: "+err.Message, true)
		}
	}()

}
func (self *StandardPage) AddLogOutNav() {
	tree :=
		s5.LI(
			s5.A(
				s5.HtmlAttrConstant(s5.HREF, "#"),
				s5.Text("Log Out"),
				s5.Event(s5.CLICK, func(evt jquery.Event) {
					evt.PreventDefault()
					self.PerformLogOut()
				}),
			),
		).Build()
	self.navRight.Dom().Append(tree)
}

func (self *StandardPage) AddNewPostNav(isActive bool) {
	self.AddLeftNav(shared.URLGen.NewPost(), "New Post", isActive)
}
func (self *StandardPage) AddBlogNav(isActive bool) {
	self.AddLeftNav(shared.URLGen.IndexPage(), "Fresno", isActive)
}

func (self *StandardPage) AddLeftNav(target string, text string, isActive bool) {
	var li *s5.ViewImpl
	var a *s5.ViewImpl

	if target == "#" {
		a = s5.A(
			s5.Text(text),
			s5.Event(s5.CLICK, func(evt jquery.Event) {
				evt.PreventDefault()
			}),
		)
	} else {
		a = s5.A(
			s5.HtmlAttrConstant(s5.HREF, target),
			s5.Text(text),
		)
	}

	if isActive {
		li =
			s5.LI(
				s5.Class(Active),
				a,
			)
	} else {
		li =
			s5.LI(
				a,
			)
	}
	tree := li.Build()
	self.navLeft.Dom().Append(tree)
}

// DisplayErrorText displays the text provided at the top of the screen for
// a configurable number of seconds (see ERR_SECS).  If the danger value is
// true, the error displayed will have a red danger background, otherwise it
// is yellow (warning).
func (self *StandardPage) DisplayErrorText(text string, danger bool) {
	add := BgWarn
	remove := BgDanger
	if danger {
		add = BgDanger
		remove = BgWarn
	}
	displayErrorText(text, add, remove, self.errText, self.errRegion)
}

// GetLoggedInUser should be called during the initialization of a page based
// on StandardPage.  This call will contact the server and see if there is a
// current user.  If there is and you supplied a userArrival function in
// #NewStandardPage, it will be called with the logged in users details.  If
// there is no logged in user, the userArrival function is still called but
// with nil as the parameter.
func (self *StandardPage) GetLoggedInUser() {
	chLoggedIn, chLoginErr := s5.AjaxGet(&shared.UserRecord{}, shared.URLGen.Me())
	go func() {
		select {
		case raw := <-chLoggedIn:
			self.user = raw.(*shared.UserRecord)
			self.AddCurrentUserNav()
			self.AddLogOutNav()
			if self.arrivalFunc != nil {
				self.arrivalFunc(self.user)
			}
		case <-chLoginErr:
			if self.includeRightNav {
				self.AddSignUpNav()
				self.AddLogInNav()
			}
			if self.arrivalFunc != nil {
				self.arrivalFunc(nil)
			}
		}
	}()
}

//
// Attempt login tries to login to the server with the credentials provided.
// If this fails, it displays its error message using the provided standard
// page.
//
func AttemptLogin(std *StandardPage, user, password string) {
	var pap s5.PasswordAuthParameters

	pap.Username = user
	pap.Password = password
	pap.Op = s5.AUTH_OP_LOGIN

	print("login attempt", user)
	contentCh, errCh := s5.AjaxPost(&pap, shared.URLGen.Auth())

	go func() {
		select {
		case <-contentCh:
			SetCurrentPage(shared.URLGen.IndexPage())
		case err := <-errCh:
			if err.StatusCode == 401 {
				std.DisplayErrorText("That's probably not your password.", false)
			} else {
				std.DisplayErrorText("Log In Trouble: "+err.Message, true)
			}
		}
	}()
}
