package approve

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/syedomair/ex-pay-gateway/lib/models"
	"github.com/syedomair/ex-pay-gateway/lib/tools/logger"
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
func (p *postgresRepo) CreateApprove(inputApprove map[string]interface{}, merchantID int) (string, error) {
	methodName := "CreateApprove"
	p.logger.Debug(p.requestID, "M:%v start", methodName)
	start := time.Now()

	ccNumber := ""
	if ccNumberValue, ok := inputApprove["cc_number"]; ok {
		ccNumber = ccNumberValue.(string)
	}
	ccExpiry := ""
	if ccExpiryValue, ok := inputApprove["cc_expiry"]; ok {
		ccExpiry = ccExpiryValue.(string)
	}
	currency := ""
	if currencyValue, ok := inputApprove["currency"]; ok {
		currency = currencyValue.(string)
	}
	amount := ""
	if amountValue, ok := inputApprove["amount"]; ok {
		amount = amountValue.(string)
	}

	key := make([]byte, 10)
	_, _ = rand.Read(key)

	newApprove := &models.Approve{}
	newApprove.MerchantID = merchantID
	newApprove.CCNumber = ccNumber
	newApprove.CCExpiry = ccExpiry
	newApprove.Currency = currency
	newApprove.Amount = amount
	newApprove.ApproveKey = fmt.Sprintf("%X", key)
	newApprove.CreatedAt = time.Now().Format(time.RFC3339)

	if err := p.client.Create(newApprove).Error; err != nil {
		return "", err
	}
	p.logger.Debug(p.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return newApprove.ApproveKey, nil
}
