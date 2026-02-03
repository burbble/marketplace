package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/burbble/marketplace/internal/exchange"
)

type rateResponse struct {
	Rate float64 `json:"rate"`
}

type ExchangeHandler struct {
	provider exchange.RateProvider
}

func NewExchangeHandler(provider exchange.RateProvider) *ExchangeHandler {
	return &ExchangeHandler{provider: provider}
}

// @Summary      Get USDT/RUB exchange rate
// @Tags         exchange
// @Produce      json
// @Success      200  {object}  rateResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /exchange/rate [get]
func (h *ExchangeHandler) GetRate(c *gin.Context) {
	rate, err := h.provider.GetUSDTRate(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to get exchange rate")
		return
	}

	c.JSON(http.StatusOK, rateResponse{Rate: rate})
}
