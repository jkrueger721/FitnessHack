import { makeAutoObservable } from 'mobx'
import UserStore from './UserStore'
import WorkoutStore from './WorkoutStore'
import ExerciseStore from './ExerciseStore'

class RootStore {
  userStore: UserStore
  workoutStore: WorkoutStore
  exerciseStore: ExerciseStore

  constructor() {
    this.userStore = new UserStore(this)
    this.workoutStore = new WorkoutStore(this)
    this.exerciseStore = new ExerciseStore(this)
    makeAutoObservable(this)
  }
}

export default RootStore 