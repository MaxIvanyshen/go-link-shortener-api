package main

import (
	"encoding/json"
	"fmt"
	"link-shortener/dto"
	"link-shortener/link"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createLink(w http.ResponseWriter, r *http.Request) {
    var req dto.LinkReq

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }

    err = db.Where("short = ?", req.Custom).First(&link.Link{}).Error
    if err != gorm.ErrRecordNotFound {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "custom link name taken")
        return
    }

    link := link.FromDto(req)
    err = db.Create(link).Error
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    jsonData, err := json.Marshal(link)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Add("Content-Type", "application/json")
    w.Write(jsonData)
}

func getLink(w http.ResponseWriter, r *http.Request) {
    var found link.Link

    err := db.Where("short = ?", r.PathValue("short")).First(&found).Error
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
    }

    if !found.Infinite {
        if (found.ExpiresAt != 0 && found.ExpiresAt < time.Now().Unix()) || found.Usages <= 0{
            w.WriteHeader(http.StatusNotFound)
            return
        } 
        if found.Usages > 0 {
            found.Usages -= 1
        }
    }

    err = db.Save(&found).Error
    if err != nil {
        fmt.Println("Error while updating link")
        w.WriteHeader(http.StatusInternalServerError)
    }

    http.Redirect(w, r, found.Original, http.StatusSeeOther)
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

    http.HandleFunc("POST /api/link", createLink)
    http.HandleFunc("GET /{short}", getLink)
    http.ListenAndServe(":42069", nil)
}
