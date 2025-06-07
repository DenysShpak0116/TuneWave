package dto

import (
	"time"

	"github.com/google/uuid"
)

type SongDTO struct {
	ID         uuid.UUID    `json:"id"`
	CreatedAt  time.Time    `json:"createdAt"`
	Duration   string       `json:"duration"`
	Title      string       `json:"title"`
	Genre      string       `json:"genre"`
	SongURL    string       `json:"songUrl"`
	CoverURL   string       `json:"coverUrl"`
	Listenings int64        `json:"listenings"`
	Likes      int64        `json:"likes"`
	Dislikes   int64        `json:"dislikes"`
	User       UserDTO      `json:"user"`
	Authors    []AuthorDTO  `json:"authors"`
	SongTags   []TagDTO     `json:"songTags"`
	Comments   []CommentDTO `json:"comments"`
}
