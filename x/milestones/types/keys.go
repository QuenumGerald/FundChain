package types

import "cosmossdk.io/collections"

const (
    // ModuleName defines the module name
    ModuleName = "milestones"

    // StoreKey defines the primary module store key
    StoreKey = ModuleName

    // GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
    // It should be synced with the gov module's name if it is ever changed.
    // See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
    GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix(0x00)

// ProjectKey is the prefix for storing projects by id
var ProjectKey = collections.NewPrefix(0x01)

// MilestoneKey is the prefix for storing milestones by (projectId, index)
// Use with collections.Pair[uint64, uint64]
var MilestoneKey = collections.NewPrefix(0x02)

// ProjectSeqKey is the prefix for the global project id sequence
var ProjectSeqKey = collections.NewPrefix(0x03)
