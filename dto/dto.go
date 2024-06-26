package dto

type LinkReq struct {
    Original string `json:"original"`
    Custom  string  `json:"custom"`
    Infinite bool   `json:"infinite"`    
    Usages uint     `json:"usages"`    
    ExpiresAt int64 `json:"expires"` 
}
