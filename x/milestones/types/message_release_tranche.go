package types

func NewMsgReleaseTranche(creator string, projectId string) *MsgReleaseTranche {
	return &MsgReleaseTranche{
		Creator:   creator,
		ProjectId: projectId,
	}
}
