package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Config struct {
	SvcHost      string
	DbDriver     string
	DbConnection string
}

type MibigService struct {
}

const (
	SqliteMibigSumbissionInit string = `CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER PRIMARY KEY NOT NULL,
		submitted TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		raw TEXT NOT NULL,
		v INTEGER );`
	SqliteGeneSumbissionInit string = `CREATE TABLE IF NOT EXISTS gene_submissions (
		id INTEGER PRIMARY KEY NOT NULL,
		bgc_id TEXT NOT NULL,
		submitted TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		raw TEXT NOT NULL,
		v INTEGER );`
	SqliteNrpsSumbissionInit string = `CREATE TABLE IF NOT EXISTS nrps_submissions (
		id INTEGER PRIMARY KEY NOT NULL,
		bgc_id TEXT NOT NULL,
		submitted TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		raw TEXT NOT NULL,
		v INTEGER );`
)

func (s *MibigService) getDb(cfg Config) (*sql.DB, error) {
	return sql.Open(cfg.DbDriver, cfg.DbConnection)
}

func (s *MibigService) Migrate(cfg Config) error {
	if cfg.DbDriver != "sqlite3" {
		log.Println("Not running Migrate for DB driver " + cfg.DbDriver)
		return nil
	}

	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(SqliteMibigSumbissionInit); err != nil {
		return err
	}
	if _, err = tx.Exec(SqliteGeneSumbissionInit); err != nil {
		return err
	}
	if _, err = tx.Exec(SqliteNrpsSumbissionInit); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *MibigService) Run(cfg Config) error {
	engine, db, err := s.CreateEngine(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	engine.Run(cfg.SvcHost)
	return nil
}

func (s *MibigService) CreateEngine(cfg Config) (*gin.Engine, *sql.DB, error) {
	db, err := s.getDb(cfg)
	if err != nil {
		return nil, nil, err
	}

	mr := &MibigResource{db: db}

	r := gin.Default()

	v1 := r.Group("/v1.0.0")
	{
		v1.GET("/bgc-registration", mr.ServiceInfo)
		v1.POST("/bgc-registration", mr.StoreMibigSubmissionV1)
	}

	v2 := r.Group("/v2.0.0")
	{
		v2.POST("/bgc-registration", mr.StoreMibigSubmissionV2)
		v2.POST("/bgc-detail-registration", mr.StoreBgcDetailSubmission)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": true, "message": "Not Found", "data": gin.H{"message": "null for uri: " + c.Request.URL.Host + c.Request.URL.Path}})
	})

	return r, db, nil
}
