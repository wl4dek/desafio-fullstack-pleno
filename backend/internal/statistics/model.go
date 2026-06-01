package statistics

type NeighborhoodAlertCount struct {
	Neighborhood     string `json:"neighborhood"`
	Health           int    `json:"health"`
	Education        int    `json:"education"`
	SocialAssistance int    `json:"social_assistance"`
}

type StatisticsResponse struct {
	Statistics []NeighborhoodAlertCount `json:"statistics"`
}

type Summary struct {
	TotalChildren int            `json:"total_children"`
	Reviewed      int            `json:"reviewed"`
	PendingReview int            `json:"pending_review"`
	AlertsByArea  map[string]int `json:"alerts_by_area"`
}
