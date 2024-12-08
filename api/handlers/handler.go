package handlers

import (
	"database/sql"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContractService interface {
	ExecContract(methodName string, params ...interface{}) error
	CallContract(methodName string, result *[]interface{}, params ...interface{}) error
}

type Handler struct {
	DB             *sql.DB
	ContractClient ContractService
}

func (h *Handler) GetContractValue(c *gin.Context) {
	var output []interface{}
	err := h.ContractClient.CallContract("get", &output)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": output[0]})
}

func (h *Handler) SetContractValue(c *gin.Context) {
	var request struct {
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	contractValue := new(big.Int)
	if _, ok := contractValue.SetString(request.Value, 10); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value, must be a numeric string"})
		return
	}

	err := h.ContractClient.ExecContract("set", contractValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "value set successfully"})
}

func (h *Handler) SyncContractValue(c *gin.Context) {
	var output []interface{}
	err := h.ContractClient.CallContract("get", &output)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	value := output[0].(*big.Int).String()

	query := `INSERT INTO simpleStorageContract (value, timestamp) VALUES ($1, $2)`
	_, err = h.DB.Exec(query, value, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to synchronize value with database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Value synchronized successfully",
		"value":   output[0],
	})
}

func (h *Handler) CheckContractValue(c *gin.Context) {
	var output []interface{}
	err := h.ContractClient.CallContract("get", &output)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dbValue int64
	query := `SELECT value FROM simpleStorageContract ORDER BY timestamp DESC LIMIT 1`
	err = h.DB.QueryRow(query).Scan(&dbValue)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"result": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch value from database"})
		return
	}

	isEqual := false
	contractValue, ok := output[0].(*big.Int)
	if ok {
		if contractValue.Int64() == dbValue {
			isEqual = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": isEqual,
	})
}
