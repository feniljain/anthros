package service

import (
	"bytes"
	"context"
	"demo-grpc/pb"
	"errors"
	"io"
	"log"

	//"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

//LaptopServer is the server struct which provides laptop services
type LaptopServer struct {
	laptopStore LaptopStore
	imageStore  ImageStore
	ratingStore RatingStore
}

//NewLaptopServer returns a new laptop server
func NewLaptopServer(
	laptopStore LaptopStore,
	imageStore ImageStore,
	ratingStore RatingStore,
) *LaptopServer {
	return &LaptopServer{
		laptopStore: laptopStore,
		imageStore:  imageStore,
		ratingStore: ratingStore,
	}
}

//CreateLaptop is a unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {

	laptop := req.GetLaptop()
	log.Printf("Receive a create laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	//if ctx.Err().Error() == "context canceled" {
	//	log.Print("Request is cancelled")
	//	return nil, status.Error(codes.Canceled, "Request is cancelled")
	//}

	//if ctx.Err().Error() == "context deadline exceeded" {
	//	log.Printf("Deadline is exceeded")
	//	return nil, status.Error(codes.DeadlineExceeded, "Deadline is exceeded")
	//}

	err := server.laptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "Cannot save laptop to the store: %v", err)
	}

	log.Printf("Saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}

//SearchLaptop is a server-streaming RPC to search for laptop
func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer,
) error {

	filter := req.GetFilter()
	log.Printf("Received a search laptop request with filter: %v", filter)

	err := server.laptopStore.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{
			Laptop: laptop,
		}

		err := stream.Send(res)
		if err != nil {
			return err
		}

		log.Printf("Sent laptop with ID: %s", laptop.GetId())
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "Unexpected error: %v", err)
	}

	return nil
}

//UploadImage is a client-streaming RPC to upload a laptop image
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		logError(err)
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("Received an upload-image request for laptop %s with image type %s", laptopID, imageType)

	laptop, err := server.laptopStore.Find(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "Cannot find laptop: %v", err))
	}

	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "Laptop %s doesn't exist", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		if err := contextError(stream.Context()); err != nil {
			return err
		}
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("No more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Couldn't receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "Image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "Cannot write chunk data: %v", err))
		}
	}

	imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "Cannot write chunk data: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "Cannot save image in the store: %v", err))
	}

	log.Printf("Image with id: %s and size: %d", imageID, imageSize)
	return nil
}

//RateLaptop is a bi-directional RPC that allows client to rate a stream of laptops with a score, and returns a stream of average score for each of them
func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("No more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Cannot receive stream request: %v", err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScore()

		log.Printf("Received a rate-laptop request: id=%s, score=%.2f", laptopID, score)

		found, err := server.laptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Cannot find laptop: %v", err))
		}

		if found == nil {
			return logError(status.Errorf(codes.NotFound, "LaptopID %s is not found", laptopID))
		}

		rating, err := server.ratingStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "Cannot add rating to the store: %v", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Cannot send stream response: %v", err))
		}
	}

	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "Request is cancelled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "Deadline is exceeded"))
	default:
		return nil
	}
}
