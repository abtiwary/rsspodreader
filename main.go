package main

import (
	"encoding/xml"
	"fmt"
	"github.com/abtiwary/rsspodreader/podrss"
	"io"
	"os"
	"path/filepath"

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

	rssFilePath := filepath.Join(curDir, "rss_examples", "cppcast.rss")

	rssFile, err := os.Open(rssFilePath)
	if err != nil {
		log.WithError(err).WithField("rss_file", rssFile).Fatal("could not open the rss file")
	}

	rssBytes, err := io.ReadAll(rssFile)
	if err != nil {
		log.WithError(err).WithField("rss_file", rssFile).Fatal("could not read the rss file")
	}

	var podChannel podrss.PodChannel
	err = xml.Unmarshal(rssBytes, &podChannel)
	if err != nil {
		log.WithError(err).Fatal("could not unmarshal the rss file")
	}

	fmt.Println(podChannel.PodChan.Title)
	fmt.Println(podChannel.PodChan.PubDate)

	fmt.Printf("parsed %v items\n", len(podChannel.PodChan.Items))

	for _, itm := range podChannel.PodChan.Items {
		fmt.Println(itm.Title)
		fmt.Println(itm.PubDate)
		//fmt.Println(itm.Description)
		fmt.Printf("%v %v\n\n", itm.Enclosure.URL, itm.Enclosure.Length)
	}

	fmt.Println("done...")
}
