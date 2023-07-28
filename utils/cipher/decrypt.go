package cipher

import (
	"crypto/aes"
	"encoding/hex"
)

// DecryptAES AES解密
//
// @Description 加密可以调用 cipher.EncryptAES(key , plainText) 函数
//
// @Param key "对称秘钥"
//
// @Param plainText "需要解密的字符串"
func DecryptAES(key string, encryptText string) (string, error) {
	decodeText, _ := hex.DecodeString(encryptText)

	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	out := make([]byte, len(decodeText))
	cipher.Decrypt(out, decodeText)

	return string(out[:]), nil
}
