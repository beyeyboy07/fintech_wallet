import type { ReactNode } from 'react'

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