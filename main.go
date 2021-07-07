package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtiwary/rsspodreader/podrss"

	log "github.com/sirupsen/logrus"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	initLogger()

	curDir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("cannot determine current working directory")
	}

	rssWorkDir := filepath.Join(curDir, "tempfiles")
	imagesDir := filepath.Join(curDir, "images")
	audioDir := filepath.Join(curDir, "audio")

	cppCast := podrss.NewPodRss("Cpp Cast",
		"https://cppcast.libsyn.com/rss",
		"74a4a440-fa80-49be-9bbb-589a3b8a7e37",
		imagesDir,
		audioDir,
		rssWorkDir,
	)

	pItems, err := cppCast.GetItems()
	if err != nil {
		log.WithError(err).Fatal("error getting podcast items")
	}

	fmt.Printf("parsed %v items\n", len(pItems))

	for _, itm := range pItems {
		fmt.Println(itm.Title)
		fmt.Println(itm.PubDate)
		//fmt.Println(itm.Description)
		fmt.Printf("%v %v\n\n", itm.Enclosure.URL, itm.Enclosure.Length)
	}

	fmt.Println("done...")
}
