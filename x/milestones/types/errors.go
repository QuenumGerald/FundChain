package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/milestones module sentinel errors
var (
	ErrInvalidSigner   = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrNotFound        = errors.Register(ModuleName, 1101, "not found")
	ErrInvalidParam    = errors.Register(ModuleName, 1102, "invalid parameter")
	ErrUnauthorized    = errors.Register(ModuleName, 1103, "unauthorized")
	ErrZeroAmount      = errors.Register(ModuleName, 1104, "amount must be greater than zero")
)
