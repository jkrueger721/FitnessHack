export interface User {
  id: string
  email: string
  username: string
  firstName?: string
  lastName?: string
  createdAt: string
  updatedAt: string
}

export interface CreateUserRequest {
  email: string
  username: string
  password: string
  firstName?: string
  lastName?: string
}

export interface LoginRequest {
  email: string
  password: string
} 