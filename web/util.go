package web

import (
	"net/http"
	"strconv"
)

func GetLimitOffset(r *http.Request, default_limit int) (int, int) {
	q := r.URL.Query()
	if q.Get("limit") != "" && q.Get("offset") != "" {
		var limit, offset int64
		limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)
		if err != nil {
			return default_limit, 0
		}

		offset, err = strconv.ParseInt(q.Get("offset"), 10, 64)
		if err != nil {
			return default_limit, 0
		}

		return int(limit), int(offset)
	} else {
		return default_limit, 0
	}
}
