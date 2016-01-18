package saga

import (
	"time"
)

type ActivityLog struct {
	ActivityID uint64         `json:"activityID"`
	Status     ActivityStatus `json:"status"`
	StartTime  time.Time      `json:"start"`
	EndTime    time.Time      `json:"end"`
}

type ActionLog struct {
	ActionID       uint64       `json:"actionID"`
	ActivityID     uint64       `json:"activityID"`
	Status         ActionStatus `json:"status"`
	StartTime      time.Time    `json:"start"`
	EndTime        time.Time    `json:"end"`
	DoFuncID       string       `json:"doFuncID"`
	DoParams       string       `json:"doParams"`
	RollbackFuncID string       `json:"rollbackFuncID"`
	RollbackParams string       `json:"rollbackParams"`
}

type actionData struct {
	actionID uint64
	data     string
}
