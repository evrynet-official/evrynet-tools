package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strconv"

	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ed25519"
)

func GenerateAccounts(num int, seed string) ([]*Account, error) {
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

func GenerateAccountsFromContext(ctx *cli.Context) ([]*Account, error) {
	var (
		num  = ctx.Int(NumAccountsFlag.Name)
		seed = ctx.String(SeedFlag.Name)
	)

	// generate accounts
	accs, err := GenerateAccounts(num, seed)
	if err != nil {
		fmt.Println("Fail to generate new account!", "Err:", err)
		return nil, err
	}
	return accs, nil

}
