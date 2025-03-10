package web

import (
	"log"
	"log/slog"
	"magnet-parser/globals"
	"os"
	"strings"
)

func appendLineToFile(magnet string) {
	f, err := os.OpenFile("magnets.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if _, err := f.Write([]byte(magnet + "\n\n")); err != nil {
		slog.Error(err.Error())
		return
	}
	if err := f.Close(); err != nil {
		slog.Error(err.Error())
		return
	}
}

func submain() {
	globals.SetupLogger("web_spider")

	os.Remove("magnets.txt")
	links := []string{"https://igg-games.ru/"}
	index := 0

	for index < len(links) {
		url := links[index]
		urlOrigin := ParseOriginUrl(url)
		html := GetHtmlByUrl(url)

		// add all found links to queue
		newLinks := ParseAllLinks(html)
		for _, newLink := range newLinks {
			if strings.HasPrefix(newLink, "http") {
				if !globals.IsItemInArray(newLink, links) {
					links = append(links, newLink)
				}
			} else if strings.HasPrefix(newLink, "/") && urlOrigin != "" {
				if !globals.IsItemInArray(newLink, links) {
					links = append(links, urlOrigin+newLink)
				}
			} else if strings.HasPrefix(newLink, "magnet:?xt=urn:") {
				appendLineToFile(newLink)
			}
		}
		index++
	}

	log.Println("Done.")
}
