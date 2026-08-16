package weather

import "math"

func ConvertFromCelsius(celsius float64) (fahrenheit, kelvin float64) {
	return round(celsius*1.8 + 32), round(celsius + 273.15)
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}
