package admin

import (
	"fmt"
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/admin"
)

type Admin interface {
	GetUsers() ([]models.User, error)
	UpgradeUser(user_id int) error
	Downgrade(user_id int) error
	ReportPost(postID int, userID int, reason string) error
	GetReports() ([]models.Report, error)
	RequestRole(user_id int) error
	ApproveRequest(user_id int) error
	RejectRequest(user_id int) error
	GetRequests() ([]models.User, error)
	CheckRequest(user_id int) (bool, error)
}

type adminService struct {
	adminRepo admin.Admin
}

func NewAdminService(adminRepo admin.Admin) *adminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (s *adminService) GetUsers() ([]models.User, error) {
	return s.adminRepo.GetUsers()
}

func (s *adminService) UpgradeUser(user_id int) error {
	if err := s.adminRepo.UpgradeUser(user_id); err != nil {
		return fmt.Errorf("failed to upgrade user: %w", err)
	}
	return nil
}

func (s *adminService) Downgrade(user_id int) error {
	if err := s.adminRepo.DowngradeUser(user_id); err != nil {
		return fmt.Errorf("failed to upgrade user: %w", err)
	}
	return nil
}

func (s *adminService) ReportPost(postID int, userID int, reason string) error {
	if err := s.adminRepo.ReportPost(postID, userID, reason); err != nil {
		return fmt.Errorf("failed to report post: %w", err)
	}
	return nil
}

func (s *adminService) GetReports() ([]models.Report, error) {
	reports, err := s.adminRepo.GetReports()
	if err != nil {
		return reports, fmt.Errorf("failed to retrieve reports: %w", err)
	}
	return reports, nil
}

func (s *adminService) RequestRole(user_id int) error {
	if err := s.adminRepo.RequestRole(user_id); err != nil {
		return fmt.Errorf("failed to request role: %w", err)
	}
	return nil
}

func (s *adminService) ApproveRequest(user_id int) error {
	if err := s.adminRepo.ApproveRequest(user_id); err != nil {
		return fmt.Errorf("failed to approve request: %w", err)
	}
	return nil
}

func (s *adminService) RejectRequest(user_id int) error {
	if err := s.adminRepo.RejectRequest(user_id); err != nil {
		return fmt.Errorf("failed to reject request: %w", err)
	}
	return nil
}

func (s *adminService) GetRequests() ([]models.User, error) {
	requests, err := s.adminRepo.GetRequests()
	if err != nil {
		return requests, fmt.Errorf("failed to retrieve requests: %w", err)
	}
	return requests, nil
}

func (s *adminService) CheckRequest(user_id int) (bool, error) {
	ok, err := s.adminRepo.CheckRequest(user_id)
	if err != nil {
		return ok, fmt.Errorf("failed to check request: %w", err)
	}
	return ok, nil
}
