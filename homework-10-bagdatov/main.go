package main

import (
	"fmt"
	"net/http"
	"time"
	web "web/crypto-farm/http"
	"web/crypto-farm/http/middleware"
	"web/crypto-farm/repository/pg"
	"web/crypto-farm/usecase"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// не бейте за это
	// контейнер не успевает запустить иногда postgres
	// нормальное решение не успел написать, поэтому time sleep
	time.Sleep(5 * time.Second)
	db, err := pg.NewSQLRepository()
	if err != nil {
		log.Print(err)
		return
	}
	defer db.CloseConnection()

	uc := usecase.NewUseCase(db)
	m := middleware.InitMiddleware(uc)
	h := web.NewHanlder(uc)

	r.Handle("/app/user/{id:[0-9]+}", m.Middleware(h.UserPage))
	r.Handle("/app/wallet/{walletName}", m.Middleware(h.WalletPage))
	r.Handle("/app/wallet/{walletName}/{option}", m.Middleware(h.WalletMining))

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("server is listening at http://localhost:8080")
	log.Print(server.ListenAndServe())
}
