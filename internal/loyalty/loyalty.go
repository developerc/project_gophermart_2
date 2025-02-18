package loyalty

import (
	"database/sql"
	"log"
	"time"

	dbstorage "github.com/developerc/project_gophermart_2/internal/db_storage"
	"github.com/developerc/project_gophermart_2/internal/general"
)

var retryAfterSec int

/*func RunLoyalty(db *sql.DB, adresAccrual string) {
	go func() {
		for {
			chanCnt := 5
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			//defer cancel()
			arrOrderNumb, err := dbstorage.GetOrderNumbs(ctx, db)
			if err != nil {
				log.Println(err)
				continue
			}
			DoRequests(db, chanCnt, arrOrderNumb, adresAccrual)
			time.Sleep(time.Duration(retryAfterSec) * time.Second)
			retryAfterSec = 0
			cancel()
			time.Sleep(2 * time.Second)
		}
	}()
}*/

func RunLoyalty2(chSignal general.ChSignal, db *sql.DB, adresAccrual string) {
	go func() {
		for {
			chanCnt := 5
			//ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			/*select {
			case <-chSignal.ChStart:
				//fmt.Println("received POST /api/user/orders")
				arrOrderNumb, err := dbstorage.GetOrderNumbs(db)
				if err != nil {
					log.Println(err)
					continue
				}
				DoRequests(db, chanCnt, arrOrderNumb, adresAccrual)
			}*/
			//cancel()
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
