package consumer

import (
	"strings"
	"time"
)

// Accept RFC3339Nano with timezone
type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}

type SensorData struct {
	SensorID      int        `json:"sensorID"`
	Name          string     `json:"name"`
	Latitude      float64    `json:"latitude"`
	Longitude     float64    `json:"longitude"`
	Temperature   float32    `json:"temperature"`
	Humidity      float32    `json:"humidity"`
	CO2           float32    `json:"co2"`
	PM25          float32    `json:"pm25"`
	PM10          float32    `json:"pm10"`
	WindDirection int        `json:"windDirection"`
	Status        string     `json:"status"`
	Timestamp     CustomTime `json:"timestamp"`
	QIA           float64    `json:"qia"`
}
