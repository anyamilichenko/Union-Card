package jsonr

type CreateTokensJson struct {
	//DeviceId     string `form:"device_id" json:"device_id" binding:"required"`
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}
