package main

import (
	"database/sql"
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//	"fyne.io/fyne/v2/data/validation"
	"strconv"

	//	"fyne.io/fyne/v2/data/binding"
	"log"
	"os"

	"fyne.io/fyne/v2/widget"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Accounts Account
}

type Account struct {
	AccountID   int64
	UserID      int64
	amount      float64
	Transaction Transaction
}

type Transaction struct {
	transactionID int64
	AccountID     int64
	ToAccountID   int64
	amount        float64
}

func main() {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "swosh",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
	myApp := app.New()
	myWindow := myApp.NewWindow("Swosh App")

	combo := widget.NewSelect([]string{"Create Account", "Deposit", "Withdraw", "Transfer", "Print account information", "Print account list", "Delete account", "Transaction history"}, func(s string) {
		log.Println(s)
	})

	account_id := widget.NewEntry()
	id := widget.NewEntry()
	to_account_id := widget.NewEntry()
	amount := widget.NewEntry()
	//	amount.Validator = validation.NewRegexp(`[0-9]+`, "not a valid amount")
	information_text := widget.NewLabel("")
	form := widget.Form{
		Items: []*widget.FormItem{
			{Text: "Please make a selection", Widget: combo},
			{Text: "Account ID", Widget: account_id},
			{Text: "ID", Widget: id},
			{Text: "To Account ID", Widget: to_account_id},
			{Text: "Amount", Widget: amount},
			{Text: "Information", Widget: information_text},
		}, OnSubmit: func() {
			log.Println("form submitted")

			switch s := combo.Selected; s {
			case "Create Account":
				id := id.Text
				idint64, _ := strconv.ParseInt(id, 10, 64)
				acc_id, err := createAccount(idint64, 0.0)
				if err != nil {
					log.Println(err)
				} else {
					information_text.SetText("Account created with id: " + fmt.Sprint(acc_id))
					information_text.Refresh()
				}
			case "Deposit":
				account_id := account_id.Text
				account_idint64, _ := strconv.ParseInt(account_id, 10, 64)
				amount := amount.Text
				amountfloat64, _ := strconv.ParseFloat(amount, 64)
				err := addMoney(account_idint64, amountfloat64)
				if err != nil {
					log.Println(err)
				} else {
					information_text.SetText("Deposited: " + amount)
					information_text.Refresh()
				}
			case "Withdraw":
				account_id := account_id.Text
				account_idint64, _ := strconv.ParseInt(account_id, 10, 64)
				amount := amount.Text
				amountfloat64, _ := strconv.ParseFloat(amount, 64)
				information_text.SetText("Withdrawn: " + amount)
				err := addMoney(account_idint64, -amountfloat64)
				if err != nil {
					log.Println(err)
				} else {
					information_text.SetText("Withdrawn: " + amount)
					information_text.Refresh()
				}
			case "Transfer":
				account_id := account_id.Text
				account_idint64, _ := strconv.ParseInt(account_id, 10, 64)
				to_account_id := to_account_id.Text
				to_account_idint64, _ := strconv.ParseInt(to_account_id, 10, 64)
				amount := amount.Text
				amountfloat64, _ := strconv.ParseFloat(amount, 64)
				_, err = createTransaction(account_idint64, to_account_idint64, amountfloat64)
				if err != nil {
					log.Println(err)
				} else {
					information_text.SetText("Transferred: " + amount + " to " + to_account_id)
					information_text.Refresh()
				}
			case "Print account information":
				account_id := account_id.Text
				account_idint64, _ := strconv.ParseInt(account_id, 10, 64)
				acc, err := accountInfo(account_idint64)
				if err != nil {
					log.Println(err)
				} else {
					information_text.SetText(acc)
					information_text.Refresh()
				}
			case "Print account list":
				id := id.Text
				idint64, _ := strconv.ParseInt(id, 10, 64)
				accounts, err := allAccounts(idint64)
				if err != nil {
					log.Println(err)
				} else {
					var s string
					for _, account := range accounts {
						s += strconv.FormatInt(account.AccountID, 10) + "\n"
					}
					information_text.SetText(s)
				}
			case "Delete account":
				account_id := account_id.Text
				account_idint64, _ := strconv.ParseInt(account_id, 10, 64)
				_, err = deleteAccountById(account_idint64)
				if err != nil {
					log.Println(err)
				} else {
					information_text.SetText("Deleted account: " + account_id)
					information_text.Refresh()
				}
			case "Transaction history":
				account_id := account_id.Text
				account_idint64, _ := strconv.ParseInt(account_id, 10, 64)
				transactions, err := transactionHistory(account_idint64)
				if err != nil {
					log.Println(err)
				} else {
					var s string
					for _, transaction := range transactions {
						s += strconv.FormatInt(transaction.transactionID, 10) + "\n"
					}
					information_text.SetText(s)
					information_text.Refresh()
				}
			}
		},
	}
	myWindow.SetContent(container.NewVBox(&form))
	myWindow.ShowAndRun()

}

func createUser(name string, email string, password string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO user(name, email, password) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}

	res, err := stmt.Exec(name, email, password)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

func userByID(id int64) (User, error) {

	var user User

	row := db.QueryRow("SELECT id, name, email, password FROM user WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found: %d", id)
		}
		return user, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func allAccounts(user_id int64) ([]Account, error) {
	var accounts []Account
	rows, err := db.Query("SELECT * FROM account WHERE user_id = ?", user_id)
	if err != nil {
		return accounts, fmt.Errorf("failed to get accounts: %v", err)
	}

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.AccountID, &account.UserID, &account.amount); err != nil {
			return accounts, fmt.Errorf("failed to scan account: %v", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func createAccount(userID int64, amount float64) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO account(user_id, amount) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}

	res, err := stmt.Exec(userID, amount)
	if err != nil {
		return 0, fmt.Errorf("failed to create account: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}

func accountById(id int64) (Account, error) {

	var account Account

	row := db.QueryRow("SELECT id, user_id, amount FROM account WHERE id = ?", id)
	if err := row.Scan(&account.AccountID, &account.UserID, &account.amount); err != nil {
		if err == sql.ErrNoRows {
			return account, fmt.Errorf("account not found: %d", id)
		}
		return account, fmt.Errorf("failed to get account: %d", id)
	}

	return account, nil
}

func addMoney(accountID int64, amount float64) error {
	stmt, err := db.Prepare("UPDATE account SET amount = amount + ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}

	res, err := stmt.Exec(amount, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found: %d", accountID)
	}
	return nil
}

func deleteAccountById(accountid int64) (int64, error) {
	var account Account

	row := db.QueryRow("SELECT id, user_id, amount FROM account WHERE id = ?", accountid)

	if err := row.Scan(&account.AccountID, &account.UserID, &account.amount); err != nil {
		if err == sql.ErrNoRows {
			return accountid, fmt.Errorf("account not found: %d", accountid)
		}
		return accountid, fmt.Errorf("failed to get account: %d", accountid)
	}
	return account.AccountID, nil
}

func accountInfo(accountID int64) (string, error) {
	var account Account
	row := db.QueryRow("SELECT id, user_id, amount FROM account WHERE id = ?", accountID)
	if err := row.Scan(&account.AccountID, &account.UserID, &account.amount); err != nil {
		if err == sql.ErrNoRows {
			return "Account not found", fmt.Errorf("account not found: %d", accountID)
		}
		return "Failed to get account", fmt.Errorf("failed to get account: %d", accountID)
	}
	return fmt.Sprintf("Account info: id: %d, user_id: %d, amount: %f", account.AccountID, account.UserID, account.amount), nil
}

func createTransaction(accountID int64, toAccountID int64, amount float64) (string, error) {

	stmt, err := db.Prepare("INSERT INTO transaction(account_id, to_account_id, amount) VALUES(?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %v", err)
	}

	res, err := stmt.Exec(accountID, toAccountID, amount)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert id: %v", err)
	}

	if err := addMoney(accountID, -amount); err != nil {
		return "", fmt.Errorf("failed to send money: %v", err)
	}
	if err := addMoney(toAccountID, amount); err != nil {
		return "", fmt.Errorf("failed to receive money: %v", err)
	}

	return fmt.Sprintf("Transaction created success with id: %d", id), nil
}

func transactionHistory(accountID int64) ([]Transaction, error) {

	var transactions []Transaction

	rows, err := db.Query("SELECT id, account_id, to_account_id, amount FROM transaction WHERE account_id = ?", accountID)
	if err != nil {
		return transactions, fmt.Errorf("failed to get transactions: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.transactionID, &transaction.AccountID, &transaction.ToAccountID, &transaction.amount); err != nil {
			return transactions, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
