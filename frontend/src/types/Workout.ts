export interface Workout {
  id: string
  name: string
  description?: string
  userId: string
  createdAt: string
  updatedAt: string
}

export interface CreateWorkoutRequest {
  name: string
  description?: string
}

export interface UpdateWorkoutRequest {
  name?: string
  description?: string
} 