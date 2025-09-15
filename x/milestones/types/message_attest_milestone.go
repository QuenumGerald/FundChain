package types

func NewMsgAttestMilestone(creator string, projectId string, milestoneHash string) *MsgAttestMilestone {
	return &MsgAttestMilestone{
		Creator:       creator,
		ProjectId:     projectId,
		MilestoneHash: milestoneHash,
	}
}
