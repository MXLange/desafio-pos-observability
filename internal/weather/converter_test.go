package weather

import "testing"

func TestConvertFromCelsius(t *testing.T) {
	tempF, tempK := ConvertFromCelsius(28.5)

	if tempF != 83.3 {
		t.Fatalf("fahrenheit = %v, want 83.3", tempF)
	}
	if tempK != 301.65 {
		t.Fatalf("kelvin = %v, want 301.65", tempK)
	}
}
