package dir

import (
	"fmt"
	"strconv"
)

type Prefix struct {
	text string
}

func (pre Prefix) getValue() int {
	i, err := strconv.Atoi(pre.text)
	if err != nil {
		return -1
	}
	return i
}

func (pre Prefix) getWidth() int {
	return len(pre.text)
}

func (pre Prefix) Increment() (string, error) {
	i := pre.getValue()
	if i < 0 {
		return "", fmt.Errorf("failed to increment index")
	}
	i += 1
	w := pre.getWidth()
	return fmt.Sprintf("%0"+strconv.Itoa(w)+"d", i), nil
}
