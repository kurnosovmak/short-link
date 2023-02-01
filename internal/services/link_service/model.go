package link_service

type Link struct {
	Link string
	Key  string
}
type CreateLinkDTO struct {
	Link string `json:"link"`
}

type RedirectLinkDTO struct {
	Key string
}
