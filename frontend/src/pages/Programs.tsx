import { useState, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Plus, Edit, Trash2, Dumbbell } from 'lucide-react'
import { useStores } from '../stores/useStores'
import { CreateProgramRequest } from '../types/Program'
import { CreateWorkoutRequest, Workout } from '../types/Workout'

const Programs = observer(() => {
  const { programStore, workoutStore } = useStores()
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [showWorkoutForm, setShowWorkoutForm] = useState(false)
  const [selectedProgramId, setSelectedProgramId] = useState<string | null>(null)
  const [formData, setFormData] = useState<CreateProgramRequest>({
    name: '',
    description: '',
    durationWeeks: undefined,
    difficulty: 'beginner'
  })
  const [workoutFormData, setWorkoutFormData] = useState<CreateWorkoutRequest>({
    name: '',
    description: ''
  })

  useEffect(() => {
    programStore.loadPrograms()
  }, [programStore])

  const handleCreateProgram = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await programStore.createProgram(formData)
      setFormData({ name: '', description: '', durationWeeks: undefined, difficulty: 'beginner' })
      setShowCreateForm(false)
    } catch (error) {
      console.error('Failed to create program:', error)
    }
  }

  const handleCreateWorkout = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedProgramId) return
    
    try {
      await workoutStore.createWorkout({
        ...workoutFormData,
        programId: selectedProgramId
      })
      setWorkoutFormData({ name: '', description: '' })
      setShowWorkoutForm(false)
      setSelectedProgramId(null)
    } catch (error) {
      console.error('Failed to create workout:', error)
    }
  }

  const handleDeleteProgram = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this program?')) {
      await programStore.deleteProgram(id)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Programs</h1>
          <p className="mt-1 text-sm text-gray-500">
            Manage your fitness programs and workouts
          </p>
        </div>
        <button
          onClick={() => setShowCreateForm(true)}
          className="btn-primary flex items-center space-x-2"
        >
          <Plus size={16} />
          <span>Create Program</span>
        </button>
      </div>

      {/* Create Program Form */}
      {showCreateForm && (
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Create New Program</h3>
          <form onSubmit={handleCreateProgram} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Name</label>
              <input
                type="text"
                required
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="input-field"
                placeholder="Enter program name"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="input-field"
                rows={3}
                placeholder="Enter program description"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Duration (weeks)</label>
                <input
                  type="number"
                  value={formData.durationWeeks || ''}
                  onChange={(e) => setFormData({ ...formData, durationWeeks: e.target.value ? parseInt(e.target.value) : undefined })}
                  className="input-field"
                  placeholder="e.g., 12"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Difficulty</label>
                <select
                  value={formData.difficulty}
                  onChange={(e) => setFormData({ ...formData, difficulty: e.target.value as any })}
                  className="input-field"
                >
                  <option value="beginner">Beginner</option>
                  <option value="intermediate">Intermediate</option>
                  <option value="advanced">Advanced</option>
                </select>
              </div>
            </div>
            <div className="flex space-x-3">
              <button type="submit" className="btn-primary" disabled={programStore.isLoading}>
                {programStore.isLoading ? 'Creating...' : 'Create Program'}
              </button>
              <button
                type="button"
                onClick={() => setShowCreateForm(false)}
                className="btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Programs List */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {programStore.programs.map((program) => (
          <div key={program.id} className="card">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="text-lg font-medium text-gray-900">{program.name}</h3>
                <p className="text-sm text-gray-500">{program.difficulty}</p>
              </div>
              <div className="flex space-x-2">
                <button
                  onClick={() => {
                    setSelectedProgramId(program.id)
                    setShowWorkoutForm(true)
                  }}
                  className="p-1 text-gray-400 hover:text-primary-600"
                  title="Add Workout"
                >
                  <Plus size={16} />
                </button>
                <button
                  onClick={() => handleDeleteProgram(program.id)}
                  className="p-1 text-gray-400 hover:text-red-600"
                  title="Delete Program"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
            
            {program.description && (
              <p className="text-sm text-gray-600 mb-3">{program.description}</p>
            )}
            
            <div className="flex items-center justify-between text-sm text-gray-500">
              <span>{program.durationWeeks ? `${program.durationWeeks} weeks` : 'No duration set'}</span>
              <span className={`px-2 py-1 rounded-full text-xs ${
                program.isActive ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
              }`}>
                {program.isActive ? 'Active' : 'Inactive'}
              </span>
            </div>

            {/* Workouts in this program */}
            <div className="mt-4 pt-4 border-t border-gray-200">
              <h4 className="text-sm font-medium text-gray-700 mb-2 flex items-center">
                <Dumbbell size={14} className="mr-1" />
                Workouts
              </h4>
              {workoutStore.workoutsByProgram(program.id).length > 0 ? (
                <ul className="space-y-1">
                  {workoutStore.workoutsByProgram(program.id).map((workout: Workout) => (
                    <li key={workout.id} className="text-sm text-gray-600">
                      {workout.name}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-sm text-gray-400">No workouts yet</p>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* Create Workout Form */}
      {showWorkoutForm && selectedProgramId && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Add Workout to Program</h3>
            <form onSubmit={handleCreateWorkout} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Workout Name</label>
                <input
                  type="text"
                  required
                  value={workoutFormData.name}
                  onChange={(e) => setWorkoutFormData({ ...workoutFormData, name: e.target.value })}
                  className="input-field"
                  placeholder="Enter workout name"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Description</label>
                <textarea
                  value={workoutFormData.description}
                  onChange={(e) => setWorkoutFormData({ ...workoutFormData, description: e.target.value })}
                  className="input-field"
                  rows={3}
                  placeholder="Enter workout description"
                />
              </div>
              <div className="flex space-x-3">
                <button type="submit" className="btn-primary" disabled={workoutStore.isLoading}>
                  {workoutStore.isLoading ? 'Creating...' : 'Create Workout'}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowWorkoutForm(false)
                    setSelectedProgramId(null)
                  }}
                  className="btn-secondary"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {programStore.error && (
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <p className="text-sm text-red-600">{programStore.error}</p>
        </div>
      )}
    </div>
  )
})

export default Programs 