package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct{
	ID int
	Name string
	Email string
	Phone_number string
	Job_post string
	Department string
	Address string
	Education_details string
}


var db *gorm.DB

func main(){

	dsn := "host=localhost user=postgres password=password dbname=usersdb port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err !=nil {
		panic(err)
	}
	db.AutoMigrate(&User{})

	e:=echo.New()

	e.POST("/users",addUser)
	e.GET("/users",listUsers)
	e.PUT("/users/:id",updateUser)
	e.DELETE("/users/:id", deleteUser)
	

	e.Start(":8080")
}

func addUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := db.Create(u).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusCreated, u)
}




func updateUser(c echo.Context) error {
    id := c.Param("id")
    var user User

	if err := db.First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
    }

	if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
	if err := db.Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, user)

}
func deleteUser(c echo.Context) error {
    id := c.Param("id")
    var user User

	if err := db.First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
    }

	if err := db.Delete(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func listUsers(c echo.Context) error {
    var users []User
	department := c.QueryParam("department")
    jobPost := c.QueryParam("job_post")
    search := c.QueryParam("search")
    page := c.QueryParam("page")
    limit := c.QueryParam("limit")

	pageNum := 1
    pageSize := 10

	if page != "" {
        fmt.Sscanf(page, "%d", &pageNum)
    }
    if limit != "" {
        fmt.Sscanf(limit, "%d", &pageSize)
    }

    offset := (pageNum - 1) * pageSize

	query := db.Model(&User{})

    if department != "" {
        query = query.Where("department = ?", department)
    }

    if jobPost != "" {
        query = query.Where("job_post = ?", jobPost)
    }

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("name ILIKE ? OR email ILIKE ?", likeSearch, likeSearch)
    }

    var total int64
    query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
        "page":       pageNum,
        "limit":      pageSize,
        "total":      total,
        "totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
        "data":       users,
    })
}