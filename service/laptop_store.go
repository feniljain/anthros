package service

import (
	"context"
	"demo-grpc/pb"
	"errors"
	"fmt"
	"log"
	"sync"
	//"time"

	"github.com/jinzhu/copier"
)

//ErrAlreadyExists is returned whena record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("Record already exists")

//LaptopStore is an interface to store laptop
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

//InMemoryLaptopStore stores laptop in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

//NewInMemoryLaptopStore returns a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

//Save saves the laptop to store
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	return nil
}

//Find finds a laptop by ID
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	if store.data[id] == nil {
		return nil, nil
	}

	return deepCopy(store.data[id])
}

//Search searches for laptop with filter, returns one by one via the found function
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		//time.Sleep(time.Second)
		log.Print("Checking laptop id: ", laptop.Id)

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("Context is cancelled")
			return errors.New("Context is cancelled")
		}

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.MinCpuCores {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.MinCpuGhz {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3
	case pb.Memory_KILOBYTE:
		return value << 13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("Cannot copy laptop data: %v", err)
	}

	return other, nil
}
