export interface Exercise {
  id: string
  name: string
  description?: string
  muscleGroup: string
  equipment?: string
  difficulty: 'beginner' | 'intermediate' | 'advanced'
  createdAt: string
  updatedAt: string
}

export interface CreateExerciseRequest {
  name: string
  description?: string
  muscleGroup: string
  equipment?: string
  difficulty: 'beginner' | 'intermediate' | 'advanced'
}

export interface UpdateExerciseRequest {
  name?: string
  description?: string
  muscleGroup?: string
  equipment?: string
  difficulty?: 'beginner' | 'intermediate' | 'advanced'
} 