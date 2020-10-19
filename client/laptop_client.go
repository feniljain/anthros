package client

import (
	"bufio"
	"context"
	"demo-grpc/pb"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//LaptopClient is a client to call laptop service RPCs
type LaptopClient struct {
	service pb.LaptopServiceClient
}

//NewLaptopClient returns a new laptop client
func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{
		service: service,
	}
}

//CreateLaptop calls create laptop RPC
func (laptopClient *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.service.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("Laptop already exists")
		} else {
			log.Fatal("Cannot create laptop: ", err)
		}
		return
	}

	log.Printf("Created laptop with id: %s", res.Id)
}

//SearchLaptop calls search laptop RPC
func (laptopClient *LaptopClient) SearchLaptop(filter *pb.Filter) {
	log.Print("Seach filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}

	stream, err := laptopClient.service.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("Cannot search laptop: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("Cannot receive response: ", err)
		}

		laptop := res.GetLaptop()
		log.Printf("Found matches: %+v\n", laptop)
	}
}

//UploadImage calls upload image RPC
func (laptopClient *LaptopClient) UploadImage(laptopID, imagePath string) {

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Cannot open image file: %v", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("Cannot Upload Image:", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("Cannot send image: ", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Cannot read chunk to buffers: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		stream.Send(req)
		if err != nil {
			log.Fatal("Cannot send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Couldn't receive response: ", err)
	}

	log.Printf("Image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

//RateLaptop calls rate laptop RPC
func (laptopClient *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot rate laptop: %v", err)
	}

	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("No more responses")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- err
				return
			}

			log.Print("Received response: ", res)
		}
	}()

	for i, laptopID := range laptopIDs {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopID,
			Score:    scores[i],
		}

		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("Cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}

		log.Print("Request sent: ", req)
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("Cannot close send: %v", err)
	}

	err = <-waitResponse
	return err
}
