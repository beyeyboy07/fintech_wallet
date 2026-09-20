interface ErrorStateProps {
  message?: string
}

function ErrorState({
  message = 'Something went wrong.',
}: ErrorStateProps) {
  return (
    <div
      className="error-state"
      role="alert"
    >
      <p>{message}</p>
    </div>
  )
}

export default ErrorState