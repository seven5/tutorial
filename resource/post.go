package resource

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/coocood/qbs"
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

	return &p, nil
}
