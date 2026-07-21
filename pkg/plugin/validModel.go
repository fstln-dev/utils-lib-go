package plugin

import (
	"context"
	"time"

	"github.com/fstln-dev/utils-lib-go/pkg/jwt"
	"github.com/go-kratos/kratos/v2/transport"
)

type JWTTenantModel struct {
	Iss    string    `json:"iss"` // 发行者
	Iat    time.Time `json:"iat"` // 签发时间
	Exp    time.Time `json:"exp"` // 失效时间
	Sub    int       `json:"sub"` // 主体，TenantId
	Aud    string    `json:"aud"` // 受众, 账户体系
	Scopes []string  `json:"sco"` // 权限
}

func GetJWTInfo(ctx context.Context) *JWTTenantModel {

	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil
	}
	jwtStr := tr.RequestHeader().Get("x-fstln-token")

	json := &JWTTenantModel{}
	err := jwt.JWTDecode(jwtStr, json)
	if err != nil {
		return nil
	}

	return json
}
