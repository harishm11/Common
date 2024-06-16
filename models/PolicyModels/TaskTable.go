package models

type TaskTable struct {
	Field     string
	Indicator string
	Task      string
}

func (TaskTable) TableName() string {
	return "policyprocessor_task_tables"
}
