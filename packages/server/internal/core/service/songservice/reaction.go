package songservice

import (
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/google/uuid"
)

func (ss *SongService) SetReaction(ctx context.Context, songID uuid.UUID, userID uuid.UUID, reactionType string) (int64, int64, error) {
	reactions, err := ss.ReactionsRepository.NewQuery(ctx).
		Where("song_id = ? AND user_id = ?", songID, userID).
		Find()
	if err != nil {
		return 0, 0, err
	}

	if len(reactions) == 0 {
		reaction := models.UserReaction{
			SongID: songID,
			UserID: userID,
			Type:   reactionType,
		}
		if err := ss.ReactionsRepository.Add(ctx, &reaction); err != nil {
			return 0, 0, err
		}
	} else {
		existingReaction := reactions[0]
		if existingReaction.Type == reactionType {
			if err := ss.ReactionsRepository.Delete(ctx, existingReaction.ID); err != nil {
				return 0, 0, err
			}
		} else {
			existingReaction.Type = reactionType
			if _, err := ss.ReactionsRepository.Update(ctx, &existingReaction); err != nil {
				return 0, 0, err
			}
		}
	}

	dislikes, err := ss.ReactionsCount(ctx, songID, "dislike")
	if err != nil {
		return 0, 0, err
	}
	likes, err := ss.ReactionsCount(ctx, songID, "like")
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
