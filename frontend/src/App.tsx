import { Routes, Route } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Workouts from './pages/Workouts'
import Exercises from './pages/Exercises'
import Profile from './pages/Profile'

const App = observer(() => {
  return (
    <div className="min-h-screen bg-gray-50">
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/workouts" element={<Workouts />} />
          <Route path="/exercises" element={<Exercises />} />
          <Route path="/profile" element={<Profile />} />
        </Routes>
      </Layout>
    </div>
  )
})

export default App 