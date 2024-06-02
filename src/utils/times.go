package utils

import "time"

type Times struct {
	StartTime time.Time `bson:"startTime" json:"startTime" validate:"required"`
	EndTime   time.Time `bson:"endTime" json:"endTime" validate:"required"`
}

// ExtractTime extracts the hour and minute from a time.Time object
func ExtractTime(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, t.Location())
}

// combineDateTime combines a date from one time.Time and a time from another time.Time
func combineDateTime(date, timeComponent time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), timeComponent.Hour(), timeComponent.Minute(), 0, 0, date.Location())
}

func AdjustDates(timesMap *map[string][]Times, date time.Time) map[string][]Times {
	result := make(map[string][]Times)

	for key, timesSlice := range *timesMap {
		var adjustedTimes []Times
		for _, t := range timesSlice {
			adjustedStartTime := time.Date(date.Year(), date.Month(), date.Day(), t.StartTime.Hour(), t.StartTime.Minute(), t.StartTime.Second(), 0, date.Location())
			adjustedEndTime := time.Date(date.Year(), date.Month(), date.Day(), t.EndTime.Hour(), t.EndTime.Minute(), t.EndTime.Second(), 0, date.Location())
			adjustedTimes = append(adjustedTimes, Times{
				StartTime: adjustedStartTime,
				EndTime:   adjustedEndTime,
			})
		}
		result[key] = adjustedTimes
	}

	return result
}

// SubtractTimes subtracts the single Times from the array of Times, considering only hour and minute
func SubtractTimes(times []Times, sub Times) []Times {
	var result []Times

	subStartTime := ExtractTime(sub.StartTime)
	subEndTime := ExtractTime(sub.EndTime)

	for _, t := range times {
		tStartTime := ExtractTime(t.StartTime)
		tEndTime := ExtractTime(t.EndTime)

		if subEndTime.Before(tStartTime) || subStartTime.After(tEndTime) {
			// No overlap
			result = append(result, t)
		} else {
			// Overlap exists
			if subStartTime.After(tStartTime) {
				result = append(result, Times{StartTime: t.StartTime, EndTime: combineDateTime(t.StartTime, subStartTime)})
			}
			if subEndTime.Before(tEndTime) {
				result = append(result, Times{StartTime: combineDateTime(t.StartTime, subEndTime), EndTime: t.EndTime})
			}
		}
	}
	return result
}
