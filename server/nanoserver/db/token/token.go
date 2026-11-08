package token

import (
	"errors"
	"log"
	"time"

	"nanoserver/pkg/constant"
	"nanoserver/pkg/my_jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// CreateUserFactory 创建 userToken 工厂
func CreateUserFactory() *userToken {
	return &userToken{
		userJwt: my_jwt.CreateMyJWT(viper.GetString("token.jwtTokenSignKey")),
	}
}

type userToken struct {
	userJwt *my_jwt.JwtSign
}

// GenerateToken 生成token
func (u *userToken) GenerateToken(userid int64, accountName, platform string, serverId int32, expireAt int64) (tokens string, err error) {
	// 根据实际业务自定义token需要包含的参数，生成token，注意：用户密码请勿包含在token
	customClaims := my_jwt.CustomClaims{
		UserId:      userid,
		AccountName: accountName,
		Platform:    platform,
		ServerId:    serverId,
		// 特别注意，针对前文的匿名结构体，初始化的时候必须指定键名，并且不带 jwt. 否则报错：Mixture of field: value and value initializers
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,       // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

// RecordLoginToken 用户login成功，记录用户token
func (u *userToken) RecordLoginToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken); err == nil {
		userId := customClaims.UserId

		ValidTokenCacheToRedis(userId, userToken, customClaims.ExpiresAt)

		return true
	} else {
		return false
	}
}

// 判断token本身是否未过期
// 参数解释：
// token： 待处理的token值
// expireAtSec： 过期时间延长的秒数，主要用于用户刷新token时，判断是否在延长的时间范围内，非刷新逻辑默认为0
func (u *userToken) isNotExpired(token string, expireAtSec int64) (*my_jwt.CustomClaims, int) {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {
		if time.Now().Unix()-(customClaims.ExpiresAt+expireAtSec) < 0 {
			// token有效
			return customClaims, constant.JwtTokenOK
		} else {
			// 过期的token
			return customClaims, constant.JwtTokenExpired
		}
	} else {
		// 无效的token
		return nil, constant.JwtTokenInvalid
	}
}

// IsEffective 判断token是否有效（未过期+数据库用户信息正常）
func (u *userToken) IsEffective(token string) bool {
	customClaims, code := u.isNotExpired(token, 0)
	if constant.JwtTokenOK == code {
		tokenRedisFact := CreateUsersTokenCacheFactory(customClaims.UserId)
		if tokenRedisFact != nil {
			defer tokenRedisFact.ReleaseRedisConn()
			if tokenRedisFact.TokenCacheIsExists(token) {
				return true
			} else {
				log.Printf("token 不存在 uid: %d  token : %s", customClaims.UserId, token)
			}
		}
	} else {
		log.Printf("token 已过期   token: %s", token)
	}

	return false
}

// ParseToken 将 token 解析为绑定时传递的参数
func (u *userToken) ParseToken(tokenStr string) (CustomClaims my_jwt.CustomClaims, err error) {
	if customClaims, err := u.userJwt.ParseToken(tokenStr); err == nil {
		return *customClaims, nil
	} else {
		return my_jwt.CustomClaims{}, errors.New("token解析失败")
	}
}

// DestroyToken 销毁token，基本用不到，因为一个网站的用户退出都是直接关闭浏览器窗口，极少有户会点击“注销、退出”等按钮，销毁token其实无多大意义
func (u *userToken) DestroyToken() {}
