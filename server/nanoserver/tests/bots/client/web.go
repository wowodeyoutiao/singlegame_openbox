package client

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"nanoserver/protocol"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type LoginResponse struct {
	Token string `json:"token"`
	Uid   int64  `json:"uid"`
}

func (c *Client) Login(name string) (resp *LoginResponse, err error) {
	// 登录
	loginResp, err := c.guestLogin(&protocol.LoginRequest{
		AppID:     "wxChat",
		ChannelID: "wx",
		IMEI:      name,
		ServerId:  viper.GetInt32("server.id"),
	})
	if err != nil {
		return nil, err
	}
	// log.Info().Msgf("客户端 %s 登录成功, Token: %s", name, loginResp.Token)
	log.Info().Msgf("客户端 %s 登录成功", name)

	// 创建角色
	var character *protocol.CharacterInfo
	if len(loginResp.CharacterList) == 0 {
		createReq := &protocol.CreateCharacterRequest{
			// Account: name,
			Sex:  2,
			Name: name,
			Job:  2,
		}
		createCharacterResp, err := c.createCharacter(createReq, loginResp.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "create character failed for account: %s, serverId: %d", name, viper.GetInt32("server.id"))
		}
		character = createCharacterResp.CharacterList
	} else {
		if len(loginResp.CharacterList) > 1 {
			log.Warn().Msgf("登录时发现多个角色 %#+v", loginResp.CharacterList)
		}
		character = loginResp.CharacterList[0]
	}
	log.Info().Msgf("客户端 %s 选择角色: %s, UID: %v", name, character.Name, character.Uid)

	// 解析JWT tokens
	// loginClaims, err := parseJWTPayload(loginResp.Token)
	// if err != nil {
	// 	log.Warn().Err(err).Msg("failed to parse login token")
	// }

	// log.Info().Msgf("客户端 %s 登录成功, UID: %v, login_token: %+#v, character: %+#v", name, character.Uid, loginClaims, character)

	// 选择角色
	selResp, err := c.selectCharacter(&protocol.SelectCharacterRequest{Uid: character.Uid}, loginResp.Token)
	if err != nil {
		return nil, err
	}

	// selectClaims, err := parseJWTPayload(selResp.Token)
	// if err != nil {
	// 	log.Warn().Err(err).Msg("failed to parse select token")
	// }
	// log.Info().Msgf("客户端 %s 选择角色成功, UID: %v, select_token: %+#v", name, character.Uid, selectClaims)

	return &LoginResponse{
		Token: selResp.Token,
		Uid:   character.Uid,
	}, nil
}

// guestLogin 游客登录
func (c *Client) guestLogin(req *protocol.LoginRequest) (*protocol.LoginResponse, error) {
	webURL := viper.GetString("server.web_url")
	if webURL == "" {
		return nil, errors.Errorf("server web_url not configured")
	}

	loginResp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(webURL + "/v1/user/login/guest")
	if err != nil {
		return nil, errors.Wrapf(err, "guest login request failed: %s", webURL+"/v1/user/login/guest")
	}

	if loginResp.IsError() {
		return nil, errors.Errorf("guest login request failed with status %s: %s", loginResp.Status(), string(loginResp.Body()))
	}

	var resp protocol.LoginResponse
	if err := json.Unmarshal(loginResp.Body(), &resp); err != nil {
		return nil, errors.Errorf("failed to unmarshal guest login response: %w, body: %s", err, string(loginResp.Body()))
	}

	if resp.Code != 0 {
		return nil, errors.Errorf("guest login failed with code: %d, message: %v", resp.Code, resp.Messages)
	}

	return &resp, nil
}

// createCharacter 创建角色
func (c *Client) createCharacter(req *protocol.CreateCharacterRequest, token string) (*protocol.CreateCharacterResponse, error) {
	webURL := viper.GetString("server.web_url")
	if webURL == "" {
		return nil, errors.Errorf("server web_url not configured")
	}

	createReq := &protocol.CreateCharacterRequest{
		Account: req.Account,
		Sex:     req.Sex,
		Name:    req.Name,
		Job:     req.Job,
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(createReq).
		Put(webURL + "/v1/user/login/create_char")
	if err != nil {
		return nil, errors.Wrapf(err, "create character request failed: %s", webURL+"/v1/user/login/create_char")
	}

	if resp.IsError() {
		return nil, errors.Errorf("create character request failed with status %s: %s", resp.Status(), string(resp.Body()))
	}

	var createResp protocol.CreateCharacterResponse
	if err := json.Unmarshal(resp.Body(), &createResp); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal create character response: %s", string(resp.Body()))
	}

	if createResp.Code != 0 {
		return nil, errors.Errorf("create character failed with code: %d, message: %s", createResp.Code, createResp.Date)
	}

	return &createResp, nil
}

// selectCharacter 选择角色
func (c *Client) selectCharacter(req *protocol.SelectCharacterRequest, token string) (*protocol.SelectCharacterResponse, error) {
	webURL := viper.GetString("server.web_url")
	if webURL == "" {
		return nil, errors.Errorf("server web_url not configured")
	}

	selectReq := &protocol.SelectCharacterRequest{
		Uid: req.Uid,
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(selectReq).
		Post(webURL + "/v1/user/login/select_char")
	if err != nil {
		return nil, errors.Wrapf(err, "select character request failed: %s", webURL+"/v1/user/login/select_char")
	}

	if resp.IsError() {
		return nil, errors.Errorf("select character request failed with status %s: %s", resp.Status(), string(resp.Body()))
	}

	var selectResp protocol.SelectCharacterResponse
	if err := json.Unmarshal(resp.Body(), &selectResp); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal select character response: %s", string(resp.Body()))
	}

	if selectResp.Code != 0 {
		return nil, errors.Errorf("select character failed with code: %d, message: %s", selectResp.Code, selectResp.Message)
	}

	return &selectResp, nil
}

// parseJWTPayload 解析JWT token的payload部分（未加密数据）
func parseJWTPayload(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT token format")
	}

	// 解码payload部分（第二部分）
	payload := parts[1]
	// 补齐base64填充
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode JWT payload")
	}

	var claims map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(string(decoded)))
	decoder.UseNumber() // 使用json.Number类型来保持数字的原始格式
	if err := decoder.Decode(&claims); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal JWT payload")
	}

	return claims, nil
}
