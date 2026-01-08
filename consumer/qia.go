package consumer

import "math"

func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func normalize(value, max float64) float64 {
	return clamp((value/max)*500.0, 0, 500)
}

func humidityPenalty(h float64) float64 {
	if h >= 30 && h <= 70 {
		return 0
	}
	if h < 30 {
		return normalize(30-h, 30)
	}
	return normalize(h-70, 30)
}

func CalculateQIA(data *SensorData) float64 {
	aqiPM25 := normalize(float64(data.PM25), 250)
	aqiPM10 := normalize(float64(data.PM10), 420)
	aqiCO2 := normalize(float64(data.CO2), 2000)
	aqiTemp := normalize(float64(data.Temperature), 50)
	aqiHum := humidityPenalty(float64(data.Humidity))

	qia :=
		0.40*aqiPM25 +
			0.25*aqiPM10 +
			0.20*aqiCO2 +
			0.10*aqiTemp +
			0.05*aqiHum

	return math.Round(qia*100) / 100
}
