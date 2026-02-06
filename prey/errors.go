package prey

import "errors"

var (
	ErrMissingClient = errors.New("prey client not available in context")
	ErrWriteDisabled = errors.New("write operations are disabled (PREY_ALLOW_WRITE=false)")
)
