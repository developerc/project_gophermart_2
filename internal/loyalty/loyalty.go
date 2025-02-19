package loyalty

import (
	"database/sql"
	"log"
	"time"

	dbstorage "github.com/developerc/project_gophermart_2/internal/db_storage"
	"github.com/developerc/project_gophermart_2/internal/general"
)

func RunLoyalty(chSignal general.ChSignal, db *sql.DB, adresAccrual string) {
	go func() {
		for {
			chanCnt := 5
			select {
			case <-chSignal.ChStart:
				arrOrderNumb, err := dbstorage.GetOrderNumbs(db)
				if err != nil {
					log.Println(err)
					continue
				}
				DoRequests(db, chanCnt, arrOrderNumb, adresAccrual, chSignal)
			case pause := <-chSignal.ChPause:
				time.Sleep(time.Duration(pause) * time.Second)
			}
		}
	}()
}
