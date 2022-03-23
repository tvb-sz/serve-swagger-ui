package entry

import (
	"fmt"
	"strings"
)

// General entry structure encapsulation

// PageReq Paging general request parameters
type PageReq struct {
	// General parameters--page number
	Page int `json:"page" form:"page" binding:"omitempty"`
	// General parameters--quantity per page
	Limit int `json:"limit" form:"limit" binding:"omitempty,max=1000"`
}

// GetOffset get offset
func (s PageReq) GetOffset() int {
	return (s.GetPage() - 1) * s.GetLimit()
}

// GetPage Get the page number, cannot be less than 0
func (s PageReq) GetPage() int {
	if s.Page <= 0 {
		return 1
	}
	return s.Page
}

// GetLimit Get the number of displayed bars cannot be less than 0 cannot be greater than 1000
func (s PageReq) GetLimit() int {
	if s.Limit <= 0 {
		return 10
	} else if s.Limit > 1000 {
		return 1000
	} else {
		return s.Limit
	}
}

// SortReq Generic sort request parameters
type SortReq struct {
	// Common parameters--sort field name
	OrderBy string `form:"order_by" json:"order_by"`
	// General parameters--sort type <asc--ascending desc--descending>
	// enum:desc,asc
	Sort string `form:"sort" json:"sort" binding:"omitempty,oneof=desc asc"`
}

// GetOrderBy 组装排序参数
func (s SortReq) GetOrderBy() string {
	s.OrderBy = strings.Replace(s.OrderBy, "`", "", -1) // Remove malicious field backticks
	if s.OrderBy != "" && s.Sort != "" {
		return fmt.Sprintf("`%s` %s", s.OrderBy, s.Sort)
	}
	return "`id` desc"
}

// PageRes paginated response data
type PageRes struct {
	// The total number is returned in full or 0 is returned if the number is not required
	Total int64 `json:"total"`
	// list data
	List interface{} `json:"list"`
}

// NumRes return a number
type NumRes struct {
	Num int64 `json:"num"`
}

// BaseRes Basic response structure
type BaseRes struct {
	// Error code ZERO success, non-ZERO failure
	Code int64 `json:"code"`
	// Error message corresponding to the error code
	Msg string `json:"msg"`
}
