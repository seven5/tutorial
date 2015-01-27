package main

import (
	"fmt"
	"strings"

	"github.com/gopherjs/jquery"
	//note that this is importing the CLIENT SIDE of seven5
	s5 "github.com/seven5/seven5/client"

	"tutorial/shared"
)

var (
	disabled = s5.NewCssClass("disabled")
)

//per page state, this lives as long as the page is displayed and there is
//only one instance per page
type signupPage struct {
	first s5.StringAttribute
	last  s5.StringAttribute
	email s5.StringAttribute
	pwd1  s5.StringAttribute
	pwd2  s5.StringAttribute

	firstFeedback s5.StringAttribute
	emailFeedback s5.StringAttribute
	pwd1Feedback  s5.StringAttribute
}

//create new signup page
func newSignupPage() signupPage {
	result := signupPage{
		first: s5.NewStringSimple(""),
		last:  s5.NewStringSimple(""),
		email: s5.NewStringSimple(""),
		pwd1:  s5.NewStringSimple(""),
		pwd2:  s5.NewStringSimple(""),

		firstFeedback: s5.NewStringSimple(""),
		emailFeedback: s5.NewStringSimple(""),
		pwd1Feedback:  s5.NewStringSimple(""),
	}

	return result
}

//convenient way to attach constraints
type feedbackInfo struct {
	input      s5.HtmlId
	output     s5.HtmlId
	constraint s5.Constraint
}

//called after the dom is ready
func (s signupPage) Start() {
	print("hello, world")

	//note that these will PANIC if the dom element cannot be found
	//usually this happens when you change your HTML and don't change
	//the corresponding code. these ids are controlled by the json in
	//signup.json
	first := s5.NewHtmlId("input", "first_name")
	firstOut := s5.NewHtmlId("label", "feedback_first_name")
	last := s5.NewHtmlId("input", "last_name")
	email := s5.NewHtmlId("input", "email")
	emailOut := s5.NewHtmlId("label", "feedback_email")
	pwd1 := s5.NewHtmlId("input", "password")
	pwd1Out := s5.NewHtmlId("label", "feedback_password")
	pwd2 := s5.NewHtmlId("input", "confirm_password")
	button := s5.NewHtmlId("button", "signup")

	//make a table that maps the _attribute_ to the feedback information so
	//we can avoid a lot of repetitive code
	attrMap := make(map[s5.StringAttribute]*feedbackInfo)
	attrMap[s.first] = &feedbackInfo{first, firstOut, s5.NewSimpleConstraint(nameFeedback, s.first, s.last)}
	attrMap[s.last] = &feedbackInfo{last, nil, nil}
	attrMap[s.email] = &feedbackInfo{email, emailOut, s5.NewSimpleConstraint(emailFeedback, s.email)}
	attrMap[s.pwd1] = &feedbackInfo{pwd1, pwd1Out, s5.NewSimpleConstraint(pwdFeedback, s.pwd1, s.pwd2)}
	attrMap[s.pwd2] = &feedbackInfo{pwd2, nil, nil}

	for attr, info := range attrMap {
		//this is just copies any value typed in the form to the attribute
		attr.Attach(s5.StringEquality(s5.NewValueAttr(info.input.Dom())))

		//add a constraint to give negative text feedback about the input
		if info.constraint != nil {
			//attach that string to the output display
			info.output.TextAttribute().Attach(info.constraint)
		}
	}

	///remove disabled when formIsBad returns false
	button.CssExistenceAttribute(disabled).Attach(s5.NewSimpleConstraint(
		formIsBad, s.first, s.last, s.email, s.pwd1, s.pwd2))

	button.Dom().On(s5.CLICK, func(evt jquery.Event) {
		evt.PreventDefault()
		var ur shared.UserRecord
		ur.EmailAddr = s.email.Value()
		ur.FirstName = s.first.Value()
		ur.LastName = s.last.Value()
		ur.Password = s.pwd1.Value()
		ur.Admin = true //hee hee hee
		s5.PostNew(&ur, "/rest",
			func(ptrToStruct interface{}) {
				user := ptrToStruct.(*shared.UserRecord)
				print("user is", user.EmailAddr, user.UserUdid)
			},
			func(httpCode int, msg string) {
				print("failed to post", httpCode, msg)
			})
	})
}

func main() {
	s5.Main(newSignupPage())
}

//
// CONSTRAINT FUNCTIONS
//
func nameFeedback(raw []s5.Equaler) s5.Equaler {
	//sadly, you have to *know* what types to expect based on the constraint
	//inputs provided to NewSimpleConstraint
	firstName := strings.TrimSpace(raw[0].(s5.StringEqualer).S)
	lastName := strings.TrimSpace(raw[1].(s5.StringEqualer).S)

	if len(firstName) == 0 && len(lastName) == 0 {
		return s5.StringEqualer{S: ""}
	}

	if len(firstName) == 0 {
		return s5.StringEqualer{S: "First name can't be blank!"}
	}
	if len(lastName) == 0 {
		return s5.StringEqualer{S: "Last name can't be blank!"}
	}
	return s5.StringEqualer{S: fmt.Sprintf("Other folks will see '%s %s'", firstName, lastName)}
}

func emailFeedback(raw []s5.Equaler) s5.Equaler {
	email := strings.TrimSpace(raw[0].(s5.StringEqualer).S)
	if len(email) == 0 {
		return s5.StringEqualer{S: ""}
	}
	if len(email) < 6 { //a@b.co
		return s5.StringEqualer{S: "That doesn't look like an email address!"}
	}
	if strings.Index(email, "@") == -1 {
		return s5.StringEqualer{S: "That doesn't look like an email address!"}
	}
	return s5.StringEqualer{S: ""} //no error
}

func pwdFeedback(raw []s5.Equaler) s5.Equaler {
	//sadly, you have to *know* what types to expect based on the constraint
	//inputs provided to NewSimpleConstraint
	pwd1 := raw[0].(s5.StringEqualer).S //don't trim space because might be significant
	pwd2 := raw[1].(s5.StringEqualer).S //don't trim space because might be significant

	if len(pwd1) == 0 && len(pwd2) == 0 {
		return s5.StringEqualer{S: ""}
	}

	if len(pwd1) < 6 {
		return s5.StringEqualer{S: "Password is too short!"}
	}
	if pwd1 != pwd2 {
		return s5.StringEqualer{S: "Passwords don't match!"}
	}
	return s5.StringEqualer{"Passwords match."}
}

//note the "logic" of this is backwards and we return true when the form is
//BAD because we are computing the existence of the "disabled" tag via this function
func formIsBad(raw []s5.Equaler) s5.Equaler {
	firstName := strings.TrimSpace(raw[0].(s5.StringEqualer).S)
	lastName := strings.TrimSpace(raw[1].(s5.StringEqualer).S)
	email := strings.TrimSpace(raw[2].(s5.StringEqualer).S)
	pwd1 := raw[3].(s5.StringEqualer).S
	pwd2 := raw[4].(s5.StringEqualer).S

	if len(firstName) == 0 {
		return s5.BoolEqualer{B: true}
	}
	if len(lastName) == 0 {
		return s5.BoolEqualer{B: true}
	}
	if len(email) < 6 {
		return s5.BoolEqualer{B: true}
	}
	if strings.Index(email, "@") == -1 {
		return s5.BoolEqualer{B: true}
	}
	if len(pwd1) < 6 {
		return s5.BoolEqualer{B: true}
	}
	if pwd1 != pwd2 {
		return s5.BoolEqualer{B: true}
	}
	//enable button
	return s5.BoolEqualer{B: false}

}
