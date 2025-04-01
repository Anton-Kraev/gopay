package typeconv

import (
	"fmt"
	"strconv"
	"strings"
)

func StringToInt64Slice(str string) ([]int64, error) {
	strSlice := strings.Split(strings.TrimSpace(str), ",")
	int64Slice := make([]int64, len(strSlice))

	var err error

	for i, el := range strSlice {
		int64Slice[i], err = strconv.ParseInt(el, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("typeconv.StringToInt64Slice: %w", err)
		}
	}

	return int64Slice, nil
}
