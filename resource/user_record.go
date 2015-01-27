package resource

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/coocood/qbs"
	s5 "github.com/seven5/seven5"

	"tutorial/shared"
)

type UserRecordResource struct {
	//stateless
}

//get the resource by id
func (self *UserRecordResource) FindQbs(udid string, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	var ur shared.UserRecord
	ur.UserUdid = udid
	err := q.Find(&ur)
	if err != nil && err == sql.ErrNoRows {
		return nil, s5.HTTPError(http.StatusNotFound, fmt.Sprintf("did not find %s", udid))
	}
	if err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("error in qbs find: %v", err))
	}
	//don't show the password
	ur.Password = ""

	return &ur, nil
}
