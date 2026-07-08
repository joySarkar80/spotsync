package zone

import (
	"spotsync/internal/domain/zone/dto"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) toResponse(z ParkingZone) (dto.ZoneResponse, error) {
	activeCount, err := s.repo.CountActiveReservations(z.ID)
	if err != nil {
		return dto.ZoneResponse{}, err
	}

	available := z.TotalCapacity - int(activeCount)
	if available < 0 {
		available = 0
	}

	return dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt.String(),
	}, nil
}

func (s *service) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.CreateZone(&zone); err != nil {
		return nil, err
	}

	resp, err := s.toResponse(zone)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *service) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.repo.GetAllZones()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		resp, err := s.toResponse(z)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *service) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}

	resp, err := s.toResponse(*zone)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
