package dto

type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateReservationResponse — POST /reservations response (flat shape, spec #6)
type CreateReservationResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	ZoneID       uint   `json:"zone_id"`
	LicensePlate string `json:"license_plate"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// MyReservationResponse — GET /reservations/my-reservations (zone nested, spec #7)
type MyReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    string      `json:"created_at"`
}

// AdminReservationResponse — GET /reservations (admin, user+zone nested, spec #9)
type AdminReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	User         UserSummary `json:"user"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    string      `json:"created_at"`
}
