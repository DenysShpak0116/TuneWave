package service

import (
	"context"
	"fmt"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
)

type UserService struct {
	*GenericService[models.User]
}

func NewUserService(repo port.Repository[models.User]) *UserService {
	return &UserService{
		GenericService: NewGenericService(repo),
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
			ProfilePicture: user.ProfilePicture,
		})
		return usersDTOs, nil
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no users found")
	}

	return usersDTOs, nil
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
