package reservation

import (
	"spotsync/internal/domain/reservation/dto"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Reserve(userID uint, req dto.CreateReservationRequest) (*dto.CreateReservationResponse, error) {
	res, err := s.repo.CreateReservation(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.CreateReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.String(),
		UpdatedAt:    res.UpdatedAt.String(),
	}, nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.repo.GetMyReservations(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneSummary{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.String(),
		})
	}
	return responses, nil
}

// CancelReservation: driver nijer reservation-i cancel korte parbe,
// admin je kono reservation cancel korte parbe.
func (s *service) CancelReservation(callerID uint, callerRole string, reservationID uint) error {
	res, err := s.repo.GetReservationByID(reservationID)
	if err != nil {
		return err
	}

	if res.UserID != callerID && callerRole != "admin" {
		return ErrForbidden
	}

	return s.repo.CancelReservation(reservationID)
}

func (s *service) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.repo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserSummary{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneSummary{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.String(),
		})
	}
	return responses, nil
}
