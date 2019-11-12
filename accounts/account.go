package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"strconv"

	"golang.org/x/crypto/ed25519"

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

func GenerateAccount(num int, seed string) ([]*Account, error) {
	var accs []*Account
	for i := 0; i < num; i++ {
		seedBytes := []byte(seed + strconv.Itoa(i))
		seedBytes = append(seedBytes, bytes.Repeat([]byte{0x00}, ed25519.SeedSize-len(seedBytes))...)

		key := ed25519.NewKeyFromSeed(seedBytes)[32:]
		privateKey, err := crypto.ToECDSA(key[:])
		if err != nil {
			return nil, err
		}

		publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		accs = append(accs,
			&Account{
				PriKey:  privateKey,
				PubKey:  publicKeyECDSA,
				Address: crypto.PubkeyToAddress(privateKey.PublicKey),
			})
	}
	return accs, nil
}
