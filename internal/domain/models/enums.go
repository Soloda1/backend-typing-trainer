package models

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type StatisticStatus string

const (
	StatisticStatusSuccess StatisticStatus = "success"
	StatisticStatusFailed  StatisticStatus = "failed"
	StatisticStatusAborted StatisticStatus = "aborted"
)
