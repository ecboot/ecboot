package config

type Config struct {
	App     App
	Jwt     Jwt
	Wechat  []WechatAccount
	Payment Payment
}

type App struct {
	Name  string `json:"name"`
	Debug bool   `json:"debug"`
}

type Jwt struct {
	Secret             string `json:"secret"`
	AccessTokenExpire  string `json:"accessTokenExpire"`
	RefreshTokenExpire string `json:"refreshTokenExpire"`
}

type WechatAccount struct {
	AppID          string `json:"appid"`
	Secret         string `json:"secret"`
	AccessTokenURL string `json:"accessTokenUrl"`
	UserInfoURL    string `json:"userInfoUrl"`
}

type WechatConfig struct {
	Accounts []WechatAccount `json:"accounts"`
}

type Payment struct {
	Wechat WechatPayment `json:"wechat"`
}

type WechatPayment struct {
	MchId      string `json:"mchId"`
	AppId      string `json:"appId"`
	ApiV3Key   string `json:"apiV3Key"`
	SerialNo   string `json:"serialNo"`
	PrivateKey string `json:"privateKey"`
	// PublicKeyId 微信支付公钥 ID（商户平台-API安全申请，形如 PUB_KEY_ID_...）。
	// PublicKey 微信支付公钥 PEM 文件路径。新商户不再签发平台证书，自动下载证书会报"无可用的平台证书"，必须改用公钥模式。
	PublicKeyId string `json:"publicKeyId"`
	PublicKey   string `json:"publicKey"`
	NotifyUrl   string `json:"notifyUrl"`
}
