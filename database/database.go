package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type User struct {
	Id      int64  `json:"id"`
	Name    string `json:"username"`
	Balance int32  `json:"balance"`
}

var (
	db    *sql.DB
	users []User
)

func init() {
	var err error

	// Подключаемя к базе данных
	db, err = sql.Open("postgres", "user=postgres password=234416 dbname=postgres sslmode=disable")
	PanicOnErr(err)

	//Пингуем базу
	err = db.Ping()
	PanicOnErr(err)
}

func runLoop(user *User) {
	task := struct {
		Id      int64  `json:"id"`
		UserId  int64  `json:"user_id"`
		Type    string `json:"type"`
		Amount  int32  `json:"amount"`
		Status  string `json:"status"`
		Created string `json:"created"`
	}{}
	for {
		tasks, err := db.Query("SELECT * FROM tasks WHERE user_id=$1 AND status=$2 ORDER BY created",
			user.Id,
			"open")
		defer tasks.Close()
		PanicOnErr(err)
		for tasks.Next() {
			err = tasks.Scan(&task.Id, &task.UserId, &task.Type, &task.Amount, &task.Status, &task.Created)
			PanicOnErr(err)
			if task.Type == "replenishment" && task.Amount > 0 {
				new_balance := user.Balance + task.Amount
				ctx := context.Background()
				tx, err := db.BeginTx(ctx, nil)
				PanicOnErr(err)
				_, err = tx.ExecContext(ctx, "UPDATE users SET balance=$1 WHERE id=$2", new_balance, user.Id)
				if err != nil {
					tx.Rollback()
					PanicOnErr(err)
				}
				_, err = tx.ExecContext(ctx, "UPDATE tasks SET status=$1 WHERE id=$2", "close", task.Id)
				if err != nil {
					tx.Rollback()
					PanicOnErr(err)
				}
				err = tx.Commit()
				PanicOnErr(err)
				user.Balance = new_balance
				fmt.Println("Task type", task.Type, "closed, user id:", user.Id, " new balance:", new_balance)
			} else if task.Type == "write-off" && task.Amount > 0 && task.Amount < user.Balance {
				new_balance := user.Balance - task.Amount
				ctx := context.Background()
				tx, err := db.BeginTx(ctx, nil)
				PanicOnErr(err)
				_, err = tx.ExecContext(ctx, "UPDATE  users SET balance=$1 WHERE id=$2", new_balance, user.Id)
				if err != nil {
					tx.Rollback()
					PanicOnErr(err)
				}
				_, err = tx.ExecContext(ctx, "UPDATE tasks SET status=$1 WHERE id=$2", "close", task.Id)
				if err != nil {
					tx.Rollback()
					PanicOnErr(err)
				}
				err = tx.Commit()
				PanicOnErr(err)
				user.Balance = new_balance
				fmt.Println("Task type", task.Type, "closed, user id:", user.Id, " new balance:", new_balance)
			} else {
				var lastId int32
				err = db.QueryRow("UPDATE tasks SET status=$1 WHERE id=$2 RETURNING id",
					"rejected",
					task.Id).Scan(&lastId)
				PanicOnErr(err)
				fmt.Println("Task rejected")
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func InitData() {
	var user User

	// Получаем данные о пользователях
	rows, err := db.Query("SELECT * FROM users")
	defer rows.Close()
	PanicOnErr(err)
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Balance)
		PanicOnErr(err)
		users = append(users, user)
	}
	for i := 0; i < len(users); i++ {
		go runLoop(&users[i])
	}
}

func GetUserById(idUser int64) []byte {
	for i := 0; i < len(users); i++ {
		if users[i].Id == idUser {
			productsJson, _ := json.Marshal(users[i])
			return productsJson
		}
	}
	return []byte("")
}

func GetAmountById(idUser int64) int32 {
	for i := 0; i < len(users); i++ {
		if users[i].Id == idUser {
			return users[i].Balance
		}
	}
	return -1
}

func GetAllUsers() []byte {
	productsJson, _ := json.Marshal(users)
	return productsJson
}

func CreateUser(name string, balance int32) int64 {
	if balance < 0 {
		return -1
	}
	var err error
	var lastId int64

	err = db.QueryRow("INSERT INTO users (username, balance) VALUES ($1, $2) RETURNING id",
		name,
		balance).Scan(&lastId)
	PanicOnErr(err)
	user := User{lastId, name, balance}
	users = append(users, user)
	go runLoop(&user)
	return lastId
}

func ChangeBalance(id int64, type_task string, amount int32) {
	db.QueryRow("INSERT INTO tasks (user_id, type_task, amount) VALUES ($1, $2, $3)", id, type_task, amount)
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
