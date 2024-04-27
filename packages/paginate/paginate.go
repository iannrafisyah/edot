package paginate

import (
	"fmt"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	NextPage   string `json:"next_page"`
	PrevPage   string `json:"prev_page"`
}

// GetOffset :
func (r *Pagination) GetOffset() int {
	return (r.GetPage() - 1) * r.GetLimit()
}

// GetLimit : Default limit 10
func (r *Pagination) GetLimit() int {
	if r.Limit == 0 {
		r.Limit = 10
	}
	return r.Limit
}

// GetPage : Default page 1
func (r *Pagination) GetPage() int {
	if r.Page == 0 {
		r.Page = 1
	}
	return r.Page
}

// Next :
func (r *Pagination) Next(c echo.Context) string {
	if r.GetPage() <= r.TotalPages || r.GetPage() >= r.TotalPages {
		if r.GetPage() >= r.TotalPages && r.TotalPages == 0 {
			r.NextPage = r.PerPage(c, 1)
		} else if r.GetPage() >= r.TotalPages && r.TotalPages > 0 {
			r.NextPage = r.PerPage(c, r.TotalPages)
		} else {
			r.NextPage = r.PerPage(c, r.GetPage()+1)
		}
	} else {
		r.NextPage = r.PerPage(c, 1)
	}

	return r.NextPage
}

// Prev :
func (r *Pagination) Prev(c echo.Context) string {
	if r.GetPage() > 1 {
		if r.GetPage() > r.TotalPages {
			prevPage := r.TotalPages - 1
			if prevPage < 1 {
				prevPage = 1
			}
			r.PrevPage = r.PerPage(c, prevPage)
		} else {
			r.PrevPage = r.PerPage(c, r.GetPage()-1)
		}
	} else {
		r.PrevPage = r.PerPage(c, 1)
	}

	return r.PrevPage
}

// PerPage :
func (r *Pagination) PerPage(c echo.Context, page int) string {
	c.QueryParams().Set("page", strconv.Itoa(page))
	url := fmt.Sprintf("%s%s?%s", "localhost", c.Request().URL.Path, c.QueryParams().Encode())
	return url
}

// Paginate :
func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
