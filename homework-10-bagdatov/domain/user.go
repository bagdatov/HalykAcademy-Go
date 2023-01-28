package domain

// User struct describes credentials.
// All fields are ready for marshaling except Password.
type User struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"-"`
	Wallets  []string `json:"wallets"`
}
