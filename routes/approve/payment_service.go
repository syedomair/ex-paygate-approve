package approve

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type PaymentService struct {
	logger    logger.Logger
	requestID string
}

// NewPaymentService Public.
func NewPaymentService(logger logger.Logger) Payment {
	return &PaymentService{logger: logger}
}

// ApprovePayment Public.
func (payWrap *PaymentService) ApprovePayment(approveObj *models.Approve) (string, error) {
	methodName := "ApprovePayment"
	payWrap.logger.Debug(payWrap.requestID, "M:%v start", methodName)
	start := time.Now()

	if approveObj.CCNumber == "4000000000000119" {
		return "", errors.New("authorisation failure")
	}

	key := make([]byte, 10)
	_, _ = rand.Read(key)

	payWrap.logger.Debug(payWrap.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return fmt.Sprintf("%X", key), nil
}
