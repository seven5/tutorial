package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/seven5/seven5/migrate"

	"github.com/seven5/tutorial/shared"
)

//
// This table records the up and down migrations for this application.
//
var defn = migrate.Definitions{
	Up: map[int]migrate.MigrationFunc{
		1: oneUp,
		2: twoUp,
	},
	Down: map[int]migrate.MigrationFunc{
		1: oneDown,
		2: twoDown,
	},
}

//
// Our main is just a wrapper around migrate.Main()
//
func main() {

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		fmt.Fprintf(os.Stderr, "failed to get DATABASE_URL from environment")
	}
	m, err := migrate.NewPostgresMigrator(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to selma: %v", err)
		return
	}
	migrate.Main(&defn, m)
	m.Close()
}

//
// MIGRATION FROM EMPTY DB -> SIMPLE POST STRUCTURE, PLUS SAMPLE POSTS
//
func oneUp(tx *sql.Tx) error {
	var err error
	//
	// USER_RECORD (has to come before POST)
	// ("user" is a bad thing in postgres, so we do "user_record" )
	//
	_, err = tx.Exec(`
        CREATE TABLE user_record (
        user_udid CHAR(36) PRIMARY KEY,
        first_name VARCHAR(63),
        last_name VARCHAR(63),
        email_addr VARCHAR(63),
        password VARCHAR(63),
        disabled BOOLEAN DEFAULT false,
        admin BOOLEAN DEFAULT false
    )`)

	if err != nil {
		return err
	}

	//
	// POST
	//
	_, err = tx.Exec(`
	CREATE TABLE post (
	id BIGSERIAL PRIMARY KEY,
	title VARCHAR(127),
	updated TIMESTAMP WITH TIME ZONE,
	created TIMESTAMP WITH TIME ZONE,
	text TEXT,
	author_udid CHAR(36) REFERENCES user_record(user_udid)
	)`)
	if err != nil {
		return err
	}

	//
	// TEST DATA USER_RECORD
	//
	_, err = tx.Exec(`
					INSERT INTO user_record (first_name, last_name, user_udid, email_addr, password, admin)  VALUES
					('Joe','Smith','df12ba96-71c7-436d-b8f6-2d157d5f8ff1','joe@example.com', 'seekret', false),
					('Mary','Jones','515f7619-8ea2-427f-8cf3-7a9201c747dd','mary@example.com', 'bigseekret', true)
					`,
	)
	if err != nil {
		return err
	}

	//
	// TEST DATA POST
	//
	_, err = tx.Exec(`
						INSERT INTO post (title, updated, created, text, author_udid)  VALUES
						('first post!',CURRENT_TIMESTAMP - INTERVAL '1 hour',CURRENT_TIMESTAMP - INTERVAL '1 hour','This is the first post on the site!', 'df12ba96-71c7-436d-b8f6-2d157d5f8ff1'),
						('apology',CURRENT_TIMESTAMP - INTERVAL '30 minutes',CURRENT_TIMESTAMP - INTERVAL '30 minutes','Please excuse Joe''s stupid post.', '515f7619-8ea2-427f-8cf3-7a9201c747dd')
						`,
	)
	if err != nil {
		return err
	}

	return nil
}

//
// MIGRATION FROM 1 to 2 (ADD COLUMN FOR SHORT VERSION OF POST)
//
func twoUp(tx *sql.Tx) error {
	var err error
	//
	// ADD A COLUMN FOR THE SHORT VERSION
	//
	text := fmt.Sprintf("ALTER TABLE post ADD COLUMN text_short VARCHAR(%d)", shared.SHORT_LIMIT)
	_, err = tx.Exec(text)
	if err != nil {
		return err
	}
	//
	// MOVE THE DEFAULT DATA'S LONG VERSION INTO SHORT VERSION
	//
	_, err = tx.Exec(`UPDATE post SET text_short = text, text=''`)
	if err != nil {
		return err
	}
	return nil
}

//
// Go from having tables to not have them (state 0)
//
func oneDown(tx *sql.Tx) error {
	//bc of foreign keys, order of these drops is siginficant
	drops := []string{
		"DROP TABLE post",
		"DROP TABLE user_record",
	}
	for _, drop := range drops {
		_, err := tx.Exec(drop)
		if err != nil {
			return err
		}
	}
	return nil
}

//
// Go from having two to 1: remove the text_short column
//
func twoDown(tx *sql.Tx) error {
	var err error
	_, err = tx.Exec(`ALTER TABLE post DROP COLUMN text_short`)
	if err != nil {
		return err
	}
	return nil
}
