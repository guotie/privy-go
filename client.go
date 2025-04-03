package privy

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gowebpki/jcs"
)

type PrivyClient struct {
	apiID           string
	secretKey       string
	verificationKey []byte
	authKey         string
}

const (
	HeaderPrivyAppID         = "privy-app-id"
	HeaderPrivyAuthSignature = "privy-authorization-signature"
	PrivyAPIBase             = "https://api.privy.io"
	PrivyAuthBase            = "https://auth.privy.io"
	PrivyWalletRPC           = "https://auth.privy.io/api/v1/wallets/rpc"
	PrivyWalletSignRPC       = "https://auth.privy.io/api/v1/wallets/%s/rpc"
)

// NewPrivyClient creates a new PrivyClient instance
func NewPrivyClient(apiID, secretKey, verificationKey string) *PrivyClient {
	if !strings.Contains(verificationKey, "-----BEGIN PUBLIC KEY-----") {
		verificationKey = "-----BEGIN PUBLIC KEY-----\n" + verificationKey + "\n" + "-----END PUBLIC KEY-----"
	}

	if apiID == "" || secretKey == "" {
		panic("AppID or app secret is empty")
	}

	return &PrivyClient{
		apiID:           apiID,
		secretKey:       secretKey,
		verificationKey: []byte(verificationKey),
	}
}

func NewPrivyClientWithAuthkey(apiID, secretKey, verificationKey string, authKey string) *PrivyClient {
	if authKey == "" {
		panic("authkey is empty")
	}

	pc := NewPrivyClient(apiID, secretKey, verificationKey)
	pc.authKey = authKey
	return pc
}

func (pc *PrivyClient) SetAuthKey(key string) {
	pc.authKey = key
}

func doSignedRequest[T any](pc *PrivyClient, method, uri string, data any, sign bool, v *T) error {
	var reader io.Reader

	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		fmt.Println("uri: " + uri)
		return err
	}

	req.SetBasicAuth(pc.apiID, pc.secretKey)
	req.Header.Add(HeaderPrivyAppID, pc.apiID)
	req.Header.Set("Content-Type", "application/json")
	// fmt.Println("apiId", pc.apiID, pc.secretKey, pc.authKey)

	if sign {
		s, err1 := pc.GetAuthorizationSignature(uri, data)
		if err1 != nil {
			return err1
		}
		// fmt.Println("signature: " + s)
		req.Header.Add(HeaderPrivyAuthSignature, s)
	}

	// fmt.Println("req:", req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = unmarshalResp(resp, v)
	return err
}

// CreateWallet create server wallet
func (pc *PrivyClient) CreateWallet(chainType string) (*ResultCreateWallet, error) {
	uri := PrivyAPIBase + "/v1/wallets"
	if chainType == "" {
		chainType = "ethereum" // solana
	}

	var res ResultCreateWallet
	err := doSignedRequest(pc, http.MethodPost, uri, map[string]string{
		"chain_type": chainType,
	}, false, &res)

	return &res, err
}

// SendETHTransaction https://docs.privy.io/guide/server-wallets/usage/ethereum
// https://api.privy.io/v1/wallets/<wallet_id>/rpc
// https://docs.privy.io/guide/delegated-actions/usage/ethereum#using-privy-io-server-auth
func (pc *PrivyClient) SignTypedDataV4(walletId string, address string, data map[string]any) (signature string, err error) {
	res := SignTypedDataV4Result{}
	err = doSignedRequest(pc,
		http.MethodPost,
		PrivyWalletRPC,
		// fmt.Sprintf(PrivyWalletSignRPC, walletId),
		map[string]any{
			"address":    address,
			"chain_type": "ethereum",
			"method":     "eth_signTypedData_v4",
			"params": map[string]any{
				"typed_data": data,
			},
		},
		true,
		&res)
	return res.Data.Signature, err
}

// SendETHTransaction https://docs.privy.io/guide/server-wallets/usage/ethereum
func (pc *PrivyClient) SendETHTransaction(chainID string, owner string, transcation map[string]any) (hash string, err error) {
	// log.Info("Send ETH Transaction", "chainId", chainID, "owner", owner, "transcation", transcation)

	res := WalletRPCResult{}
	err = doSignedRequest(pc,
		http.MethodPost,
		PrivyWalletRPC,
		map[string]any{
			"method":     "eth_sendTransaction",
			"caip2":      "eip155:" + chainID,
			"chain_type": "ethereum",
			"address":    owner,
			"params": map[string]any{
				"transaction": transcation,
			},
		},
		true,
		&res)
	hash = res.Data.Hash
	return
}

// GetUser get privy user info
func (pc *PrivyClient) GetUser(did string) (*UserInfo, error) {
	uri := PrivyAuthBase + "/api/v1/users/" + did
	var res UserInfo
	err := doSignedRequest(pc, http.MethodGet, uri, nil, false, &res)

	// fmt.Println(res)
	return &res, err
}

func (pc *PrivyClient) GetAuthorizationSignature(url string, body any) (string, error) {
	if pc.authKey == "" {
		return "", errors.New("no auth key")
	}

	payload := map[string]any{
		"version": 1,
		"method":  "POST",
		"url":     url,
		"body":    body,
		"headers": map[string]string{
			"privy-app-id": pc.apiID,
		},
	}

	// JSON-canonicalize the payload
	serializedPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	canonicalPayload, err := jcs.Transform(serializedPayload)
	if err != nil {
		return "", err
	}
	// Replace this with your app's authorization key
	privateKeyAsString := strings.Replace(pc.authKey, "wallet-auth:", "", 1)

	// Convert your private key to PEM format
	privateKeyAsPem := fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", privateKeyAsString)

	// Parse the PEM encoded private key
	block, _ := pem.Decode([]byte(privateKeyAsPem))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	switch pk := privateKey.(type) {
	case *ecdsa.PrivateKey:
		// fmt.Println("private key is ecdsa", pk.Curve.Params())

		// Sign the payload buffer with your private key
		h := sha256.New()
		h.Write(canonicalPayload)
		hash := h.Sum(nil)
		signature, err := ecdsa.SignASN1(rand.Reader, pk, hash[:])
		if err != nil {
			return "", err
		}

		// Combine r and s into a single byte slice
		// println(len(r.Bytes()), len(s.Bytes()))
		// println(len(signature))
		// signature := append(r.Bytes(), s.Bytes()...)
		// Serialize the signature to a base64 string
		return base64.StdEncoding.EncodeToString(signature), nil

	default:
		panic(fmt.Sprintf("privy signature: unsupport private key format: %T", pk))
	}
}
