import { makeAutoObservable, runInAction } from 'mobx'
import RootStore from './RootStore'
import { Exercise, CreateExerciseRequest, UpdateExerciseRequest } from '../types/Exercise'

class ExerciseStore {
  rootStore: RootStore
  exercises: Exercise[] = []
  currentExercise: Exercise | null = null
  isLoading = false
  error: string | null = null

  constructor(rootStore: RootStore) {
    this.rootStore = rootStore
    makeAutoObservable(this)
  }

  setExercises(exercises: Exercise[]) {
    this.exercises = exercises
  }

  setCurrentExercise(exercise: Exercise | null) {
    this.currentExercise = exercise
  }

  setLoading(loading: boolean) {
    this.isLoading = loading
  }

  setError(error: string | null) {
    this.error = error
  }

  async loadExercises() {
    this.setLoading(true)
    this.setError(null)
    
    try {
      // TODO: Implement exercise service
      const exercises: Exercise[] = []
      runInAction(() => {
        this.exercises = exercises
        this.isLoading = false
      })
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Failed to load exercises'
        this.isLoading = false
      })
    }
  }
}

export default ExerciseStore 