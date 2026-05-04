package others

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

var TZ_JAKARTA *time.Location

func init() {
	var v_err error
	TZ_JAKARTA, v_err = time.LoadLocation("Asia/Jakarta")
	if v_err != nil {
		TZ_JAKARTA = time.FixedZone("WIB", 7*3600)
	}
}

func GetTTLUntilMidnightJakarta() int {
	v_now      := time.Now().In(TZ_JAKARTA)
	v_midnight := time.Date(v_now.Year(), v_now.Month(), v_now.Day()+1, 0, 0, 0, 0, TZ_JAKARTA)
	return int(v_midnight.Sub(v_now).Seconds())
}

func ValidateRequestDatetime(p_datetime string) (bool, string) {
	v_layout         := "2006-01-02 15:04:05"
	v_request_dt, v_err := time.ParseInLocation(v_layout, p_datetime, TZ_JAKARTA)
	if v_err != nil {
		return false, "Format datetime tidak valid. Gunakan YYYY-MM-DD HH:MM:SS"
	}
	v_now_jakarta   := time.Now().In(TZ_JAKARTA)
	v_diff          := v_now_jakarta.Sub(v_request_dt)
	v_abs_diff      := time.Duration(math.Abs(float64(v_diff)))
	if v_abs_diff > 24*time.Hour {
		return false, "Datetime request harus dalam rentang ±24 jam dari waktu sekarang"
	}
	return true, "OK"
}

func NowJakartaStr() string {
	return time.Now().In(TZ_JAKARTA).Format("2006-01-02 15:04:05")
}

func IntToStr(p_n int) string {
	return fmt.Sprintf("%d", p_n)
}

func Float64ToStr(p_f float64, p_decimal int) string {
	return strconv.FormatFloat(p_f, 'f', p_decimal, 64)
}

func ParsePagination(p_page, p_limit string) (int, int, int) {
	v_page, v_err := strconv.Atoi(p_page)
	if v_err != nil || v_page < 1 {
		v_page = 1
	}
	v_limit, v_err := strconv.Atoi(p_limit)
	if v_err != nil || v_limit < 1 {
		v_limit = 20
	}
	if v_limit > 100 {
		v_limit = 100
	}
	v_offset := (v_page - 1) * v_limit
	return v_page, v_limit, v_offset
}

func CalcDaysRemaining(p_deadline string) int {
	v_deadline, v_err := time.ParseInLocation("2006-01-02", p_deadline, TZ_JAKARTA)
	if v_err != nil {
		return 0
	}
	v_now   := time.Now().In(TZ_JAKARTA).Truncate(24 * time.Hour)
	v_diff  := v_deadline.Sub(v_now)
	v_days  := int(v_diff.Hours() / 24)
	if v_days < 0 {
		return 0
	}
	return v_days
}

func CalcDailySavingNeeded(p_remaining_amount, p_days_remaining float64) float64 {
	if p_days_remaining <= 0 {
		return p_remaining_amount
	}
	return math.Ceil(p_remaining_amount / p_days_remaining)
}

func CalcProgressPercent(p_current, p_target float64) float64 {
	if p_target <= 0 {
		return 0
	}
	v_pct := (p_current / p_target) * 100
	if v_pct > 100 {
		return 100
	}
	return math.Round(v_pct*10) / 10
}
