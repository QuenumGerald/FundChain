package types

func NewMsgSubmitProject(creator string, title string, budget string, ipfsHash string) *MsgSubmitProject {
	return &MsgSubmitProject{
		Creator:  creator,
		Title:    title,
		Budget:   budget,
		IpfsHash: ipfsHash,
	}
}
