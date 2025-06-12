package dto

import "github.com/google/uuid"

type CountUserFollowersFunc func(userID uuid.UUID) int64
type CountSongLikesFunc func(songID uuid.UUID) int64
type CountSongDislikesFunc func(songID uuid.UUID) int64

type dtoContext struct {
	CountUserFollowers CountUserFollowersFunc
	CountSongLikes     CountSongLikesFunc
	CountSongDislikes  CountSongDislikesFunc
}

type DTOBuilder struct {
	dtoCtx dtoContext
}

func NewDTOBuilder() *DTOBuilder {
	return &DTOBuilder{
		dtoCtx: dtoContext{
			CountUserFollowers: func(userID uuid.UUID) int64 { return 0 },
			CountSongLikes:     func(songID uuid.UUID) int64 { return 0 },
			CountSongDislikes:  func(songID uuid.UUID) int64 { return 0 },
		},
	}
}

func (b *DTOBuilder) SetCountUserFollowersFunc(
	countUserFollowersFunc CountUserFollowersFunc,
) *DTOBuilder {
	b.dtoCtx.CountUserFollowers = countUserFollowersFunc
	return b
}

func (b *DTOBuilder) SetCountSongLikesFunc(
	countSongLikesFunc CountSongLikesFunc,
) *DTOBuilder {
	b.dtoCtx.CountSongLikes = countSongLikesFunc
	return b
}

func (b *DTOBuilder) SetCountSongDislikesFunc(
	countSongDislikesFunc CountSongDislikesFunc,
) *DTOBuilder {
	b.dtoCtx.CountSongDislikes = countSongDislikesFunc
	return b
}

func (b *DTOBuilder) Build() (*dtoContext, error) {
	// Possible some validations

	return &b.dtoCtx, nil
}
