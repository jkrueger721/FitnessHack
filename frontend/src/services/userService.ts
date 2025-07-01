import api from './api'
import { User, CreateUserRequest, LoginRequest } from '../types/User'

export const userService = {
  async login(credentials: LoginRequest): Promise<User> {
    const response = await api.post('/auth/login', credentials)
    const { user, token } = response.data
    localStorage.setItem('authToken', token)
    return user
  },

  async register(userData: CreateUserRequest): Promise<User> {
    const response = await api.post('/auth/register', userData)
    const { user, token } = response.data
    localStorage.setItem('authToken', token)
    return user
  },

  async getCurrentUser(): Promise<User> {
    const response = await api.get('/auth/me')
    return response.data
  },

  async updateProfile(userData: Partial<User>): Promise<User> {
    const response = await api.put('/users/profile', userData)
    return response.data
  },

  async logout(): Promise<void> {
    localStorage.removeItem('authToken')
    await api.post('/auth/logout')
  }
} 