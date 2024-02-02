package _const

const (
	CONFIG                                   = "application.yaml"   // 系统配置文件名
	TokenType                                = "act"                // token类型(自定义)
	TokenSignature                           = "link_ease_platform" // token令牌签名
	DbLinkEncryptKey                         = "201dd1f39f184638"   // 数据库链接加密key
	TokenName                                = "Authorization"      // 令牌Header存放key名称
	ConsulEndId                              = "consulServiceId"
	ClientID                                 = "_clientId"               // 租户id
	ClientTx                                 = "_client_db_transaction_" // 租户带事务的数据源
	MasterTx                                 = "_master_db_transaction_" // 平台库带事务的数据源
	CustomTx                                 = "_custom_db_transaction_" // 自定义带事务数据库标识
	JwtData                                  = "_jwt_data"               // 解析后的token信息
	UserId                                   = "_userId"                 // 用户id
	UserCode                                 = "_user_code"              // 用户唯一编码获取key名称
	ClientCode                               = "_client_code"            // 租户唯一编码获取key名称
	TokenOriginal                            = "_token_original"         // token原文
	WebSocketTokenStoreHttpRequestHeaderName = "Sec-Websocket-Protocol"  // webSocket存储token请求头名称
)
