package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ed25519"
)

type Account struct {
	PriKey  *ecdsa.PrivateKey
	PubKey  *ecdsa.PublicKey
	Address common.Address
}

func (a *Account) PrivateKey() string {
	return hex.EncodeToString(crypto.FromECDSA(a.PriKey))
}

func (a *Account) PublicKey() string {
	return hex.EncodeToString(crypto.FromECDSAPub(a.PubKey))
}

func GenerateAccount(seed string) (*Account, error) {
	seedBytes := []byte(seed)
	seedBytes = append(seedBytes, bytes.Repeat([]byte{0x00}, ed25519.SeedSize-len(seedBytes))...)

	key := ed25519.NewKeyFromSeed(seedBytes)[32:]
	privateKey, err := crypto.ToECDSA(key[:])
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return &Account{
		PriKey:  privateKey,
		PubKey:  publicKeyECDSA,
		Address: crypto.PubkeyToAddress(privateKey.PublicKey),
	}, nil
}
