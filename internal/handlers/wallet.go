package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	DB *sql.DB
}

func NewWalletHandler(db *sql.DB) *WalletHandler {
	return &WalletHandler{DB: db}
}

// GetWallets godoc
// @Summary Get user wallets
// @Description Get all cryptocurrency wallets for the authenticated user
// @Tags Wallets
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Success 200 {object} []models.Wallet "List of wallets"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/wallets [get]
// GetWallets godoc
// @Summary Get user wallets
// @Description Get all cryptocurrency wallets for the authenticated user, including virtual ones for coins without a DB record
// @Tags Wallets
// @Accept json
// @Produce json
// @Security BackendSecret
// @Security BearerAuth
// @Success 200 {object} []models.Wallet "List of wallets"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/wallets [get]
func (h *WalletHandler) GetWallets(c *gin.Context) {
	// 1. Get user ID from token (set by middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "User ID not found in context"})
		return
	}
	// We assume userID is string or uuid, let's treat it as string for query params if driver handles it,
	// but better to cast if it's uuid.UUID. Middleware usually sets it as string or UUID.
	// routes.go middleware likely sets it. models.User has ID as uuid.UUID.

	// 2. Query coins using LEFT JOIN to get all potential wallets
	// We use COALESCE to provide defaults for minimal required fields
	query := `
		SELECT 
			w.id, 
            c.id as coin_id,
			COALESCE(w.balance, '0'), 
            COALESCE(w.frozen_balance, '0'),
			c.name, c.ticker, c.logo, c.decimal
		FROM coins c
		LEFT JOIN wallets w ON c.id = w.coin_id AND w.user_id = $1
		WHERE c.status = 1
		ORDER BY c.ticker ASC
	`

	rows, err := h.DB.Query(query, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database_error", "message": "Failed to fetch wallets"})
		return
	}
	defer rows.Close()

	// 3. Parse results
	// We define a transient struct to hold the combined data
	type VirtualWalletResponse struct {
		ID            *string `json:"id"` // Pointer to allow null, or we can use string "0000.."
		CoinID        int     `json:"coin_id"`
		Balance       string  `json:"balance"`
		FrozenBalance string  `json:"frozen_balance"`
		CoinName      string  `json:"coin_name"`
		CoinTicker    string  `json:"coin_ticker"`
		CoinLogo      *string `json:"coin_logo"`
		CoinDecimal   int     `json:"coin_decimal"`
	}

	wallets := []VirtualWalletResponse{}
	for rows.Next() {
		var id sql.NullString
		var coinID int
		var balance, frozenBalance string
		var name, ticker string
		var logo *string
		var decimal int

		if err := rows.Scan(
			&id, &coinID, &balance, &frozenBalance,
			&name, &ticker, &logo, &decimal,
		); err != nil {
			continue // Skip erroneous rows
		}

		// Construct response
		var idStr *string
		if id.Valid {
			s := id.String
			idStr = &s
		}

		w := VirtualWalletResponse{
			ID:            idStr,
			CoinID:        coinID,
			Balance:       balance,
			FrozenBalance: frozenBalance,
			CoinName:      name,
			CoinTicker:    ticker,
			CoinLogo:      logo,
			CoinDecimal:   decimal,
		}
		wallets = append(wallets, w)
	}

	c.JSON(http.StatusOK, wallets)
}
