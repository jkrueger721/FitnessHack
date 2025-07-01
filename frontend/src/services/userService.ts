import api from './api'
import { User, CreateUserRequest, LoginRequest } from '../types/User'

export const userService = {
  async login(email: string, password: string): Promise<User> {
    const response = await api.post('/v1/auth/login', { email, password })
    const { data } = response.data
    localStorage.setItem('authToken', data.token)
    return data.user
  },

  async register(userData: CreateUserRequest): Promise<User> {
    const response = await api.post('/v1/users', userData)
    const { data } = response.data
    return data
  },

  async getCurrentUser(): Promise<User> {
    const response = await api.get('/v1/users/me')
    return response.data
  },

  async updateProfile(userData: Partial<User>): Promise<User> {
    const response = await api.put('/v1/users/profile', userData)
    return response.data
  },

  async logout(): Promise<void> {
    localStorage.removeItem('authToken')
  }
} 