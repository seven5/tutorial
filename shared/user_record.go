package shared

type UserRecord struct {
	UserUdid  string `qbs:"pk"`
	FirstName string
	LastName  string
	EmailAddr string
	Password  string
	Admin     bool
}
