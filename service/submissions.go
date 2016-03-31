package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kblin/mibig-api/api"
	"log"
	"strconv"
	"time"
)

type MibigResource struct {
	db *sql.DB
}

func (mr *MibigResource) ServiceInfo(c *gin.Context) {
	c.String(200, "bgc registration up and running")
}

func (mr *MibigResource) StoreMibigSubmissionV1(c *gin.Context) {
	mr.storeMibigSubmissionGeneric(c, "//thankyou.html", false)
}

func (mr *MibigResource) StoreMibigSubmissionV2(c *gin.Context) {
	mr.storeMibigSubmissionGeneric(c, "//genes_form.html", true)
}

func (mr *MibigResource) storeMibigSubmissionGeneric(c *gin.Context, forward_to string, redirect bool) {
	mibigJson := c.PostForm("json")
	versionString := c.PostForm("version")

	if mibigJson == "" {
		log.Println("no json")
		c.JSON(400, gin.H{"error": true, "message": "json not provided"})
		return
	}

	if versionString == "" {
		log.Println("no version")
		c.JSON(400, gin.H{"error": true, "message": "Version parameter not provided. Need a version parameter greater than 0"})
		return
	}

	version, err := strconv.ParseInt(versionString, 10, 32)
	if err != nil {
		log.Println(versionString)
		c.JSON(400, gin.H{"error": true, "message": "Version parameter not a valid number"})
		return
	}

	if version <= 0 {
		log.Println(version)
		c.JSON(400, gin.H{"error": true, "message": "Need a version parameter greater than 0"})
		return
	}

	mibig := &api.MibigSubmission{
		Submitted: time.Now().UTC(),
		Modified:  time.Now().UTC(),
		Raw:       mibigJson,
		Version:   int(version),
	}

	if err := mibig.Create(mr.db); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": true, "message": "could not store submission"})
		return
	}

	if redirect {
		c.Redirect(303, forward_to)
	} else {
		c.JSON(200, gin.H{"error": false, "message": "bgc registration successful.", "redirect_url": forward_to})
	}
}

func (mr *MibigResource) StoreBgcDetailSubmission(c *gin.Context) {
	data := c.PostForm("data")
	target := c.PostForm("target")
	versionString := c.DefaultPostForm("version", "1")
	bgc_id := c.DefaultPostForm("bgc_id", "BGC00000")

	if data == "" {
		c.JSON(400, gin.H{"error": true, "message": "json not provided"})
		return
	}

	if target == "" {
		c.JSON(400, gin.H{"error": true, "message": "target not provided"})
		return
	}

	version, err := strconv.ParseInt(versionString, 10, 32)
	if err != nil {
		log.Println(versionString)
		c.JSON(400, gin.H{"error": true, "message": "Version parameter not a valid number"})
		return
	}

	if version <= 0 {
		c.JSON(400, gin.H{"error": true, "message": "Need a version parameter greater than 0"})
		return
	}

	bgc := &api.BgcDetailSubmission{
		BgcId:     bgc_id,
		Submitted: time.Now().UTC(),
		Modified:  time.Now().UTC(),
		Raw:       data,
		Version:   int(version),
	}

	if target == "gene_info" {
		err = bgc.Create(mr.db, "gene_submissions")
	} else if target == "nrps_info" {
		err = bgc.Create(mr.db, "nrps_submissions")
	} else {
		c.JSON(400, gin.H{"error": true, "message": "target parameter not matching. Must be one of 'gene_info' or 'nrps_info'"})
		return
	}

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": true, "message": "could not store submission"})
		return
	}

	c.AbortWithStatus(204)
}
