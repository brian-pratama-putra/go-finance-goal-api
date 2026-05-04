package main

import (
	"net/http"
	"os"

	"go-finance-goal-api/config"
	controllers "go-finance-goal-api/controllers/api"
	"go-finance-goal-api/middleware"
	"go-finance-goal-api/others"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.GetPgPool()
	config.GetRedisClient()

	if os.Getenv("STATUS_APP") != "DEV" {
		gin.SetMode(gin.ReleaseMode)
	}

	v_router := gin.New()
	v_router.Use(gin.Recovery())
	v_router.Use(middleware.SecurityMiddleware())
	v_router.Use(middleware.SecurityHeaders())

	v_router.NoRoute(func(p_c *gin.Context) {
		p_c.JSON(http.StatusOK, gin.H{
			"err_code": 404,
			"err_msg":  "Route not found.",
			"datetime": others.NowJakartaStr(),
		})
	})

	v_router.GET("/", func(p_c *gin.Context) {
		p_c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	v_router.GET("/health", func(p_c *gin.Context) {
		p_c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	v_inquiry := v_router.Group("/")
	v_inquiry.Use(middleware.ParseBody())
	v_inquiry.Use(middleware.RateLimiter())
	{
		v_inquiry.POST("/inquiry", InquiryHandler)
	}

	v_router.Run(":" + config.GetPort())
}

func InquiryHandler(p_c *gin.Context) {
	v_body, _ := p_c.Get("parsed_body")
	v_data     := v_body.(map[string]interface{})

	v_method    := ""
	v_datetime  := ""
	if v_m, v_ok := v_data["method"].(string); v_ok {
		v_method = v_m
	}
	if v_dt, v_ok := v_data["datetime"].(string); v_ok {
		v_datetime = v_dt
	}

	if v_datetime != "" {
		v_is_valid, v_msg := others.ValidateRequestDatetime(v_datetime)
		if !v_is_valid {
			others.ResponseService(p_c, v_method, 405, v_msg, nil)
			return
		}
	}

	switch v_method {
	case "create_goal":
		controllers.CreateGoal(p_c)
	case "get_goal_list":
		controllers.GetGoalList(p_c)
	case "get_goal_detail":
		controllers.GetGoalDetail(p_c)
	case "update_goal":
		controllers.UpdateGoal(p_c)
	case "delete_goal":
		controllers.DeleteGoal(p_c)
	case "add_saving":
		controllers.AddSaving(p_c)
	case "withdraw_saving":
		controllers.WithdrawSaving(p_c)
	case "get_saving_history":
		controllers.GetSavingHistory(p_c)
	case "get_goal_progress":
		controllers.GetGoalProgress(p_c)
	case "get_finance_summary":
		controllers.GetFinanceSummary(p_c)
	case "get_milestone_list":
		controllers.GetMilestoneList(p_c)
	default:
		others.ResponseService(p_c, v_method, 405, "Invalid Method", nil)
	}
}
