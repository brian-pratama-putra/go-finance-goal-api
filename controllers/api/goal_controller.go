package controllers

import (
	"os"
	"strconv"

	dao "go-finance-goal-api/dao/api"
	"go-finance-goal-api/others"

	"github.com/gin-gonic/gin"
)

func CreateGoal(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_name          := others.GetStr(v_data, "name")
	v_description   := others.GetStr(v_data, "description")
	v_target_amount := others.GetStr(v_data, "target_amount")
	v_deadline      := others.GetStr(v_data, "deadline")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_name == "" || v_target_amount == "" || v_deadline == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_target_float, v_err := strconv.ParseFloat(v_target_amount, 64)
	if v_err != nil || v_target_float <= 0 {
		others.ResponseService(p_c, v_method, 422, "Target amount harus angka positif", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_name + "#" + v_target_amount + "#" + v_deadline + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.InsertGoal(v_user_id, v_name, v_description, v_deadline, v_target_float)

	if v_hasil["status"] == "T" {
		dao.CacheDeletePattern("goal:list:" + v_user_id + ":*")
		dao.CacheDelete("goal:summary:" + v_user_id)
		v_result := v_hasil["result"].(map[string]interface{})
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"goal_id": v_result["goal_id"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetGoalList(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_status        := others.GetStr(v_data, "status")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetGoalList(v_user_id, v_status)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetGoalDetail(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetGoalDetail(v_goal_id, v_user_id)

	if v_hasil["status_code"] == 204 {
		others.ResponseService(p_c, v_method, 204, "Data tidak ditemukan", nil)
	} else if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func UpdateGoal(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_update := map[string]interface{}{}
	if v_name := others.GetStr(v_data, "name"); v_name != "" {
		v_update["name"] = v_name
	}
	if v_desc := others.GetStr(v_data, "description"); v_desc != "" {
		v_update["description"] = v_desc
	}
	if v_ta := others.GetStr(v_data, "target_amount"); v_ta != "" {
		if v_ta_float, v_err := strconv.ParseFloat(v_ta, 64); v_err == nil && v_ta_float > 0 {
			v_update["target_amount"] = v_ta_float
		}
	}
	if v_dl := others.GetStr(v_data, "deadline"); v_dl != "" {
		v_update["deadline"] = v_dl
	}

	if len(v_update) == 0 {
		others.ResponseService(p_c, v_method, 422, "Tidak ada data yang diupdate", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.UpdateGoal(v_goal_id, v_user_id, v_update)

	if v_hasil["status"] == "T" {
		dao.CacheDelete("goal:detail:" + v_goal_id)
		dao.CacheDeletePattern("goal:list:" + v_user_id + ":*")
		dao.CacheDelete("goal:summary:" + v_user_id)
		others.ResponseService(p_c, v_method, 200, "Success", nil)
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func DeleteGoal(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.DeleteGoal(v_goal_id, v_user_id)

	if v_hasil["status"] == "T" {
		dao.CacheDelete("goal:detail:" + v_goal_id)
		dao.CacheDeletePattern("goal:list:" + v_user_id + ":*")
		dao.CacheDelete("goal:summary:" + v_user_id)
		others.ResponseService(p_c, v_method, 200, "Success", nil)
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func AddSaving(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_amount        := others.GetStr(v_data, "amount")
	v_note          := others.GetStr(v_data, "note")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_amount == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_amount_float, v_err := strconv.ParseFloat(v_amount, 64)
	if v_err != nil || v_amount_float <= 0 {
		others.ResponseService(p_c, v_method, 422, "Amount harus angka positif", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_amount + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.AddSaving(v_goal_id, v_user_id, v_note, v_amount_float)

	if v_hasil["status"] == "T" {
		dao.CacheDelete("goal:detail:" + v_goal_id)
		dao.CacheDeletePattern("goal:list:" + v_user_id + ":*")
		dao.CacheDeletePattern("goal:history:" + v_goal_id + ":*")
		dao.CacheDelete("goal:summary:" + v_user_id)
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, v_hasil["status_code"].(int), v_hasil["message"].(string), nil)
	}
}

func WithdrawSaving(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_amount        := others.GetStr(v_data, "amount")
	v_note          := others.GetStr(v_data, "note")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_amount == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_amount_float, v_err := strconv.ParseFloat(v_amount, 64)
	if v_err != nil || v_amount_float <= 0 {
		others.ResponseService(p_c, v_method, 422, "Amount harus angka positif", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_amount + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.WithdrawSaving(v_goal_id, v_user_id, v_note, v_amount_float)

	if v_hasil["status"] == "T" {
		dao.CacheDelete("goal:detail:" + v_goal_id)
		dao.CacheDeletePattern("goal:list:" + v_user_id + ":*")
		dao.CacheDeletePattern("goal:history:" + v_goal_id + ":*")
		dao.CacheDelete("goal:summary:" + v_user_id)
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, v_hasil["status_code"].(int), v_hasil["message"].(string), nil)
	}
}

func GetSavingHistory(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_page          := others.GetStr(v_data, "page")
	v_limit         := others.GetStr(v_data, "limit")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	if v_page == "" {
		v_page = "1"
	}
	if v_limit == "" {
		v_limit = "20"
	}

	v_user_id               := v_session_data["user_id"].(string)
	_, v_limit_int, v_offset_int := others.ParsePagination(v_page, v_limit)
	v_hasil, _              := dao.GetSavingHistory(v_goal_id, v_user_id, v_limit_int, v_offset_int)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetGoalProgress(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetGoalProgress(v_goal_id, v_user_id)

	if v_hasil["status_code"] == 204 {
		others.ResponseService(p_c, v_method, 204, "Data tidak ditemukan", nil)
	} else if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetFinanceSummary(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetFinanceSummary(v_user_id)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetMilestoneList(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_goal_id       := others.GetStr(v_data, "goal_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_goal_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_goal_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetMilestoneList(v_goal_id, v_user_id)

	if v_hasil["status_code"] == 204 {
		others.ResponseService(p_c, v_method, 204, "Data tidak ditemukan", nil)
	} else if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}
