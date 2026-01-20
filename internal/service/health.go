package service

import (
	"context"
	"time"

	"github.com/emresahna/url-shortener-app/internal/models"
)

// HealthCheck returns the health status of the service
func (s *service) HealthCheck(ctx context.Context) (models.HealthResponse, error) {
	return models.HealthResponse{
		Status:    "ok",
		Version:   "1.0.0", // This should be fetched from a version package in a real app
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// ReadinessCheck checks if the service is ready to accept requests
func (s *service) ReadinessCheck(ctx context.Context) (models.HealthResponse, error) {
	// In a real application, you would check database connections,
	// cache availability, and other dependencies

	// Check database connection
	if s.db == nil {
		return models.HealthResponse{}, models.InternalServerErr()
	}

	// Check Redis cache connection
	if s.rcc == nil {
		return models.HealthResponse{}, models.InternalServerErr()
	}

	return models.HealthResponse{
		Status:    "ready",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// LivenessCheck checks if the service is alive
func (s *service) LivenessCheck(ctx context.Context) (models.HealthResponse, error) {
	return models.HealthResponse{
		Status:    "alive",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
