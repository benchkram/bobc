package application

import (
	"errors"
)

var (
	ErrArtifactNotFound      = errors.New("artifact not found")
	ErrArtifactAlreadyExists = errors.New("artifact already exists")
	ErrProjectNotFound       = errors.New("project not found")
	ErrProjectAlreadyExists  = errors.New("project already exists")
	ErrTokenAlreadyExists    = errors.New("access token already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidUsername       = errors.New("invalid username")
	ErrInvalidProjectName    = errors.New("invalid project name")
)
