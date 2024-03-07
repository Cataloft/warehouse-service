package errors

const (
	WarehouseError             = "unable to SUM(amount): Warehouse data not found"
	WarehouseAvailabilityError = "unable to reserve good : Warehouse is not available"
	GoodError                  = "unable to reserve good: Current amount of available warehouse will be < 0"
)
