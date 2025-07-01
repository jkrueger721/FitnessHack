import { observer } from 'mobx-react-lite'

const Profile = observer(() => {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Profile</h1>
        <p className="mt-1 text-sm text-gray-500">
          Manage your account settings
        </p>
      </div>
      
      <div className="card">
        <p className="text-gray-500">Profile settings coming soon!</p>
      </div>
    </div>
  )
})

export default Profile 