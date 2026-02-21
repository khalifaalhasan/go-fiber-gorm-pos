package core

import "errors"

// Sentinel errors — satu-satunya error yang boleh muncul di HTTP response.
// Service layer memetakan semua DB/raw error ke salah satu di bawah ini.
var (
	ErrNotFound          = errors.New("data tidak ditemukan")
	ErrAlreadyExists     = errors.New("data sudah ada")
	ErrInsufficientStock = errors.New("stok produk tidak mencukupi")
	ErrVoucherInvalid    = errors.New("voucher tidak valid atau sudah kadaluarsa")
	ErrVoucherMinOrder   = errors.New("total pesanan tidak memenuhi minimum untuk voucher ini")
	ErrOrderAlreadyPaid  = errors.New("pesanan sudah dibayar")
	ErrInvalidSignature  = errors.New("signature webhook tidak valid")
	ErrInternalServer    = errors.New("terjadi kesalahan pada server")
)
