package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	pb "tictactoe/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type MyServer struct {
	Name        string
	AIscore     int64
	ClientScore int64
	pb.UnimplementedTictactoeAIGameServer
}

func (s *MyServer) GetScore(ctx context.Context, req *pb.RequestScore) (*pb.Score, error) {
	return &pb.Score{
		Super_AIScore: s.AIscore,
		HandsomeScore: s.ClientScore,
	}, nil
}

func (s *MyServer) StartGame(stream pb.TictactoeAIGame_StartGameServer) error {

	status := newBoard()

	if err := stream.Send(status); err != nil {
		fmt.Println("Response failed. Error: ", err)
		return err
	}

	ctx := stream.Context()

	var gameEnded bool

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled: ", ctx.Err())
			return ctx.Err()
		default:
			break
		}

		if gameEnded {
			gameEnded = false
			status = newBoard()

			if err := stream.Send(status); err != nil {
				fmt.Println("Response failed. Error: ", err)
				return err
			}
			continue
		}

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Recv error:", err)
			return err
		}

		if action := req.GetErr(); action == "restart" {
			gameEnded = true
			continue
		}

		if s.Win(status) {
			if err := stream.Send(req); err != nil {
				fmt.Println("Response failed. Error: ", err)
				return err
			}
			gameEnded = true
			continue
		}

		req = MakeMove(req, status)

		if s.Win(status) {
			if err := stream.Send(req); err != nil {
				fmt.Println("Response failed. Error: ", err)
				return err
			}
			gameEnded = true
			continue
		}

		if err := stream.Send(req); err != nil {
			fmt.Println("Response failed. Error: ", err)
			return err
		}

	}
	return nil
}

func newBoard() *pb.Status {
	return &pb.Status{Reply: &pb.Status_Board{Board: &pb.Board{
		Line1: []int64{-1, -1, -1},
		Line2: []int64{-1, -1, -1},
		Line3: []int64{-1, -1, -1},
	}}}
}

func (s *MyServer) Win(status *pb.Status) bool {
	b := status.GetBoard()
	if b == nil {
		return false
	}

	var winner string

	if ok, name := checkSum(b.Line1...); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line2...); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line3...); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line1[0], b.Line2[0], b.Line3[0]); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line1[1], b.Line2[1], b.Line3[1]); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line1[2], b.Line2[2], b.Line3[2]); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line1[0], b.Line2[1], b.Line3[2]); ok {
		winner = name
	}

	if ok, name := checkSum(b.Line1[2], b.Line2[1], b.Line3[0]); ok {
		winner = name
	}

	switch winner {
	case "server":
		s.AIscore++
		status.GetBoard().Win = &pb.Board_IsWin{IsWin: false}

	case "client":
		s.ClientScore++
		status.GetBoard().Win = &pb.Board_IsWin{IsWin: true}
	default:
		return false
	}

	return true
}

func checkSum(nums ...int64) (bool, string) {
	if len(nums) != 3 {
		return false, ""
	}
	if nums[0] != -1 && nums[0] == nums[1] && nums[0] == nums[2] {
		switch nums[0] {
		case 0:
			return true, "server"
		case 1:
			return true, "client"
		}
	}
	return false, ""
}

func MakeMove(req *pb.Status, status *pb.Status) *pb.Status {
	m := req.GetMove()
	if m == nil {
		return &pb.Status{Reply: &pb.Status_Err{
			Err: "non 'move' stream from client",
		}}
	}

	if x := m.GetX(); !(x >= 0 && x <= 2) {
		return &pb.Status{Reply: &pb.Status_Err{
			Err: "incorrect 'move' from client",
		}}
	}

	if y := m.GetY(); !(y >= 0 && y <= 2) {
		return &pb.Status{Reply: &pb.Status_Err{
			Err: "incorrect 'move' from client",
		}}
	}

	if err := enemyMove(status, m.GetX(), m.GetY()); err != nil {
		return &pb.Status{Reply: &pb.Status_Err{
			Err: "incorrect 'move' from client",
		}}
	}

	for {
		x, y := rand.Intn(3), rand.Intn(3)
		if isSave(status.GetBoard(), x, y) {
			return status
		}
	}

	return nil
}

func enemyMove(status *pb.Status, x, y int64) error {
	b := status.GetBoard()
	if b == nil {
		return nil
	}

	if y == 0 && b.Line1[x] == -1 {
		b.Line1[x] = 1
		return nil

	} else if y == 1 && b.Line2[x] == -1 {
		b.Line2[x] = 1
		return nil

	} else if y == 2 && b.Line3[x] == -1 {
		b.Line3[x] = 1
		return nil
	}
	return fmt.Errorf("incorrect move")
}

func isSave(b *pb.Board, x, y int) bool {
	if b == nil {
		return false
	}

	if y == 0 && b.Line1[x] == -1 {
		b.Line1[x] = 0
		return true

	} else if y == 1 && b.Line2[x] == -1 {
		b.Line2[x] = 0
		return true

	} else if y == 2 && b.Line3[x] == -1 {
		b.Line3[x] = 0
		return true
	}

	return false
}

func main() {
	port := ":8081"

	rand.Seed(time.Now().Unix())

	listener, err := net.Listen("tcp", port)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
		return
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	server := &MyServer{
		Name: "DeepBlue",
	}

	pb.RegisterTictactoeAIGameServer(grpcServer, server)

	fmt.Println("Start gRPC service...")
	fmt.Println(grpcServer.Serve(listener))
}
