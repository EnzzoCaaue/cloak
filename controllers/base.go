package controllers

import (
	"github.com/Cloakaac/cloak/util"
)

// BaseController is the main controller that all controllers extend
type BaseController struct {
	Session *util.Session
}