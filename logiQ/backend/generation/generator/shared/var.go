package shared

import "errors"

var (
	ErrGenerationBudgetExceeded = errors.New("generation: generation budget exceeded")
	ErrRngRequired              = errors.New("generation: rng is required")
)

const MAX_ATTEMPTS = 64
const MAX_EXPR_LENGTH = 60
