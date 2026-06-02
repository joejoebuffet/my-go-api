package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountStatusRequest struct {
	AccountNo string `json:"accountNo"`
	Status    string `json:"status"`
}

type HelloController struct {
	DB *sql.DB // Replaces @Autowired private JdbcTemplate jdbcTemplate;
}

func NewHelloController(db *sql.DB) *HelloController {
	return &HelloController{DB: db}
}

// 1. GET /hello?accountNo=xxxx
func (ctrl *HelloController) GetAccountHolder(c *gin.Context) {
	accountNo := c.Query("accountNo") // Replaces @RequestParam("accountNo")
	fmt.Printf("DEBUG INPUT: Received accountNo = [%s]\n", accountNo)

	// Change this line:
	sqlQuery := "SELECT gl_segment FROM bill_media WHERE rec_id = '3' AND account_no = $1"

	var holderName string
	// QueryRow handles single-record lookups. Parameter types are handled implicitly by Go.
	err := ctrl.DB.QueryRow(sqlQuery, accountNo).Scan(&holderName)

	if err != nil {
		fmt.Printf("DEBUG DATABASE ERROR: %v\n", err)
		c.String(http.StatusOK, "Error: Account number "+accountNo+" not found.")
		return
	}

	c.String(http.StatusOK, "Account Holder: "+holderName)
}

// 2. POST /updateStatus
func (ctrl *HelloController) UpdateAccountStatus(c *gin.Context) {
	var request AccountStatusRequest
	// Replaces @RequestBody mapping payload to your class
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "FAILURE: Invalid payload format")
		return
	}

	sqlQuery := "UPDATE cim_account SET status = $1 WHERE account_no = $2"
	// Exec runs updates/deletes and returns performance metadata
	result, err := ctrl.DB.Exec(sqlQuery, request.Status, request.AccountNo)
	if err != nil {
		c.String(http.StatusOK, "FAILURE: Database error: "+err.Error())
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		c.String(http.StatusOK, fmt.Sprintf("SUCCESS: Account %s updated to status '%s'", request.AccountNo, request.Status))
	} else {
		c.String(http.StatusOK, "WARNING! No record found for account_no "+request.AccountNo)
	}
}

// 3. POST /updateStatusBulk
func (ctrl *HelloController) UpdateAccountStatusBulk(c *gin.Context) {
	var payloadList []AccountStatusRequest
	if err := c.ShouldBindJSON(&payloadList); err != nil {
		c.String(http.StatusBadRequest, "FAILURE: Invalid payload format")
		return
	}

	sqlQuery := "UPDATE cim_account SET status = $1 WHERE account_no = $2"
	// Go does batch processing safely via Database Transactions (Tx) and Prepared Statements
	tx, err := ctrl.DB.Begin()
	if err != nil {
		c.String(http.StatusOK, "FAILURE: Batch transaction start error: "+err.Error())
		return
	}
	// If anything crashes mid-loop, safely roll back the whole batch
	defer tx.Rollback()

	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		c.String(http.StatusOK, "FAILURE: Batch preparation error: "+err.Error())
		return
	}
	defer stmt.Close()

	totalRowsUpdated := 0
	for _, recordItem := range payloadList {
		result, err := stmt.Exec(recordItem.Status, recordItem.AccountNo)
		if err != nil {
			c.String(http.StatusOK, "FAILURE: Batch database error: "+err.Error())
			return
		}
		rows, _ := result.RowsAffected()
		totalRowsUpdated += int(rows)
	}

	// Commit the entire transaction block to the DB
	if err := tx.Commit(); err != nil {
		c.String(http.StatusOK, "FAILURE: Transaction commit failed: "+err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("SUCCESS: Record-batch complete. Total records updated: %d", totalRowsUpdated))
}
