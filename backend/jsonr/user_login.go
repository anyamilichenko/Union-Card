package jsonr

type UserLoginJson struct {
	Email    string `form:"email" json:"email"  binding:"required"`
	Password string `form:"password" json:"password"  binding:"required"`
}
