export interface Program {
  id: string
  name: string
  description?: string
  userId: string
  durationWeeks?: number
  difficulty: 'beginner' | 'intermediate' | 'advanced'
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateProgramRequest {
  name: string
  description?: string
  durationWeeks?: number
  difficulty: 'beginner' | 'intermediate' | 'advanced'
}

export interface UpdateProgramRequest {
  name?: string
  description?: string
  durationWeeks?: number
  difficulty?: 'beginner' | 'intermediate' | 'advanced'
  isActive?: boolean
} 