package model

const (
	UserClubStatusApply = 1
	UserClubStatusAgree = 2
)

type DBOpt int32

const (
	Add DBOpt = iota + 1
	Del
	Update
)
