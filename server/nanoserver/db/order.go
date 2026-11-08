package db

// import (
// 	"fmt"
// 	"strings"

// 	"nanoserver/db/model"
// 	"nanoserver/pkg/algoutil"
// 	"nanoserver/pkg/errutil"
// )

// const (
// 	noLimitFlag  = -1 // 如果count == -1则表示返回所有数据
// 	noTimeFilter = -1 // 如果start/end == -1则表示无时间筛选
// )

// func QueryOrder(orderID string) (*model.Order, error) {
// 	order := &model.Order{OrderId: orderID}
// 	has, err := database.Get(order)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !has {
// 		return nil, errutil.ErrOrderNotFound
// 	}
// 	return order, nil
// }

// func QueryOrdersByUserID(userID int64) ([]*model.Order, error) {
// 	if userID <= 0 {
// 		return nil, fmt.Errorf("invalid userID: %d", userID)
// 	}

// 	order := &model.Order{Uid: userID, Status: OrderStatusNotified}
// 	orders := make([]*model.Order, 0)
// 	// err := database.Where("uid = ? AND status IN (?)", userID, statuses).Find(&orders)
// 	err := database.Find(&orders, order)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query orders for userID %d: %w", userID, err)
// 	}

// 	return orders, nil
// }

// // 通过订单号和Nonce查找订单
// func QueryOrderByOpenIDAndNonce(openID, nonce string) (*model.Order, error) {
// 	if openID == "" || nonce == "" {
// 		return nil, fmt.Errorf("invalid parameters: nonce=%s, openid=%s", nonce, openID)
// 	}

// 	order := &model.Order{
// 		OpenID: openID,
// 		Nonce:  nonce,
// 	}
// 	has, err := database.Get(order)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query order: %w", err)
// 	}
// 	if !has {
// 		return nil, fmt.Errorf("order not found with openID=%s and nonce=%s", openID, nonce)
// 	}
// 	return order, nil
// }

// // 查询并且需改Nonce
// func QueryAndUpdateOrderNonce(orderID, nonce string) (*model.Order, error) {
// 	if orderID == "" || nonce == "" {
// 		return nil, fmt.Errorf("invalid parameters: orderID=%s, nonce=%s", orderID, nonce)
// 	}

// 	session := database.NewSession()
// 	defer session.Close()

// 	if err := session.Begin(); err != nil {
// 		return nil, fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	order := &model.Order{}
// 	has, err := session.Where("order_id = ?", orderID).ForUpdate().Get(order)
// 	if err != nil {
// 		session.Rollback()
// 		return nil, fmt.Errorf("failed to query order: %w", err)
// 	}
// 	if !has {
// 		session.Rollback()
// 		return nil, fmt.Errorf("order not found: %s", orderID)
// 	}

// 	if order.Status != OrderStatusCreated {
// 		session.Rollback()
// 		return nil, fmt.Errorf("order not in created status: %s, current status: %d", orderID, order.Status)
// 	}

// 	order.Nonce = nonce
// 	_, err = session.ID(order.Id).Update(order) // 使用主键更新，更符合 ORM 思想
// 	if err != nil {
// 		session.Rollback()
// 		return nil, fmt.Errorf("failed to update order nonce: %w", err)
// 	}

// 	if err := session.Commit(); err != nil {
// 		return nil, fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return order, nil
// }

// func InsertOrder(order *model.Order) error {
// 	if order == nil {
// 		return errutil.ErrInvalidParameter
// 	}
// 	_, err := database.Insert(order)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func UpdateOrder(order *model.Order) error {
// 	if order == nil {
// 		return errutil.ErrInvalidParameter
// 	}
// 	_, err := database.Where("order_id=?", order.OrderId).Update(order)
// 	if err != nil {
// 		return errutil.ErrDBOperation
// 	}
// 	return nil
// }

// func YXPayOrderList(uid int64, appid, channelID, orderID string, start, end int64, typ, offset, count int) ([]model.Order, int, error) {
// 	order := &model.Order{
// 		AppId:     appid,
// 		ChannelId: channelID,
// 		Uid:       uid,
// 		OrderId:   orderID,
// 		Type:      typ,
// 	}

// 	start, end = algoutil.TimeRange(start, end)

// 	// println("uid", uid, "appid", appid, "channelid", channelID, "start", start, "end", end, "offset", offset, "count", count)

// 	total, err := database.Where("created_at BETWEEN ? AND ?", start, end).Count(order)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, 0, errutil.ErrDBOperation
// 	}

// 	result := make([]model.Order, 0)
// 	if count == noLimitFlag {
// 		err = database.Where("created_at BETWEEN ? AND ?", start, end).
// 			Desc("id").Find(&result, order)
// 	} else {
// 		err = database.Where("created_at BETWEEN ? AND ?", start, end).
// 			Desc("id").Limit(count, offset).Find(&result, order)
// 	}

// 	if err != nil {
// 		logger.Error(err)
// 		return nil, 0, errutil.ErrDBOperation
// 	}

// 	return result, int(total), nil
// }

// func OrderList(uid int64, appid, channelID, orderID, payBy string, start, end int64, status, offset, count int) ([]model.Order, int, error) {
// 	order := &model.Order{
// 		AppId:       appid,
// 		ChannelId:   channelID,
// 		Uid:         uid,
// 		OrderId:     orderID,
// 		PayPlatform: payBy,
// 		Status:      status,
// 	}

// 	start, end = algoutil.TimeRange(start, end)

// 	// println("uid", uid, "appid", appid, "channelid", channelID, "start", start, "end", end, "offset", offset, "count", count)

// 	total, err := database.Where("created_at BETWEEN ? AND ?", start, end).Count(order)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, 0, errutil.ErrDBOperation
// 	}

// 	result := make([]model.Order, 0)
// 	if count == noLimitFlag {
// 		err = database.Where("created_at BETWEEN ? AND ?", start, end).
// 			Desc("id").Find(&result, order)
// 	} else {
// 		err = database.Where("created_at BETWEEN ? AND ?", start, end).
// 			Desc("id").Limit(count, offset).Find(&result, order)
// 	}

// 	if err != nil {
// 		logger.Error(err)
// 		return nil, 0, errutil.ErrDBOperation
// 	}

// 	return result, int(total), nil
// }

// func BalanceList(uids []string) (map[string]string, error) {
// 	if uids == nil {
// 		return nil, errutil.ErrIllegalParameter
// 	}

// 	sql := "SELECT  uid, coin from `user` WHERE uid IN ( " + strings.Join(uids, ",") + ")"
// 	results, err := database.Query(sql)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, errutil.ErrDBOperation
// 	}

// 	m := make(map[string]string)

// 	for _, result := range results {
// 		m[string(result["uid"])] = string(result["coin"])
// 	}
// 	return m, nil
// }
