package orders

import "errors"

var (
	ErrOrderNotFound              = errors.New("orden no encontrada")
	ErrProductNotFound            = errors.New("producto no encontrado")
	ErrInvalidQuantity            = errors.New("cantidad inválida")
	ErrIdempotencyKeyConflict     = errors.New("idempotency key ya usada")
	ErrProductsServiceUnavailable = errors.New("servicio de productos no disponible")
	ErrEmptyOrder                 = errors.New("la orden no puede estar vacía")
	ErrOrderNotPending            = errors.New("la orden no está pendiente")
)
