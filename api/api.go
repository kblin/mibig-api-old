package api

import (
	"database/sql"
	"fmt"
	"time"
)

type MibigSubmission struct {
	Id        int
	Submitted time.Time
	Modified  time.Time
	Raw       string
	Version   int
}

func (s *MibigSubmission) Create(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO submissions(submitted, modified, raw, v) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(s.Submitted, s.Modified, s.Raw, s.Version); err != nil {
		return err
	}
	return nil
}

type BgcDetailSubmission struct {
	Id        int
	BgcId     string
	Submitted time.Time
	Modified  time.Time
	Raw       string
	Version   int
}

func (s *BgcDetailSubmission) Create(db *sql.DB, target string) error {
	statement := fmt.Sprintf("INSERT INTO %s(bgc_id, submitted, modified, raw, v) VALUES($1, $2, $3, $4, $5)", target)
	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(s.BgcId, s.Submitted, s.Modified, s.Raw, s.Version); err != nil {
		return err
	}

	return nil
}
