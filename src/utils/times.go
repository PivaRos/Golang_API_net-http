package utils

import "time"

type Times struct {
	StartTime time.Time `bson:"startTime" json:"startTime"`
	EndTime   time.Time `bson:"endTime" json:"endTime"`
}

// extractTime extracts the hour and minute from a time.Time object
func extractTime(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, t.Location())
}

// combineDateTime combines a date from one time.Time and a time from another time.Time
func combineDateTime(date, timeComponent time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), timeComponent.Hour(), timeComponent.Minute(), 0, 0, date.Location())
}

// SubtractTimes subtracts the single Times from the array of Times, considering only hour and minute
func SubtractTimes(times []Times, sub Times) []Times {
	var result []Times

	subStartTime := extractTime(sub.StartTime)
	subEndTime := extractTime(sub.EndTime)

	for _, t := range times {
		tStartTime := extractTime(t.StartTime)
		tEndTime := extractTime(t.EndTime)

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
