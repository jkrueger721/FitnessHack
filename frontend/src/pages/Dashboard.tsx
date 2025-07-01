import { observer } from 'mobx-react-lite'

const Dashboard = observer(() => {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-sm text-gray-500">
          Welcome to your fitness tracking dashboard
        </p>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Recent Workouts</h3>
          <p className="mt-2 text-sm text-gray-500">No recent workouts</p>
        </div>
        
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Total Exercises</h3>
          <p className="mt-2 text-sm text-gray-500">0 exercises</p>
        </div>
        
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">This Week</h3>
          <p className="mt-2 text-sm text-gray-500">0 workouts completed</p>
        </div>
      </div>
    </div>
  )
})

export default Dashboard 