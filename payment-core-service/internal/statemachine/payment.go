package statemachine

type PaymentStatus string

const (
	StatusCreated    PaymentStatus = "CREATED"
	StatusProcessing PaymentStatus = "PROCESSING"
	StatusPaying     PaymentStatus = "PAYING"
	StatusSuccess    PaymentStatus = "SUCCESS"
	StatusFailed     PaymentStatus = "FAILED"
	StatusExpired    PaymentStatus = "EXPIRED"
	StatusClosed     PaymentStatus = "CLOSED"
	StatusUnknown    PaymentStatus = "UNKNOWN"
)

var paymentTransitions = map[PaymentStatus]map[PaymentStatus]bool{
	StatusCreated: {
		StatusProcessing: true,
	},
	StatusProcessing: {
		StatusPaying:  true,
		StatusSuccess: true,
		StatusFailed:  true,
		StatusUnknown: true,
	},
	StatusPaying: {
		StatusSuccess: true,
		StatusFailed:  true,
		StatusExpired: true,
		StatusUnknown: true,
	},
	StatusUnknown: {
		StatusSuccess: true,
		StatusFailed:  true,
	},
}

func CanTransit(from PaymentStatus, to PaymentStatus) bool {
	return paymentTransitions[from][to]
}

