export interface Workout {
  id: string
  name: string
  description?: string
  userId: string
  programId?: string
  createdAt: string
  updatedAt: string
}

export interface CreateWorkoutRequest {
  name: string
  description?: string
  programId?: string
}

export interface UpdateWorkoutRequest {
  name?: string
  description?: string
} 