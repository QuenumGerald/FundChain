package types

func NewMsgSubmitProject(creator string, title string, budget string, ipfsHash string, reviewers []string, attestThreshold uint32) *MsgSubmitProject {
    return &MsgSubmitProject{
        Creator:         creator,
        Title:           title,
        Budget:          budget,
        IpfsHash:        ipfsHash,
        Reviewers:       reviewers,
        AttestThreshold: attestThreshold,
    }
}
