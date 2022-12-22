package main

import (
	"log"
  "net/http"
  "strconv"
  "strings"
  "github.com/gin-gonic/gin"
  "time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TodoItem struct {
  Id int `json:"id" gorm:"column:id;"`
  Title string `json:"title" gorm:"column:title;"`
  Status string `json:"status" gorm:"column:status"`
  CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
  UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;"`
}
func (TodoItem) TableName() string {
  return "todo_items"
}
func main() {
	dsn := "root:passroot@tcp(localhost:3303)/todo_db?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("cannot connect to mysql", err)
	}
	log.Println("Conected", db)
  router := gin.Default()
  v1 := router.Group("/v1")
  {
    v1.POST("/items", createItem(db))
    v1.GET("/items", getListOfItems(db))
    v1.GET("/items/:id", readItemById(db))
    v1.PUT("/items/:id", editItemById(db))
    v1.DELETE("/items/:id", deleteItemById(db))
    v1.GET("/check-health", checkHealth())
  }
  router.Run()
}

func createItem(db *gorm.DB) gin.HandlerFunc {
  return func(c *gin.Context){
    var dataItem TodoItem
    if err := c.ShouldBind(&dataItem); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    dataItem.Title = strings.TrimSpace(dataItem.Title)
    if dataItem.Title == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "title can't blank"})
      return
    }
    dataItem.Status = "Doing"
    if err:=db.Create(&dataItem).Error;err!=nil{
      c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
      return
    }
    c.JSON(http.StatusOK,gin.H{"data": dataItem.Id})
  }
}

func readItemById(db *gorm.DB) gin.HandlerFunc {
  return func(c *gin.Context){
    var dataItem TodoItem
    id, err := strconv.Atoi(c.Param("id"))
    if err!=nil{
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    if err:=db.Where("id=?", id).First(&dataItem).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
      return 
    }
    c.JSON(http.StatusOK, gin.H{"data": dataItem})
  }
}

func getListOfItems(db *gorm.DB) gin.HandlerFunc {
  return func(c *gin.Context){
    type DataPaging struct {
      Page int `json:"page" form:"page"`
      Limit int `json:"limit" form:"limit"`
      Total int64 `json:"total" form:"-"`
    }
    var paging DataPaging
    if err:=c.ShouldBind(&paging);err!=nil{
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    if paging.Page <=0 {
      paging.Page = 1
    }
    if paging.Limit<=0 {
      paging.Limit = 10
    }
    offset := (paging.Page -1)*paging.Limit
    var result []TodoItem
    if err:=db.Table(TodoItem{}.TableName()).
      Count(&paging.Total).
      Offset(offset).
      Order("id desc").
      Find(&result).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": result})
  }
}
func editItemById(db *gorm.DB) gin.HandlerFunc{
  return func(c *gin.Context){
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil{
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    var dataItem TodoItem
    if err:=c.ShouldBind(&dataItem); err!=nil {
      c.JSON(http.StatusBadRequest, gin.H{"erorr": err.Error()})
      return
    }
    if err:=db.Where("id=?", id).Updates(&dataItem).Error;err!=nil{
      c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
      return
    }
    c.JSON(http.StatusOK, gin.H{"data":true})
  }

}
func deleteItemById(db *gorm.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
      id, err := strconv.Atoi(c.Param("id"))
      if err!=nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      }
      if  err:=db.Table(TodoItem{}.TableName()).
        Where("id = ?", id).
        Delete(nil).Error; err!=nil{
          c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
          return
      }
      c.JSON(http.StatusOK, gin.H{"data": true})
    }
}

func checkHealth() gin.HandlerFunc {
  return func(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{"data": true})
  }
}
