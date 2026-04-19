package dto

type SongSearchResponse struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
	Image  string `json:"image"`
}

type SongDetailResponse struct {
	Title           string               `json:"title"`
	Artist          string               `json:"artist"`
	Key             string               `json:"key"`
	Chords          []string             `json:"chords"`
	Content         string               `json:"content"`
	SimplifiedUrl   string               `json:"simplified_url,omitempty"`
	PrincipalUrl    string               `json:"principal_url,omitempty"`
	KeyboardUrl     string               `json:"keyboard_url,omitempty"`
	Recommendations []SongSearchResponse `json:"recommendations"`
}

