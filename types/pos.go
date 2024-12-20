package types

import (
	"fmt"
	"strconv"
	"strings"

	mt "github.com/minetest-go/types"
)

// parses a comma separated string to a pos object
func ParsePos(pos string) (*mt.Pos, error) {
	parts := strings.Split(pos, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid part-count: %d, should be 3", len(parts))
	}

	x, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing x part: %v", err)
	}

	y, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing y part: %v", err)
	}

	z, err := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing z part: %v", err)
	}

	return mt.NewPos(int(x), int(y), int(z)), nil
}
