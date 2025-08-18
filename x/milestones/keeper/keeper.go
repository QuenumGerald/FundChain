package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"fundchain/x/milestones/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	// collections
	Projects   collections.Map[uint64, types.Project]
	Milestones collections.Map[collections.Pair[uint64, string], types.Milestone]
	ProjectSeq collections.Sequence

	bankKeeper types.BankKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		bankKeeper: bankKeeper,
		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	// Initialize collections
	k.Projects = collections.NewMap(sb, types.ProjectKey, "projects", collections.Uint64Key, codec.CollValue[types.Project](cdc))
	k.Milestones = collections.NewMap(
		sb,
		types.MilestoneKey,
		"milestones",
		collections.PairKeyCodec(collections.Uint64Key, collections.StringKey),
		codec.CollValue[types.Milestone](cdc),
	)
	k.ProjectSeq = collections.NewSequence(sb, types.ProjectSeqKey, "project_seq")

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// AppendProject stores a new project, auto-assigning an ID, and returns it.
func (k Keeper) AppendProject(ctx context.Context, p types.Project) (uint64, error) {
	id, err := k.ProjectSeq.Next(ctx)
	if err != nil {
		return 0, err
	}
	p.Id = id
	if err := k.Projects.Set(ctx, id, p); err != nil {
		return 0, err
	}
	return id, nil
}

// GetProject retrieves a project by ID.
func (k Keeper) GetProject(ctx context.Context, id uint64) (types.Project, bool, error) {
	p, err := k.Projects.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Project{}, false, nil
		}
		return types.Project{}, false, err
	}
	return p, true, nil
}

// SetProject updates a project.
func (k Keeper) SetProject(ctx context.Context, p types.Project) error {
	return k.Projects.Set(ctx, p.Id, p)
}

// AppendMilestone stores or updates a milestone keyed by (project_id, hash).
func (k Keeper) AppendMilestone(ctx context.Context, m types.Milestone) error {
	key := collections.Join(m.ProjectId, m.Hash)
	return k.Milestones.Set(ctx, key, m)
}
