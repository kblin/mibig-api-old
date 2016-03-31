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
	statement := fmt.Sprintf("INSERT INTO submissions(submitted, modified, raw, v) VALUES($1, $2, '%s', %d)",
		s.Raw, s.Version)
	if _, err := db.Exec(statement, s.Submitted, s.Modified); err != nil {
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
	statement := fmt.Sprintf("INSERT INTO %s(bgc_id, submitted, modified, raw, v) VALUES('%s', $1, $2, '%s', %d)",
		target, s.BgcId, s.Raw, s.Version)
	if _, err := db.Exec(statement, s.Submitted, s.Modified); err != nil {
		return err
	}

	return nil
}
