package podrss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

type PodcastFeedReader interface {
	GetItems() ([]PodItem, error)
}

type PodRss struct {
	Title      string
	URL        string
	FileName   string
	WorkingDir string
}

func NewPodRss(title, url, filenm, workdir string) PodRss {
	return PodRss{
		Title:      title,
		URL:        url,
		FileName:   filenm,
		WorkingDir: workdir,
	}
}

func (p PodRss) GetItems() ([]PodItem, error) {
	flPath := filepath.Join(p.WorkingDir, p.FileName)
	fexists, err := CheckIfFileExists(flPath)
	if err != nil {
		return nil, errors.Annotatef(err, "could not get items")
	}

	if !fexists {
		err = DownloadFromURL(p.URL, p.FileName, p.WorkingDir)
		if err != nil {
			return nil, err
		}
	}

	rssFile, err := os.Open(flPath)
	if err != nil {
		return nil, errors.Annotatef(err, "could not open the rss file %v", flPath)
	}

	rssBytes, err := io.ReadAll(rssFile)
	if err != nil {
		return nil, errors.Annotatef(err, "could not read the rss file %v", flPath)
	}

	var podChannel PodChannel
	err = xml.Unmarshal(rssBytes, &podChannel)
	if err != nil {
		return nil, errors.Annotatef(err, "could not unmarshal the rss file")
	}

	return podChannel.PodChan.Items, nil
}

func CheckIfFileExists(pathToFile string) (bool, error) {
	_, err := os.Stat(pathToFile)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, errors.Errorf("could not determine if file exists: %v", err.Error())
	}

	return true, nil
}

func DownloadFromURL(url, filenm, workDir string) error {
	flPath := filepath.Join(workDir, filenm)

	out, err := os.Create(flPath)
	if err != nil {
		return errors.Annotatef(err, "could not create file (%v) in directory (%v)", filenm, workDir)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return errors.Annotatef(err, "could not get file (%v) from url (%v)", filenm, url)
	}
	defer resp.Body.Close()

	nBytes, err := io.Copy(out, resp.Body)
	fmt.Printf("a total of %v bytes were written to %v", nBytes, flPath)
	return nil
}
