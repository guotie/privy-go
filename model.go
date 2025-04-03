package privy

type BaseResp struct {
	AuthorizationThreshold int `json:"authorization_threshold"`
}

// ResultCreateWallet server wallet response
type ResultCreateWallet struct {
	BaseResp
	ID        string `json:"id"`
	Address   string `json:"address"`
	ChainType string `json:"chain_type"`
}

// LinkedAccount struct to represent the linked account
type LinkedAccount struct {
	Type             string `json:"type"`
	Address          string `json:"address"`
	ChainType        string `json:"chain_type"`
	ChainID          string `json:"chain_id"`
	WalletClient     string `json:"wallet_client"`
	WalletClientType string `json:"wallet_client_type"`
	ConnectorType    string `json:"connector_type"`
	VerifiedAt       int64  `json:"verified_at"`
	FirstVerifiedAt  int64  `json:"first_verified_at"`
	LatestVerifiedAt int64  `json:"latest_verified_at"`
	Delegated        bool   `json:"delegated"`
	Imported         bool   `json:"imported"`
	RecoveryMethod   string `json:"recovery_method"`
}

// UserInfo privy user info
type UserInfo struct {
	ID               string          `json:"id"`
	CreatedAt        int64           `json:"created_at"`
	LinkedAccounts   []LinkedAccount `json:"linked_accounts"`
	MfaMethods       []string        `json:"mfa_methods"` // Adjust type as needed
	HasAcceptedTerms bool            `json:"has_accepted_terms"`
	IsGuest          bool            `json:"is_guest"`
}

type TxData struct {
	Hash  string `json:"hash"`
	Caip2 string `json:"caip2"`
}
type WalletRPCResult struct {
	Method string `json:"method"`
	Data   TxData `json:"data"`
}

type SignTypedDataV4Result struct {
	Method string `json:"method"`
	Data   struct {
		Signature string `json:"signature"`
		Encoding  string `json:"encoding"`
	} `json:"data"`
}
