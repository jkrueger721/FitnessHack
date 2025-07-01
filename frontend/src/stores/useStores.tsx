import { useContext, createContext } from 'react'
import RootStore from './RootStore'

const StoreContext = createContext<RootStore | null>(null)

export const StoreProvider = ({ children, store }: { children: React.ReactNode; store: RootStore }) => {
  return (
    <StoreContext.Provider value={store}>
      {children}
    </StoreContext.Provider>
  )
}

export const useStores = () => {
  const store = useContext(StoreContext)
  if (!store) {
    throw new Error('useStores must be used within a StoreProvider')
  }
  return store
} 