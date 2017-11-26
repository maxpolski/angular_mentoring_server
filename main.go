package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/*Course represents course*/
type Course struct {
	ID                uint   `json:"courseId" gorm:"AUTO_INCREMENT;primary_key"`
	CourseName        string `json:"courseName"`
	CourseType        string `json:"courseType"`
	CourseEndDate     uint64 `json:"courseEndDate"`
	CourseStartDate   uint64 `json:"courseStartDate"`
	CourseDescription string `json:"courseDescription,omitempty"`
}

var db *gorm.DB
var err error

func initializeDB() {
	db, err = gorm.Open("mysql", "angular_mentor:123@/angular_mentoring?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		handleError(nil, err)
	}

	db.AutoMigrate(&Course{})
}

func getIndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func getCourseHandler(c *gin.Context) {
	var course Course
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&course).Error; err == nil {
		c.JSON(http.StatusOK, course)
	} else {
		handleError(c, err)
	}
}

func addCourseHandler(c *gin.Context) {
	var course Course
	if err := c.ShouldBindJSON(&course); err == nil {
		db.Create(&course)
		c.JSON(200, course)
	} else {
		handleError(c, err)
	}
}

func updateCourseHandler(c *gin.Context) {
	var course Course
	id := c.Param("id")

	if err := db.Where("id = ?", id).First(&course).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(404)
	}

	if err := c.ShouldBindJSON(&course); err == nil {
		db.Save(&course)
		c.JSON(http.StatusOK, course)
	} else {
		handleError(c, err)
	}
}

func handleError(c *gin.Context, err error) {
	if c != nil {
		c.String(http.StatusBadGateway, err.Error())
	}

	fmt.Println(err.Error())
}

func main() {
	initializeDB()
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFS("/static", http.Dir("/Users/maksymkorabelskyi/dev/js/build/angular-mentoring/dist"))

	router.GET("/", getIndexHandler)
	router.GET("/courses/:id", getCourseHandler)
	router.POST("/courses", addCourseHandler)
	router.PUT("/courses", updateCourseHandler)

	router.Run(":8080")
}
