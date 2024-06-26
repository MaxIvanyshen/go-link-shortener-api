package link

import (
	"math/rand"
    "link-shortener/dto"

	"gorm.io/gorm"
)

var linkLength = 5;

type Link struct {
    gorm.Model
    Original string `json:"original"`
    Short string    `json:"short"`
    Usages uint     `json:"usages"`
    Infinite bool   `json:"infinite"`
    ExpiresAt int64 `json:"expires"`
}

func NewLink(link string) *Link {
    return &Link{Original: link, Short: generateLink()}
}

func FromDto(dto dto.LinkReq) *Link {
    short := generateLink()
    if dto.Custom != "" {
        short = dto.Custom
    }
    return &Link{
        Original: dto.Original,
        Short: short,
        Infinite: dto.Infinite,
        Usages: dto.Usages, 
        ExpiresAt: dto.ExpiresAt,
    }
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateLink() string {
    b := make([]rune, linkLength)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}
