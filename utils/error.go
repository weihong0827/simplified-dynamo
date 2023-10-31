package utils

import "errors"

// Define error variables for standardized errors
var (
	ErrNoNodesAvailable = errors.New("no nodes available")
	ErrInvalidParameter = errors.New("invalid parameter")
)
