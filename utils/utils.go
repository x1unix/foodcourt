package utils

import (
	"time"
	"github.com/jinzhu/now"
	"fmt"
	"foodcourt/logger"
)

// Gets next day to process the order
func GetNextDay(date time.Time) time.Time {
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