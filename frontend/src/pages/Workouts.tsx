import { observer } from 'mobx-react-lite'

const Workouts = observer(() => {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Workouts</h1>
        <p className="mt-1 text-sm text-gray-500">
          Manage your workout routines
        </p>
      </div>
      
      <div className="card">
        <p className="text-gray-500">No workouts yet. Create your first workout!</p>
      </div>
    </div>
  )
})

export default Workouts 