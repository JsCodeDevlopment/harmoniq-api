package dto

type SongSearchResponse struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
}

type SongDetailResponse struct {
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	Key      string   `json:"key"`
	Chords   []string `json:"chords"`
	Content  string   `json:"content"`
}
