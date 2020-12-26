package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"go-kit-demo/model"
)

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVO struct {
	Username string
	Password string
	Email    string
}

type UserService interface {
	// 登录接口
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)
	// 注册接口
	Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error)
}

type UserServiceImpl struct {
	userDao model.User
}

func MakeUserServiceImpl(userDAO model.User) UserService {
	return &UserServiceImpl{
		userDao: userDAO,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {
	user, err := userService.userDao.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		fmt.Sprintf("err : %s", err)
	}
	return nil, err
}

func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {
	existUser, err := userService.userDao.SelectByEmail(vo.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &model.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}
		err = userService.userDao.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err
}