package resource

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coocood/qbs"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	s5 "github.com/seven5/seven5"

	"tutorial/shared"
)

type PostResource struct {
	//stateless
}

//get the resource by id
func (self *PostResource) FindQbs(id int64, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	var p shared.Post
	p.Id = id
	err := q.Find(&p)
	if err != nil && err == sql.ErrNoRows {
		return nil, s5.HTTPError(http.StatusNotFound, fmt.Sprintf("did not find %d", id))
	}
	if err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("error in qbs find: %v", err))
	}
	//don't show the password
	p.Author.Password = ""

	_, wantMarkdown := pb.Query("markdown")
	if wantMarkdown {
		p.Text = markdown(p.TextShort + p.Text)
		p.TextShort = ""
	}
	return &p, nil
}

func (self *PostResource) IndexQbs(pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	limit := int(pb.IntQueryParameter(shared.LIMIT_PARAM, 10))
	offset := int(pb.IntQueryParameter(shared.OFFSET_PARAM, 0))
	q, err := qbs.GetQbs()
	if err != nil {
		log.Printf("unable to get db connection:%v", err)
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to get db connection:%v", err))
	}
	var posts []*shared.Post
	err = q.Limit(limit).Offset(offset).OrderByDesc("created").FindAll(&posts)
	if err != nil && err == sql.ErrNoRows {
		return []*shared.Post{}, nil
	}
	if err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("error in qbs find: %v", err))
	}
	for _, p := range posts {
		author := p.Author
		author.Password = ""
		p.TextShort = markdown(p.TextShort)
		p.Text = markdown(p.Text)
	}
	return posts, nil
}

func (self *PostResource) PutQbs(id int64, i interface{}, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	//we ONLY trust the id from the client at this point
	p := &shared.Post{Id: id}
	if err := q.Find(p); err != nil {
		if err == sql.ErrNoRows {
			return nil, s5.HTTPError(http.StatusNotFound, "Id"+fmt.Sprint(id))
		}
		return nil, s5.HTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to find post to change: %v", err.Error()))
	}

	//we know they have a session because we checked in allow
	user := pb.Session().UserData().(*shared.UserRecord)
	if !user.Admin && user.UserUdid != p.AuthorUdid {
		return nil, s5.HTTPError(http.StatusUnauthorized, "you lose")
	}

	fromClient := i.(*shared.Post)
	//copy in the fields that they can change
	p.Title = fromClient.Title
	p.Text = fromClient.Text
	p.TextShort = fromClient.TextShort

	if _, err := q.Save(p); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to update post: %v", err.Error()))
	}
	p.Author.Password = "" //make sure they cant see that
	return p, nil
}

func (self *PostResource) PostQbs(i interface{}, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	//this is ok because of the AllowWrite check
	user := pb.Session().UserData().(*shared.UserRecord)
	proposed := i.(*shared.Post)
	var actual shared.Post
	actual.TextShort = strings.TrimSpace(proposed.TextShort)
	actual.Text = strings.TrimSpace(proposed.Text)
	actual.Title = strings.TrimSpace(proposed.Title)
	if actual.Title == "" || actual.TextShort == "" {
		return nil, s5.HTTPError(http.StatusBadRequest, "title and at content are required")
	}
	actual.AuthorUdid = user.UserUdid
	if _, err := q.Save(&actual); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to save post: %v", err.Error()))
	}
	return &actual, nil
}

func (self *PostResource) DeleteQbs(id int64, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	p := &shared.Post{Id: id}

	if err := q.Find(p); err != nil {
		if err == sql.ErrNoRows {
			return nil, s5.HTTPError(http.StatusNotFound, "Id"+fmt.Sprint(id))
		}
		return nil, s5.HTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to find post to delete: %v", err.Error()))
	}

	//we know they have a session because we checked in allow
	user := pb.Session().UserData().(*shared.UserRecord)
	if !user.Admin && user.UserUdid != p.AuthorUdid {
		return nil, s5.HTTPError(http.StatusUnauthorized, "you lose")

	}
	if _, err := q.Delete(p); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to delete post: %v", err.Error()))
	}
	return p, nil
}

//You may only post to this resource if you have a session
func (self *PostResource) AllowWrite(pb s5.PBundle) bool {
	return pb.Session() != nil
}

//You may only change this resource or delete if you are an admin.
func (self *PostResource) Allow(id int64, method string, pb s5.PBundle) bool {
	if method == "PUT" || method == "DELETE" {
		if pb.Session() == nil {
			return false
		}

		//we are going to allow this to go through and check in PUT or DELETE
		//to see if they have the proper authority

	}
	return true //find is always ok
}

func markdown(input string) string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return ""
	}
	unsafe := blackfriday.MarkdownCommon([]byte(input))
	return string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}
