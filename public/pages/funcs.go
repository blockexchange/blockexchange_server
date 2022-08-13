package pages

import (
	"fmt"
	"time"
)

func prettysize(num int) string {
	if num > (1000 * 1000) {
		return fmt.Sprintf("%d MB", num/(1000*1000))
	} else if num > 1000 {
		return fmt.Sprintf("%d kB", num/(1000))
	} else {
		return fmt.Sprintf("%d bytes", num)
	}
}

func formattime(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("2006-02-01")
}

