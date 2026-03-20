package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	repo *Repository
}

func NewUserService(repo *Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(userID int64) (*model.UserResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		OpenID:    user.OpenID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Phone:     user.Phone,
	}, nil
}

func (s *UserService) UpdateProfile(userID int64, nickname, avatarURL string) error {
	return s.repo.UpdateUserProfile(userID, nickname, avatarURL)
}
