package main

import (
	"context"
	"fmt"
	"gRPC-tutorial/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFileRequest) (*pb.ListFileResponse, error) {
	fmt.Println("ListenFiles was invoked")

	dir := "/Users/isawashun/Desktop/gRPC-tutorial/storage"

	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFileResponse{
		Filename: filenames,
	}
	return res, nil
}

func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := "/Users/isawashun/Desktop/gRPC-tutorial/storage"

	file, err := os.Open(filepath.Join(path, filename))
	if err != nil {
			return err
	}
	defer file.Close()

	buf := make([]byte, 1024) // または他の適切なサイズ
	for {
			n, err := file.Read(buf)
			if err == io.EOF {
					break
			}
			if err != nil {
					return err
			}

			res := &pb.DownloadResponse{Data: buf[:n]}
			if err := stream.Send(res); err != nil {
					return err
			}
			time.Sleep(1 * time.Second) // 必要に応じて調整または削除
	}

	return nil
}


func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is running ... ")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}