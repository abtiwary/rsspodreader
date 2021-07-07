package podrss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/juju/errors"
)

type PodcastFeedReader interface {
	GetItems() ([]PodItem, error)
}

type PodRss struct {
	Title    string
	URL      string
	FileName string
	ImageDir string
	AudioDir string
	WorkDir  string
}

func NewPodRss(title, url, filenm, imgdir, audiodir, workdir string) PodRss {
	return PodRss{
		Title:    title,
		URL:      url,
		FileName: filenm,
		ImageDir: imgdir,
		AudioDir: audiodir,
		WorkDir:  workdir,
	}
}

// GetItems gets the rss file from the URL, if a recent copy does not exist,
// and parses it. It downloads a podcast image if needed, and returns a list
// of all the podcast items.
func (p PodRss) GetItems() ([]PodItem, error) {
	flPath := filepath.Join(p.WorkDir, p.FileName)
	fexists, err := CheckIfFileExists(flPath)
	if err != nil {
		return nil, errors.Annotatef(err, "could not get items")
	}

	// download the rss file from the URL since we don't have one locally
	if !fexists {
		err = DownloadFromURL(p.URL, p.FileName, p.WorkDir)
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

	// channel image
	_, imgFile := filepath.Split(podChannel.PodChan.Image.URL)
	imgFileExt := strings.Split(imgFile, ".")[1]
	imgFileName := fmt.Sprintf("%v.%v", p.FileName, imgFileExt)
	podImageFile := filepath.Join(p.ImageDir, imgFileName)
	imgExists, _ := CheckIfFileExists(podImageFile)
	if !imgExists {
		err = DownloadFromURL(podChannel.PodChan.Image.URL, imgFileName, p.ImageDir)
	}

	podRssAge, podDateParsed := GetAgeRelativeToNow(&podChannel.PodChan.PubDate)
	if podDateParsed {
		fmt.Printf("the published RSS is %v seconds old\n", podRssAge)
	} else {
		fmt.Printf("could not determine the age of the RSS\n")
	}

	return podChannel.PodChan.Items, nil
}

// GetAgeRelativeToNow attempts to parse a given date time string
// and return the difference relative to now in seconds
func GetAgeRelativeToNow(dateTime *string) (float64, bool) {
	dtNow := time.Now()

	testFormats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.RFC3339Nano,
		time.RFC3339,
	}

	var testParseDate time.Time
	var testParseErr error
	parseSuccess := false
	for _, tfmt := range testFormats {
		testParseDate, testParseErr = time.Parse(tfmt, *dateTime)
		if testParseErr != nil {
			fmt.Printf("tried parsing as %v, moving on...\n", tfmt)
		} else {
			parseSuccess = true
			break
		}
	}

	diffDur := dtNow.Sub(testParseDate)
	return diffDur.Seconds(), parseSuccess
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
