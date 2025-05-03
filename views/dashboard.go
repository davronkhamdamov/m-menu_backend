package views

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
)

type WeekReport struct {
	Date  string `json:"date"`
	Total uint   `json:"total"`
}
type MostPopularFood struct {
	Id          string `json:"id"`
	TotalQty    uint   `json:"total_quantity"`
	Name        string `json:"name"`
	Image       string `json:"image_url"`
	Price       uint   `json:"price"`
	Weight      uint   `json:"weight"`
	WeightType  string `json:"weight_type"`
	Description string `json:"description"`
}

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	var (
		ordersCount     int64
		tablesCount     int64
		categoriesCount int64
		foodsCount      int64
		feedbacksCount  int64
		todayRevenue    int64
		ordersThisWeek  []models.Order
	)

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	startOfWeek := now.AddDate(0, 0, -6)

	models.DB.Model(&models.Order{}).
		Where("created_at >= ? AND created_at < ? AND user_id IS NOT NULL AND status = 'done'", startOfDay, endOfDay).
		Count(&ordersCount)

	models.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total), 0)").
		Where("created_at >= ? AND created_at < ? AND user_id IS NOT NULL AND status = 'done'", startOfDay, endOfDay).
		Scan(&todayRevenue)

	models.DB.Model(&models.Feedback{}).Count(&feedbacksCount)
	models.DB.Model(&models.Table{}).Count(&tablesCount)
	models.DB.Model(&models.Category{}).Count(&categoriesCount)
	models.DB.Model(&models.Food{}).Count(&foodsCount)

	models.DB.Where("created_at BETWEEN ? AND ? AND user_id IS NOT NULL AND status = 'done'", startOfWeek, endOfDay).Find(&ordersThisWeek)

	dailyTotals := make(map[string]uint)
	for _, order := range ordersThisWeek {
		dateStr := order.CreatedAt.Format("2006-01-02")
		dailyTotals[dateStr] += order.Total
	}

	var oneWeekReport []WeekReport
	for i := range 7 {
		date := startOfWeek.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		oneWeekReport = append(oneWeekReport, WeekReport{
			Date:  dateStr,
			Total: dailyTotals[dateStr],
		})
	}
	response := map[string]interface{}{
		"today_orders":     ordersCount,
		"total_feedback":   feedbacksCount,
		"total_tables":     tablesCount,
		"total_categories": categoriesCount,
		"total_foods":      foodsCount,
		"today_revenue":    todayRevenue,
		"one_week_report":  oneWeekReport,
	}
	utils.RespondWithSuccess(w, http.StatusOK, "OK", response)
}

func GetMostCommonFood(w http.ResponseWriter, r *http.Request) {
	var results []MostPopularFood
	lang := r.URL.Query().Get("lang")
	fmt.Println(lang)
	nameCol := "order_foods.name_uz"
	nameColDesc := "order_foods.description_uz"
	if lang == "en" {
		nameCol = "order_foods.name_en"
		nameColDesc = "order_foods.description_en"
	} else if lang == "ru" {
		nameCol = "order_foods.name_ru"
		nameColDesc = "order_foods.description_ru"
	}

	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	models.DB.Table("order_foods").
		Select(fmt.Sprintf("order_foods.food_id as id, SUM(order_foods.quantity) as total_qty, %s as name, order_foods.image as image, order_foods.price as price, order_foods.weight as weight, order_foods.weight_type as weight_type, %s as description", nameCol, nameColDesc)).
		Where("order_foods.created_at >= ?", oneWeekAgo).
		Group(fmt.Sprintf("order_foods.food_id, %s, order_foods.image, order_foods.price, order_foods.weight, order_foods.weight_type, %s", nameCol, nameColDesc)).
		Order("total_qty DESC").
		Limit(6).
		Scan(&results)
	utils.RespondWithSuccess(w, http.StatusOK, "OK", results)
}
