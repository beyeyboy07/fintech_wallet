import { useAuthStore } from '../stores/authStore'
import PageHeader from '../components/PageHeader'
import Card from '../components/Card'
import ErrorState from '../components/ErrorState'

function ProfilePage() {
  const user = useAuthStore((state) => state.user)

  if (!user) {
    return <ErrorState message="User data not available." />
  }

  return (
    <section className="profile-page">
      <PageHeader
        title="Profile"
        description="View your account information."
      />

      <Card className="profile-card">
        <p>
          <strong>Name:</strong> {user.name}
        </p>

        <p>
          <strong>Email:</strong> {user.email}
        </p>

        <p>
          <strong>Role:</strong> {user.role}
        </p>
      </Card>
    </section>
  )
}

export default ProfilePage