package accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccount(t *testing.T) {
	type args struct {
		number int
		seed   string
	}

	type account struct {
		PriKey  string `json:"private_key"`
		PubKey  string `json:"public_key"`
		Address string `json:"address"`
	}

	tests := []struct {
		name string
		args []args
		want []account
	}{
		{
			name: "The privatekey of account 1,2 is same if seed like previous input after many running time",
			args: []args{
				{
					number: 2,
					seed:   "test_1",
				},
				{
					number: 3,
					seed:   "test_1",
				},
			},
			want: []account{
				{
					PriKey:  "c5ef913e7c1ec3e0cee5d56aba9cbd12be7bdd1a1b3364047a84a2477ba3a569",
					PubKey:  "0480c430dfc951e020341fd9ef07bff2c897ead65682ba49c54733f639e151c3d6452d071a004001f8e173e40d04ca73d78068291c815c1c7bad44509fcdf3af7e",
					Address: "0x1289709BFaE305Fb7BE040b710A97c97672068bE",
				},
				{
					PriKey:  "04f795d32a98316d9fb417f7ab0f71cabb0382956d2a6135396a8aab77c25a20",
					PubKey:  "045374464df78a486f13c4cbd8dbda069f11ff899371e7bef1bdbbb330da6635315f274fdd9ba09be7db92abd95da7b28028a48cece98e5265680f15d23f1d2a24",
					Address: "0xfD568b69259Ed4b124A3F176c94F4A11Ef8b4d49",
				},
				{
					PriKey:  "34abc1637903c9b304682fb115ffb9618721c8a1d66d2782e8e94da0951cd6c0",
					PubKey:  "04facf7444da89287866f6f641342900bf515a1eb04175e0217a9d588a77788319d135e2df73a8bd6ca2fa2758de06ccb98358ed7fc1ed122663036023fb10aafd",
					Address: "0xD56Cae1315B31a5e7A5E6e1274B8960372E80220",
				},
			},
		},
		{
			name: "The privatekey of account 1,2 is different from the prevous test if seed is different",
			args: []args{
				{
					number: 2,
					seed:   "test_change",
				},
			},
			want: []account{
				{
					PriKey:  "7a7c049240033b1d4ed4d03877e9ab33d3c43d22cffda5132f476fda48529721",
					PubKey:  "04194419aa99c81f30d09c6c50519d258377a26cfa61222fa4ebc596919138dcd744d913978e4ee50cc3c4ad28dd2fa1d78d7ed4066dac2d231494a771f96a8f2b",
					Address: "0x0AB5bCa03792c104C3B2D5D9EbEF48A1234f30A9",
				},
				{
					PriKey:  "78d454164dcc867b47f762f826308c70f48420572aceb430cdd88a1657a7d33d",
					PubKey:  "040885baea472b3e24231d4c992384ae54e122e0bc15a0fa095767200d07409371f1c263ebbae4804fd306c84db4e4bf5580223e324ada0f8bea7478a98f65cb2c",
					Address: "0x1638084DAfC21cA4e8484e528F33C24611321d2A",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, arg := range tt.args {
				accs, err := GenerateAccounts(arg.number, arg.seed)
				for i, acc := range accs {
					assert.NoError(t, err)
					assert.Equal(t, tt.want[i].PriKey, acc.PrivateKeyStr())
					assert.Equal(t, tt.want[i].PubKey, acc.PublicKeyStr())
					assert.Equal(t, tt.want[i].Address, acc.Address.Hex())
				}
			}
		})
	}
}
