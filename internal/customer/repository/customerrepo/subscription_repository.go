package customerrepo

// import (
// 	"github.com/fathimasithara01/tradeverse/pkg/models"
// 	"gorm.io/gorm"
// )

// type SubscriptionRepository interface {
// 	CreateSubscription(sub *models.TraderSubscription) error
// 	GetActiveSubscription(userID, traderPlanID uint) (*models.TraderSubscription, error)
// }

// type subscriptionRepository struct {
// 	db *gorm.DB
// }

// func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
// 	return &subscriptionRepository{db: db}
// }

// func (r *subscriptionRepository) CreateSubscription(sub *models.TraderSubscription) error {
// 	return r.db.Create(sub).Error
// }

// func (r *subscriptionRepository) GetActiveSubscription(userID, traderPlanID uint) (*models.TraderSubscription, error) {
// 	var sub models.TraderSubscription
// 	err := r.db.Where("user_id = ? AND trader_subscription_plan_id = ? AND is_active = ?", userID, traderPlanID, true).
// 		First(&sub).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &sub, nil
// }
