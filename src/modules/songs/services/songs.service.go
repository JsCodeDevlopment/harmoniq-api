package services

import (
	"api/src/modules/songs/dto"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SongsService struct{}

func NewSongsService() *SongsService {
	return &SongsService{}
}

func (s *SongsService) Search(query string) ([]dto.SongSearchResponse, error) {
	url := fmt.Sprintf("https://www.cifraclub.com.br/?q=%s", strings.ReplaceAll(query, " ", "+"))
	
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cifraclub returned status %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var songs []dto.SongSearchResponse
	doc.Find(".gsc-result").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".gs-title").Text()
		link, _ := s.Find(".gs-title a").Attr("href")
		if title != "" && link != "" {
			songs = append(songs, dto.SongSearchResponse{
				Title: title,
				Url:   link,
			})
		}
	})

	// Alternative search if the above doesn't work (CifraClub uses GCS for search in some pages)
	// But let's try to parse the direct results if available.
	// Actually, cifraclub often has a list of results in a different format.
	
	return songs, nil
}

func (s *SongsService) GetSong(url string) (*dto.SongDetailResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cifraclub returned status %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find(".t1").Text()
	artist := doc.Find(".t3").Text()
	key := doc.Find("#cifra_tom a").Text()
	
	content, _ := doc.Find("pre").Html()
	
	var chords []string
	doc.Find("#cifra_capo").NextAll().Find("b").Each(func(i int, s *goquery.Selection) {
		chords = append(chords, s.Text())
	})

	return &dto.SongDetailResponse{
		Title:   strings.TrimSpace(title),
		Artist:  strings.TrimSpace(artist),
		Key:     strings.TrimSpace(key),
		Chords:  chords,
		Content: content,
	}, nil
}
