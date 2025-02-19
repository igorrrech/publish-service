package handlers

import "errors"

var (
	ErrNotAllowedByRol             = errors.New("operation not allowed by user rol")
	ErrUnknownUser                 = errors.New("can not find user in context")
	ErrCanNotParseUser             = errors.New("can not parse user")
	ErrUserNotInGroup              = errors.New("user not in the group")
	ErrUserHaveNotAccessToTheGroup = errors.New("user have not access to the group")
)
