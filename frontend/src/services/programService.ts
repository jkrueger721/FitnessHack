import api from './api'
import { Program, CreateProgramRequest, UpdateProgramRequest } from '../types/Program'

export const programService = {
  async createProgram(programData: CreateProgramRequest): Promise<Program> {
    const response = await api.post('/programs', programData)
    return response.data
  },

  async getProgram(id: string): Promise<Program> {
    const response = await api.get(`/programs/${id}`)
    return response.data
  },

  async listPrograms(): Promise<Program[]> {
    const response = await api.get('/programs')
    return response.data
  },

  async updateProgram(id: string, programData: UpdateProgramRequest): Promise<Program> {
    const response = await api.put(`/programs/${id}`, programData)
    return response.data
  },

  async deleteProgram(id: string): Promise<void> {
    await api.delete(`/programs/${id}`)
  },

  async getProgramWorkouts(programId: string): Promise<any[]> {
    const response = await api.get(`/programs/${programId}/workouts`)
    return response.data
  }
} 