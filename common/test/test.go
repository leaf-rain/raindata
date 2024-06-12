package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	key := []byte("gDyVgzwa0mFz9uUP7M6GQQ==")
	data := []byte("1559129713_67297598215591297148AP6DT9ybtniUJfbwx20afc706a711eefc")

	h := hmac.New(sha256.New, key)
	h.Write(data)
	hash := h.Sum(nil)

	// 将hash结果转换为16进制字符串
	hashHex := hex.EncodeToString(hash)
	fmt.Println("HMAC-SHA256 (Hex):", hashHex)
}
