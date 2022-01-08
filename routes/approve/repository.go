package approve

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Repository interface
type Repository interface {
	SetRequestID(requestID string)
	GetMerchantID(merchantKey string) (int, error)
	CreateApprove(approveModel *models.Approve, merchantID int, approveKey string) (*models.Approve, error)
}
