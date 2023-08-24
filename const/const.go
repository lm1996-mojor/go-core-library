package _const

const (
	CONFIG           = "application.yaml"      // 系统配置文件名
	ClientID         = "clientId"              // 租户id
	ClientTx         = "client_db_transaction" // 租户带事务的数据源
	MasterTx         = "master_db_transaction" // 平台库带事务的数据源
	CustomTx         = "custom_db_transaction" // 自定义带事务数据库标识
	JwtData          = "jwt_data"              // 原token信息（解析后的）
	UserId           = "userId"                // 用户id
	TokenType        = "act"                   // token类型(自定义)
	TokenSignature   = "link_ease_platform"    // token令牌签名
	TokenName        = "Authorization"         // 令牌Header存放key名称
	UserCode         = "user_code"             // 用户唯一编码获取key名称
	ClientCode       = "client_code"           // 租户唯一编码获取key名称
	DbLinkEncryptKey = "201dd1f39f184638"      // 数据库链接加密key
)
