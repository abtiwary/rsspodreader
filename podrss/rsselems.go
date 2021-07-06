package podrss

type ChannelImage struct {
	URL string `xml:"url"`
}

type PodEnclosure struct {
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
	URL    string `xml:"url,attr"`
}

type PodItem struct {
	Title       string       `xml:"title"`
	PubDate     string       `xml:"pubDate"`
	Description string       `xml:"description"`
	Enclosure   PodEnclosure `xml:"enclosure"`
}

type Channel struct {
	Title   string       `xml:"title"`
	PubDate string       `xml:"pubDate"`
	Image   ChannelImage `xml:"image"`
	Items   []PodItem    `xml:"item"`
}

type PodChannel struct {
	PodChan Channel `xml:"channel"`
}
