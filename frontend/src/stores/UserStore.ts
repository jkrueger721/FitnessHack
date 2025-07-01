import { makeAutoObservable, runInAction } from 'mobx'
import { RootStore } from './RootStore'
import { User } from '../types/User'
import { userService } from '../services/userService'

class UserStore {
  rootStore: RootStore
  currentUser: User | null = null
  isLoading = false
  error: string | null = null

  constructor(rootStore: RootStore) {
    this.rootStore = rootStore
    makeAutoObservable(this)
  }

  setCurrentUser(user: User | null) {
    this.currentUser = user
  }

  setLoading(loading: boolean) {
    this.isLoading = loading
  }

  setError(error: string | null) {
    this.error = error
  }

  async login(email: string, password: string) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      const user = await userService.login(email, password)
      runInAction(() => {
        this.currentUser = user
        this.isLoading = false
      })
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Login failed'
        this.isLoading = false
      })
    }
  }

  async register(userData: Partial<User>) {
    this.setLoading(true)
    this.setError(null)
    
    try {
      const user = await userService.register(userData)
      runInAction(() => {
        this.currentUser = user
        this.isLoading = false
      })
    } catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : 'Registration failed'
        this.isLoading = false
      })
    }
  }

  logout() {
    this.currentUser = null
    this.error = null
  }
}

export default UserStore 