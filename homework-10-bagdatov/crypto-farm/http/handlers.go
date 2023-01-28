package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web/crypto-farm/http/middleware"
	"web/domain"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type handler struct {
	domain.UseCase
}

func NewHanlder(uc domain.UseCase) *handler {
	return &handler{uc}
}

func (handle *handler) UserPage(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:

		user, err := handle.GetUserByID(r.Context(), id)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		reply, err := json.Marshal(user)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		w.Write(reply)

	case http.MethodPost:

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user := &domain.User{
			ID:       id,
			Username: username,
			Password: password,
		}

		err := handle.CreateUser(r.Context(), user)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (handle *handler) WalletPage(w http.ResponseWriter, r *http.Request) {

	walletName := mux.Vars(r)["walletName"]
	user := r.Context().Value(middleware.CtxKeyUser).(*domain.User)

	switch r.Method {
	case http.MethodGet:

		wallet, err := handle.FindWallet(r.Context(), user, walletName)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		reply, err := json.Marshal(wallet)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		w.Write(reply)

	case http.MethodPost:

		wallet := &domain.CryptoWallet{
			Name:      walletName,
			OwnerName: user.Username,
			OwnerID:   user.ID,
			Amount:    0,
		}

		err := handle.CreateWallet(r.Context(), wallet)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Wallet created"))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (handle *handler) WalletMining(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	option := mux.Vars(r)["option"]
	walletName := mux.Vars(r)["walletName"]
	user := r.Context().Value(middleware.CtxKeyUser).(*domain.User)

	switch option {
	case "start":

		err := handle.StartMine(r.Context(), user, walletName)

		if err != nil {
			ErrorHandler(w, err)
			return
		}

		w.Write([]byte("Mining started"))

	case "stop":

		err := handle.StopMine(r.Context(), user, walletName)
		if err != nil {
			ErrorHandler(w, err)
			return
		}

		w.Write([]byte("Mining stopped"))

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func ErrorHandler(w http.ResponseWriter, err error) {

	switch err {
	case domain.ErrInitilized, domain.ErrNotInitilized:
		http.Error(w, err.Error(), http.StatusBadRequest)

	case domain.ErrNotFound:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	case domain.ErrExists:
		http.Error(w, err.Error(), http.StatusBadRequest)

	case domain.ErrNotAuthorized:
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

	default:
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
