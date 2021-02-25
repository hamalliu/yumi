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

// Deploy 程序配置
type Deploy struct {
	SysName string    // 系统名称
	Env     DeployEnv // 部署环境
	Region  string    //服务器所在地区
	Zone    string    //服务器所在分区
	Version string    //版本号
}
