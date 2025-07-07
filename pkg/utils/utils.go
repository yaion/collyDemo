package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

func Decrypt(urlStr, text string) (string, error) {
	// 1. 参数检查
	if urlStr == "" || text == "" {
		return "", fmt.Errorf("URL and text must not be empty")
	}

	// 2. URL编码
	str := getStr(urlStr)

	// 5. 提取密钥和IV
	orgKey := str[:16]
	orgIV := str[12:28]

	// 6. 解码Base64密文
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %v", err)
	}

	// 7. 创建AES解密器
	block, err := aes.NewCipher([]byte(orgKey))
	if err != nil {
		return "", fmt.Errorf("cipher creation error: %v", err)
	}

	// 8. 检查IV长度
	if len(orgIV) != block.BlockSize() {
		return "", fmt.Errorf("IV length must equal block size")
	}

	// 9. 创建CBC模式解密器
	mode := cipher.NewCBCDecrypter(block, []byte(orgIV))

	// 10. 解密数据
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// 11. 移除PKCS7填充
	unpadded, err := unpadPKCS7(decrypted)
	if err != nil {
		return "", fmt.Errorf("padding removal error: %v", err)
	}

	// 12. 返回解密后的字符串
	return string(unpadded), nil
}

// 移除PKCS7填充
func unpadPKCS7(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("input data is empty")
	}

	padding := int(data[len(data)-1])
	if padding < 1 || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding size")
	}

	if len(data) < padding {
		return nil, fmt.Errorf("data length is less than padding size")
	}

	// 检查填充字节是否有效
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}

	return data[:len(data)-padding], nil
}

func getStr(url string) string {
	// 1. 执行encodeURI
	encoded := encodeURI(url)

	// 2. 进行Base64编码
	base64Str := base64.StdEncoding.EncodeToString([]byte(encoded))

	// 3. 重复3次并返回
	return strings.Repeat(base64Str, 3)
}

func encodeURI(s string) string {
	safe := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.!~*'();,/?:@&=+$#"
	var result strings.Builder
	for _, r := range s {
		if strings.ContainsRune(safe, r) {
			result.WriteRune(r)
		} else {
			// UTF-8编码非安全字符
			for _, b := range []byte(string(r)) {
				result.WriteString(fmt.Sprintf("%%%02X", b))
			}
		}
	}
	return result.String()
}
