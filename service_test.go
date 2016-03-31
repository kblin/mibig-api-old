package main_test

import (
	"database/sql"
	"github.com/kblin/mibig-api/client"
	"github.com/kblin/mibig-api/service"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MibigSuite struct {
	svc    *service.MibigService
	cfg    service.Config
	port   int
	dbFile *os.File
	db     *sql.DB
	ts     *httptest.Server
}

var _ = Suite(&MibigSuite{})

// Actual tests here

func (s *MibigSuite) TestServiceInfo(c *C) {
	client := client.MibigClient{Host: s.ts.URL}

	message, err := client.ServiceInfo()
	c.Assert(err, IsNil)
	c.Assert(message, Matches, "bgc registration up and running")
}

func (s *MibigSuite) TestStoreMibigSubmissionOnlyValid(c *C) {
	client := client.MibigClient{Host: s.ts.URL}
	resp, err := client.StoreMibigSubmissionOnly(`{"foo":"bar"}`, 1)
	c.Assert(err, IsNil)
	c.Assert(resp.Error, Equals, false)
	c.Assert(resp.Message, Matches, "bgc registration successful.")
}

func (s *MibigSuite) TestStoreMibigSubmissionOnlyNoJson(c *C) {
	client := client.MibigClient{Host: s.ts.URL}
	resp, err := client.StoreMibigSubmissionOnly("", 1)
	c.Assert(err, NotNil)
	c.Assert(resp.Error, Equals, true)
	c.Assert(resp.Message, Matches, "json not provided")
}

func (s *MibigSuite) TestStoreMibigSubmissionOnlyInvalidNumber(c *C) {
	client := client.MibigClient{Host: s.ts.URL}
	resp, err := client.StoreMibigSubmissionOnly(`{"foo":"bar"}`, -1)
	c.Assert(err, NotNil)
	c.Assert(resp.Error, Equals, true)
	c.Assert(resp.Message, Matches, "Need a version parameter greater than 0")
}

func (s *MibigSuite) TestStoreMibigSubmission(c *C) {
	client := client.MibigClient{Host: s.ts.URL}
	status, err := client.StoreMibigSubmission(`{"foo":"bar"}`, 1)
	c.Assert(err, IsNil)
	c.Assert(status, Equals, 303)
}

// Test Boilerplate
func (s *MibigSuite) SetUpSuite(c *C) {
	// Set up a SQLite3 DB for testing
	file, err := ioutil.TempFile("", "mibig_test_db")
	c.Assert(err, IsNil)
	s.dbFile = file
	s.cfg = service.Config{DbDriver: "sqlite3", DbConnection: s.dbFile.Name()}
	s.svc = &service.MibigService{}
	c.Assert(s.svc.Migrate(s.cfg), IsNil)

	engine, db, err := s.svc.CreateEngine(s.cfg)
	c.Assert(err, IsNil)
	s.db = db
	s.ts = httptest.NewServer(engine)
}

func (s *MibigSuite) TearDownSuite(c *C) {
	s.db.Close()
	s.ts.Close()
	os.Remove(s.dbFile.Name())
}

func (s *MibigSuite) TearDownTest(c *C) {
	s.db.Exec("TRUNCATE TABLE submissions; TRUNCATE TABLE nrps_submissions; TRUNCATE TABLE gene_submissions")
}
