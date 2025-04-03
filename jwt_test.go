package privy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAppID           = "xxx"
	testSecretKey       = "xxxx" //nolint
	testVerificationKey = `-----BEGIN PUBLIC KEY-----
xxxxxxxx==
-----END PUBLIC KEY-----`
	testAuthKey = `wallet-auth:xxxxxx+xxxx+`
)

func TestJWT(t *testing.T) {
	SkipCI(t)

	ak := "xxxxx.xxxxx.xxxxx"
	pc := NewPrivyClient(testAppID, testSecretKey, testVerificationKey)
	claim, err := pc.ValidAccessToken(ak)
	assert.Nil(t, err)
	fmt.Println(claim)
}

func TestGetUser(t *testing.T) {
	SkipCI(t)

	did := "did:privy:xxxxxx"
	pc := NewPrivyClient(testAppID, testSecretKey, testVerificationKey)
	_, err := pc.GetUser(did)
	assert.Nil(t, err)
}
