import { makeAutoObservable, runInAction } from 'mobx'
import RootStore from './RootStore'
import { Workout, CreateWorkoutRequest, UpdateWorkoutRequest } from '../types/Workout'

class WorkoutStore {
  rootStore: RootStore
  workouts: Workout[] = []
  currentWorkout: Workout | null = null
  isLoading = false
  error: string | null = null

  constructor(rootStore: RootStore) {
    this.rootStore = rootStore
    makeAutoObservable(this)
  }

  setWorkouts(workouts: Workout[]) {
    this.workouts = workouts
  }

  setCurrentWorkout(workout: Workout | null) {
    this.currentWorkout = workout
  }

  setLoading(loading: boolean) {
    this.isLoading = loading
  }

  setError(error: string | null) {
    this.error = error
  }

  async createWorkout(workoutData: CreateWorkoutRequest) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      // TODO: Implement workout service
      const workout: Workout = {
        id: 'temp-id',
        name: workoutData.name,
        description: workoutData.description,
        userId: 'temp-user-id',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
      
      runInAction(() => {
        this.workouts.unshift(workout)
        this.isLoading = false
      })
      return workout
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to create workout'
        this.isLoading = false
      })
      throw error
    }
  }

  get workoutsByProgram() {
    return (programId: string) => this.workouts.filter(w => w.programId === programId)
  }
}

export default WorkoutStore 