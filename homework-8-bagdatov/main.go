package main

import (
	"log"
	"net/http"
	"time"
	"web/internal"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	ub := internal.NewUserBase()

	r.Handle("/app/user/{id:[0-9]+}", ub.Middleware(ub.UserPage))
	r.Handle("/app/wallet/{walletName}", ub.Middleware(ub.WalletPage))
	r.Handle("/app/wallet/{walletName}/{option}", ub.Middleware(ub.WalletMining))

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("server is listening at http://localhost:8080")
	log.Println(server.ListenAndServe())
}
