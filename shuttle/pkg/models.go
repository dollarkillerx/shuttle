package pkg

import (
	"encoding/json"

	"google.dev/google/shuttle/utils"
)

type VerifyToken struct {
	UserJWT     string      `json:"user_jwt"`     // 用户jwt
	NodeID      string      `json:"node_id"`      // 授权node id
	Expiration  int64       `json:"expiration"`   // token 过期时间
	MountSocks5 MountSocks5 `json:"mount_socks5"` // 挂载socks5
}

type MountSocks5 struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (v *VerifyToken) ToToken(aesKey string) string {
	marshal, _ := json.Marshal(v)

	encrypt, err := utils.AESEncrypt(marshal, []byte(aesKey))
	if err != nil {
		panic(err)
	}
	return encrypt
}

func (v *VerifyToken) FromToken(token string, aesKey string) {

	decrypt, err := utils.AESDecrypt(token, []byte(aesKey))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(decrypt, v)
	if err != nil {
		panic(err)
	}
}
