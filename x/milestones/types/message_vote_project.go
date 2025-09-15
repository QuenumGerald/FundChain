package types

func NewMsgVoteProject(creator string, projectId string, support bool) *MsgVoteProject {
	return &MsgVoteProject{
		Creator:   creator,
		ProjectId: projectId,
		Support:   support,
	}
}
