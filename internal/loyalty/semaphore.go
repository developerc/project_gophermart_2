package loyalty

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	dbstorage "github.com/developerc/project_gophermart_2/internal/db_storage"
	"github.com/developerc/project_gophermart_2/internal/general"
)

/*type Semaphore struct {
	semaCh chan struct{}
}

func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaCh
}

func DoRequests(db *sql.DB, chanCnt int, arrOrderNumb []int, adresAccrual string, chSignal general.ChSignal) {
	var wg sync.WaitGroup
	semaphore := NewSemaphore(chanCnt)
	for idx := 0; idx < len(arrOrderNumb); idx++ {
		select {
		case pause := <-chSignal.ChPause:
			time.Sleep(time.Duration(pause) * time.Second)
		default:
			wg.Add(1)
			go func(orderNumb int) {
				semaphore.Acquire()
				defer wg.Done()
				defer semaphore.Release()
				if err := ReqLoyalty(db, adresAccrual, orderNumb, chSignal); err != nil {
					log.Println(err)
				}
			}(arrOrderNumb[idx])
		}
		wg.Wait()

	}
}*/

func ReqLoyalty(db *sql.DB, adresAccrual string, orderNumb int, chSignal general.ChSignal) error {
	response, err := http.Get(adresAccrual + "/api/orders/" + strconv.FormatInt(int64(orderNumb), 10))
	if err != nil {
		log.Println(err)
		return err
	}
	stCode := response.StatusCode
	//fmt.Println("ReqLoyalty stCode: ", stCode)
	if stCode == 204 {
		return errors.New("response Status code: 204")
	}
	if stCode == 500 {
		return errors.New("response Status code: 500")
	}
	if stCode == 429 {
		retryAfter := response.Header.Get("Retry-After")
		retryAfterSec, err := strconv.Atoi(retryAfter)
		if err != nil {
			retryAfterSec = 0
		}
		chSignal.ChPause <- retryAfterSec
		return errors.New("response Status code: 429")
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	loyaltyOrder := general.LoyaltyOrder{}
	if err = json.Unmarshal(body, &loyaltyOrder); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err = dbstorage.SetStatusAccrual(ctx, db, loyaltyOrder.Order, loyaltyOrder.Status, loyaltyOrder.Accrual)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// --
func worker(db *sql.DB, adresAccrual string, chSignal general.ChSignal, jobs <-chan int) {
	for orderNumb := range jobs {
		//fmt.Println("worker: ", orderNumb)
		//проверяем атомик с временем задержки
		if atomic.LoadInt32(&sleepSec) > 0 {
			time.Sleep(time.Duration(atomic.LoadInt32(&sleepSec)) * time.Second)
		}
		if err := ReqLoyalty(db, adresAccrual, orderNumb, chSignal); err != nil {
			log.Println(err)
		}
	}
}
