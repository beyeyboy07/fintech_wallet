
import { useNavigate } from 'react-router-dom'
import '../assets/css/BackButton.css'

function BackButton() {
  const navigate = useNavigate()

  function handleBack() {
    navigate('/')
  }

  return (
    <button
      type="button"
      className="back-button"
      onClick={handleBack}
      aria-label="Back to dashboard"
    >
      <span className="back-button-icon" aria-hidden="true">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M15 18l-6-6 6-6"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </span>

      <span>Back</span>
    </button>
  )
}

export default BackButton
