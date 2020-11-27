package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/powerqueue/fitque-users-api/models"
)

//ILoginService - service interface
type ILoginService interface {
	GetLogin(location string, username string, memberID string, loginDate *time.Time) ([]*models.LoginDefinition, error)
	CreateLogin(Login models.LoginDefinition) (*models.LoginDefinition, error)
	TermLogin(AddressHierarchy *models.LoginDefinition) (*models.LoginDefinition, error)
}

//LoginService - Service Repo struct
type LoginService struct {
	LoginRepo models.ILoginRepository
}

//NewLoginService - service getter
func NewLoginService(loginRepo models.ILoginRepository) *LoginService {
	return &LoginService{loginRepo}
}

//ValidateCreateLoginRequest - service method definition
func (loginSvc *LoginService) ValidateCreateLoginRequest(CreateLoginRequest *models.LoginDefinition) error {

	/*put here any validations*/

	if len(CreateLoginRequest.LocationID) <= 0 {
		return errors.New("The Create Login Request Location is required to be defined")
	}
	return nil
}

//GetLogin - service method definition
func (loginSvc *LoginService) GetLogin(Login *models.LoginDefinition) ([]*models.LoginDefinition, error) {
	if err := loginSvc.ValidateCreateLoginRequest(Login); err != nil {
		return nil, err
	}

	return loginSvc.LoginRepo.GetLogin(Login.LocationID, Login.UserName, Login.MemberID, Login.EfctvStartDt)
}

//CreateLogin - service method definition
func (loginSvc *LoginService) CreateLogin(Login *models.LoginDefinition) (*models.LoginDefinition, error) {
	fmt.Println("Inside CreateLogin Service method definition")
	if err := loginSvc.ValidateCreateLoginRequest(Login); err != nil {
		return nil, err
	}

	return loginSvc.LoginRepo.CreateLogin(Login)
}

//TermLogin - service method definition
func (loginSvc *LoginService) TermLogin(Login *models.LoginDefinition) (*models.LoginDefinition, error) {
	if err := loginSvc.ValidateCreateLoginRequest(Login); err != nil {
		return nil, err
	}

	return loginSvc.LoginRepo.TermLogin(Login)
}
