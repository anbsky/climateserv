package climateserv

import "math"

var descriptions = map[string]string{}

/*
CurrentPM25toAQI converts current PM2.5 concentration to AQI.

Code based on AirNow calculator (https://airnow.gov/index.cfm?action=airnow.calculator)
*/
func CurrentPM25toAQI(concentration float64) (AQI int, category string) {
	fConc := math.Floor(concentration*10) / 10
	if 12.1 > fConc && fConc >= 0 {
		AQI = convertPM25toAQI(50, 0, 12, 0, fConc)
		category = "good"
	} else if 35.5 > fConc && fConc >= 12.1 {
		AQI = convertPM25toAQI(100, 51, 35.4, 12.1, fConc)
		category = "mdoerate"
	} else if 55.5 > fConc && fConc >= 35.5 {
		AQI = convertPM25toAQI(150, 101, 55.4, 35.5, fConc)
		category = "unhealthy for sensitive groups"
	} else if 150.5 > fConc && fConc >= 55.5 {
		AQI = convertPM25toAQI(200, 151, 150.4, 55.5, fConc)
		category = "unhealthy"
	} else if 250.5 > fConc && fConc >= 150.5 {
		AQI = convertPM25toAQI(300, 201, 250.4, 150.5, fConc)
		category = "very unhealthy"
	} else if 350.5 > fConc && fConc >= 250.5 {
		AQI = convertPM25toAQI(400, 301, 350.4, 250.5, fConc)
		category = "hazardous"
	} else if 500.5 > fConc && fConc >= 350.5 {
		AQI = convertPM25toAQI(500, 401, 500.4, 350.5, fConc)
		category = "hazardous"
	} else {
		AQI = 999
		category = "out of range"
	}
	return AQI, category
}

/*
convertPM25toAQI calculates AQI based on current concentration
and previous high and low values using EPA formula
(https://en.wikipedia.org/wiki/Air_quality_index#Computing_the_AQI)
*/
func convertPM25toAQI(AQIhigh uint, AQIlow uint, concHigh float64, concLow float64, conc float64) (AQI int) {
	resultAQI := ((conc-concLow)/(concHigh-concLow))*float64((AQIhigh-AQIlow)) + float64(AQIlow)
	return int(math.Round(resultAQI))
}
