import { makeAutoObservable } from 'mobx'
import UserStore from '../stores/UserStore'
import WorkoutStore from '../stores/WorkoutStore'
import ExerciseStore from '../stores/ExerciseStore'
import ProgramStore from '../stores/ProgramStore'

class RootStore {
  userStore: UserStore
  workoutStore: WorkoutStore
  exerciseStore: ExerciseStore
  programStore: ProgramStore

  constructor() {
    this.userStore = new UserStore(this)
    this.workoutStore = new WorkoutStore(this)
    this.exerciseStore = new ExerciseStore(this)
    this.programStore = new ProgramStore(this)
    makeAutoObservable(this)
  }
}

export default RootStore 