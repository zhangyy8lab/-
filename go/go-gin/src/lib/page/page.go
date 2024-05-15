package page

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// Page 分页
func Page(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	order := c.GetBool("order_by")
	_, paramsOrder := c.GetQuery("order_by")
	return func(db *gorm.DB) *gorm.DB {
		//q := c.URL.Query()
		q := c.Request.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("pageSize"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		if order || paramsOrder {
			return db.Offset(offset).Limit(pageSize)
		} else {
			return db.Offset(offset).Limit(pageSize).Order("id desc")
		}

	}
}
