package accounts

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/crypto"
)

type Account struct {
	PriKey  *ecdsa.PrivateKey
	PubKey  *ecdsa.PublicKey
	Address common.Address
}

func (a *Account) PrivateKeyStr() string {
	return hex.EncodeToString(crypto.FromECDSA(a.PriKey))
}

func (a *Account) PublicKeyStr() string {
	return hex.EncodeToString(crypto.FromECDSAPub(a.PubKey))
}
