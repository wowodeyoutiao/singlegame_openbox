package web

import (
	"context"
	"nanoserver/dispatcher"
	"nanoserver/internal/web/api"
)

func CreateCharacterHandler(eventData *dispatcher.CreateCharacter) {
	paymentClinet := api.NewPaymentClient()
	resp, err := paymentClinet.CreatePayCharacterRequest(context.TODO(), eventData)
	if err != nil {
		logger.Errorf("CreateCharacterHandler: 创建支付角色失败: %v", err)
		return
	}
	logger.Infof("CreateCharacterHandler: 创建支付角色成功: %+v", resp)
}
