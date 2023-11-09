package cipher

import (
	"crypto/aes"
)

// AesDecryptECB AES_ECB解密
//
// @Description 加密可以调用 cipher.EncryptAES(key , plainText) 函数
//
// @Param []byte key "对称秘钥",必须要16/32位
//
// @Param []byte encrypted "需要解密的字符串"
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
