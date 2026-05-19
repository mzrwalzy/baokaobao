package model

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrAdminNotFound      = errors.New("admin user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBanned         = errors.New("user is banned")
	ErrQuestionNotFound   = errors.New("question not found")
	ErrBankNotFound       = errors.New("question bank not found")
	ErrNoAccessToBank     = errors.New("no access to this question bank")
	ErrBankNotPurchased   = errors.New("please purchase this bank first")
	ErrInvalidAnswer      = errors.New("invalid answer format")
)
