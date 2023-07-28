package cipher

import (
	"crypto/aes"
	"encoding/hex"
)

// EncryptAES AES加密
//
// @Description 解密可以调用 cipher.DecryptAES(key ,encryptText) 函数
//
// @Param key "对称秘钥"
//
// @Param plainText "需要加密的字符串内容"
func EncryptAES(key string, plainText string) (string, error) {
	cipher, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	out := make([]byte, len(plainText))

	cipher.Encrypt(out, []byte(plainText))

	return hex.EncodeToString(out), nil
}
