package service

import (
	"errors"
	"time"

	"baokaobao/internal/config"
	"baokaobao/internal/model"
	"baokaobao/internal/pkg/jwt"
	"baokaobao/internal/pkg/wechat"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo      *Repository
	jwt       *jwt.JWT
	wechatSDK *wechat.WechatSDK
}

func NewAuthService(repo *Repository, jwt *jwt.JWT, wechatSDK *wechatSDK) *AuthService {
	return &AuthService{
		repo:      repo,
		jwt:       jwt,
		wechatSDK: wechatSDK,
	}
}

type LoginByWechatRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginResponse struct {
	Token string             `json:"token"`
	User  model.UserResponse `json:"user"`
}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string              `json:"token"`
	User  model.AdminResponse `json:"user"`
}

func (s *AuthService) LoginByWechat(code string) (*LoginResponse, error) {
	result, err := s.wechatSDK.Code2Session(code)
	if err != nil {
		return nil, errors.New("微信登录失败: " + err.Error())
	}

	user, err := s.repo.GetUserByOpenID(result.OpenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &model.User{
				OpenID:  result.OpenID,
				UnionID: result.UnionID,
				Status:  1,
			}
			if err := s.repo.CreateUser(user); err != nil {
				return nil, errors.New("创建用户失败")
			}
		} else {
			return nil, errors.New("查询用户失败")
		}
	}

	if user.Status != 1 {
		return nil, model.ErrUserBanned
	}

	user.LastLogin = time.Now()
	s.repo.UpdateUser(user)

	token, err := s.jwt.GenerateToken(user.ID, user.OpenID, "mini")
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &LoginResponse{
		Token: token,
		User: model.UserResponse{
			ID:        user.ID,
			OpenID:    user.OpenID,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarURL,
			Phone:     user.Phone,
		},
	}, nil
}

func (s *AuthService) AdminLogin(req *AdminLoginRequest) (*AdminLoginResponse, error) {
	admin, err := s.repo.GetAdminByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, errors.New("查询用户失败")
	}

	if admin.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return nil, model.ErrInvalidCredentials
	}

	admin.LastLogin = time.Now()
	s.repo.UpdateAdmin(admin)

	token, err := s.jwt.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &AdminLoginResponse{
		Token: token,
		User: model.AdminResponse{
			ID:       admin.ID,
			Username: admin.Username,
			Nickname: admin.Nickname,
			Role:     admin.Role,
		},
	}, nil
}

func (s *AuthService) DecryptPhone(userID int64, code string) (string, error) {
	result, err := s.wechatSDK.GetPhoneNumber(code)
	if err != nil {
		return "", errors.New("获取手机号失败: " + err.Error())
	}

	phone := result.PhoneInfo.PhoneNumber
	if err := s.repo.UpdateUserPhone(userID, phone); err != nil {
		return "", errors.New("更新手机号失败")
	}

	return phone, nil
}

func (s *AuthService) GetWechatSDK() *wechatSDK {
	return s.wechatSDK
}
