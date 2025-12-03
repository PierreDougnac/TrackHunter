package server

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/PierreDougnac/TrackHunter/proto"
)

// TrackHunterService implements the gRPC server
type TrackHunterServer struct {
	pb.UnimplementedTrackHunterServiceServer

	mu     sync.Mutex
	tracks map[string]*pb.Track
}

func NewTrackHunterServer() *TrackHunterServer {
	return &TrackHunterServer{
		tracks: make(map[string]*pb.Track),
	}
}

// IdentifySong is a placeholder that returns a dummy match
func (s *TrackHunterServer) IdentifySong(ctx context.Context, req *pb.IdentifySongRequest) (*pb.IdentifySongResponse, error) {

	// TODO: audio fingerprinting logic

	// For now: always return the first track found (demo mode)
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, track := range s.tracks {
		return &pb.IdentifySongResponse{
			Track:      track,
			Confidence: 0.42,
		}, nil
	}

	return nil, fmt.Errorf("no tracks in database")
}

func (s *TrackHunterServer) AddFingerprint(ctx context.Context, req *pb.AddFingerprintRequest) (*pb.AddFingerprintResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.Metadata.Id
	if id == "" {
		return nil, fmt.Errorf("track id is required")
	}

	s.tracks[id] = req.Metadata

	return &pb.AddFingerprintResponse{
		Sucess:  true,
		Message: "Fingerprint stored successfully",
	}, nil
}

func (s *TrackHunterServer) GetTrackInfo(ctx context.Context, req *pb.GetTrackInfoRequest) (*pb.GetTrackInfoResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	track, ok := s.tracks[req.Id]
	if !ok {
		return nil, fmt.Errorf("track %s not found", req.Id)
	}

	return &pb.GetTrackInfoResponse{
		Track: track,
	}, nil
}
