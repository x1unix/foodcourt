package control

import (
	"foodcourt/logger"
	"foodcourt/menu"
	"github.com/jinzhu/now"
	"time"
	"fmt"
	"strconv"
)

const dateFmt = "20060102"

func LockMenu() (bool, error) {
	log := logger.GetLogger()
	today := now.New(time.Now()).BeginningOfDay()

	dateToBlock, _ := strconv.Atoi(GetDayToBlock(today).Format(dateFmt))

	log.Info(fmt.Sprintf("Blocking the menu for %d", dateToBlock))

	err := menu.SetMenuLockStatus(true, dateToBlock)

	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	return true, nil
}

// Gets next week day to block
func GetDayToBlock(date time.Time) time.Time {
	now.FirstDayMonday = true

	// Reset the time to 00:00
	today := now.New(date).BeginningOfDay()

	// Get today's week of day
	weekDay := today.Weekday()

	var daysToAdd int

	// Today is the weekend or friday - skip to next week
	if weekDay > time.Thursday {
		daysToAdd = 7 - int(weekDay) + 1
	} else {
		daysToAdd = 1
	}

	dayToBlock := today.AddDate(0, 0, daysToAdd)

	nexWeekDay := dayToBlock.Weekday()

	logMsg := fmt.Sprintf("%d %v : %v  =>  %d %v : %v", weekDay, weekDay, today, nexWeekDay, nexWeekDay, dayToBlock)
	logger.GetLogger().Debug(logMsg)

	return dayToBlock
}
