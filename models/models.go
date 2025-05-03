package models

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// dsn := "host=ep-long-violet-a43f1tu5-pooler.us-east-1.aws.neon.tech user=restaurant_owner password=npg_IijM68WRqnUz dbname=restaurant port=5432 TimeZone=Asia/Tashkent"
	dsn := "host=localhost user=postgres password=j24xt200 dbname=restaraunt_backend port=5432 TimeZone=Asia/Tashkent"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected to PostgreSQL database!")
}

func MigrateDB() {
	err := DB.AutoMigrate(&User{}, &Table{}, &Category{}, &Food{}, &Order{}, &OrderFood{}, &Feedback{})
	if err != nil {
		panic("failed to migrate database")
	}
	fmt.Println("Database migrated!")
}

type UserRole int

const (
	Admin UserRole = iota
	Staff
)

func (r UserRole) String() string {
	return [...]string{"Admin", "Staff"}[r]
}

type User struct {
	ID        string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Login     string    `gorm:"unique;not null" json:"login"`
	Password  string    `json:"password" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated"`
}

type Table struct {
	ID        string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Number    uint      `gorm:"unique; not null" json:"number" validate:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated"`
}
type Category struct {
	ID        string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	NameUz    string    `json:"name_uz" validate:"required"`
	NameRu    string    `json:"name_ru" validate:"required"`
	NameEn    string    `json:"name_en" validate:"required"`
	Name      string    `json:"name" gorm:"-"`
	Foods     []Food    `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"foods"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated"`
}
type Food struct {
	ID            string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	NameUz        string    `json:"name_uz" validate:"required"`
	Name          string    `json:"name" gorm:"-"`
	NameRu        string    `json:"name_ru" validate:"required"`
	NameEn        string    `json:"name_en" validate:"required"`
	DescriptionUz string    `json:"description_uz" validate:"required"`
	Description   string    `json:"description" gorm:"-"`
	DescriptionRu string    `json:"description_ru" validate:"required"`
	DescriptionEn string    `json:"description_en" validate:"required"`
	Price         uint      `json:"price" validate:"required"`
	ImageUrl      string    `json:"image_url" validate:"required"`
	Weight        uint      `json:"weight" validate:"required"`
	WeightType    string    `json:"weight_type" validate:"required"`
	Available     bool      `json:"available" gorm:"default:true" validate:"-"`
	CategoryID    string    `gorm:"not null" json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" validate:"-"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated"`
}
type Order struct {
	ID        string      `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	TableID   string      `gorm:"not null" json:"table_id" validate:"required"`
	OrderId   string      `gorm:"" json:"order_id" validate:"-"`
	Table     Table       `gorm:"foreignKey:TableID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"table" validate:"-"`
	UserID    *string     `json:"user_id"`
	User      User        `gorm:"foreignKey:UserID" json:"-" validate:"-"`
	Total     uint        `gorm:"not null" json:"total"`
	Status    string      `gorm:"not null" json:"status"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated"`
	OrderFood []OrderFood `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"foods" validate:"-"`
}
type OrderFood struct {
	ID            string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID       string    `gorm:"not null" json:"order_id"`
	FoodID        string    `gorm:"not null" json:"food_id" validate:"required"`
	Quantity      uint      `gorm:"not null" json:"quantity" validate:"required"`
	NameUz        string    `json:"name_uz"`
	NameRu        string    `json:"name_ru"`
	NameEn        string    `json:"name_en"`
	DescriptionUz string    `json:"description_uz" validate:"required"`
	Description   string    `json:"description" gorm:"-"`
	DescriptionRu string    `json:"description_ru" validate:"required"`
	DescriptionEn string    `json:"description_en" validate:"required"`
	Name          string    `json:"name" gorm:"-"`
	Price         uint      `json:"price"`
	Image         string    `json:"image"`
	Weight        uint      `json:"weight"`
	WeightType    string    `json:"weight_type"`
	Food          Food      `gorm:"foreignKey:FoodID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-" validate:"-"`
	Order         Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-" validate:"-"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated"`
}
type Feedback struct {
	ID        string    `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	TableID   string    `gorm:"not null" json:"table_id" validate:"required"`
	Table     Table     `gorm:"foreignKey:TableID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"table" validate:"-"`
	Feedback  string    `json:"feedback"`
	Region    string    `json:"region" validate:"required"`
	Star      uint      `gorm:"type:int;check:star >= 1 AND star <= 5" json:"star" validate:"required,min=1,max=5"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created"`
}
type Login struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
