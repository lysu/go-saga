package activity

type ActivityStatus uint8

const (
	ActivityStarted ActivityStatus = iota + 1
	ActivitySuccess
	ActivityFailure
)

type ActionStatus uint8

const (
	ActionStarted ActionStatus = iota + 1
	ActionSuccess
	ActionFailure
)