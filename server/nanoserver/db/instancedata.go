package db

import (
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
)

// AddInstanceData adds a new instance data record
func AddInstanceData(data *model.InstanceData) error {
	if data.Uid <= 0 {
		return errutil.ErrInvalidParameter
	}

	_, err := database.Insert(data)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// GetInstanceData retrieves instance data by Uid
func QueryInstanceData(uid int64) (*model.InstanceData, error) {
	if uid <= 0 {
		return nil, errutil.ErrInvalidParameter
	}

	var data model.InstanceData
	has, err := database.Where("uid = ?", uid).Get(&data)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	if !has {
		return nil, errutil.ErrUserNotFound
	}
	return &data, nil
}

// UpdateInstanceData updates an existing instance data record
func UpdateInstanceData(data *model.InstanceData) error {
	if data.Uid <= 0 {
		return errutil.ErrInvalidParameter
	}

	_, err := database.Where("uid = ?", data.Uid).Cols("cur_battle_stage, reward_stage," +
		" battle_cur_stage, battle_reward_stage, boss_cur_stage, boss_challenge_time, trial_tower_cur_stage," +
		" relic_cur_stage, relic_challenge_time, relic_evaluation, relic_chapter_reward,fight_end_cur_stage,fight_end_evaluation,fight_end_refresh_time," +
		" element_cur_stage, element_challenge_time, element_need_sweep,main_battle_id,challenge_record,kill_mon_control,pushing_wall_count,pushing_wall_ticket_purchase_count").Update(data)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// DeleteInstanceData deletes an instance data record by Uid
func DeleteInstanceData(uid int64) error {
	if uid <= 0 {
		return errutil.ErrInvalidParameter
	}

	_, err := database.Where("uid = ?", uid).Delete(&model.InstanceData{})
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}
