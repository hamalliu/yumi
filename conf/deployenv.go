package conf

// DeployEnv ...
type DeployEnv string

// deploy env.
const (
	DeployEnvDev  DeployEnv = "dev"
	DeployEnvFat  DeployEnv = "fat"
	DeployEnvUat  DeployEnv = "uat"
	DeployEnvPre  DeployEnv = "pre"
	DeployEnvProd DeployEnv = "prod"
)
