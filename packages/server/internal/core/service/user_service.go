package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type UserService struct {
	*GenericService[models.User]
	SongService             services.SongService
	FileStorage             port.FileStorage
	UserFollowersRepository port.Repository[models.UserFollower]
}

func NewUserService(
	repo port.Repository[models.User],
	songService services.SongService,
	fileStorage port.FileStorage,
	userFollowersRepository port.Repository[models.UserFollower],
) services.UserService {
	return &UserService{
		GenericService:          NewGenericService(repo),
		SongService:             songService,
		FileStorage:             fileStorage,
		UserFollowersRepository: userFollowersRepository,
	}
}

func (us *UserService) GetUsers(
	ctx context.Context,
	page int,
	limit int,
) ([]models.User, error) {
	users, err := us.Repository.NewQuery(ctx).Take(limit).Skip((page - 1) * limit).Preload("Followers").Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (us *UserService) UpdateUserPassword(email string, hashedPassword string) error {
	users, err := us.Repository.NewQuery(context.Background()).Where("email = ?", email).Find()
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("user not found")
	}
	user := &users[0]

	user.PasswordHash = hashedPassword

	if err := us.Update(context.TODO(), user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}

func (us *UserService) UpdateUserPfp(ctx context.Context, pfpParams services.UpdatePfpParams) error {
	users, err := us.Repository.NewQuery(ctx).
		Where("id = ?", pfpParams.UserID).
		Find()
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if len(users) == 0 {
		return fmt.Errorf("user not found")
	}
	user := &users[0]

	if pfpParams.Pfp != nil && pfpParams.PfpHeader != nil {
		oldUserPfpKey := helpers.ExtractS3Key(user.ProfilePicture)
		if err := us.FileStorage.Remove(ctx, oldUserPfpKey); err != nil {
			return fmt.Errorf("failed to remove old user file: %w", err)
		}

		key := fmt.Sprintf("pfp/%s/%d-%s", user.ID, time.Now().Unix(), pfpParams.PfpHeader.Filename)
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, pfpParams.Pfp); err != nil {
			return err
		}

		url, err := us.FileStorage.Save(ctx, key, buf)
		if err != nil {
			return err
		}

		user.ProfilePicture = url
	}

	if err := us.Repository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetUserFollowersCount(ctx context.Context, userID uuid.UUID) int64 {
	count, _ := us.UserFollowersRepository.NewQuery(ctx).
		Where(&models.UserFollower{UserID: userID}).Count()

	return count
}
