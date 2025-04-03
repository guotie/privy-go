package privy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SkipCI(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}
}

func TestCreateWallet(t *testing.T) {
	SkipCI(t)

	c := NewPrivyClient(testAppID, testSecretKey, testVerificationKey)
	res, err := c.CreateWallet("")
	assert.Nil(t, err)
	fmt.Printf("wallet id: %v wallet address: %v\n", res.ID, res.Address)
}

func TestAuthSign(t *testing.T) {
	SkipCI(t)

	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
	signature, err := c.GetAuthorizationSignature(PrivyWalletRPC, map[string]any{
		"chain_type": "ethereum",
		// "method":     "personal_sign",
		// "params": map[string]string{
		// 	"message":  "Hello, Ethereum.",
		// 	"encoding": "utf-8",
		// },
	})
	assert.Nil(t, err)
	fmt.Println("signature: " + signature)
}

func TestPost(t *testing.T) {
	SkipCI(t)

	reader := bytes.NewBuffer([]byte("{}"))
	resp, err := http.Post(PrivyWalletRPC, "application/json", reader)
	assert.Nil(t, err)

	fmt.Println(resp)
	fmt.Println(resp.Body)
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err == nil {
		fmt.Println(string(buf))
	} else {
		fmt.Println(err)
	}
	fmt.Println(err)
}

func TestSendTx(t *testing.T) {
	SkipCI(t)

	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)

	res := map[string]any{}
	err := doSignedRequest(c,
		http.MethodPost,
		PrivyWalletRPC,
		// "https://auth.privy.io/api/v1/wallets/xxxxxx/rpc",
		map[string]any{
			"method":     "eth_sendTransaction",
			"caip2":      "eip155:8453",
			"chain_type": "ethereum",
			"address":    "",
			"params": map[string]any{
				"transaction": map[string]any{
					"chain_id": 8453,
					"to":       "",
					"value":    2000000000,
				},
			},
		},
		true,
		&res)

	// map[data:map[caip2:eip155:8453 hash:] method:eth_sendTransaction]
	fmt.Println(res)
	assert.Nil(t, err)
}

func TestApproveTx(t *testing.T) {
	SkipCI(t)

	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
	res := map[string]any{}
	err := doSignedRequest(c,
		http.MethodPost,
		PrivyWalletRPC,
		map[string]any{
			"method":     "eth_sendTransaction",
			"caip2":      "eip155:8453",
			"chain_type": "ethereum",
			"address":    "",
			"params": map[string]any{
				"transaction": map[string]any{
					"chain_id": 8453,
					// "value":    0,
					"to":   "",
					"from": "",
					"data": "",
				},
			},
		},
		true,
		&res)

	// map[data:map[caip2:eip155:8453 hash:] method:eth_sendTransaction]
	fmt.Println(res)
	assert.Nil(t, err)
}

func TestSwapTx(t *testing.T) {
	SkipCI(t)

	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
	res := map[string]any{}
	err := doSignedRequest(c,
		http.MethodPost,
		PrivyWalletRPC,
		map[string]any{
			"method":     "eth_sendTransaction",
			"caip2":      "eip155:8453",
			"chain_type": "ethereum",
			"address":    "",
			"params": map[string]any{
				"transaction": map[string]any{
					"chain_id": 8453,
					"to":       "",
					"from":     "",
					"data":     "",
				},
			},
		},
		true,
		&res)

	// map[data:map[caip2:eip155:8453 hash:] method:eth_sendTransaction]
	fmt.Println(res)
	assert.Nil(t, err)
}

func TestSwapTx2(t *testing.T) {
	SkipCI(t)

	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
	hash, err := c.SendETHTransaction("8453", "", map[string]any{
		"chain_id":  8453,
		"to":        "",
		"from":      "",
		"data":      "",
		"value":     "", //
		"gas_limit": 2000000,
	})
	assert.Nil(t, err)
	fmt.Println("txhash:", hash)
}
