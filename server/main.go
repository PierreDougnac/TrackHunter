package main

import (
	"log"
	"net"

	pb "github.com/PierreDougnac/TrackHunter/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterTrackHunterServiceServer(server, NewTrackHunterServer())

	log.Println("TrackHunter gRPC server running on :50052")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
