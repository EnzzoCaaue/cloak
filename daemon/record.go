package daemon

import (
	"github.com/Cloakaac/cloak/models"
	"log"
	"time"
)

type recordDaemon struct{}

func init() {
	daemons.Add("record", time.Minute, &recordDaemon{})
}

func (r *recordDaemon) tick() {
	total := models.GetOnlineCount()
	err := models.AddOnlineRecord(total, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
}
