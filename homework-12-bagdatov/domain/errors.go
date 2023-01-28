package domain

import "errors"

// ErrNotAuthorized - user credentials are incorrect.
var ErrNotAuthorized = errors.New("invalid user credentials")

// ErrExists - object with these credentials already exists.
var ErrExists = errors.New("object already exists")

// ErrNotFound - object with these credentials does not exist.
var ErrNotFound = errors.New("object not exists")

// ErrInitilized - mining process is already started.
var ErrInitilized = errors.New("mining is alredy started")

// ErrNotInitilized - mining process is not run.
var ErrNotInitilized = errors.New("mining is not run")

var ErrNoHeader = errors.New("authorization header not found")

var ErrInvalidHeader = errors.New("invalid authorization header")

var ErrInvalidToken = errors.New("invalid token")

var ErrExpiredToken = errors.New("expired token")

var ErrTokenNotCreated = errors.New("failed to create token")
