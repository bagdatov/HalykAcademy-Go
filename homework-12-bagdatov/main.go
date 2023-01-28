package main

import (
	"fmt"
	"net/http"
	"time"
	middleware "web/auth/http"
	rd "web/auth/repository/redis"
	auc "web/auth/usecase"
	web "web/crypto-farm/http"
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

	// Connect to Redis and set 'alive' duration
	cache, err := rd.NewRedisClient(2*time.Minute, 5*time.Minute)
	if err != nil {
		log.Print(err)
		return
	}

	uc := usecase.NewCryptoUseCase(db)
	auth := auc.NewAuthUseCase(db, cache)
	m := middleware.InitMiddleware(auth)
	h := web.NewHanlder(uc)

	r.Handle("/login", m.LoginHandler()).Methods("POST")
	r.Handle("/update", m.UpdateHanlder()).Methods("POST")
	r.Handle("/app/user/{id:[0-9]+}", m.CheckAuthMiddleware(h.GetUserHandler)).Methods("GET")
	r.Handle("/app/user/{id:[0-9]+}", m.CheckAuthMiddleware(h.CreateUserHanlder)).Methods("POST")
	r.Handle("/app/wallet/{walletName}", m.CheckAuthMiddleware(h.GetWalletHanlder)).Methods("GET")
	r.Handle("/app/wallet/{walletName}", m.CheckAuthMiddleware(h.CreateWalletHanlder)).Methods("POST")
	r.Handle("/app/wallet/{walletName}/start", m.CheckAuthMiddleware(h.StartMiningHandler)).Methods("OPTIONS")
	r.Handle("/app/wallet/{walletName}/stop", m.CheckAuthMiddleware(h.StopMiningHandler)).Methods("OPTIONS")

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("server is listening at http://localhost:8080")
	log.Print(server.ListenAndServe())
}
