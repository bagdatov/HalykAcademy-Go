package main

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	pb "tictactoe/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	// TODO: здесь писать код gRPC клиента

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:8081", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	fmt.Println("gRPC connect created")

	client := pb.NewTictactoeAIGameClient(conn)
	ctx := context.Background()

	stream, err := client.StartGame(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	status, err := stream.Recv()
	if err != nil {
		fmt.Println("Recv err:", err)
	}

	PrintBoard(status.GetBoard())

	for {

		m := Move(ctx, client)

		if err := stream.Send(m); err != nil {
			fmt.Println("Send error:", err)
		}

		status, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Recv err:", err)
		}

		if win := status.GetBoard().Win; win != nil {
			if status.GetBoard().GetIsWin() {
				fmt.Println("\033[92m", "You win!", "\033[37m")
			} else {
				fmt.Println("\033[91m", "You lose!", "\033[37m")
			}

			status, err = stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Recv err:", err)
			}
		}

		PrintBoard(status.GetBoard())
	}
}

func PrintBoard(b *pb.Board) {
	if b == nil {
		return
	}
	fmt.Println("Current game status:\n")
	fmt.Println("   0     1     2\n")

	fmt.Printf("%v   %v |  %v  | %v\n", 0, getSymbol(b.Line1[0]), getSymbol(b.Line1[1]), getSymbol(b.Line1[2]))
	fmt.Println("   ——————————————")
	fmt.Printf("%v   %v |  %v  | %v\n", 1, getSymbol(b.Line2[0]), getSymbol(b.Line2[1]), getSymbol(b.Line2[2]))
	fmt.Println("   ——————————————")
	fmt.Printf("%v   %v |  %v  | %v\n\n", 2, getSymbol(b.Line3[0]), getSymbol(b.Line3[1]), getSymbol(b.Line3[2]))
}

func getSymbol(n int64) string {
	switch n {
	case 0:
		return "o"
	case 1:
		return "x"
	}
	return " "
}

func Move(ctx context.Context, client pb.TictactoeAIGameClient) *pb.Status {
	fmt.Println("Please make your move in the integer format: `x:y`")
	fmt.Println("You can see current score by typing 'score'")
	fmt.Println("You restart game by typing 'restart'")

	for {
		var str string
		fmt.Scan(&str)

		if str == "score" {
			res, err := client.GetScore(ctx, &pb.RequestScore{})
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("server: %v || you: %v\n", res.Super_AIScore, res.HandsomeScore)
			continue
		}

		if str == "restart" {
			return &pb.Status{
				Reply: &pb.Status_Err{
					Err: "restart",
				},
			}
		}

		s := strings.Split(str, ":")
		if len(s) != 2 {
			fmt.Println("Incorrect format, try again")
			continue
		}

		x, err1 := strconv.Atoi(s[0])
		y, err2 := strconv.Atoi(s[1])

		if err1 != nil || err2 != nil {
			fmt.Println("Incorrect format, try again")
			continue
		}

		return &pb.Status{
			Reply: &pb.Status_Move{
				Move: &pb.Move{
					X: int64(x),
					Y: int64(y),
				},
			},
		}
	}
	return nil
}
