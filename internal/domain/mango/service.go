package mango

import "spotsync/internal/domain/mango/dto"

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateMango(req dto.CreateRequest) (*dto.Response, error) {
	mango := Mango{
		Name:        req.Name,
		Description: req.Description,
		Variety:     req.Variety,
		PricePerKg:  req.PricePerKg,
		StockKg:     req.StockKg,
		ImageURL:    req.ImageURL,
	}

	if err := s.repo.Create(&mango); err != nil {
		return nil, err
	}

	return mango.ToResponse(), nil
}

func (s *service) GetMangoes() ([]dto.Response, error) {
	mangoes, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.Response
	for _, m := range mangoes {
		responses = append(responses, *m.ToResponse())
	}
	return responses, nil
}

func (s *service) GetMangoByID(mangoId uint) (*dto.Response, error) {
	mango, err := s.repo.GetByID(mangoId)
	if err != nil {
		return nil, err
	}
	return mango.ToResponse(), nil
}

func (s *service) UpdateMango(mangoId uint, req dto.UpdateRequest) (*dto.Response, error) {
	mango, err := s.repo.GetByID(mangoId)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		mango.Name = req.Name
	}
	if req.Description != "" {
		mango.Description = req.Description
	}
	if req.Variety != "" {
		mango.Variety = req.Variety
	}
	if req.PricePerKg != 0 {
		mango.PricePerKg = req.PricePerKg
	}
	if req.StockKg != 0 {
		mango.StockKg = req.StockKg
	}
	if req.ImageURL != "" {
		mango.ImageURL = req.ImageURL
	}

	if err := s.repo.Update(mango); err != nil {
		return nil, err
	}

	return mango.ToResponse(), nil
}
