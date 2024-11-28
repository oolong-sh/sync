package main

import (
	"fmt"
	"log"

	"github.com/oolong-sh/sync"
)

func main() {
	// parse config file
	configPath := "./oolong-sync.toml"
	cfg, err := sync.ReadConfig(configPath)
	if err != nil {
		log.Println(err)
		return
	}

	// initialize a new sync client
	c, err := sync.NewClient(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	// operations can now be performed
	fp := fmt.Sprintf("/home/%s/test.txt", cfg.User)
	if err := c.Update(fp); err != nil {
		log.Println(err)
	}

	// close sftp and ssh connections when finished
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}
