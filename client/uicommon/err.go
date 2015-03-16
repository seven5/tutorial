package uicommon

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"
)

const (
	ERR_SECS = 15 // how long to leave error on screen
)

// SlideDown is a utility for causing a particular HTML element to gradually
// appear (jquery "slideDown").
func SlideDown(id s5.HtmlId) {
	selector := id.TagName() + "#" + id.Id() //awful, should be done through API
	if !s5.TestMode {
		jq := jquery.NewJQuery(selector)
		jq.Underlying().Call("slideDown", "slow")
	}
}

// SlideDown is a utility for causing a particular HTML element to gradually
// disappear (jquery "slideUp").
func SlideUp(id s5.HtmlId) {
	selector := id.TagName() + "#" + id.Id() //awful, should be done through API
	if !s5.TestMode {
		jq := jquery.NewJQuery(selector)
		jq.Underlying().Call("slideUp", "slow")
	}
}

//displayErrorText is the lower-level routine that actually causes the animation
//transitions, adds/removes CSS classes, and handles the timeout for removing
//the error message after a delay.
func displayErrorText(text string, clazzAdd s5.CssClass, clazzRemove s5.CssClass, errText s5.HtmlId, errRegion s5.HtmlId) {
	errText.Dom().SetText(text)

	//remove any old classes
	errRegion.Dom().RemoveClass(clazzRemove.ClassName()) //should be through API
	//now add the new class
	errRegion.Dom().AddClass(clazzAdd.ClassName()) //should be through API

	SlideDown(errRegion)

	js.Global.Call("setTimeout", func() {
		SlideUp(errRegion)
	}, ERR_SECS*1000)
}
