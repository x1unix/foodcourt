package sender

import (
	//"time"
	"../../shared/logger"
	"../../shared/config"
	"fmt"
	"errors"
	"strconv"
	"runtime"
)


const logMsgSend = "(T%d): Sending mail to %s ..."
const logMsgSendOk = "(GO#%d): Send success to %s"
const logMsgSendFail = "(GO#%d): Send fail to %s: %s"
const logDbgTaskAlloc = "(TS): [%d:%d] => GOROUTINE#%d %v"

func SendLunchOrders() (bool, error) {
	log := logger.GetLogger()

	mails := []string{"ascii@live.ru", "foo@bar.com", "user@mail.com", "denis0051@gmail.com", "jdoe@mail.com", "u1@gmai.com", "fff", "aa"}
	done := make(chan SendResult)
	total := len(mails)

	// Total threads for goroutines
	threadCount, err := strconv.Atoi(config.Get(config.MAX_THREAD_COUNT, "-1"))

	if (err != nil) || (threadCount <= 0) {
		// Take all cores if MAX_THREAD_COUNT is invalid
		log.Warning("invalid MAX_THREAD_COUNT value, a default value used")
		threadCount = runtime.NumCPU()
	}

	// Don't spawn more goroutines than mails we have
	if threadCount >= total {
		if total > 1 {
			threadCount = total / 2
		} else {
			threadCount = 1
		}
	}

	// Allocate CPU cores
	runtime.GOMAXPROCS(threadCount)

	var taskPerGoroutine int

	// Allocate mails to send per goroutine
	taskPerGoroutine = total / threadCount

	restTasks := total % threadCount

	failed := 0


	lastGoroutine := threadCount - 1

	/**
		Debug info:

		TOD0 - Total mails
		GOR  - Avail goroutines
		TPG  - Mail per goroutine
		REST - Rest mails (%)
	 */
	log.Debug(fmt.Sprintf("(TS): TODO: %d; GOR: %d; TPG: %d; REST: %d", total, threadCount, taskPerGoroutine, restTasks))

	var offset int
	var limit int

	for t := 0; t < threadCount; t++ {
		switch t {
		case 0:
			offset = 0
			limit = taskPerGoroutine
			break
		case lastGoroutine:
			offset = limit
			limit = (taskPerGoroutine * 2) + restTasks
			break;
		default:
			offset = limit
			limit = taskPerGoroutine * 2
			break
		}

		arr := mails[offset:limit]
		log.Debug(fmt.Sprintf(logDbgTaskAlloc, offset, limit, t, arr))

		go sendLunchOrder(t, arr, done)



	}

	//for t := 0; t < threadCount; t++ {
	//
	//	var offset int
	//	var arr []string
	//	var end int
	//
	//	if t == 0 {
	//		offset = 0
	//	} else {
	//		offset = taskPerGoroutine * t
	//	}
	//
	//	if t == lastGoroutine {
	//		arr = mails[offset:]
	//		end = total - 1
	//	} else {
	//		end = taskPerGoroutine + 1
	//		arr = mails[offset:end]
	//	}
	//
	//	go sendLunchOrder(t, arr, done)
	//	log.Debug(fmt.Sprintf(logDbgTaskAlloc, offset, end, t, arr))
	//}

	for i := 0; i < total; i++ {
		res := <- done

		if res.Success {
			log.Info(fmt.Sprintf(logMsgSendOk, res.GoRoutine, res.Email))
		} else {
			failed++
			log.Error(fmt.Sprintf(logMsgSendFail, res.GoRoutine, res.Email, res.Error))
		}
	}

	if failed > 0 {
		return false, errors.New(fmt.Sprintf("%d mails were not send", failed))
	} else {
		return true, nil
	}
}

func sendLunchOrder(routineId int, emails []string, ch chan SendResult) {
	for _, email := range emails {
		err := "Foo"

		//time.Sleep(1 * time.Second)


		r := SendResult{routineId,email, true, err}
		ch <- r
	}
}