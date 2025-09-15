package types

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// NewParams creates a new Params instance.
func NewParams() Params {
    return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
    return Params{
        Treasury: "",
        Denom:    "ufund",
    }
}

// Validate validates the set of params.
func (p Params) Validate() error {
    // Treasury: optional bech32 address
    if len(p.Treasury) > 0 {
        if _, err := sdk.AccAddressFromBech32(p.Treasury); err != nil {
            return fmt.Errorf("invalid treasury address: %w", err)
        }
    }
    // Denom: required, non-empty and valid denom
    if len(p.Denom) > 0 {
        if err := sdk.ValidateDenom(p.Denom); err != nil {
            return fmt.Errorf("invalid denom: %w", err)
        }
    }
    return nil
}

// Legacy params module compatibility (optional)
// Keys
var (
    ParamStoreKeyTreasury = []byte("Treasury")
    ParamStoreKeyDenom    = []byte("Denom")
)

// ParamKeyTable returns the KeyTable for params module
func ParamKeyTable() paramtypes.KeyTable {
    return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns the param pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
    return paramtypes.ParamSetPairs{
        paramtypes.NewParamSetPair(ParamStoreKeyTreasury, &p.Treasury, validateTreasury),
        paramtypes.NewParamSetPair(ParamStoreKeyDenom, &p.Denom, validateDenom),
    }
}

func validateTreasury(v interface{}) error {
    s, ok := v.(string)
    if !ok {
        return fmt.Errorf("invalid type for treasury: %T", v)
    }
    if s == "" {
        return nil
    }
    if _, err := sdk.AccAddressFromBech32(s); err != nil {
        return fmt.Errorf("invalid treasury address: %w", err)
    }
    return nil
}

func validateDenom(v interface{}) error {
    s, ok := v.(string)
    if !ok {
        return fmt.Errorf("invalid type for denom: %T", v)
    }
    if s != "" {
        if err := sdk.ValidateDenom(s); err != nil {
            return fmt.Errorf("invalid denom: %w", err)
        }
    }
    return nil
}
