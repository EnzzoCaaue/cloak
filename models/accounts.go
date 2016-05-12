package models

import (
    "github.com/Cloakaac/cloak/database"
	"time"
)

// Account struct for accounts tables
type Account struct {
    ID int64
    Name string
    Password string
    Type int
    Premdays int
    Lastday int64
    Email string
    SecretKey string
}

// CloakaAccount struct for cloaka_account tables
type CloakaAccount struct {
	ID      int64 
	Account *Account 
	Token   string 
	Points  int    
	Admin   int
    TwoFactor int
    RecoveryKey string    
}

// NewAccount creates and return a new cloaka account
func NewAccount() *CloakaAccount {
    account := &Account{
        -1,
        "",
        "",
        0,
        0,
        0,
        "",
        "",
    }
    cloakaAccount := &CloakaAccount{
        -1,
        account,
        "",
        0,
        0,
        0,
        "",
    }
    return cloakaAccount
}

// Save registers an account
func (account *Account) Save() error {
    result, err := database.Connection.Exec("INSERT INTO accounts (name, password, type, premdays, lastday, email, creation) VALUES (?, ?, 0, ?, 0, ?, ?)", 
    account.Name, 
    account.Password, 
    account.Premdays, 
    account.Email, 
    time.Now().Unix())
	if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    account.ID = id
    return nil
}

// Save registers a cloaka account
func (account *CloakaAccount) Save() error {
    _, err := database.Connection.Exec("INSERT INTO cloaka_accounts (account, token, admin, twofactor, recovery) VALUES (?, ?, ?, ?, ?)", account.Account.ID, account.Token, account.Admin, 0, "")
	return err
}

// NameExists checks if an account name is already in use
func (account *CloakaAccount) NameExists() bool {
    row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE name = ?)", account.Account.Name)
    exists := false
    row.Scan(&exists)
    return exists
}

// EmailExists checks if an account email is already in use
func (account *CloakaAccount) EmailExists() bool {
    row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE email = ?)", account.Account.Email)
    exists := false
    row.Scan(&exists)
    return exists
}

// GetAccountByToken gets an account with the given cookie token
func GetAccountByToken(token string) *CloakaAccount {
    if !TokenExists(token) {
        return nil
    }
    account := NewAccount()
    row := database.Connection.QueryRow("SELECT a.id, a.name, a.password, a.email, a.premdays, b.points, b.admin, b.twofactor, b.recovery FROM accounts a, cloaka_accounts b WHERE a.id = b.account AND b.token = ?", token)
    row.Scan(&account.Account.ID, &account.Account.Name, &account.Account.Password, &account.Account.Email, &account.Account.Premdays, &account.Points, &account.Admin, &account.TwoFactor, &account.RecoveryKey)
    return account
}

// TokenExists checks if a token is already in use by an account
func TokenExists(token string) bool {
    row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM cloaka_accounts WHERE token = ?)", token)
    exists := false
    row.Scan(&exists)
    return exists
}

// UpdateToken sets an account cookie token
func (account *CloakaAccount) UpdateToken(token string) error {
    _, err := database.Connection.Exec("UPDATE cloaka_accounts SET token = ? WHERE account = ?", token, account.Account.ID)
    return err
}

// SignIn checks if a given username and password exists
func (account *CloakaAccount) SignIn() bool {
    row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE name = ? AND password = ?)", account.Account.Name, account.Account.Password)
    success := false
    row.Scan(&success)
    if !success {
        return false
    }
    row = database.Connection.QueryRow("SELECT a.id, a.secret, b.twofactor FROM accounts a, cloaka_accounts b WHERE a.name = ? AND a.id = b.account", account.Account.Name)
    row.Scan(&account.Account.ID, &account.Account.SecretKey, &account.TwoFactor)
    return true
}

// UpdateRecoveryKey sets an account recovery key
func (account *CloakaAccount) UpdateRecoveryKey(key string) error {
    _, err := database.Connection.Exec("UPDATE cloaka_accounts SET recovery = ? WHERE account = ?", key, account.Account.ID)
    return err
}

// EnableTwoFactor enables the two-factor google auth system on a given account
func (account *CloakaAccount) EnableTwoFactor(secret string) error {
    _, err := database.Connection.Exec("UPDATE accounts a, cloaka_accounts b SET b.twofactor = 1, a.secret = ? WHERE a.id = ? AND b.account = a.id", secret, account.Account.ID)
    return err
}