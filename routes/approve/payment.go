package approve

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

type Payment interface {
	ApprovePayment(*models.Approve) (string, error)
}

type PaymentWrapper struct {
	logger    logger.Logger
	requestID string
}

// NewPayment Public.
func NewPayment(logger logger.Logger) Payment {
	return &PaymentWrapper{logger: logger}
}

// ApprovePayment Public.
func (payWrap *PaymentWrapper) ApprovePayment(*models.Approve) (string, error) {
	methodName := "ApprovePayment"
	payWrap.logger.Debug(payWrap.requestID, "M:%v start", methodName)
	start := time.Now()

	key := make([]byte, 10)
	_, _ = rand.Read(key)

	payWrap.logger.Debug(payWrap.requestID, "M:%v ts %+v", methodName, time.Since(start))
	return fmt.Sprintf("%X", key), nil
}
