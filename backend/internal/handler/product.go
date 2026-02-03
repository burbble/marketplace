package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/internal/service"
	"github.com/burbble/marketplace/pkg/pagination"
)

var allowedProductSortFields = map[string]bool{
	"name":       true,
	"price":      true,
	"created_at": true,
	"brand":      true,
}

type productListQuery struct {
	Page       uint64 `form:"page"`
	PageSize   uint64 `form:"page_size"`
	SortFields string `form:"sort_fields"`
	CategoryID string `form:"category_id"`
	Brand      string `form:"brand"`
	MinPrice   *int   `form:"min_price"`
	MaxPrice   *int   `form:"max_price"`
	Search     string `form:"search"`
}

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// @Summary      List products
// @Tags         products
// @Produce      json
// @Param        page         query     int     false  "Page number"               default(1)
// @Param        page_size    query     int     false  "Page size"                 default(24)
// @Param        sort_fields  query     string  false  "Sort (e.g. price:asc,name:desc)"
// @Param        category_id  query     string  false  "Category UUID"
// @Param        brand        query     string  false  "Brand filter"
// @Param        min_price    query     int     false  "Min price (RUB)"
// @Param        max_price    query     int     false  "Max price (RUB)"
// @Param        search       query     string  false  "Search by name"
// @Success      200  {object}  domain.ProductList
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var q productListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 24
	}

	pag := pagination.PagePagination{Page: q.Page, PageSize: q.PageSize}

	sfr := pagination.SortFieldsRequest{SortFields: q.SortFields}
	sortClauses, err := sfr.ParseSortFields()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for _, clause := range sortClauses {
		field := strings.Fields(clause)[0]
		if !allowedProductSortFields[field] {
			errorResponse(c, http.StatusBadRequest, "invalid sort field: "+field)
			return
		}
	}

	filter := domain.ProductFilter{
		Limit:    pag.GetLimit(),
		Offset:   pag.GetOffset(),
		SortBy:   sortClauses,
		MinPrice: q.MinPrice,
		MaxPrice: q.MaxPrice,
	}

	if q.CategoryID != "" {
		id, err := uuid.Parse(q.CategoryID)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "invalid category_id")
			return
		}
		filter.CategoryID = &id
	}
	if q.Brand != "" {
		filter.Brand = &q.Brand
	}
	if q.Search != "" {
		filter.Search = &q.Search
	}

	result, err := h.svc.GetByFilter(c.Request.Context(), filter)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to get products")
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product UUID"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "product not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to get product")
		return
	}

	c.JSON(http.StatusOK, product)
}

// @Summary      List all brands
// @Tags         products
// @Produce      json
// @Success      200  {array}   string
// @Failure      500  {object}  ErrorResponse
// @Router       /brands [get]
func (h *ProductHandler) GetBrands(c *gin.Context) {
	brands, err := h.svc.GetBrands(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to get brands")
		return
	}

	c.JSON(http.StatusOK, brands)
}
