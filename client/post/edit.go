package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	s5 "github.com/seven5/seven5/client"

	"tutorial/client/uicommon"
	"tutorial/shared"
)

func main() {
	//
	// sanity check the current URL
	//
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

	//
	// we have a sensible url, lets create an instance of the new/edit page
	// with our choices for edit
	//
	ep := uicommon.NewEditPostPage(shared.URLGen.IndexPage(), func() {
	}, func() {
		uicommon.SetCurrentPage(shared.URLGen.PostView(id))
	}, func(p *shared.Post) (chan interface{}, chan s5.AjaxError) {
		p.Id = id
		return s5.AjaxPut(p, shared.URLGen.Post(id))
	})

	//
	// our SetupFunc does the work of calling the server to get the content
	// to initialize the form
	//
	ep.SetupFunc = func() {
		content, errorChan := s5.AjaxGet(&shared.Post{}, shared.URLGen.Post(id))
		go func() {
			select {
			case raw := <-content:
				p := raw.(*shared.Post)
				title := s5.NewHtmlId("input", "title")
				title.Dom().SetVal(p.Title)
				content := s5.NewHtmlId("textarea", "content")
				content.Dom().SetVal(p.TextShort + p.Text)
			case err := <-errorChan:
				ep.DisplayErrorText("failed to read post:"+err.Message, true)
			}
		}()

	}

	//
	// NAV BAR
	//
	ep.StandardPage.GetLoggedInUser()
	ep.StandardPage.AddBlogNav(false)
	ep.StandardPage.AddLeftNav("#", "Edit #"+fmt.Sprint(id), true)

	s5.Main(ep)
}
