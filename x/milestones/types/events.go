package types

// Event types for the milestones module
const (
	EventSubmitProject    = "milestones_submit_project"
	EventVoteProject      = "milestones_vote_project"
	EventAttestMilestone  = "milestones_attest_milestone"
	EventReleaseTranche   = "milestones_release_tranche"
)

// Attribute keys for events
const (
	AttrProjectID = "project_id"
	AttrHash      = "hash"
	AttrAmount    = "amount"
	AttrDenom     = "denom"
	AttrOwner     = "owner"
)
