package activity

import (
	"time"
)

type ActivityRecord struct {
	ID        uint64         `json:"id"`
	Status    ActivityStatus `json:"status"`
	StartTime time.Time      `json:"start"`
	EndTime   time.Time      `json:"end"`
}

type ActionRecord struct {
	ID             uint64       `json:"id"`
	ActivityID     uint64       `json:"activityID"`
	Status         ActionStatus `json:"status"`
	StartTime      time.Time    `json:"start"`
	EndTime        time.Time    `json:"end"`
	DoFuncID       string       `json:"doFuncID"`
	DoParams       string       `json:"doParams"`
	RollbackFuncID string       `json:"rollbackFuncID"`
	RollbackParams string       `json:"rollbackParams"`
}
