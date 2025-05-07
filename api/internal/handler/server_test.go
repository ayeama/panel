package handler

import (
	"testing"
)

type mockServerService struct{}

func (m *mockServerService) Create() error {
	return nil
}

func TestServerCreate(t *testing.T) {
	// service := &mockServerService{}
	// handler := NewServerHandler(service)

	// req := httptest.NewRequest(http.MethodPost, "/servers", nil)
	// rec := httptest.NewRecorder()
}

func TestServerRead(t *testing.T) {}

func TestServerUpdate(t *testing.T) {}

func TestServerDelete(t *testing.T) {}
