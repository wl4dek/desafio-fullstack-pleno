export interface Child {
  id: string
  name: string
  age: number
  neighborhood: string
  has_alert: boolean
  reviewed: boolean
  reviewed_by?: string | null
  reviewed_at?: string | null
  notes: string
  created_at: string
}

export interface Areas {
  health: Health
  socialAssistance: SocialAssistance
  education: Education
}

export interface Health {
  type?: string
  vaccinationsUpToDate: boolean
  alerts: string[]
  lastConsultation: string
}

export interface SocialAssistance {
  type?: string
  alerts: string[]
  cadUnico: boolean
  activeBenefit: boolean
}

export interface Education {
  type?: string
  alerts: string[]
  schoolName: string
  frequenciaPercent: number
}

export interface Pagination {
  page: number
  per_page: number
  total: number
  total_pages: number
}

export interface PaginatedResponse {
  data: Child[]
  pagination: Pagination
}

export interface Summary {
  total_children: number
  reviewed: number
  pending_review: number
  alerts_by_area: Record<string, number>
}

export interface LoginResponse {
  access_token: string
  token_type: string
  expires_in: number
}

export interface ChildFilters {
  childName?: string
  neighborhood?: string
  has_alert?: string
  reviewed?: string
  page?: string
  per_page?: string
}
