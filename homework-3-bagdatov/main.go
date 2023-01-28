package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

type Target struct {
	InputFilePath  string // Путь до файла с паролями
	OutputFilePath string // Путь до файла, куда должны записываться результаты
	*Connection
}

func HackServer(ctx context.Context, target *Target) {
	// Открываем файл с паролями
	fileInput, err := os.Open(target.InputFilePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer fileInput.Close()

	scanner := bufio.NewScanner(fileInput)

	// Считываем и отправляем запросы на лету
	for scanner.Scan() {
		req := &Request{
			Ctx:      ctx,
			Password: scanner.Text(),
		}

		go SendRequest(target.Connection, req)
	}
}

func main() {
	requestChan := make(chan *Request)
	responseChan := make(chan *Response)
	defer close(requestChan)
	defer close(responseChan)

	connection := &Connection{
		RequestConn:  requestChan,
		ResponseConn: responseChan,
	}

	target := &Target{
		InputFilePath:  "darkweb2017-top10000.txt",
		OutputFilePath: "output.txt",
		Connection:     connection,
	}
	// Заменить "Password" на один из 10000 паролей
	server := NewVulnerableServer("aezakmi", connection)

	go server.Run()

	// Пробовать запускать с разными контекстами
	ctx, cancel := context.WithCancel(context.Background())
	HackServer(ctx, target)

	// Открываем или создаем файл для сохранения результата
	fileOutput, err := os.OpenFile(target.OutputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Println(err)
		return
	}

	// Очищаем файл
	if err := os.Truncate(target.OutputFilePath, 0); err != nil {
		log.Printf("Failed to truncate: %v\n", err)
	}

	// Считываем результаты и записываем
	for {
		response := <-target.ResponseConn
		fileOutput.WriteString(fmt.Sprintf("%s:%v\n", response.Password, response.Pass))

		// Как только необходимый пароль найден останавливаем подбор
		if response.Pass {
			log.Println("\033[32m", "SUCCESS", "\033[0m")
			cancel()
			break
		}
	}
	fileOutput.Close()
}
