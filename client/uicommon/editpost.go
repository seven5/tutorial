// +build js

package uicommon

import (
	_ "fmt"
	"strings"

	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client" //client side

	"tutorial/shared"
)

//This is per page state.  This is common between the "new post" and
//"edit post" functionality.
type EditPostPage struct {
	*StandardPage

	Title       s5.StringAttribute
	Content     s5.StringAttribute
	SetupFunc   func()
	SuccessFunc func()
	NetworkFunc func(*shared.Post) (chan interface{}, chan s5.AjaxError)

	//private stuff
	errRegion       s5.HtmlId
	errText         s5.HtmlId
	button          s5.HtmlId
	cancel          s5.HtmlId
	titleInput      s5.HtmlId
	contentTextArea s5.HtmlId
	disabled        s5.CssClass
}

// FormIsBad returns true if either the title or the textual content of the
// blog post is empty.  This is used to insure you can't press the post
// or update button if you don't have any content to send.
func (self *EditPostPage) FormIsBad(raw []s5.Equaler) s5.Equaler {
	title := strings.TrimSpace(raw[0].(s5.StringEqualer).S)
	content := strings.TrimSpace(raw[1].(s5.StringEqualer).S)
	if len(title) == 0 || len(content) == 0 {
		return s5.BoolEqualer{B: true}
	}
	return s5.BoolEqualer{B: false}
}

//NewEditPostPage creates new post or edit post page.  Note that the three functions
//here provide a callback-style API for instantiators of this type. The setup
//func is called at the time the page is loaded, to allow the caller to
//start any network loads that are needed.  The success func is called after
//a successful PUT or POST is completed.  The network func allows the caller
//to control PUT or POST ajax calls.  Other than these functions, the UI for
//new post and edit post are the same.
func NewEditPostPage(curr string, setup func(), success func(),
	network func(*shared.Post) (chan interface{}, chan s5.AjaxError)) *EditPostPage {
	result := &EditPostPage{
		Title:        s5.NewStringSimple(""),
		Content:      s5.NewStringSimple(""),
		SetupFunc:    setup,
		SuccessFunc:  success,
		NetworkFunc:  network,
		StandardPage: NewStandardPage(curr, true, nil),
		//private stuff
		errRegion:       s5.NewHtmlId("div", "err-region"),
		errText:         s5.NewHtmlId("h5", "err-text"),
		button:          s5.NewHtmlId("button", "post"),
		cancel:          s5.NewHtmlId("button", "cancel"),
		titleInput:      s5.NewHtmlId("input", "title"),
		contentTextArea: s5.NewHtmlId("textarea", "content"),
		disabled:        s5.NewCssClass("disabled"),
	}
	return result
}

// SubmitPost is used to actually send the content from this client to the
// server.  It uses the NetworkFunc supplied at the creation of this
// EditPostPage to allow either POST or PUT to be used.  If the server
// responds with success, it calls the SuccessFunc to know what action should
// be performed.
func (self *EditPostPage) SubmitPost() {
	var p shared.Post
	p.Title = self.Title.Value()
	raw := self.Content.Value()
	var short, rest string
	if len(raw) < shared.SHORT_LIMIT {
		short = raw
	} else {
		short = raw[:shared.SHORT_LIMIT]
		index := strings.LastIndex(short, "\n")
		if index != -1 {
			short = raw[:index]
			rest = raw[index:]
		} else {
			//didn't find a \n
			rest = raw[shared.SHORT_LIMIT:]
		}
	}
	p.TextShort = short
	p.Text = rest
	contentChan, errorChan := self.NetworkFunc(&p)
	go func() {
		select {
		case <-contentChan:
			self.SuccessFunc()
		case err := <-errorChan:
			self.DisplayErrorText(err.Message, true)
		}
	}()
}

//called after the dom is ready
func (self *EditPostPage) Start() {

	//hook up the constraints FROM the UI to the attributes
	self.Title.Attach(s5.StringEquality(s5.NewValueAttr(self.titleInput.Dom())))
	self.Content.Attach(s5.StringEquality(s5.NewValueAttr(self.contentTextArea.Dom())))

	///remove disabled when formIsBad returns false, computed from attributes
	self.button.CssExistenceAttribute(self.disabled).Attach(s5.NewSimpleConstraint(
		self.FormIsBad, self.Title, self.Content))

	self.cancel.Dom().On(s5.CLICK, func(evt jquery.Event) {
		evt.PreventDefault()
		SetCurrentPage(shared.URLGen.IndexPage())
	})

	self.button.Dom().On(s5.CLICK, func(evt jquery.Event) {
		evt.PreventDefault()
		self.SubmitPost()
	})

	self.SetupFunc()
}
