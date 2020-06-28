package messages

type BusinessMessage struct {
	Type string `json:"type"`
}

const (
	P2PMessageTypeRFP              = "rfp"
	P2PMessageTypeProposal         = "proposal"
	P2PMessageTypeContractRequest  = "contractrequest"
	P2PMessageTypeContractAccepted = "contractaccepted"
	P2PMessageTypePO               = "po"
)

const (
	BlockchainMessageTypeContract = "contract"
	BlockchainMessageTypePO       = "po"
)
