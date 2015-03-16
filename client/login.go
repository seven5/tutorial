// +build js

package main

import (
	"github.com/gopherjs/jquery"
	//note that this is importing the CLIENT SIDE of seven5
	s5 "github.com/seven5/seven5/client"

	"tutorial/client/uicommon"
	"tutorial/shared"
)

var (
	emailInput = s5.NewHtmlId("input", "email")
	pwdInput   = s5.NewHtmlId("input", "pwd")
	button     = s5.NewHtmlId("button", "login")

	primary   = s5.NewHtmlId("div", "primary")
	secondary = s5.NewHtmlId("div", "secondary")
)

//This is per page state, this lives as long as the page is displayed and there is
//only one instance per page.  The login page needs to have attributes that
//are connected to the type-in fields on the screen.  alreadyLoggedIn is used
//to trigger a completely different display if you reach this page already
//logged in, since this page would make no sense.
type loginPage struct {
	*uicommon.StandardPage
	email           s5.StringAttribute
	pwd             s5.StringAttribute
	alreadyLoggedIn s5.BooleanAttribute
}

func formIsBad(raw []s5.Equaler) s5.Equaler {
	e := raw[0].(s5.StringEqualer).S
	p := raw[1].(s5.StringEqualer).S
	if len(e) < 1 || len(p) < 6 {
		return s5.BoolEqualer{B: true}
	}
	return s5.BoolEqualer{B: false}
}

// Prepare the primary UI, the ui for the case where the user needs to login.
func (self *loginPage) preparePrimary() {
	self.email.Attach(s5.StringEquality(s5.NewValueAttr(emailInput.Dom())))
	self.pwd.Attach(s5.StringEquality(s5.NewValueAttr(pwdInput.Dom())))

	///remove disabled when formIsBad returns false
	button.CssExistenceAttribute(uicommon.Disabled).Attach(s5.NewSimpleConstraint(
		formIsBad, self.email, self.pwd))

	button.Dom().On(s5.CLICK, func(evt jquery.Event) {
		evt.PreventDefault()
		uicommon.AttemptLogin(self.StandardPage, self.email.Value(), self.pwd.Value())
	})
}

// selectPrimaryOrSecondary uses constraints to display either the primary UI or
// the secondary UI (which is just an error page).
func (self *loginPage) selectPrimaryOrSecondary() {
	s5.Equality(primary.CssExistenceAttribute(uicommon.Hide), self.alreadyLoggedIn)
	secondary.CssExistenceAttribute(uicommon.Hide).Attach(s5.NewBooleanInverter(self.alreadyLoggedIn))
}

// newLoginPage return a new instance of the loginPage structure.
func newLoginPage() *loginPage {
	result := &loginPage{
		email:           s5.NewStringSimple(""),
		pwd:             s5.NewStringSimple(""),
		alreadyLoggedIn: s5.NewBooleanSimple(false),
	}
	result.preparePrimary()
	result.selectPrimaryOrSecondary()

	return result
}

//called after the dom is ready
func (self *loginPage) Start() {
	self.StandardPage = uicommon.NewStandardPage(
		shared.URLGen.LoginPage(),
		false,
		func(u *shared.UserRecord) {
			if u != nil {
				self.alreadyLoggedIn.Set(true)
			}
		})
	self.StandardPage.AddBlogNav(false)
	self.StandardPage.GetLoggedInUser()
}

func main() {
	s5.Main(newLoginPage())
}
