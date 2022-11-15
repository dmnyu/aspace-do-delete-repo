package main

import (
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"log"
	"os"
	"strings"
	"time"
)

var (
	repoid      int
	config      string
	environment string
	client      *aspace.ASClient
	test        bool
)

func init() {
	flag.IntVar(&repoid, "repo", 0, "the repository id of to delete to delete DOS from")
	flag.StringVar(&config, "config", "", "the location of a go-aspace config file")
	flag.StringVar(&environment, "environment", "", "the environment to delete files from")
	flag.BoolVar(&test, "test", false, "")
}

func main() {
	//parse flags
	flag.Parse()

	//setup the log
	t := time.Now()
	tf := t.Format("20060102T15:04")
	logFile, err := os.Create(fmt.Sprintf("aspace-do-delete-repo-%s-%s.log", environment, tf))
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	client, err = aspace.NewClient(config, environment, 20)
	if err != nil {
		panic(err)
	}

	doids, err := client.GetDigitalObjectIDs(repoid)
	if err != nil {
		panic(err)
	}

	for _, doid := range doids {
		//get the DO metadata
		domd, err := client.GetDigitalObject(repoid, doid)
		if err != nil {
			doErrMsg := fmt.Sprintf("[ERROR] repo-id: %d do-id: %d error: %s", repoid, doid, strings.ReplaceAll(err.Error(), "\n", " "))
			fmt.Println(doErrMsg)
			log.Println(doErrMsg)
			continue
		}

		//get the uris from the file version
		fileversionUris := ""
		for i, fv := range domd.FileVersions {
			if i > 0 {
				fileversionUris = fileversionUris + ", "
			}
			fileversionUris = fileversionUris + fv.FileURI
		}

		infoMsg := fmt.Sprintf("[INFO] DO-URI: %s, TITLE: %s, FILE-URIS: %s", domd.URI, domd.Title, fileversionUris)
		fmt.Println(infoMsg)
		log.Println(infoMsg)

		//delete the do
		if test == false {
			msg, err := client.DeleteDigitalObject(repoid, doid)
			if err != nil {
				errMsg := fmt.Sprintf("[ERROR] repo-id: %d do-id: %d error: %s", strings.ReplaceAll(err.Error(), "\n", " "))
				fmt.Println(errMsg)
				log.Println(errMsg)
				continue
			} else {
				infoMsg2 := fmt.Sprintf("[INFO] DELETED %s %s", domd.URI, strings.ReplaceAll(msg, "\n", " "))
				fmt.Println(infoMsg2)
				log.Println(infoMsg2)
			}
		}
	}
}
