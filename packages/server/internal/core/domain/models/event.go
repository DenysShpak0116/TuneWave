package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type EventType string

const (
	PlayEvent               EventType = "play"
	AddToPlaylistEvent      EventType = "add_to_playlist"
	RemoveFromPlaylistEvent EventType = "remove_from_playlist"
	CreatePlaylistEvent     EventType = "create_playlist"
	CommentEvent            EventType = "comment"
	LikeEvent               EventType = "like"
	ViewPlaylistEvent       EventType = "view_playlist"
)

type Event struct {
	BaseModel

	UserID uuid.UUID `json:"userId"`
	User   User      `gorm:"foreignKey:UserID"`

	EventType EventType `json:"eventType"`

	TrackID *uuid.UUID `json:"trackId,omitempty"`
	Track   *Song      `gorm:"foreignKey:TrackID"`

	PlaylistID *uuid.UUID  `json:"playlistId,omitempty"`
	Playlist   *Collection `gorm:"foreignKey:PlaylistID"`

	CommentID *uuid.UUID `json:"commentId,omitempty"`

	Metadata datatypes.JSONMap `json:"metadata,omitempty"`
}
