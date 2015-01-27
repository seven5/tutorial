package shared

import (
	"time"
)

type Post struct {
	Id         int64 `qbs:"pk"`
	Title      string
	Updated    time.Time
	Created    time.Time
	Text       string
	AuthorUdid string `qbs:"fk:Author"`
	Author     *UserRecord
}
