package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

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

func (s *UserService) UploadAvatar(userID int64, file io.Reader, filename string) (string, error) {
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("create upload dir failed: %w", err)
	}

	dst := filepath.Join(uploadDir, filename)
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create file failed: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("save file failed: %w", err)
	}

	url := fmt.Sprintf("/uploads/avatars/%s", filename)

	s.repo.UpdateUserProfile(userID, "", url)

	return url, nil
}

func (s *UserService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.repo.ListUsers(page, pageSize)
}

func (s *UserService) GetUser(id int64) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUserStatus(id int64, status int8) error {
	return s.repo.UpdateUserStatus(id, status)
}
