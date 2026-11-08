package utils

// GetChapterID 从战斗ID获取章节ID, 格式: 1701010, 17(副本ID) 01(章节ID) 010(关卡ID)
func GetChapterID(battleID int32) int32 {
	return (battleID / 1000) % 100
}
