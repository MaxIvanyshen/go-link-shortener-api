package main

import (
	//"encoding/json"
	"fmt"
	"link-shortener/dto"
	"link-shortener/link"
	"net/http"
	"os"
	"time"
    "strconv"

    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createLink(c *gin.Context) {
    var req dto.LinkReq

    c.BindJSON(&req)

    err := db.Where("short = ?", req.Custom).First(&link.Link{}).Error
    if err != gorm.ErrRecordNotFound {
        c.Status(http.StatusConflict)
        c.Writer.Write([]byte("custom link name taken"))
        return
    }

    link := link.FromDto(req)
    err = db.Create(link).Error
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return
    }

    c.IndentedJSON(http.StatusOK, link)
}

func getAllLinks(c *gin.Context) {
    var allLinks []link.Link
    db.Where("").Scopes(Paginate(c.Request)).Find(&allLinks)
    c.IndentedJSON(http.StatusOK, allLinks)
}

func getLink(c *gin.Context) {
    var found link.Link

    err := db.Where("short = ?", c.Param("short")).First(&found).Error
    if err != nil {
        c.Status(http.StatusNotFound)
        return
    }

    if !found.Infinite {
        if (found.ExpiresAt != 0 && found.ExpiresAt < time.Now().Unix()) || found.Usages <= 0 {
            c.Status(http.StatusNotFound)
            return
        } 
        if found.Usages > 0 {
            found.Usages -= 1
        }
    }

    err = db.Save(&found).Error
    if err != nil {
        fmt.Println("Error while updating link")
        c.Status(http.StatusInternalServerError)
        return
    }

    c.Redirect(http.StatusMovedPermanently, found.Original)
}


func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
  return func (db *gorm.DB) *gorm.DB {
    q := r.URL.Query()
    page, _ := strconv.Atoi(q.Get("page"))
    if page <= 0 {
      page = 1
    }

    pageSize, _ := strconv.Atoi(q.Get("page_size"))
    switch {
    case pageSize > 100:
      pageSize = 100
    case pageSize <= 0:
      pageSize = 10
    }

    offset := (page - 1) * pageSize
    return db.Offset(offset).Limit(pageSize)
  }
}

var db *gorm.DB

func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("couldn't parse .env")
    }

    db, err = gorm.Open(postgres.Open(os.Getenv("POSTGRES_CONN_STRING")), &gorm.Config{})
    if err != nil {
        fmt.Println("couldn't connect to db")
    }

    db.AutoMigrate(&link.Link{})

    router := gin.Default();

    router.POST("/api/link", createLink)
    router.GET("/api/link", getAllLinks)
    router.GET("/:short", getLink)

    router.Run("localhost:42069")
}
