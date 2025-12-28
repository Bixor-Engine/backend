package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Bixor-Engine/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

// GetTransactions godoc
// @Summary Get user transactions
// @Description Get transaction history for the authenticated user
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "List of transactions with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// 1. Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "User ID not found in context"})
		return
	}

	// 2. Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 3. Query transactions with joins for better context (e.g. coin ticker)
	// Note: We need to join wallets -> coins to get the ticker
	query := `
		SELECT 
			t.id, t.user_id, t.wallet_id, t.type, t.amount, t.fee, 
			t.description, t.reference_id, t.payment_method, t.status, t.created_at, t.updated_at,
			c.ticker, c.name
		FROM transactions t
		JOIN wallets w ON t.wallet_id = w.id
		JOIN coins c ON w.coin_id = c.id
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.DB.Query(query, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database_error", "message": "Failed to fetch transactions"})
		return
	}
	defer rows.Close()

	type TransactionWithMeta struct {
		models.Transaction
		CoinTicker string `json:"coin_ticker"`
		CoinName   string `json:"coin_name"`
	}

	transactions := []TransactionWithMeta{}
	for rows.Next() {
		var t TransactionWithMeta
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.WalletID, &t.Type, &t.Amount, &t.Fee,
			&t.Description, &t.ReferenceID, &t.PaymentMethod, &t.Status, &t.CreatedAt, &t.UpdatedAt,
			&t.CoinTicker, &t.CoinName,
		); err != nil {
			continue
		}
		transactions = append(transactions, t)
	}

	// 4. Get total count for pagination
	var total int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM transactions WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		total = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
