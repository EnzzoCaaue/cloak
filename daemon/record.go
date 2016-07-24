package daemon

import (
	"log"
	"time"

	"github.com/Cloakaac/cloak/models"
)

type RecordDaemon struct{}

func (r *RecordDaemon) tick() {
	total := models.GetOnlineCount()
	err := models.AddOnlineRecord(total, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
}
