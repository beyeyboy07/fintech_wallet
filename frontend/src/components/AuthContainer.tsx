import type { ReactNode } from 'react'
import '../assets/css/AuthContainer.css'



interface AuthContainerProps {
  children: ReactNode
}

function AuthContainer({
  children,
}: AuthContainerProps) {
  return (
    <div className="auth-container">
      {children}
    </div>
  )
}

export default AuthContainer