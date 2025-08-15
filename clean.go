package main

import (
	"net/url"
)

func clean(host string, href []string) []string {
	var cleaned []string
	for _, h := range href {
		cleaned = append(cleaned, cleanOne(host, h))
	}
	return cleaned
}

func cleanOne(host string, href string) string {
	hostUrl, err := url.Parse(host)
	if err != nil {
		return "Error parsing host"
	}

	hrefUrl, err := url.Parse(href)
	if err != nil {
		return "Error parsing href"
	}

	if hrefUrl.Host != "" {
		return hrefUrl.String()
	}

	hostUrl.Path = hrefUrl.Path // works online
	// hostUrl.Path = path.Clean(path.Join(hostUrl.Path, hrefUrl.Path)) // works on localhost
	// hostUrl.RawQuery = hrefUrl.RawQuery

	return hostUrl.String()
}
