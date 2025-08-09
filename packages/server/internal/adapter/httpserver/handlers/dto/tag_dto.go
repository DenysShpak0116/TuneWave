package dto

import "github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"

type TagDTO struct {
	Name string `json:"name"`
}

func (b *DTOBuilder) BuildTagDTO(songTag *models.SongTag) *TagDTO {
	return &TagDTO{
		Name: songTag.Tag.Name,
	}
}
