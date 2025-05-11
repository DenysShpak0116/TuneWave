package service

import (
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type UserService struct {
	*GenericService[models.User]
	SongService services.SongService
}

func NewUserService(repo port.Repository[models.User], songService services.SongService) services.UserService {
	return &UserService{
		GenericService: NewGenericService(repo),
		SongService:    songService,
	}
}

func (us *UserService) GetUsers(
	ctx context.Context,
	page int,
	limit int,
) ([]dtos.UserDTO, error) {
	users, err := us.Repository.NewQuery(ctx).Take(limit).Skip((page - 1) * limit).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var usersDTOs []dtos.UserDTO
	for _, user := range users {
		usersDTOs = append(usersDTOs, dtos.UserDTO{
			ID:             user.ID,
			Username:       user.Username,
			Role:           user.Role,
			ProfilePicture: user.ProfilePicture,
		})
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no users found")
	}

	return usersDTOs, nil
}

func (us *UserService) GetFullDTOByID(ctx context.Context, id uuid.UUID) (*dtos.UserExtendedDTO, error) {
	users, err := us.Repository.NewQuery(ctx).Where("id = ?", id).
		Preload("Songs").Preload("Collections").Preload("Collections.User").
		Preload("Chats1").Preload("Chats2").Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := users[0]

	songsDTOs := make([]dtos.SongDTO, len(user.Songs))
	for i, song := range user.Songs {
		songExtDTO, err := us.SongService.GetFullDTOByID(ctx, song.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get song by ID: %w", err)
		}
		songsDTOs[i] = dtos.SongDTO{
			ID:         songExtDTO.ID,
			Duration:   songExtDTO.Duration,
			Title:      songExtDTO.Title,
			SongURL:    songExtDTO.SongURL,
			CoverURL:   songExtDTO.CoverURL,
			Listenings: songExtDTO.Listenings,
			Likes:      songExtDTO.Likes,
			Dislikes:   songExtDTO.Dislikes,
			User:       songExtDTO.User,
		}
	}

	collectinsDTOs := make([]dtos.CollectionDTO, len(user.Collections))
	for i, collection := range user.Collections {
		collectinsDTOs[i] = dtos.CollectionDTO{
			ID:       collection.ID,
			Title:    collection.Title,
			CoverURL: collection.CoverURL,
			User: dtos.UserDTO{
				ID:             collection.User.ID,
				Username:       collection.User.Username,
				ProfilePicture: collection.User.ProfilePicture,
			},
		}
	}

	chatDTOs := make([]dtos.ChatDTO, len(user.Chats1))
	for i, chat := range user.Chats1 {
		chatDTOs[i] = dtos.ChatDTO{
			ID: chat.ID,
			User1: dtos.UserDTO{
				ID:             chat.User1.ID,
				Username:       chat.User1.Username,
				ProfilePicture: chat.User1.ProfilePicture,
			},
		}
	}

	userDTO := &dtos.UserExtendedDTO{
		ID:             user.ID,
		Username:       user.Username,
		Role:           user.Role,
		ProfileInfo:    user.ProfileInfo,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
		Songs:          songsDTOs,
		Collections:    collectinsDTOs,
		Chats:          chatDTOs,
	}

	return userDTO, nil
}

func (us *UserService) UpdateUserPassword(email string, hashedPassword string) error {
	user, err := us.Repository.NewQuery(context.Background()).Where("email = ?", email).Find()
	if err != nil {
		return err
	}

	if len(user) == 0 {
		return fmt.Errorf("user not found")
	}

	user[0].PasswordHash = hashedPassword

	_, err = us.Update(context.TODO(), &user[0])
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}
