package mnc_solo

import (
	"errors"
)

var (
	ErrVerifyHeaderFailed = errors.New("verify failed")
	ErrNotSoloNode        = errors.New("not solo note.")
)
