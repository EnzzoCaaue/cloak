package app

import (
	"time"

	"github.com/yaimko/yaimko"
	"github.com/yaimko/yaimko/cache"
	"github.com/yaimko/yaimko/session"
)

// Start executes and loads the whole cloak app
func Start() {
	cacheInstance, err := cache.Drivers.NewCache("memory", 5*time.Minute)
	if err != nil {
		yaimko.ERROR.Fatal(err)
	}
	sess, err := session.Storages.NewSession("memory", "testing", 160, yaimko.Config.StringDefault("cookie.key", "AES256Key-32Characters1234567890"))
	if err != nil {
		yaimko.ERROR.Fatal(err)
	}
	yaimko.Filters.Add(func(c *yaimko.Controller) *yaimko.Result {
		sess, err := sess.Open(c.Request, c.Response)
		if err != nil {
			yaimko.ERROR.Fatal(err)
			return nil
		}
		c.Session = sess
		c.Cache = cacheInstance
		return nil
	})
	//	yaimko.Route.Get("/", controllers.LoginController.Login)
}
