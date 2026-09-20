package consts

const (
	AppName = "ecboot"
	Version = "1.0.0"
	Release = "20260920"

	// EnvJwtSecret JWT 签名密钥的环境变量名：真实密钥不入库（git），
	// 部署时经环境变量注入（gcfg 不支持 ${ENV} 展开，故在代码层做覆盖）。
	EnvJwtSecret = "JWT_SECRET"
)
