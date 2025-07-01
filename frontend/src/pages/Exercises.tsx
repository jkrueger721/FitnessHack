import { observer } from 'mobx-react-lite'

const Exercises = observer(() => {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Exercises</h1>
        <p className="mt-1 text-sm text-gray-500">
          Browse and manage exercises
        </p>
      </div>
      
      <div className="card">
        <p className="text-gray-500">No exercises yet. Add your first exercise!</p>
      </div>
    </div>
  )
})

export default Exercises 