package main

import "errors"

// Define error variables for standardized errors
var (
	NodeDead  = errors.New("Node is already Dead")
	NodeALive = errors.New("Node is already Alive")
)
