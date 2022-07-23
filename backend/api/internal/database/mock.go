package database

import (
	"github.com/stretchr/testify/mock"
)

// MockDb mocks of ReadWrite
type MockDb struct {
	mock.Mock
}

func (m *MockDb) Read(source, orderBy string, limit int) ([]map[string]interface{}, UpdatedAt, error) {
	args := m.Called(source, orderBy, limit)
	data, ok := args.Get(0).([]map[string]interface{})
	if !ok {
		data = nil
	}
	return data, args.Get(1).(UpdatedAt), args.Error(2)
}

func (m *MockDb) Write(source string, models ...Model) error {
	args := m.Called(source, models)
	return args.Error(0)
}

func (m *MockDb) Close() {
	m.Called()
}
