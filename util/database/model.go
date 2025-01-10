package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MySQLFilter struct {
	Where   gin.H
	Like    *Like
	Not     gin.H
	Or      gin.H
	Preload []string
	Limit   int64
	Offset  int64
	Order   string
	Sort    string
}

type Like struct {
	Columns []string
	Value   string
}

func (f *MySQLFilter) SetPagination(page, limit int64) {
	f.Limit = limit
	f.Offset = (page - 1) * limit
}

func (f *MySQLFilter) SetLike(req map[string][]string) {
	if len(req) > 0 {
		for k, v := range req {
			f.Like = &Like{Columns: v, Value: k}
		}
	}
}

func BuildMySQLFilter(dbMySQL *gorm.DB, filter MySQLFilter) *gorm.DB {
	if len(filter.Where) > 0 {
		for k, v := range filter.Where {
			dbMySQL = dbMySQL.Where(k, v)
		}
	}

	if filter.Like != nil {
		for _, column := range filter.Like.Columns {
			dbMySQL = dbMySQL.Or(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", filter.Like.Value))
		}
	}

	if len(filter.Not) > 0 {
		for k, v := range filter.Not {
			dbMySQL = dbMySQL.Not(k, v)
		}
	}

	if len(filter.Or) > 0 {
		for _, cond := range filter.Or {
			dbMySQL = dbMySQL.Or(cond)
		}
	}

	if len(filter.Preload) > 0 {
		for _, v := range filter.Preload {
			dbMySQL = dbMySQL.Preload(v)
		}
	}

	if filter.Limit > 0 {
		dbMySQL = dbMySQL.Limit(int(filter.Limit))
	}

	if filter.Offset > 0 {
		dbMySQL = dbMySQL.Offset(int(filter.Offset))
	}

	if filter.Order != "" && filter.Sort != "" {
		dbMySQL = dbMySQL.Order(fmt.Sprintf("%s %s", filter.Order, filter.Sort))
	}

	return dbMySQL
}

type Filter struct {
	Equal  gin.H
	Search map[string][]string
	Page   int64
	Limit  int64
	Sort   string
	Order  string
}

func (f *Filter) SetSearch(request map[string][]string) {
	if len(request) > 0 {
		f.Search = request
	}
}

func (f *Filter) SetPagination(Page, Limit int64) {
	f.Page = Page
	f.Limit = Limit
}

func (f *Filter) SetSortAndOrder(Sort, Order string) {
	f.Sort = Sort
	f.Order = Order
}

func (f *Filter) ToMySQLFilter() MySQLFilter {
	mFilter := MySQLFilter{
		Where:  f.Equal,
		Limit:  f.Limit,
		Offset: (f.Page - 1) * f.Limit,
		Order:  f.Sort,
		Sort:   f.Order,
	}
	mFilter.SetLike(f.Search)

	return mFilter
}
