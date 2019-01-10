package climateserv

import "math"

var descriptions = map[string]string{}

// PM25toAQI convert PM2.5 concentration to AQI value
func PM25toAQI(concentration float64) (AQI int, category string) {
	fConc := math.Floor(concentration*10) / 10
	if 12.1 > fConc && fConc >= 0 {
		AQI = linear(50, 0, 12, 0, fConc)
		category = "good"
	} else if 35.5 > fConc && fConc >= 12.1 {
		AQI = linear(100, 51, 35.4, 12.1, fConc)
		category = "mdoerate"
	} else if 55.5 > fConc && fConc >= 35.5 {
		AQI = linear(150, 101, 55.4, 35.5, fConc)
		category = "unhealthy for sensitive groups"
	} else if 150.5 > fConc && fConc >= 55.5 {
		AQI = linear(200, 151, 150.4, 55.5, fConc)
		category = "unhealthy"
	} else if 250.5 > fConc && fConc >= 150.5 {
		AQI = linear(300, 201, 250.4, 150.5, fConc)
		category = "very unhealthy"
	} else if 350.5 > fConc && fConc >= 250.5 {
		AQI = linear(400, 301, 350.4, 250.5, fConc)
		category = "hazardous"
	} else if 500.5 > fConc && fConc >= 350.5 {
		AQI = linear(500, 401, 500.4, 350.5, fConc)
		category = "hazardous"
	} else {
		AQI = 999
		category = "out of range"
	}
	return AQI, category
}

func linear(AQIhigh uint, AQIlow uint, concHigh float64, concLow float64, conc float64) (AQI int) {
	resultAQI := ((conc-concLow)/(concHigh-concLow))*float64((AQIhigh-AQIlow)) + float64(AQIlow)
	return int(math.Round(resultAQI))
}
