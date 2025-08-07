package service

import (
	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/utils"
)

type AdminService struct {
	Repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{Repo: repo}
}

func (s *AdminService) Register(admin models.Admin) (models.Admin, error) {
	hashedPwd, _ := utils.HashPassword(admin.Password)
	admin.Password = hashedPwd
	return s.Repo.Create(admin)
}

func (s *AdminService) Login(email, password string) (string, error) {
	admin, err := s.Repo.GetByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, admin.Password) {
		return "", err
	}
	return utils.GenerateJWT(admin.ID, admin.Email, "admin")
}
