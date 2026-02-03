package pagination

import (
	"errors"
	"strings"
)

type sortOrdering = string

const (
	defaultPageSize  = 24
	DefaultSortField = "created_at"

	ASC  sortOrdering = "ASC"
	DESC sortOrdering = "DESC"

	defaultSortOrder = DESC
)

type PagePagination struct {
	Page     uint64 `json:"page" form:"page" binding:"required,min=1"`
	PageSize uint64 `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100"`
	Limit    uint64 `json:"limit" form:"limit" binding:"omitempty,min=0,max=100"`
}

func (p PagePagination) GetOffset() uint64 {
	return (p.Page - 1) * p.GetLimit()
}

func (p PagePagination) GetLimit() uint64 {
	if p.Limit > 0 {
		return p.Limit
	}

	if p.PageSize > 0 {
		return p.PageSize
	}

	return defaultPageSize
}

type SortRequest struct {
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

func (s SortRequest) GetSortClause() string {
	field := s.SortBy
	if field == "" {
		field = DefaultSortField
	}

	order := strings.ToUpper(s.SortOrder)
	if order != ASC && order != DESC {
		order = defaultSortOrder
	}

	return field + " " + order
}

type SortFieldsRequest struct {
	SortFields string `json:"sort_fields" form:"sort_fields"`
}

func (s SortFieldsRequest) ParseSortFields() ([]string, error) {
	if s.SortFields == "" {
		return []string{DefaultSortField + " " + defaultSortOrder}, nil
	}

	parts := strings.Split(s.SortFields, ",")
	clauses := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		items := strings.Split(part, ":")
		field := strings.TrimSpace(items[0])
		if field == "" {
			return nil, errors.New("invalid sort field")
		}

		order := defaultSortOrder
		if len(items) > 1 {
			orderCandidate := strings.ToUpper(strings.TrimSpace(items[1]))
			if orderCandidate == ASC || orderCandidate == DESC {
				order = orderCandidate
			} else {
				return nil, errors.New("invalid sort order: " + items[1])
			}
		}

		clauses = append(clauses, field+" "+order)
	}

	return clauses, nil
}

func BuildSortClause(requests []SortRequest) string {
	clauses := make([]string, 0, len(requests))
	for _, r := range requests {
		clauses = append(clauses, r.GetSortClause())
	}

	return strings.Join(clauses, ", ")
}
