package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"fmt"
	"strconv"
)

func main() {
	db := Database()
	db.AutoMigrate(&Todo{})

	router := gin.Default()

	v1 := router.Group("/api/v1/todos")

	{
		v1.POST("/", CreateTodo)
		v1.GET("/", FetchAllTodo)
		v1.GET("/:id", FetchSingleTodo)
		v1.PUT("/:id", UpdateTodo)
		v1.DELETE("/:id", DeleteTodo)
	}

	router.Run()
}

func CreateTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := Todo{Title: c.PostForm("title"), Completed: completed}
	db := Database()
	db.Save(&todo)

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
	fmt.Printf("Received completed=%s and title=%s", c.PostForm("completed"), c.PostForm("title"))
}
func FetchAllTodo(c *gin.Context) {
	var todos []Todo
	db := Database()
	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": mapTodosToDtos(todos)})
}

func mapTodosToDtos(todos []Todo) (result []TransformedTodo) {
	for _, todo := range todos {
		completed := false
		if todo.Completed == 1 {
			completed = true
		}
		result = append(result, TransformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed})
	}
	return
}

func FetchSingleTodo(c *gin.Context) {

}
func UpdateTodo(c *gin.Context) {

}
func DeleteTodo(c *gin.Context) {

}

func Database() *gorm.DB {
	db, err := gorm.Open("mysql", "gotest:1234@tcp(db:3306)/demo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed int `json:"completed"`
}

type TransformedTodo struct {
	ID        uint `json:"id"`
	Title     string `json:"title"`
	Completed bool `json:"completed"`
}
