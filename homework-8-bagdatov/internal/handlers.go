package internal

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ub *UserBase) UserPage(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:

		username, err := ub.FindUsernameByID(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		wallets := ub.FindWalletsByID(id)

		data := struct {
			Username string   `json:"username"`
			Wallets  []string `json:"wallets"`
		}{Username: username, Wallets: wallets}

		reply, err := json.Marshal(data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

		err := ub.CreateUser(id, username, password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (ub *UserBase) WalletPage(w http.ResponseWriter, r *http.Request) {

	walletName := mux.Vars(r)["walletName"]
	user := r.Context().Value(ctxKeyUser).(*User)

	switch r.Method {
	case http.MethodGet:

		sum, err := ub.WalletSum(walletName, user.ID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		data := struct {
			WalletName string `json:"walletName"`
			WalletsSum int64  `json:"walletSum"`
		}{WalletName: walletName, WalletsSum: sum}

		reply, err := json.Marshal(data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(reply)

	case http.MethodPost:

		err := ub.CreateWallet(walletName, user.ID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Wallet created"))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (ub *UserBase) WalletMining(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	option := mux.Vars(r)["option"]
	walletName := mux.Vars(r)["walletName"]
	user := r.Context().Value(ctxKeyUser).(*User)

	switch option {
	case "start":

		err := ub.StartMine(walletName, user.ID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Write([]byte("Mining started"))

	case "stop":

		err := ub.StopMine(walletName, user.ID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Write([]byte("Mining stopped"))

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
