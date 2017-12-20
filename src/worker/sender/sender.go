package sender

import (
	"time"
	"math/rand"
	"../../shared/logger"
	"fmt"
	"errors"
)


const logMsgSend = "Sending mail to %s ..."
const logMsgSendOk = "Send success to %s"
const logMsgSendFail = "Send fail to %s: %s"

func SendLunchOrders() (bool, error) {
	log := logger.GetLogger()

	mails := []string{"ascii@live.ru", "foo@bar.com", "user@mail.com"}
	done := make(chan SendResult)
	sz := len(mails)
	failed := 0

	for _, mail := range mails {
		log.Info(fmt.Sprintf(logMsgSend, mail))
		go sendLunchOrder(mail, done)
	}

	for i := 0; i < sz; i++ {
		res := <- done

		if res.Success {
			log.Info(fmt.Sprintf(logMsgSendOk, res.Email))
		} else {
			failed++
			log.Error(fmt.Sprintf(logMsgSendFail, res.Email, res.Error))
		}
	}

	if failed > 0 {
		return false, errors.New(fmt.Sprintf("%d mails were not send", failed))
	} else {
		return true, nil
	}
}

func sendLunchOrder(email string, ch chan SendResult) {

	seed := rand.Intn(7)
	success := !(seed > 5)
	err := "Foo"

	time.Sleep(time.Duration(seed) * time.Second)


	r := SendResult{email, success, err}
	ch <- r
}