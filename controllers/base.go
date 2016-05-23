package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
)

// BaseController is the main controller that all controllers extend
type BaseController struct {
	Session *util.Session
	Data map[interface{}]interface{}
	Template string
	Error string
	Redirect string
	Account *models.CloakaAccount
}