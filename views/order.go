package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davronkhamdamov/restaraunt_backend/middleware"
	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var your_secret_key = []byte("your_secret_key")
var HubInstance = utils.NewHub()

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type OrderRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}

var Clients = make(map[*websocket.Conn]bool)

func NewOrder(w http.ResponseWriter, r *http.Request) {
	var request models.Order

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	order := models.Order{
		TableID: request.TableID,
		Status:  "pending",
	}
	tx := models.DB.Begin()
	if tx.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to start transaction", tx.Error.Error())
		return
	}

	if err := createOrder(tx, &order); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order", err.Error())
		return
	}

	if err := processOrderFoods(tx, &order, request); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order foods", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}
	var table models.Table
	if dbResult := models.DB.Where("ID = ?", order.TableID).Find(&table); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve orders", dbResult.Error.Error)
	}
	order.Table = table
	HubInstance.BroadcastToAll(utils.WebSocketMessage{
		Event: "new_order",
		Data:  order,
	})
	utils.RespondWithSuccess(w, http.StatusCreated, "Order created successfully", nil)
}

func createOrder(tx *gorm.DB, order *models.Order) error {
	order.OrderId = utils.GenerateOrderID()
	dbResult := tx.Create(order)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func processOrderFoods(tx *gorm.DB, order *models.Order, request models.Order) error {
	var total uint = 0

	for i := range request.OrderFood {
		orderFood := models.OrderFood{
			OrderID:  order.ID,
			FoodID:   request.OrderFood[i].FoodID,
			Quantity: request.OrderFood[i].Quantity,
		}
		var food models.Food
		if err := tx.First(&food, "ID = ?", orderFood.FoodID).Error; err != nil {
			return err
		}
		orderFood.Weight = food.Weight
		orderFood.NameUz = food.NameUz
		orderFood.NameRu = food.NameRu
		orderFood.NameEn = food.NameEn
		orderFood.WeightType = food.WeightType
		orderFood.Price = food.Price
		orderFood.Image = food.ImageUrl
		orderFood.DescriptionUz = food.DescriptionUz
		orderFood.DescriptionRu = food.DescriptionRu
		orderFood.DescriptionEn = food.DescriptionEn
		total += food.Price * orderFood.Quantity
		if err := tx.Create(&orderFood).Error; err != nil {
			return err
		}
	}

	order.Total = total
	if err := tx.Save(order).Error; err != nil {
		return err
	}

	return nil
}
func GetOrder(w http.ResponseWriter, r *http.Request) {
	var orders models.Order
	vars := mux.Vars(r)
	orderID := vars["id"]
	lang := r.URL.Query().Get("lang")
	if err := models.DB.Where("ID = ?", orderID).Preload("OrderFood").Preload("Table").First(&orders).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get orders", err.Error())
		return
	}
	for i, food := range orders.OrderFood {
		var foodNameToReturn string
		switch lang {
		case "uz":
			foodNameToReturn = food.NameUz
		case "ru":
			foodNameToReturn = food.NameRu
		case "en":
			foodNameToReturn = food.NameEn
		default:
			foodNameToReturn = food.NameEn
		}
		orders.OrderFood[i].Name = foodNameToReturn
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Orders retrieved successfully", orders)
}
func getOrdersFromDB() ([]models.Order, error) {
	var orders []models.Order
	if err := models.DB.
		Order("Created_At DESC").
		Preload("Feedback").
		Preload("Table").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
func GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := getOrdersFromDB()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get orders", err.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Orders retrieved successfully", orders)
}
func GetOrdersForStaff(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	var orders []models.Order
	if err := models.DB.
		Where("status IN ?", []string{"in_process", "pending"}).
		Where("user_id = ? OR user_id IS NULL", userID).
		Order("Created_At DESC").
		Preload("Table").
		Find(&orders); err.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get orders", err.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Orders retrieved successfully", orders)
}

func authenticateWebSocket(r *http.Request) (string, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		return "", fmt.Errorf("authorization header is missing")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return your_secret_key, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	return claims.UserID, nil
}

func Orders(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	userId, errAuth := authenticateWebSocket(r)
	if errAuth != nil {
		log.Println("Authentication failed:", err)
	}
	tableID := r.URL.Query().Get("table_id")

	client := &utils.Client{
		Conn:   conn,
		UserID: userId,
		Rooms:  make(map[string]bool),
	}
	if tableID != "" {
		client.UserID = tableID
	}
	HubInstance.Clients[client] = true
	HubInstance.JoinRoom(client, client.UserID)
	go readPump(client)
}
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	userID := r.Context().Value(middleware.UserIDKey)

	var order models.Order
	if err := models.DB.First(&order, "ID = ?", orderID).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found", err.Error())
		return
	}
	if order.Status == "done" {
		utils.RespondWithError(w, http.StatusBadRequest, "Order already has the specified status", nil)
		return
	}
	if userID != *order.UserID {
		utils.RespondWithError(w, http.StatusBadRequest, "Unauthorized status update attempt", "You must claim the order before changing its status")
		return
	}
	order.Status = "done"
	if err := models.DB.Save(&order).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update order status", err.Error())
		return
	}

	HubInstance.BroadcastToRoom(order.TableID, utils.WebSocketMessage{
		Event: "status_updated",
		Data:  order,
	})
	if order.UserID != nil {
		HubInstance.BroadcastToRoom(*order.UserID, utils.WebSocketMessage{
			Event: "status_updated",
			Data:  order,
		})
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Status updated", nil)
}
func ReceiveOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	order := models.Order{}
	userID := r.Context().Value(middleware.UserIDKey)

	if dbResult := models.DB.First(&order, "ID = ?", orderID).Error; dbResult != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found", dbResult.Error())
		return
	}
	if order.UserID != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Order already received", nil)
		return
	}
	if val, ok := userID.(string); ok {
		order.UserID = &val
		order.Status = "in_process"
		if dbResult := models.DB.Save(&order).Error; dbResult != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update order", dbResult.Error())
			return
		}
		HubInstance.BroadcastToAll(utils.WebSocketMessage{Event: "status_updated", Data: order})
		utils.RespondWithSuccess(w, http.StatusOK, "Order received", order)
		return
	}
	utils.RespondWithError(w, http.StatusBadRequest, "Failed to order", nil)
}
func DeleteAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := models.Table{}

	if dbResult := models.DB.Where("1 = 1").Delete(&orders); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Orders not found", dbResult.Error.Error())
		return
	}
	if dbResult := models.DB.Delete(&orders); dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", dbResult.Error.Error())
		return
	}
	utils.RespondWithSuccess(w, http.StatusOK, "Orders deleted successfully", nil)
}

func DownloadOrderExcel(w http.ResponseWriter, r *http.Request) {
	type PopularOrder struct {
		OrderId     string
		NameUz      string
		Price       uint
		Quantity    uint
		Weight      uint
		Region      string
		WeightType  string
		TableNumber string
		CreatedAt   time.Time
	}

	var results []PopularOrder

	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	models.DB.Table("order_foods").
		Select(`orders.order_id,
			order_foods.name_uz,
	        order_foods.price,
	        feedbacks.region,
	        order_foods.quantity,
	        order_foods.weight,
	        order_foods.weight_type,
	        tables.number as table_number,
	        orders.created_at`).
		Joins("JOIN orders ON orders.id = order_foods.order_id").
		Joins("JOIN tables ON tables.id = orders.table_id").
		Joins("LEFT JOIN feedbacks ON orders.id = feedbacks.order_id").
		Where("order_foods.created_at >= ?", oneWeekAgo).
		Order("orders.created_at DESC").
		Scan(&results)

	exporter := utils.NewTaomExcelExporter()
	headers := []string{"Buyurtma raqami", "Ovqat nomi", "Mijoz davlati", "Narxi (so'm)", "Soni", "O'girligi", "Stol raqami", "Vaqt"}
	exporter.SetHeaders(headers)

	var rows [][]any
	for _, r := range results {
		rows = append(rows, []any{
			r.OrderId,
			r.NameUz,
			r.Region,
			r.Price,
			r.Quantity,
			fmt.Sprintf("%d %s", r.Weight, r.WeightType),
			r.TableNumber,
			r.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	exporter.SetRows(rows)

	buf, err := exporter.Generate()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="weekly_orders.xlsx"`)
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
func readPump(client *utils.Client) {
	defer func() {
		client.Conn.Close()
	}()
	for {
		var msg utils.WebSocketMessage
		if err := client.Conn.ReadJSON(&msg); err != nil {
			break
		}
		event := msg.Event
		data := msg.Data

		switch event {
		case "status_updated":
			HubInstance.BroadcastToRoom(client.UserID, utils.WebSocketMessage{
				Event: "status_updated",
				Data:  data,
			})
		}
	}
}
