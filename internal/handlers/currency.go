package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Bixor-Engine/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CurrencyHandler struct {
	DB *sql.DB
}

func NewCurrencyHandler(db *sql.DB) *CurrencyHandler {
	return &CurrencyHandler{
		DB: db,
	}
}

// GetCoins godoc
// @Summary Get list of all supported coins
// @Description Retrieve a list of all supported cryptocurrencies
// @Tags Currency
// @Accept json
// @Produce json
// @Success 200 {object} models.CoinListResponse "List of coins"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/currency [get]
func (h *CurrencyHandler) GetCoins(c *gin.Context) {
	query := `
		SELECT 
			id, name, ticker, decimal, price_decimal, logo, price,
			deposit_gateway, withdraw_gateway, deposit_fee, withdraw_fee,
			deposit_fee_type, withdraw_fee_type, confirmation, status,
			withdraw_status, deposit_status, website, explorer, explorer_tx,
			explorer_address, created_at, updated_at
		FROM coins
		WHERE status = 1
		ORDER BY name ASC
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve coins",
		})
		return
	}
	defer rows.Close()

	var coins []models.Coin
	for rows.Next() {
		var coin models.Coin
		var depositGateway, withdrawGateway pq.StringArray
		var price sql.NullString
		var depositFee, withdrawFee sql.NullString

		err := rows.Scan(
			&coin.ID, &coin.Name, &coin.Ticker, &coin.Decimal, &coin.PriceDecimal,
			&coin.Logo, &price, &depositGateway, &withdrawGateway,
			&depositFee, &withdrawFee, &coin.DepositFeeType, &coin.WithdrawFeeType,
			&coin.Confirmation, &coin.Status, &coin.WithdrawStatus, &coin.DepositStatus,
			&coin.Website, &coin.Explorer, &coin.ExplorerTx, &coin.ExplorerAddress,
			&coin.CreatedAt, &coin.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "database_error",
				"message": "Failed to scan coin data",
			})
			return
		}

		// Convert arrays
		coin.DepositGateway = []string(depositGateway)
		coin.WithdrawGateway = []string(withdrawGateway)

		// Convert price and fees
		if price.Valid {
			coin.Price = price.String
		} else {
			coin.Price = "0"
		}

		if depositFee.Valid {
			coin.DepositFee = &depositFee.String
		}

		if withdrawFee.Valid {
			coin.WithdrawFee = &withdrawFee.String
		}

		coins = append(coins, coin)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Error iterating coins",
		})
		return
	}

	c.JSON(http.StatusOK, models.CoinListResponse{
		Coins: coins,
		Total: len(coins),
	})
}

// GetCoinByTicker godoc
// @Summary Get coin information by ticker
// @Description Retrieve detailed information about a specific cryptocurrency by its ticker symbol
// @Tags Currency
// @Accept json
// @Produce json
// @Param ticker path string true "Coin ticker symbol (e.g., BTC, ETH)"
// @Success 200 {object} models.CoinResponse "Coin information"
// @Failure 404 {object} map[string]interface{} "Coin not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/currency/{ticker} [get]
func (h *CurrencyHandler) GetCoinByTicker(c *gin.Context) {
	ticker := strings.ToUpper(strings.TrimSpace(c.Param("ticker")))
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_ticker",
			"message": "Ticker parameter is required",
		})
		return
	}

	query := `
		SELECT 
			id, name, ticker, decimal, price_decimal, logo, price,
			deposit_gateway, withdraw_gateway, deposit_fee, withdraw_fee,
			deposit_fee_type, withdraw_fee_type, confirmation, status,
			withdraw_status, deposit_status, website, explorer, explorer_tx,
			explorer_address, created_at, updated_at
		FROM coins
		WHERE UPPER(ticker) = $1 AND status = 1
		LIMIT 1
	`

	var coin models.Coin
	var depositGateway, withdrawGateway pq.StringArray
	var price sql.NullString
	var depositFee, withdrawFee sql.NullString

	err := h.DB.QueryRow(query, ticker).Scan(
		&coin.ID, &coin.Name, &coin.Ticker, &coin.Decimal, &coin.PriceDecimal,
		&coin.Logo, &price, &depositGateway, &withdrawGateway,
		&depositFee, &withdrawFee, &coin.DepositFeeType, &coin.WithdrawFeeType,
		&coin.Confirmation, &coin.Status, &coin.WithdrawStatus, &coin.DepositStatus,
		&coin.Website, &coin.Explorer, &coin.ExplorerTx, &coin.ExplorerAddress,
		&coin.CreatedAt, &coin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "coin_not_found",
				"message": "Coin with ticker " + ticker + " not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "Failed to retrieve coin",
		})
		return
	}

	// Convert arrays
	coin.DepositGateway = []string(depositGateway)
	coin.WithdrawGateway = []string(withdrawGateway)

	// Convert price and fees
	if price.Valid {
		coin.Price = price.String
	} else {
		coin.Price = "0"
	}

	if depositFee.Valid {
		coin.DepositFee = &depositFee.String
	}

	if withdrawFee.Valid {
		coin.WithdrawFee = &withdrawFee.String
	}

	c.JSON(http.StatusOK, models.CoinResponse{
		Coin: coin,
	})
}
