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

// HostEnv 程序配置
type HostEnv struct {
	SysName     string    // 系统名称
	Environment DeployEnv // 部署环境
	Region      string    //服务器所在地区
	Zone        string    //服务器所在分区
	Version     string    //版本号
}
