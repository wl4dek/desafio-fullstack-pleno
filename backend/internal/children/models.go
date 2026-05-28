package children

import "time"

type Child struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Age          int        `json:"age"`
	Neighborhood string     `json:"neighborhood"`
	HasAlert     bool       `json:"has_alert"`
	Reviewed     bool       `json:"reviewed"`
	ReviewedBy   *string    `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Health struct {
	VaccinationsUpToDate bool       `json:"vaccinationsUpToDate"`
	Alerts               []string   `json:"alerts"`
	LastConsultation     *time.Time `json:"lastConsultation"`
}

type SocialAssistance struct {
	Alerts        []string `json:"alerts"`
	CadUnico      bool     `json:"cadUnico"`
	ActiveBenefit bool     `json:"activeBenefit"`
}

type Education struct {
	Alerts            []string `json:"alerts"`
	SchoolName        *string  `json:"schoolName"`
	FrequenciaPercent int      `json:"frequenciaPercent"`
}

type Areas struct {
	Health           Health           `json:"health"`
	SocialAssistance SocialAssistance `json:"socialAssistance"`
	Education        Education        `json:"education"`
}

type PaginatedResponse struct {
	Data       []Child    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type Summary struct {
	TotalChildren int            `json:"total_children"`
	Reviewed      int            `json:"reviewed"`
	PendingReview int            `json:"pending_review"`
	AlertsByArea  map[string]int `json:"alerts_by_area"`
}
