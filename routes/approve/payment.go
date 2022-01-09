package approve

import (
	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// Payment Interface
type Payment interface {
	ApprovePayment(*models.Approve) (string, error)
}
