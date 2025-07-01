import { Routes, Route, Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { StoreProvider } from './stores/useStores'
import RootStore from './stores/RootStore'
import Login from './pages/Login'
import Register from './pages/Register'
import Layout from './components/Layout.tsx'
import Dashboard from './pages/Dashboard.tsx'
import Workouts from './pages/Workouts.tsx'
import Programs from './pages/Programs.tsx'
import Exercises from './pages/Exercises.tsx'
import Profile from './pages/Profile.tsx'

const rootStore = new RootStore()

// Protected Route component
const ProtectedRoute = observer(({ children }: { children: React.ReactNode }) => {
  const { userStore } = rootStore
  return userStore.currentUser ? <>{children}</> : <Navigate to="/login" replace />
})

const App = observer(() => {
  return (
    <StoreProvider store={rootStore}>
      <div className="min-h-screen bg-gray-50">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route
            path="/*"
            element={
              <ProtectedRoute>
                <Layout>
                  <Routes>
                    <Route path="/" element={<Dashboard />} />
                    <Route path="/workouts" element={<Workouts />} />
                    <Route path="/exercises" element={<Exercises />} />
                    <Route path="/programs" element={<Programs />} />
                    <Route path="/profile" element={<Profile />} />
                  </Routes>
                </Layout>
              </ProtectedRoute>
            }
          />
        </Routes>
      </div>
    </StoreProvider>
  )
})

export default App 