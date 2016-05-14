package shared

import (
	"time"
)

const SHORT_LIMIT = 480

type Post struct {
	Id         int64 `qbs:"pk"`
	Title      string
	Updated    time.Time
	Created    time.Time
	Text       string
	TextShort  string
	AuthorUdid string `qbs:"fk:Author"`
	Author     *UserRecord
}
