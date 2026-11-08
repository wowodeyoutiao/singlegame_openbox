package db

// import (
// 	"nanoserver/db/model"
// 	"nanoserver/pkg/errutil"

// 	log "github.com/sirupsen/logrus"
// )

// func InsertRecoder(r *model.ItemControlRefresh) error {
// 	if r == nil {
// 		return errutil.ErrInvalidParameter
// 	}
// 	_, err := database.Insert(r)
// 	if err != nil {
// 		return errutil.ErrDBOperation
// 	}
// 	return nil
// }

// func QueryRecoder(id int32) (*model.ItemControlRefresh, error) {
// 	r := &model.ItemControlRefresh{Uid: id}
// 	has, err := database.Get(r)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	}
// 	if !has {
// 		return nil, errutil.ErrOrderNotFound
// 	}
// 	return r, nil
// }

// func UpdateRecoder(r *model.ItemControlRefresh) error {
// 	if r == nil {
// 		return errutil.ErrInvalidParameter
// 	}
// 	_, err := database.Where("uid=?", r.Uid).Cols("recoder_time").Update(r)
// 	if err != nil {
// 		return errutil.ErrDBOperation
// 	}
// 	return nil
// }

// // item_control_last
// func QueryControlLast(d *model.ItemControlLast) (*model.ItemControlLast, error) {

// 	has, err := database.Get(d)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	}
// 	if !has {
// 		return nil, errutil.ErrOrderNotFound
// 	}
// 	return d, nil
// }

// func UpdateControlLast(r *model.ItemControlLast) error {
// 	if r == nil {
// 		return errutil.ErrInvalidParameter
// 	}
// 	_, err := database.Where("uid=? and item_cfg_id=?", r.Uid, r.ItemCfgId).Cols("item_num").Update(r)
// 	if err != nil {
// 		return errutil.ErrDBOperation
// 	}
// 	return nil
// }

// func InsertControlLast(r *model.ItemControlLast) error {
// 	if r == nil {
// 		return errutil.ErrInvalidParameter
// 	}
// 	_, err := database.Insert(r)
// 	if err != nil {
// 		return errutil.ErrDBOperation
// 	}
// 	return nil
// }

// func QueryAllControlLast() ([]model.ItemControlLast, error) {
// 	d := []model.ItemControlLast{}
// 	err := database.Find(&d)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	}

// 	return d, nil
// }
