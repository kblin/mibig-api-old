package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/kblin/mibig-api/service"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "mibig"
	app.Usage = "work with the MIBiG API service"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "config.yml", Usage: "config file to use", EnvVar: "MIBIG_CONFIG"},
	}

	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "Run the http API server",
			Action: serve,
		},
		{
			Name:   "migratedb",
			Usage:  "Perform database migrations",
			Action: migrateDb,
		},
	}

	app.Run(os.Args)
}

func getConfig(c *cli.Context) (service.Config, error) {
	yamlPath := c.GlobalString("config")
	config := service.Config{}

	if _, err := os.Stat(yamlPath); err != nil {
		return config, errors.New("Invalid config file path")
	}

	yamlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(yamlData), &config)

	return config, err
}

func serve(c *cli.Context) {
	cfg, err := getConfig(c)
	if err != nil {
		log.Fatal(err)
		return
	}

	svc := service.MibigService{}

	if err = svc.Run(cfg); err != nil {
		log.Fatal(err)
		return
	}
}

func migrateDb(c *cli.Context) {
	cfg, err := getConfig(c)
	if err != nil {
		log.Fatal(err)
		return
	}

	svc := service.MibigService{}

	if err = svc.Migrate(cfg); err != nil {
		log.Fatal(err)
		return
	}
}
