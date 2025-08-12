package repository

// import (
// 	"errors"

// 	"github.com/fathimasithara01/tradeverse/models" // Make sure this path is correct
// 	"gorm.io/gorm"
// )

// // type UserRepository struct {
// // 	DB *gorm.DB
// // }

// func NewCustomerRepository(db *gorm.DB) *UserRepository {
// 	return &UserRepository{DB: db}
// }

// // FindAll retrieves all customers from the database.
// func (r *UserRepository) FindAll() ([]models.CustomerProfile, error) {
// 	var customers []models.CustomerProfile
// 	// Order by ID to ensure a consistent list
// 	if err := r.DB.Order("id asc").Find(&customers).Error; err != nil {
// 		return nil, err
// 	}
// 	return customers, nil
// }

// func (r *UserRepository) Create(customer *models.CustomerProfile) error {
// 	return r.DB.Create(customer).Error
// }

// func (r *UserRepository) FindByID(id uint) (models.CustomerProfile, error) {
// 	var customer models.CustomerProfile
// 	if err := r.DB.First(&customer, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return models.CustomerProfile{}, errors.New("customer not found")
// 		}
// 		return models.CustomerProfile{}, err
// 	}
// 	return customer, nil
// }

// // Update saves the changes to an existing customer
// func (r *UserRepository) Update(customer *models.CustomerProfile) error {
// 	return r.DB.Save(customer).Error
// }

// // Delete removes a customer from the database by their ID
// func (r *UserRepository) Delete(id uint) error {
// 	return r.DB.Delete(&models.CustomerProfile{}, id).Error
// }
