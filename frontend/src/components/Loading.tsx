import '../assets/css/Loading.css'

interface LoadingProps {
  message?: string
}

function Loading({
  message = 'Loading...',
}: LoadingProps) {
  return (
    <div
        className="loading"
        role="status"
    >
      <span>{message}</span>
    </div>
  )
}

export default Loading