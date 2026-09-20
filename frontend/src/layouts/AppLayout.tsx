import { useState, type ReactNode } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { logout } from '../services/authService'
import { clearTokens } from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import PageContainer from '../components/PageContainer'
import Button from '../components/Button'






interface AppLayoutProps {
  children: ReactNode
}

function AppLayout({ children }: AppLayoutProps) {
  const navigate = useNavigate()
  const clearAuth = useAuthStore((state) => state.clearAuth)
  const [loggingOut, setLoggingOut] = useState(false)

  async function handleLogout() {
  if (loggingOut) {
    return
  }

  setLoggingOut(true)
    try {
      await logout()
    } finally {
      clearTokens()
      clearAuth()
      navigate('/login', { replace: true })
    }
  }

  return (
    <div>
      <header className="app-header">
        <h1 className="app-logo">Fintech Wallet</h1>

        <nav className="app-nav">
          <Link to="/">Dashboard</Link>
          <Link to="/profile">Profile</Link>
          <Link to="/top-up">Top Up</Link>
          <Link to="/transfer">Transfer</Link>
          <Link to="/transactions">Transactions</Link>

          <Button
            type="button"
            variant="secondary"
            loading={loggingOut}
            onClick={handleLogout}
          >
            {loggingOut ? 'Logging out...' : 'Logout'}
          </Button>
        </nav>
      </header>

      <main>
        <PageContainer>{children}</PageContainer>
      </main>
    </div>
  )
}

export default AppLayout