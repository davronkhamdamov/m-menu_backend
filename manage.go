package main

import (
	"fmt"
	"net/http"

	"github.com/davronkhamdamov/restaraunt_backend/middleware"
	"github.com/davronkhamdamov/restaraunt_backend/models"
	"github.com/davronkhamdamov/restaraunt_backend/views"
	"github.com/gorilla/mux"
)

func main() {
	runserver()
}

func runserver() {
	router := mux.NewRouter()
	models.ConnectDB()
	models.MigrateDB()

	// Auth
	router.HandleFunc("/v1/login", views.Login).Methods("POST")
	// Users
	router.Handle("/v1/staff", middleware.AuthMiddleware(http.HandlerFunc(views.CreateStaff))).Methods("POST")
	router.Handle("/v1/staff", middleware.AuthMiddleware(http.HandlerFunc(views.GetStaffs))).Methods("GET")
	router.Handle("/v1/staff/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.GetStaff))).Methods("GET")
	router.Handle("/v1/staff/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.UpdateStaff))).Methods("PUT")
	// Tables
	router.HandleFunc("/v1/table/{id}", views.GetTable).Methods("GET")
	router.HandleFunc("/v1/table/one/{id}", views.GetOneTable).Methods("GET")
	router.Handle("/v1/table", middleware.AuthMiddleware(http.HandlerFunc(views.CreateTable))).Methods("POST")
	router.Handle("/v1/table", middleware.AuthMiddleware(http.HandlerFunc(views.GetTables))).Methods("GET")
	router.Handle("/v1/table/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.UpdateTable))).Methods("PUT")
	router.Handle("/v1/table/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.DeleteTable))).Methods("DELETE")
	// Food
	router.HandleFunc("/v1/food/{id}", views.GetFood).Methods("GET")
	router.HandleFunc("/v1/food-with-category", views.GetCategoriesAndFoods).Methods("GET")
	router.HandleFunc("/v1/food", views.GetAllFood).Methods("GET")
	router.Handle("/v1/food", middleware.AuthMiddleware(http.HandlerFunc(views.CreateFood))).Methods("POST")
	router.Handle("/v1/food/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.UpdateFood))).Methods("PUT")
	router.Handle("/v1/food/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.DeleteFood))).Methods("DELETE")
	// Category
	router.HandleFunc("/v1/category/{id}", views.GetCategory).Methods("GET")
	router.HandleFunc("/v1/category", views.GetAllCategory).Methods("GET")
	router.Handle("/v1/category", middleware.AuthMiddleware(http.HandlerFunc(views.CreateCategory))).Methods("POST")
	router.Handle("/v1/category/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.UpdateCategory))).Methods("PUT")
	router.Handle("/v1/category/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.DeleteCategory))).Methods("DELETE")
	// Order
	router.HandleFunc("/ws", views.Orders)
	router.HandleFunc("/v1/order", views.NewOrder).Methods("POST")
	router.HandleFunc("/v1/order/xlsx", views.DownloadOrderExcel).Methods("GET")
	router.HandleFunc("/v1/order/{id}", views.GetOrder).Methods("GET")
	router.HandleFunc("/v1/order", views.GetOrders).Methods("GET")
	router.Handle("/v1/order_staff", middleware.AuthMiddleware(http.HandlerFunc(views.GetOrdersForStaff))).Methods("GET")
	router.Handle("/v1/order/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.UpdateOrderStatus))).Methods("PUT")
	router.Handle("/v1/order/receive/{id}", middleware.AuthMiddleware(http.HandlerFunc(views.ReceiveOrder))).Methods("PUT")
	router.Handle("/v1/orders", middleware.AuthMiddleware(http.HandlerFunc(views.DeleteAllOrders))).Methods("DELETE")
	// Feedback
	router.HandleFunc("/v1/feedback", views.GetAllFeedback).Methods("GET")
	router.HandleFunc("/v1/feedback/xlsx", views.DownloadFeedbackExcel).Methods("GET")
	router.HandleFunc("/v1/feedback", views.CreateFeedback).Methods("POST")
	// Dashboard
	router.Handle("/v1/dashboard", middleware.AuthMiddleware(http.HandlerFunc(views.GetDashboard))).Methods("GET")
	router.HandleFunc("/v1/common_food", views.GetMostCommonFood).Methods("GET")

	fmt.Println("Starting Server http://localhost:8080/")
	http.ListenAndServe("0.0.0.0:8080", middleware.CorsMiddleware(router))
}
