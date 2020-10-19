package sample

import (
	"demo-grpc/pb"

	"github.com/golang/protobuf/ptypes"
)

//NewKeyboard returns a new sample keyboard
func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}

//NewCPU returns a new sample CPU
func NewCPU() *pb.CPU {
	brand := randomCPUBrand()

	numberCores := randomInt(2, 8)

	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	return &pb.CPU{
		Brand:       brand,
		Name:        randomCPUName(brand),
		NumberCores: uint32(numberCores),
		MinGhz:      minGhz,
		MaxGhz:      maxGhz,
	}
}

//NewGPU returns a new sample GPU
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)

	return &pb.GPU{
		Brand:  brand,
		Name:   randomGPUName(brand),
		MinGhz: minGhz,
		MaxGhz: maxGhz,
	}
}

//NewRAM returns a new sample RAM
func NewRAM() *pb.Memory {
	memGB := randomInt(4, 64)

	return &pb.Memory{
		Value: uint64(memGB),
		Unit:  pb.Memory_GIGABYTE,
	}
}

//NewSSD returns a new sample SSD
func NewSSD() *pb.Storage {
	memGB := randomInt(128, 1024)

	return &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(memGB),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

//NewHDD returns a new sample HDD
func NewHDD() *pb.Storage {
	memTB := randomInt(1, 6)

	return &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(memTB),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

//NewScreen returns a sample screen
func NewScreen() *pb.Screen {
	return &pb.Screen{
		SizeInch:   randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
}

//NewLaptop returns a sample laptop
func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	return &pb.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     randomLaptopName(brand),
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*pb.GPU{NewGPU(), NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3500),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:   ptypes.TimestampNow(),
	}
}

//RandomLaptopScore returns a random laptop score
func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
