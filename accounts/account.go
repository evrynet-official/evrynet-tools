package accounts

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/crypto"
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
