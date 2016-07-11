package services

import (
	"github.com/stretchr/testify/mock"
)

type MockBitmarkPay struct {
	mock.Mock
}

func (mock *MockBitmarkPay) Encrypt(bitmarkPayType BitmarkPayType) error {
	// mock.Called(bitmarkPayType)
	return nil
}

func (mock *MockBitmarkPay) Info(bitmarkPayType BitmarkPayType) error {
	return nil
}

func (mock *MockBitmarkPay) Pay(bitmarkPayType BitmarkPayType) error {
	return nil
}

func (mock *MockBitmarkPay) Restore(bitmarkPayType BitmarkPayType) error {
	return nil
}

func (mock *MockBitmarkPay) Status(hashString string) string {
	return "running"
}

func (mock *MockBitmarkPay) Kill() error {
	return nil
}

func (mock *MockBitmarkPay) GetBitmarkPayJobHash() string {
	return "12345678"
}

func (mock *MockBitmarkPay) GetBitmarkPayJobResult(bitmarkPayType BitmarkPayType) ([]byte, error) {
	return []byte("test"), nil
}

func (mock *MockBitmarkPay) GetBitmarkPayJobType(hashString string) string {
	return "info"
}
