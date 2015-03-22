// +build js

package main

import (
	s5 "github.com/seven5/seven5/client"

	"tutorial/client/uicommon"
	"tutorial/shared"
)

func main() {

	// create a new instance of the edit/new post page with our choices
	// that are right for new
	ep := uicommon.NewEditPostPage(shared.URLGen.IndexPage(), func() {
		//nothing to do
	}, func() {
		uicommon.SetCurrentPage(shared.URLGen.IndexPage())
	}, func(p *shared.Post) (chan interface{}, chan s5.AjaxError) {
		return s5.AjaxPost(p, shared.URLGen.PostResource())
	})

	//navbar
	ep.StandardPage.GetLoggedInUser()
	ep.StandardPage.AddBlogNav(false)
	ep.StandardPage.AddNewPostNav(true)

	s5.Main(ep)
}
