package model

type URL struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Short   string `json:"short"`
	Created string `json:"created"`
	UserID  string `json:"user_id"`
}
