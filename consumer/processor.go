package consumer

import (
	"encoding/json"
	"fmt"
	"log"
)

func analyzeSensorStatus(data *SensorData) string {
	if data.Humidity == 0 {
		return "faulty"
	}
	if data.Temperature > 50 || data.PM25 > 300 || data.PM10 > 400 {
		return "fire"
	}
	if data.Temperature > 35 || data.PM25 > 100 || data.PM10 > 150 {
		return "affected"
	}
	return "normal"
}

func DefaultProcessMessage(msg []byte) {
	var data SensorData
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Printf("[Processor] Error decoding JSON: %v", err)
		return
	}

	// Status
	data.Status = analyzeSensorStatus(&data)

	// QIA
	data.QIA = CalculateQIA(&data)

	fmt.Printf(
		"[Processed] Sensor %d (%s) | Temp: %.1f°C | Hum: %.1f%% | CO2: %.1f | PM2.5: %.1f | PM10: %.1f | QIA: %.2f | Status: %s\n",
		data.SensorID, data.Name, data.Temperature, data.Humidity,
		data.CO2, data.PM25, data.PM10, data.QIA, data.Status,
	)

	SendToAPI(data)
}
