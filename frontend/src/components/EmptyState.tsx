import '../assets/css/EmptyState.css'

interface EmptyStateProps {
  message?: string
}

function EmptyState({
  message = 'No data found.',
}: EmptyStateProps) {
  return (
    <div
      className="empty-state"
      role="status"
    >
      <p>{message}</p>
    </div>
  )
}

export default EmptyState