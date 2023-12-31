package main

import (
	"context"
	"fmt"
	"gRPC-tutorial/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	callListFiles(client)
	// callDownload(client)
	// CallUpload(client)
	// CallUploadAndNotifyProgress(client)
}

func callListFiles(client pb.FileServiceClient) {
	md := metadata.New(map[string]string{"authorization": "Bearer bad-token"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	res, err := client.ListFiles(ctx, &pb.ListFileRequest{})
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println(res.GetFilename())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Response from Download(butes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}

func CallUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/isawashun/Desktop/gRPC-tutorial/storage" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}

// func CallUploadAndNotifyProgress(client pb.FileServiceClient) {
// 	filename := "sports.txt"
// 	path := filepath.Join("/Users/isawashun/Desktop/gRPC-tutorial/storage", filename)

// 	file, err := os.Open(path)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer file.Close()

// 	stream, err := client.UploadAndNotifyProgress(context.Background())
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	buf := make([]byte, 5)
// 	go func() {
// 		for {
// 			n, err := file.Read(buf)
// 			if n == 0 || err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				log.Fatalln(err)
// 			}

// 			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
// 			sendErr := stream.Send(req)
// 			if sendErr != nil {
// 				log.Fatalln(sendErr)
// 			}
// 			time.Sleep(1 * time.Second)
// 		}

// 		err := stream.CloseSend()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	}()

// 	// Response
// 	ch := make(chan struct{})
// 	go func() {
// 		for {
// 			res, err := stream.Recv()
// 			if err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				log.Fatalln(err)
// 			}

// 			log.Printf("received message: %v", res.GetMsg())
// 		}
// 		close(ch)
// 	}()
// 	<-ch
// }