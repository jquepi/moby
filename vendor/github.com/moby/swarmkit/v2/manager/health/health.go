// Package health provides some utility functions to health-check a server. The implementation
// is based on protobuf. Users need to write their own implementations if other IDLs are used.
//
// See original source: https://github.com/grpc/grpc-go/blob/master/health/health.go
//
// We use our own implementation of grpc server health check to include the authorization
// wrapper necessary for the Managers.
package health

import (
	"context"
	"sync"

	"github.com/moby/swarmkit/v2/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents a Health Check server to check
// if a service is running or not on some host.
type Server struct {
	mu sync.Mutex
	// statusMap stores the serving status of the services this HealthServer monitors.
	statusMap map[string]api.HealthCheckResponse_ServingStatus
}

// NewHealthServer creates a new health check server for grpc services.
func NewHealthServer() *Server {
	return &Server{
		statusMap: make(map[string]api.HealthCheckResponse_ServingStatus),
	}
}

// Check checks if the grpc server is healthy and running.
func (s *Server) Check(ctx context.Context, in *api.HealthCheckRequest) (*api.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if in.Service == "" {
		// check the server overall health status.
		return &api.HealthCheckResponse{
			Status: api.HealthCheckResponse_SERVING,
		}, nil
	}
	if status, ok := s.statusMap[in.Service]; ok {
		return &api.HealthCheckResponse{
			Status: status,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "unknown service")
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *Server) SetServingStatus(service string, status api.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	s.statusMap[service] = status
	s.mu.Unlock()
}
