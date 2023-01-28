package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Request struct {
	Ctx      context.Context // Так делать плохо, но для домашки - можно :)
	Password string
}

type Response struct {
	Password string
	Pass     bool
}

type Connection struct {
	RequestConn  chan *Request  // "Входящий" коннекшен для нашего сервера
	ResponseConn chan *Response // "Выходящий" коннекшен для нашего сервера
}

type VulnerableServer struct {
	*Connection

	secretPassword string
}

func (vs *VulnerableServer) Run() {
	defer fmt.Println("Server stoped")
	timer := time.NewTimer(2 * time.Second)

LOOP:
	for {
		select {
		case req, ok := <-vs.RequestConn: // Обработка запросов на авторизацию
			if !ok {
				break LOOP
			}
			if req != nil {
				go func(r *Request) {
					timeCh := time.After(time.Second)
					select {
					case <-timeCh:
						select {
						case vs.ResponseConn <- &Response{
							Password: r.Password,
							Pass:     r.Password == vs.secretPassword,
						}:
						case <-r.Ctx.Done():
							fmt.Println("VulnerableServer.Run: Request canceled")
							return
						}

					case <-r.Ctx.Done():
						fmt.Println("VulnerableServer.Run: Request canceled")
						return
					}
				}(req)
			}
		case <-timer.C:
			fmt.Println("Server lifetime end. Goodbye:)")
			os.Exit(1)
		}

	}
}

func SendRequest(conn *Connection, req *Request) {
	conn.RequestConn <- req
	// Имитируем отправку данных по интернету. Предположим, что отправка идет через старый dial up модем с задержкой в секунду :)
	time.Sleep(time.Second)

}

func NewVulnerableServer(SecretPassword string, conn *Connection) *VulnerableServer {
	return &VulnerableServer{
		secretPassword: SecretPassword,
		Connection:     conn,
	}
}
