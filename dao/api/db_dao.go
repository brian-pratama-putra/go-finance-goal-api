package dao

import (
	"context"
	"fmt"
	"math"

	"go-finance-goal-api/config"
	"go-finance-goal-api/others"
)

func ResponseJson(p_status_code int, p_flag string, p_message string, p_result interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status_code": p_status_code,
		"status":      p_flag,
		"message":     p_message,
		"result":      p_result,
	}
}

func InsertGoal(p_user_id, p_name, p_description, p_deadline string, p_target_amount float64) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		INSERT INTO goals (user_id, name, description, target_amount, current_amount, deadline, created_at)
		VALUES($1, $2, $3, $4, 0, $5, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
		RETURNING goal_id
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_user_id, p_name, p_description, p_target_amount, p_deadline)
	var v_goal_id string
	if v_err := v_row.Scan(&v_goal_id); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", map[string]interface{}{"goal_id": v_goal_id}), nil
}

func GetGoalList(p_user_id, p_status string) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("goal:list:%s:%s", p_user_id, p_status)
	v_cached, v_err := CacheGetList(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool      := config.GetPgPool()
	v_ctx       := context.Background()
	v_where     := "WHERE g.user_id = $1 AND g.is_deleted = false"
	v_args      := []interface{}{p_user_id}
	v_idx       := 2

	if p_status != "" {
		v_where += fmt.Sprintf(" AND g.status = $%d", v_idx)
		v_args  = append(v_args, p_status)
	}

	v_query := fmt.Sprintf(`
		SELECT
			g.goal_id::varchar,
			g.name,
			coalesce(g.description, '') description,
			g.target_amount::varchar,
			g.current_amount::varchar,
			ROUND((g.current_amount / NULLIF(g.target_amount, 0) * 100)::numeric, 1)::varchar progress_percent,
			g.status,
			to_char(g.deadline, 'YYYY-MM-DD') deadline,
			to_char(g.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at
		FROM goals g
		%s
		ORDER BY g.deadline ASC NULLS LAST, g.created_at DESC
	`, v_where)

	v_rows, v_err := v_pool.Query(v_ctx, v_query, v_args...)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_rows.Close()

	v_result := []map[string]interface{}{}
	for v_rows.Next() {
		var (
			v_goal_id, v_name, v_description                       string
			v_target_amount, v_current_amount, v_progress_percent  string
			v_status, v_deadline, v_created_at                     string
		)
		if v_err := v_rows.Scan(&v_goal_id, &v_name, &v_description, &v_target_amount, &v_current_amount, &v_progress_percent, &v_status, &v_deadline, &v_created_at); v_err != nil {
			continue
		}
		v_result = append(v_result, map[string]interface{}{
			"goal_id":          v_goal_id,
			"name":             v_name,
			"description":      v_description,
			"target_amount":    v_target_amount,
			"current_amount":   v_current_amount,
			"progress_percent": v_progress_percent,
			"status":           v_status,
			"deadline":         v_deadline,
			"created_at":       v_created_at,
		})
	}

	if len(v_result) > 0 {
		CacheSetList(v_cache_key, v_result, 120)
	}
	return ResponseJson(200, "T", "Success", v_result), nil
}

func GetGoalDetail(p_goal_id, p_user_id string) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("goal:detail:%s", p_goal_id)
	v_cached, v_err := CacheGet(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			g.goal_id::varchar,
			g.name,
			coalesce(g.description, '') description,
			g.target_amount::varchar,
			g.current_amount::varchar,
			ROUND((g.current_amount / NULLIF(g.target_amount, 0) * 100)::numeric, 1)::varchar progress_percent,
			g.status,
			to_char(g.deadline, 'YYYY-MM-DD') deadline,
			to_char(g.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			to_char(g.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at
		FROM goals g
		WHERE g.goal_id = $1
		AND g.user_id = $2
		AND g.is_deleted = false
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_goal_id, p_user_id)
	var (
		v_goal_id, v_name, v_description                       string
		v_target_amount, v_current_amount, v_progress_percent  string
		v_status, v_deadline, v_created_at, v_updated_at       string
	)
	if v_err := v_row.Scan(&v_goal_id, &v_name, &v_description, &v_target_amount, &v_current_amount, &v_progress_percent, &v_status, &v_deadline, &v_created_at, &v_updated_at); v_err != nil {
		return ResponseJson(204, "T", "Data tidak ditemukan", nil), nil
	}

	v_data := map[string]interface{}{
		"goal_id":          v_goal_id,
		"name":             v_name,
		"description":      v_description,
		"target_amount":    v_target_amount,
		"current_amount":   v_current_amount,
		"progress_percent": v_progress_percent,
		"status":           v_status,
		"deadline":         v_deadline,
		"created_at":       v_created_at,
		"updated_at":       v_updated_at,
	}
	CacheSet(v_cache_key, v_data, 120)
	return ResponseJson(200, "T", "Success", v_data), nil
}

func UpdateGoal(p_goal_id, p_user_id string, p_data map[string]interface{}) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_set   := ""
	v_args  := []interface{}{}
	v_idx   := 1

	for v_k, v_v := range p_data {
		if v_set != "" {
			v_set += ", "
		}
		v_set += fmt.Sprintf("%s = $%d", v_k, v_idx)
		v_args = append(v_args, v_v)
		v_idx++
	}
	v_set  += ", updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'"
	v_args  = append(v_args, p_goal_id, p_user_id)

	v_query := fmt.Sprintf(`
		UPDATE goals SET %s
		WHERE goal_id = $%d AND user_id = $%d AND is_deleted = false
	`, v_set, v_idx, v_idx+1)

	_, v_err := v_pool.Exec(v_ctx, v_query, v_args...)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", nil), nil
}

func DeleteGoal(p_goal_id, p_user_id string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		UPDATE goals
		SET is_deleted = true, updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
		WHERE goal_id = $1 AND user_id = $2 AND is_deleted = false
	`
	_, v_err := v_pool.Exec(v_ctx, v_query, p_goal_id, p_user_id)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", nil), nil
}

func AddSaving(p_goal_id, p_user_id, p_note string, p_amount float64) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()

	v_tx, v_err := v_pool.Begin(v_ctx)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_tx.Rollback(v_ctx)

	v_check_query := `
		SELECT target_amount, current_amount, status
		FROM goals
		WHERE goal_id = $1 AND user_id = $2 AND is_deleted = false
		FOR UPDATE
	`
	v_row := v_tx.QueryRow(v_ctx, v_check_query, p_goal_id, p_user_id)
	var v_target_amount, v_current_amount float64
	var v_status string
	if v_err := v_row.Scan(&v_target_amount, &v_current_amount, &v_status); v_err != nil {
		return ResponseJson(404, "F", "Goal tidak ditemukan", nil), v_err
	}

	if v_status == "completed" {
		return ResponseJson(400, "F", "Goal sudah selesai", nil), fmt.Errorf("goal already completed")
	}

	v_insert_query := `
		INSERT INTO saving_logs (goal_id, user_id, type, amount, note, created_at)
		VALUES($1, $2, 'deposit', $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	`
	if _, v_err := v_tx.Exec(v_ctx, v_insert_query, p_goal_id, p_user_id, p_amount, p_note); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	v_new_amount    := v_current_amount + p_amount
	v_new_status    := v_status
	if v_new_amount >= v_target_amount {
		v_new_amount    = v_target_amount
		v_new_status    = "completed"
	}

	v_update_query := `
		UPDATE goals
		SET current_amount = $1, status = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
		WHERE goal_id = $3
	`
	if _, v_err := v_tx.Exec(v_ctx, v_update_query, v_new_amount, v_new_status, p_goal_id); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	if v_err := v_tx.Commit(v_ctx); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	v_progress := others.CalcProgressPercent(v_new_amount, v_target_amount)
	return ResponseJson(200, "T", "Success", map[string]interface{}{
		"current_amount":   others.Float64ToStr(v_new_amount, 2),
		"target_amount":    others.Float64ToStr(v_target_amount, 2),
		"progress_percent": others.Float64ToStr(v_progress, 1),
		"status":           v_new_status,
	}), nil
}

func WithdrawSaving(p_goal_id, p_user_id, p_note string, p_amount float64) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()

	v_tx, v_err := v_pool.Begin(v_ctx)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_tx.Rollback(v_ctx)

	v_check_query := `
		SELECT target_amount, current_amount, status
		FROM goals
		WHERE goal_id = $1 AND user_id = $2 AND is_deleted = false
		FOR UPDATE
	`
	v_row := v_tx.QueryRow(v_ctx, v_check_query, p_goal_id, p_user_id)
	var v_target_amount, v_current_amount float64
	var v_status string
	if v_err := v_row.Scan(&v_target_amount, &v_current_amount, &v_status); v_err != nil {
		return ResponseJson(404, "F", "Goal tidak ditemukan", nil), v_err
	}

	if p_amount > v_current_amount {
		return ResponseJson(400, "F", "Jumlah withdraw melebihi tabungan saat ini", nil), fmt.Errorf("insufficient saving")
	}

	v_insert_query := `
		INSERT INTO saving_logs (goal_id, user_id, type, amount, note, created_at)
		VALUES($1, $2, 'withdraw', $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	`
	if _, v_err := v_tx.Exec(v_ctx, v_insert_query, p_goal_id, p_user_id, p_amount, p_note); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	v_new_amount    := v_current_amount - p_amount
	v_new_status    := "on_track"
	if v_new_amount <= 0 {
		v_new_amount    = 0
		v_new_status    = "on_track"
	}

	v_update_query := `
		UPDATE goals
		SET current_amount = $1, status = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
		WHERE goal_id = $3
	`
	if _, v_err := v_tx.Exec(v_ctx, v_update_query, v_new_amount, v_new_status, p_goal_id); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	if v_err := v_tx.Commit(v_ctx); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	v_progress := others.CalcProgressPercent(v_new_amount, v_target_amount)
	return ResponseJson(200, "T", "Success", map[string]interface{}{
		"current_amount":   others.Float64ToStr(v_new_amount, 2),
		"target_amount":    others.Float64ToStr(v_target_amount, 2),
		"progress_percent": others.Float64ToStr(v_progress, 1),
		"status":           v_new_status,
	}), nil
}

func GetSavingHistory(p_goal_id, p_user_id string, p_limit, p_offset int) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("goal:history:%s:%d:%d", p_goal_id, p_offset, p_limit)
	v_cached, v_err := CacheGetList(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			sl.log_id::varchar,
			sl.type,
			sl.amount::varchar,
			coalesce(sl.note, '') note,
			to_char(sl.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at
		FROM saving_logs sl
		JOIN goals g ON g.goal_id = sl.goal_id
		WHERE sl.goal_id = $1
		AND sl.user_id = $2
		AND g.is_deleted = false
		ORDER BY sl.created_at DESC
		LIMIT $3 OFFSET $4
	`
	v_rows, v_err := v_pool.Query(v_ctx, v_query, p_goal_id, p_user_id, p_limit, p_offset)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_rows.Close()

	v_result := []map[string]interface{}{}
	for v_rows.Next() {
		var v_log_id, v_type, v_amount, v_note, v_created_at string
		if v_err := v_rows.Scan(&v_log_id, &v_type, &v_amount, &v_note, &v_created_at); v_err != nil {
			continue
		}
		v_result = append(v_result, map[string]interface{}{
			"log_id":     v_log_id,
			"type":       v_type,
			"amount":     v_amount,
			"note":       v_note,
			"created_at": v_created_at,
		})
	}

	if len(v_result) > 0 {
		CacheSetList(v_cache_key, v_result, 120)
	}
	return ResponseJson(200, "T", "Success", v_result), nil
}

func GetGoalProgress(p_goal_id, p_user_id string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			g.name,
			g.target_amount,
			g.current_amount,
			g.status,
			to_char(g.deadline, 'YYYY-MM-DD') deadline
		FROM goals g
		WHERE g.goal_id = $1
		AND g.user_id = $2
		AND g.is_deleted = false
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_goal_id, p_user_id)
	var v_name, v_status, v_deadline string
	var v_target_amount, v_current_amount float64
	if v_err := v_row.Scan(&v_name, &v_target_amount, &v_current_amount, &v_status, &v_deadline); v_err != nil {
		return ResponseJson(204, "T", "Data tidak ditemukan", nil), nil
	}

	v_remaining_amount  := math.Max(0, v_target_amount-v_current_amount)
	v_days_remaining    := others.CalcDaysRemaining(v_deadline)
	v_daily_needed      := others.CalcDailySavingNeeded(v_remaining_amount, float64(v_days_remaining))
	v_progress         := others.CalcProgressPercent(v_current_amount, v_target_amount)

	return ResponseJson(200, "T", "Success", map[string]interface{}{
		"name":                 v_name,
		"target_amount":        others.Float64ToStr(v_target_amount, 2),
		"current_amount":       others.Float64ToStr(v_current_amount, 2),
		"remaining_amount":     others.Float64ToStr(v_remaining_amount, 2),
		"progress_percent":     others.Float64ToStr(v_progress, 1),
		"days_remaining":       others.IntToStr(v_days_remaining),
		"daily_saving_needed":  others.Float64ToStr(v_daily_needed, 2),
		"status":               v_status,
		"deadline":             v_deadline,
	}), nil
}

func GetFinanceSummary(p_user_id string) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("goal:summary:%s", p_user_id)
	v_cached, v_err := CacheGet(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			COUNT(*)::varchar                                                                    total_goal,
			COUNT(CASE WHEN status = 'on_track' THEN 1 END)::varchar                            total_on_track,
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::varchar                           total_completed,
			COUNT(CASE WHEN deadline < CURRENT_DATE AND status != 'completed' THEN 1 END)::varchar total_overdue,
			COALESCE(SUM(current_amount), 0)::varchar                                           total_saved,
			COALESCE(SUM(target_amount), 0)::varchar                                            total_target,
			ROUND(COALESCE(SUM(current_amount) / NULLIF(SUM(target_amount), 0) * 100, 0)::numeric, 1)::varchar overall_progress
		FROM goals
		WHERE user_id = $1
		AND is_deleted = false
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_user_id)
	var (
		v_total_goal, v_total_on_track, v_total_completed, v_total_overdue string
		v_total_saved, v_total_target, v_overall_progress                  string
	)
	if v_err := v_row.Scan(&v_total_goal, &v_total_on_track, &v_total_completed, &v_total_overdue, &v_total_saved, &v_total_target, &v_overall_progress); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	v_data := map[string]interface{}{
		"total_goal":       v_total_goal,
		"total_on_track":   v_total_on_track,
		"total_completed":  v_total_completed,
		"total_overdue":    v_total_overdue,
		"total_saved":      v_total_saved,
		"total_target":     v_total_target,
		"overall_progress": v_overall_progress,
	}
	CacheSet(v_cache_key, v_data, 120)
	return ResponseJson(200, "T", "Success", v_data), nil
}

func GetMilestoneList(p_goal_id, p_user_id string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			g.target_amount,
			g.current_amount,
			ROUND((g.current_amount / NULLIF(g.target_amount, 0) * 100)::numeric, 1) progress_pct
		FROM goals g
		WHERE g.goal_id = $1
		AND g.user_id = $2
		AND g.is_deleted = false
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_goal_id, p_user_id)
	var v_target_amount, v_current_amount, v_progress_pct float64
	if v_err := v_row.Scan(&v_target_amount, &v_current_amount, &v_progress_pct); v_err != nil {
		return ResponseJson(204, "T", "Data tidak ditemukan", nil), nil
	}

	v_milestones := []map[string]interface{}{}
	v_pct_list   := []float64{25, 50, 75, 100}
	for _, v_pct := range v_pct_list {
		v_milestone_amount  := v_target_amount * v_pct / 100
		v_reached           := v_current_amount >= v_milestone_amount
		v_milestones = append(v_milestones, map[string]interface{}{
			"milestone":        fmt.Sprintf("%.0f%%", v_pct),
			"amount_needed":    others.Float64ToStr(v_milestone_amount, 2),
			"reached":          v_reached,
		})
	}

	return ResponseJson(200, "T", "Success", v_milestones), nil
}
