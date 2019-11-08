package deployment

/**
Deployment represents a collection of evrynet deploy (can be mainnet, testnet etc...)
There might be multiple deployments in same network used for different purpose.
Deployment is a separated concept from running mode to allow developers to run any deployment in debug mode.
*/

//Deployment is a enum type for checking valid DeploymentMode
//go:generate stringer -type=Deployment -linecomment
type Deployment int

const (
	//Production is production mode for deployment
	Mainnet Deployment = iota //production
	//Staging is staging mode for deployment
	TestNet //staging
)

