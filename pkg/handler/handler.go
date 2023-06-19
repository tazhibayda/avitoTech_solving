package handler

import (
	"avitoTech_solving/pkg/database"
	"avitoTech_solving/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func balance(c *gin.Context) {
	var id struct {
		Id int `json:"user_id"`
	}
	jsonData, err := c.GetRawData()
	fmt.Println("parse")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := json.Unmarshal(jsonData, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User

	q := fmt.Sprintf("SELECT balance from users where user_id = %d", id.Id)
	err = database.DB.QueryRow(q).Scan(&user.Balance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("balance not found for this user")})
		return
	}

	c.JSON(http.StatusOK, user.Balance)

	//Тело запроса:            }
	//user_id - уникальный идентификатор пользователя.
	//	Параметры запроса:
	//currency - валюта баланса.
}

func transaction(c *gin.Context) {
	var id struct {
		Id int `json:"user_id"`
	}
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(jsonData, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transactions []model.Transaction

	q := fmt.Sprintf("SELECT * from transactions where user_id = %d", id.Id)
	rows, err := database.DB.Queryx(q)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "balance has no transactions"})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.TransactionId, &t.UserId, &t.Amount, &t.Operation, &t.Date); err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		transactions = append(transactions, t)
	}

	c.IndentedJSON(http.StatusOK, &transactions)
	//Тело запроса:
	//user_id - уникальный идентификатор пользователя.
	//	Параметры запроса:
	//sort - сортировка списка транзакций.
}
func topUp(c *gin.Context) {
	var req struct {
		ID     int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(jsonData, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you can't do this operation with negative value"})
		return
	}
	var user model.User

	q := fmt.Sprintf("SELECT * from users where user_id = %d", req.ID)
	err = database.DB.QueryRow(q).Scan(&user.ID, &user.Balance)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			user.Balance = user.Balance + req.Amount
			tx := database.DB.MustBegin()
			fmt.Println(user.ID, user.Balance, req.ID)
			tx.MustExec("INSERT INTO users values ($1,$2)", req.ID, user.Balance)
			err = tx.Commit()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusOK, &user)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Balance += req.Amount
	_, err = database.DB.NamedExec("Update users SET balance=:balance WHERE user_id =:user_id", user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, &user)
	return
	//Тело запроса:
	//user_id - идентификатор пользователя,
	//	amount - сумма пополнения в RUB.
}
func debit(c *gin.Context) {
	var req struct {
		ID     int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(jsonData, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you can't do this operation with negative value"})
		return
	}
	var user model.User

	q := fmt.Sprintf("SELECT * from users where user_id = %d", req.ID)
	err = database.DB.QueryRow(q).Scan(&user.ID, &user.Balance)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("balance not found for this user")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount > user.Balance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds for this operation"})
		return
	}
	user.Balance -= req.Amount
	_, err = database.DB.NamedExec("Update users SET balance=:balance WHERE user_id =:user_id", user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, &user)
	return
	//Тело запроса:
	//user_id - идентификатор пользователя,
	//	amount - сумма списания в RUB.
}
func transfer(c *gin.Context) {
	var req struct {
		ID     int     `json:"user_id"`
		ToID   int     `json:"to_id"`
		Amount float64 `json:"amount"`
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(jsonData, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user model.User
	var userTo model.User
	q := fmt.Sprintf("SELECT * from users where user_id = %d", req.ID)
	qTo := fmt.Sprintf("SELECT * from users where user_id = %d", req.ToID)
	err = database.DB.QueryRow(q).Scan(&user.ID, &user.Balance)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("balance not found for this user")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.QueryRow(qTo).Scan(&userTo.ID, &userTo.Balance)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("balance not found for this user")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount > user.Balance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds for this operation"})
		return
	}
	if userTo.Balance > user.Balance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds for this operation"})
		return
	}

	user.Balance -= req.Amount
	_, err = database.DB.NamedExec("Update users SET balance=:balance WHERE user_id =:user_id", user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userTo.Balance += req.Amount
	_, err = database.DB.NamedExec("Update users SET balance=:balance WHERE user_id =:user_id", userTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, &userTo)
	return
	//Тело запроса:
	//user_id - идентификатор пользователя, с баланса которого списываются средства,
	//	to_id - идентификатор пользователя, на баланс которого начисляются средства,
	//	amount - сумма перевода в RUB.
}
