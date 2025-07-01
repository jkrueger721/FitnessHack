import { makeAutoObservable, runInAction } from 'mobx'
import RootStore from './RootStore'
import { Program, CreateProgramRequest, UpdateProgramRequest } from '../types/Program'
import { programService } from '../services/programService'

class ProgramStore {
  rootStore: RootStore
  programs: Program[] = []
  currentProgram: Program | null = null
  isLoading = false
  error: string | null = null

  constructor(rootStore: RootStore) {
    this.rootStore = rootStore
    makeAutoObservable(this)
  }

  setPrograms(programs: Program[]) {
    this.programs = programs
  }

  setCurrentProgram(program: Program | null) {
    this.currentProgram = program
  }

  setLoading(loading: boolean) {
    this.isLoading = loading
  }

  setError(error: string | null) {
    this.error = error
  }

  async createProgram(programData: CreateProgramRequest) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      const program = await programService.createProgram(programData)
      runInAction(() => {
        this.programs.unshift(program)
        this.isLoading = false
      })
      return program
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to create program'
        this.isLoading = false
      })
      throw error
    }
  }

  async loadPrograms() {
    this.setLoading(true)
    this.setError(null)
    
    try {
      const programs = await programService.listPrograms()
      runInAction(() => {
        this.programs = programs
        this.isLoading = false
      })
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to load programs'
        this.isLoading = false
      })
    }
  }

  async loadProgram(id: string) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      const program = await programService.getProgram(id)
      runInAction(() => {
        this.currentProgram = program
        this.isLoading = false
      })
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to load program'
        this.isLoading = false
      })
    }
  }

  async updateProgram(id: string, programData: UpdateProgramRequest) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      const updatedProgram = await programService.updateProgram(id, programData)
      runInAction(() => {
        const index = this.programs.findIndex(p => p.id === id)
        if (index !== -1) {
          this.programs[index] = updatedProgram
        }
        if (this.currentProgram?.id === id) {
          this.currentProgram = updatedProgram
        }
        this.isLoading = false
      })
      return updatedProgram
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to update program'
        this.isLoading = false
      })
      throw error
    }
  }

  async deleteProgram(id: string) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      await programService.deleteProgram(id)
      runInAction(() => {
        this.programs = this.programs.filter(p => p.id !== id)
        if (this.currentProgram?.id === id) {
          this.currentProgram = null
        }
        this.isLoading = false
      })
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to delete program'
        this.isLoading = false
      })
    }
  }

  get activePrograms() {
    return this.programs.filter(p => p.isActive)
  }

  get programById() {
    return (id: string) => this.programs.find(p => p.id === id)
  }
}

export default ProgramStore 