package sample

import (
	"demo-grpc/pb"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(s ...string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	return s[rand.Intn(n)]
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon",
			"Core i3",
			"Core i5",
			"Core i7",
			"Core i9",
		)
	}

	return randomStringFromSet(
		"Ryzen 3",
		"Ryzen 5",
		"Ryzen 7",
	)
}

func randomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "Nvidia" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"GTX 1660",
			"GTX 1070",
		)
	}
	return randomStringFromSet(
		"RX 590",
		"RX 580",
		"RX 5700",
		"RX Vega",
	)
}

func randomScreenResolution() *pb.Screen_Resolution {

	height := randomInt(1080, 4320)
	width := height * 16 / 9

	return &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomID() string {
	return uuid.New().String()
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53")
	}

}
