package ctxm

import (
	"context"
	"github.com/2908755265/mutil/typem"
)

const (
	infoKey userInfoKey = "userInfo"
)

type userInfoKey string

func WithUserInfo(ctx context.Context, info *typem.UserInfo) context.Context {
	return context.WithValue(ctx, infoKey, info)
}

func GetUserInfo(ctx context.Context) *typem.UserInfo {
	var r *typem.UserInfo
	v := ctx.Value(infoKey)
	if v == nil {
		return r
	}
	return v.(*typem.UserInfo)
}
