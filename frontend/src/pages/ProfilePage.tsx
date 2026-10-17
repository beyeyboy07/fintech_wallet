import { useAuthStore } from '../stores/authStore'
import PageHeader from '../components/PageHeader'
import Card from '../components/Card'
import ErrorState from '../components/ErrorState'
import BackButton from '../components/BackButton'

function ProfilePage() {
  const user = useAuthStore((state) => state.user)

  if (!user) {
    return <ErrorState message="User data not available." />
  }

  return (
    <section className="profile-page">
      <BackButton />

      <PageHeader
        title="Profile"
        description="View your account information."
      />

      {/* <Card className="profile-card"> */}
  
      <Card className="balance-card">
        <div className="balance-card-header">
          <p className="balance-label">DANQIVA Cash</p>
          <span className="balance-brand">DANQIVA</span>
        </div>

        <p>
          <strong>Name:</strong> {user.name}
        </p>

        <p>
          <strong>Email:</strong> {user.email}
        </p>

        <p>
          <strong>Role:</strong> {user.role}
        </p>

         <p>
          <strong>Account Number:</strong> {user.account_number}
        </p>
      </Card>
    </section>
  )
}

export default ProfilePage