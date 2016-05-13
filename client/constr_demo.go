package main

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/jquery"
	//note that this is importing the CLIENT SIDE of seven5
	s5 "github.com/seven5/seven5/client"
)

var (
	tl       = s5.NewHtmlId("button", "top-left")
	br       = s5.NewHtmlId("button", "bottom-right")
	bordered = s5.NewHtmlId("div", "bordered")

	dragInProgress = false
	dragOriginX    = 0
	dragOriginY    = 0
	startX         = 0
	startY         = 0
	dragging       s5.HtmlId
)

type constraintDemoPage struct {
	tlLeft   s5.IntegerAttribute
	tlTop    s5.IntegerAttribute
	brTop    s5.IntegerAttribute
	brLeft   s5.IntegerAttribute
	brHeight s5.IntegerAttribute
	brWidth  s5.IntegerAttribute
}

func newConstraintDemoPage() *constraintDemoPage {
	return &constraintDemoPage{
		tlLeft:   s5.NewIntegerSimple(0),
		tlTop:    s5.NewIntegerSimple(0),
		brTop:    s5.NewIntegerSimple(0),
		brLeft:   s5.NewIntegerSimple(0),
		brHeight: s5.NewIntegerSimple(0),
		brWidth:  s5.NewIntegerSimple(0),
	}
}

func (self *constraintDemoPage) Start() {
	s5.Equality(self.tlLeft, tl.StyleAttributeAsInt(s5.LEFT.String()))
	s5.Equality(self.tlTop, tl.StyleAttributeAsInt(s5.TOP.String()))
	s5.Equality(self.brTop, br.StyleAttributeAsInt(s5.TOP.String()))
	s5.Equality(self.brLeft, br.StyleAttributeAsInt(s5.LEFT.String()))
	s5.Equality(self.brHeight, br.StyleAttributeAsInt(s5.HEIGHT.String()))
	s5.Equality(self.brWidth, br.StyleAttributeAsInt(s5.WIDTH.String()))

	bordered.StyleAttribute(s5.WIDTH.String()).Attach(
		s5.NewSimpleConstraint(overallWidth, self.tlLeft, self.brTop, self.brWidth))
	bordered.StyleAttribute(s5.HEIGHT.String()).Attach(
		s5.NewSimpleConstraint(overallHeight, self.tlTop, self.brTop, self.brHeight))

	//opposite order for other
	tl.Dom().On(s5.MOUSE_DOWN, func(evt jquery.Event) {
		evt.PreventDefault()
		mouseDown(tl, evt)
	})
	br.Dom().On(s5.MOUSE_DOWN, func(evt jquery.Event) {
		evt.PreventDefault()
		mouseDown(br, evt)
	})
	bordered.Dom().On(s5.MOUSE_MOVE, func(evt jquery.Event) {
		evt.PreventDefault()
		mouseMove(evt)
	})
	bordered.Dom().On(s5.MOUSE_UP, func(evt jquery.Event) {
		evt.PreventDefault()
		mouseUp(evt)
	})

}

func main() {
	s5.Main(newConstraintDemoPage())
}

//
// DRAG FUNCTIONS
//
func mouseDown(button s5.HtmlId, evt jquery.Event) {
	dragInProgress = true
	startX, startY := button.Dom().OffsetOnPage()
	dragOriginX = evt.PageX
	dragOriginY = evt.PageY

	dragging = button
	print("start ", dragOriginX-startX, dragOriginY-startY)
}
func mouseMove(evt jquery.Event) {
	if !dragInProgress {
		return
	}
	deltaX := evt.PageX - dragOriginX
	deltaY := evt.PageY - dragOriginY
	proposedX := startX + deltaX
	proposedY := startY + deltaY
	if proposedX < 0 || proposedY < 0 {
		return
	}
	if proposedX > 1000 || proposedY > 800 {
		return
	}
	print("proposed ", proposedX, proposedY, "vs", startX, startY)
	dragging.Dom().SetCss("left", fmt.Sprint(proposedX)+"px")
	dragging.Dom().SetCss("top", fmt.Sprint(proposedY)+"px")
}
func mouseUp(evt jquery.Event) {
	mouseMove(evt)
	dragInProgress = false
}

func toIntOrPanic(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}

//
// CONSTRAINT FUNCTIONS
//
func overallWidth(raw []s5.Equaler) s5.Equaler {
	//we know these are integers
	left := raw[0].(s5.IntEqualer).I
	right := raw[1].(s5.IntEqualer).I
	width := raw[2].(s5.IntEqualer).I
	return s5.IntEqualer{I: ((right + width) - left) + 2}
}
func overallHeight(raw []s5.Equaler) s5.Equaler {
	//we know these are integers
	top := raw[0].(s5.IntEqualer).I
	right := raw[1].(s5.IntEqualer).I
	height := raw[2].(s5.IntEqualer).I
	return s5.IntEqualer{I: ((right + height) - top) + 2}
}
