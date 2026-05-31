export interface Child {
  id: string
  name: string
  age: number
  neighborhood: string
  alert_categories: string[]
  reviewed: boolean
  reviewed_by?: string | null
  reviewed_at?: string | null
  notes: string
  created_at: string
}

export interface ChildById extends Child {
  health: Health
  social_assistance: SocialAssistance
  education: Education
}

export const AlertsCategories = {
  health: "Saúde",
  social_assistance: "Assistência Social",
  education: "Educação",
} as const

export type AlertCategoryType = keyof typeof AlertsCategories

export interface Alert {
  category: AlertCategoryType
  code: string
  message: string
}

export interface Health {
  alerts: string[]
  vaccinationsUpToDate: boolean
  lastConsultation: string
}

export interface SocialAssistance {
  alerts: string[]
  cadUnico: boolean
  activeBenefit: boolean
}

export interface Education {
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
  alert?: string
  has_alert?: string
  reviewed?: string
  page?: string
  per_page?: string
}
