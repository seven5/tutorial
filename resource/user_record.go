package resource

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

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

func (self *UserRecordResource) PostQbs(i interface{}, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	var ur shared.UserRecord
	proposed := i.(*shared.UserRecord)
	e := strings.ToLower(strings.TrimSpace(proposed.EmailAddr))

	err := q.WhereEqual("email_addr", e).Find(&ur)
	if err != nil && err != sql.ErrNoRows {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("couldn't find: %v", err))
	}
	if err == nil {
		return nil, s5.HTTPError(http.StatusBadRequest, fmt.Sprintf("email address already registered %s: %s", proposed.EmailAddr, ur.UserUdid))
	}
	//just to make doubly sure we don't inadvently trust data from the client we
	//use the newly created ur which is a zero value at this point
	ur.Admin = false    //nuke the site from orbit
	ur.Disabled = false //its the only way to be sure
	if strings.Index(e, "@") == -1 {
		return nil, s5.HTTPError(http.StatusBadRequest, fmt.Sprintf("email address not ok: %s", proposed.EmailAddr))
	}
	ur.EmailAddr = e //copy it over
	if len(proposed.Password) < 6 {
		return nil, s5.HTTPError(http.StatusBadRequest, fmt.Sprintf("password too short: %s", proposed.Password))
	}
	ur.Password = proposed.Password //copy it over
	if len(proposed.FirstName) == 0 || len(proposed.LastName) == 0 {
		return nil, s5.HTTPError(http.StatusBadRequest, fmt.Sprintf("bad first or last name"))
	}
	ur.FirstName = proposed.FirstName
	ur.LastName = proposed.LastName
	//XXX this has a race condition which could cause two or more users with same
	//XXX email, right way to fix it is a DB constraint "unique"
	//values are ok, write it
	proposed.UserUdid = s5.UDID() //generate a random UDID
	if _, err := q.Save(proposed); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("couldn't save: %v", err))
	}
	return proposed, nil
}
