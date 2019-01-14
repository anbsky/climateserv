package climateserv

import (
	"testing"
)

func TestCurrentPM25toAQI(t *testing.T) {
	values := [][3]interface{}{
		{111.5, 180, "unhealthy"},
		{10.0, 42, "good"},
		{12.5, 52, "moderate"},
		{53.0, 144, "unhealthy for sensitive groups"},
		{99.5, 174, "unhealthy"},
		{9.8, 41, "good"},
		{201.0, 251, "very unhealthy"},
		{400.0, 434, "hazardous"},
		{501.99, 999, "out of range"},
	}
	for _, expected := range values {
		concentration, _ := expected[0].(float64)
		resultAQI, resultDescription := CurrentPM25toAQI(concentration)
		if resultAQI != expected[1] {
			t.Errorf(
				"AQI for %v should be %v but we got %v",
				expected[0], expected[1], resultAQI)
		}
		if resultDescription != expected[2] {
			t.Errorf(
				"description for %v should be %v but we got %v",
				expected[0], expected[2], resultDescription)
		}
	}
}
