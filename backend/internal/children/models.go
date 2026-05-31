package children

import (
	"database/sql"
	"time"
)

type Child struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Age             int            `json:"age"`
	Neighborhood    string         `json:"neighborhood"`
	AlertCategories []string `json:"alert_categories"`
	Reviewed        bool           `json:"reviewed"`
	ReviewedBy      *string        `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time     `json:"reviewed_at,omitempty"`
	Notes           sql.NullString `json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
}

type ChildById struct {
	Child
	Health           Health           `json:"health"`
	Education        Education        `json:"education"`
	SocialAssistance SocialAssistance `json:"social_assistance"`
}

type Alerts struct {
	Category string         `json:"category"`
	Code     string         `json:"code"`
	Message  sql.NullString `json:"message"`
}

type ChildResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Age             int        `json:"age"`
	Neighborhood    string     `json:"neighborhood"`
	AlertCategories []string   `json:"alert_categories"`
	Reviewed        bool       `json:"reviewed"`
	ReviewedBy      *string    `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ChildByIdResponse struct {
	ChildResponse
	Health           Health           `json:"health"`
	Education        Education        `json:"education"`
	SocialAssistance SocialAssistance `json:"social_assistance"`
}

type Health struct {
	VaccinationsUpToDate bool       `json:"vaccinationsUpToDate"`
	LastConsultation     *time.Time `json:"lastConsultation"`
	Alerts               []string   `json:"alerts"`
}

type SocialAssistance struct {
	CadUnico      bool     `json:"cadUnico"`
	ActiveBenefit bool     `json:"activeBenefit"`
	Alerts        []string `json:"alerts"`
}

type Education struct {
	SchoolName        *string  `json:"schoolName"`
	FrequenciaPercent int      `json:"frequenciaPercent"`
	Alerts            []string `json:"alerts"`
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
