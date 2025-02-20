package loyalty

import (
	"database/sql"
	"log"
	"sync/atomic"
	"time"

	dbstorage "github.com/developerc/project_gophermart_2/internal/db_storage"
	"github.com/developerc/project_gophermart_2/internal/general"
)

/*func RunLoyalty(chSignal general.ChSignal, db *sql.DB, adresAccrual string) {
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
}*/

var sleepSec int32

func RunLoyalty(chSignal general.ChSignal, db *sql.DB, adresAccrual string) {
	atomic.StoreInt32(&sleepSec, 0)
	jobs := make(chan int)
	for i := 0; i < 5; i++ {
		go worker(db, adresAccrual, chSignal, jobs)
	}
	go func() {
		for {
			select {
			case <-chSignal.ChStart:
				arrOrderNumb, err := dbstorage.GetOrderNumbs(db)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, orderNumb := range arrOrderNumb {
					jobs <- orderNumb
				}
			case pause := <-chSignal.ChPause:
				//устанавливаем атомик
				atomic.StoreInt32(&sleepSec, int32(pause))
				time.Sleep(time.Duration(pause) * time.Second)
				//сбрасываем атомик
				atomic.StoreInt32(&sleepSec, 0)
			}
		}
	}()

}
