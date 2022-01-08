package approve

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type postgresRepo struct {
	client    *gorm.DB
	logger    logger.Logger
	requestID string
}

// NewPostgresRepository Public.
func NewPostgresRepository(c *gorm.DB, logger logger.Logger) Repository {
	return &postgresRepo{client: c, logger: logger, requestID: ""}
}

func (p *postgresRepo) SetRequestID(requestID string) {
	p.requestID = requestID
}

// GetMerchantID Public
func (p *postgresRepo) GetMerchantID(merchantKey string) (int, error) {
	methodName := "GetMerchantID"
	p.logger.Debug(p.requestID, "M:%v start", methodName)
	start := time.Now()

	merchant := models.Merchant{}

	if err := p.client.Table("merchant").
		Where("key = ?", merchantKey).
		Scan(&merchant).Error; err != nil {
		return 0, errors.New("merchant not found")
	}

	p.logger.Debug(p.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return merchant.ID, nil

}

// CreateApprove Public
func (p *postgresRepo) CreateApprove(approveObj *models.Approve, merchantID int, approveKey string) (*models.Approve, error) {
	methodName := "CreateApprove"
	p.logger.Debug(p.requestID, "M:%v start", methodName)
	start := time.Now()

	approveObj.MerchantID = merchantID
		approveObj.ApproveKey = approveKey
	if err := p.client.Create(approveObj).Error; err != nil {
		return &models.Approve{}, err
	}
	p.logger.Debug(p.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return approveObj, nil
}
