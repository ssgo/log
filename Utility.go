package log

import "time"

func MakeLogTime(time time.Time) float64 {
	return float64(time.UnixNano()) / 1e9
}

func MakeUesdTime(startTime, endTime time.Time) float32 {
	return float32(endTime.UnixNano()-startTime.UnixNano()) / 1e6
}

//func Parse(data string, to interface{}) error {
//
//}