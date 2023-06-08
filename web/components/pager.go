package components

import (
	"net/http"
	"strconv"
)

type PagerModel struct {
	Offset      int
	Count       int
	Limit       int
	Pages       []any
	Pagecount   int
	CurrentPage int //1-indexed page number
	Previous    int
	Next        int
}

func Pager(r *http.Request, limit int, count int) *PagerModel {
	m := &PagerModel{
		Limit: limit,
		Count: count,
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err == nil {
		m.Offset = int(page-1) * limit
		m.CurrentPage = int(page)
	} else {
		m.CurrentPage = 1
	}

	m.Previous = m.CurrentPage - 1
	m.Next = m.CurrentPage + 1

	m.Pagecount = (count / limit) + 1

	m.Pages = make([]any, m.Pagecount)
	for i := 0; i < m.Pagecount; i++ {
		m.Pages[i] = i + 1
	}

	return m
}
