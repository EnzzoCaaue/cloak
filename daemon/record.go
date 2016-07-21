package daemon

import (
	"log"
	"time"

	"github.com/Cloakaac/cloak/models"
)

type recordDaemon struct{}

func init() {
	daemons.Add("record", 5*time.Minute, &recordDaemon{})
}

func (r *recordDaemon) tick() {
	total := models.GetOnlineCount()
	err := models.AddOnlineRecord(total, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
}
