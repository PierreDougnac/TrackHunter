package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/PierreDougnac/TrackHunter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command> [args]")
		fmt.Println("Commands: add, list, get, delete")
		return
	}

	cmd := os.Args[1]

	// Connexion gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTrackHunterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			log.Fatal("Usage: add <title>")
		}
		title := os.Args[2]

		req := &pb.AddFingerprintRequest{
			Metadata: &pb.Track{
				Id:     fmt.Sprintf("%d", time.Now().UnixNano()), // simple ID
				Title:  title,
				Artist: "Unknown",
			},
			// Audio can be added later
		}

		resp, err := client.AddFingerprint(ctx, req)
		if err != nil {
			log.Fatalf("AddFingerprint failed: %v", err)
		}
		fmt.Println("Success:", resp.Success, "Message:", resp.Message)

	case "list":
		// On suppose que tu ajoutes ListFingerprintsRequest/Response dans ton proto
		resp, err := client.ListFingerprints(ctx, &pb.ListFingerprintsRequest{})
		if err != nil {
			log.Fatalf("ListFingerprints failed: %v", err)
		}
		fmt.Println("Fingerprints:")
		for _, t := range resp.Tracks {
			fmt.Printf("- ID: %s, Title: %s, Artist: %s\n", t.Id, t.Title, t.Artist)
		}

	case "get":
		if len(os.Args) < 3 {
			log.Fatal("Usage: get <id>")
		}
		id := os.Args[2]
		resp, err := client.GetFingerprint(ctx, &pb.GetFingerprintRequest{Id: id})
		if err != nil {
			log.Fatalf("GetFingerprint failed: %v", err)
		}
		t := resp.Track
		fmt.Printf("Track: ID: %s, Title: %s, Artist: %s\n", t.Id, t.Title, t.Artist)

	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("Usage: delete <id>")
		}
		id := os.Args[2]
		resp, err := client.DeleteFingerprint(ctx, &pb.DeleteFingerprintRequest{Id: id})
		if err != nil {
			log.Fatalf("DeleteFingerprint failed: %v", err)
		}
		if resp.Success {
			fmt.Println("Track deleted successfully")
		} else {
			fmt.Println("Track not found")
		}

	default:
		fmt.Println("Unknown command:", cmd)
	}
}
