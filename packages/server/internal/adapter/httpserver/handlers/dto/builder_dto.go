package dto

import "github.com/google/uuid"

type CountUserFollowersFunc func(userID uuid.UUID) int64
type CountSongLikesFunc func(songID uuid.UUID) int64
type CountSongDislikesFunc func(songID uuid.UUID) int64

type DTOBuilder struct {
	CountUserFollowers CountUserFollowersFunc
	CountSongLikes     CountSongLikesFunc
	CountSongDislikes  CountSongDislikesFunc
}

func NewDTOBuilder() *DTOBuilder {
	return &DTOBuilder{
		CountUserFollowers: func(userID uuid.UUID) int64 { return 0 },
		CountSongLikes:     func(songID uuid.UUID) int64 { return 0 },
		CountSongDislikes:  func(songID uuid.UUID) int64 { return 0 },
	}
}

func (b *DTOBuilder) SetCountUserFollowersFunc(
	countUserFollowersFunc CountUserFollowersFunc,
) *DTOBuilder {
	b.CountUserFollowers = countUserFollowersFunc
	return b
}

func (b *DTOBuilder) SetCountSongLikesFunc(
	countSongLikesFunc CountSongLikesFunc,
) *DTOBuilder {
	b.CountSongLikes = countSongLikesFunc
	return b
}

func (b *DTOBuilder) SetCountSongDislikesFunc(
	countSongDislikesFunc CountSongDislikesFunc,
) *DTOBuilder {
	b.CountSongDislikes = countSongDislikesFunc
	return b
}
