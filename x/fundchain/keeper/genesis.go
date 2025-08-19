package keeper

import (
	"context"
	"errors"

	"fundchain/x/fundchain/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	genesis := types.DefaultGenesis()

	params, err := k.Params.Get(ctx)
	if err != nil {
		// If params are not yet set in store, fall back to defaults instead of erroring.
		if errors.Is(err, collections.ErrNotFound) {
			params = types.DefaultParams()
		} else {
			return nil, err
		}
	}
	genesis.Params = params

	return genesis, nil
}
