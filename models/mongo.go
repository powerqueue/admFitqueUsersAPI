package models

//DBConfigs - struc definition
type DBConfigs1 struct {
	Host     string `json:"Host" validate:"required"`
	Port     int64  `json:"Port" validate:"required"`
	User     string `json:"User" validate:"required"`
	Password string `json:"Password" validate:"required"`
}
