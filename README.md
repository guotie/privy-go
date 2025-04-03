# privy-go

## Install
```bash
go get github.com/guotie/privy-go
```

## Create client

```go  
client := privy.NewPrivyClient(apiID, secretKey, verificationKey)
// or 
client := privy.NewPrivyClientWithAuthkey(apiID, secretKey, verificationKey, authKey)
```

`apiID`, `secretKey`, `verificationKey`, `authKey` should be provided by Privy.

## Create wallet
```go
	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
    wallet, err := c.CreateWallet("")
```

## SignTypedDataV4
```go 
	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
	signature, err := c.SignTypedDataV4(
		walletID,
		aaddress,
		map[string]interface{}{
			"types":        types,
			"primary_type": primaryType,
			"domain":       domain,
			"message":      message,
		},
	)
```

## SendETHTransaction

```go
	c := NewPrivyClientWithAuthkey(testAppID, testSecretKey, testVerificationKey, testAuthKey)
	hash, err := c.SendETHTransaction("8453", "the privy wallet address", map[string]any{
		"chain_id":  8453,
		"to":        "your to address",
		"from":      "your from address",
		"data":      "0x" + "your tx data",
		"value":     "0x" + "your value", //
		"gas_limit": 2000000,
	})
```

