package models

type CreateGoalRequest struct {
	Method      string `json:"method"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TargetAmount string `json:"target_amount"`
	Deadline    string `json:"deadline"`
	SessionKey  string `json:"session_key"`
	Datetime    string `json:"datetime"`
	Checksum    string `json:"checksum"`
}

type GetGoalListRequest struct {
	Method     string `json:"method"`
	Status     string `json:"status"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetGoalDetailRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type UpdateGoalRequest struct {
	Method       string `json:"method"`
	GoalID       string `json:"goal_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	TargetAmount string `json:"target_amount"`
	Deadline     string `json:"deadline"`
	SessionKey   string `json:"session_key"`
	Datetime     string `json:"datetime"`
	Checksum     string `json:"checksum"`
}

type DeleteGoalRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type AddSavingRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	Amount     string `json:"amount"`
	Note       string `json:"note"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type WithdrawSavingRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	Amount     string `json:"amount"`
	Note       string `json:"note"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetSavingHistoryRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	Page       string `json:"page"`
	Limit      string `json:"limit"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetGoalProgressRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetFinanceSummaryRequest struct {
	Method     string `json:"method"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetMilestoneListRequest struct {
	Method     string `json:"method"`
	GoalID     string `json:"goal_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}
