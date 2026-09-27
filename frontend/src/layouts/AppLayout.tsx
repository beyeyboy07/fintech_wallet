import { useState, type ReactNode } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { logout } from '../services/authService'
import { clearTokens } from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import PageContainer from '../components/PageContainer'
import Button from '../components/Button'
import '../assets/css/AppLayout.css'





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
    <div className="app-layout">
      <header className="app-header">
        <h1 className="app-logo">Fintech Wallet</h1>

        <nav className="app-nav">
          <NavLink
            to="/"
            end
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Dashboard
          </NavLink>

          <NavLink
            to="/profile"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Profile
          </NavLink>

          <NavLink
            to="/top-up"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Top Up
          </NavLink>

          <NavLink
            to="/transfer"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Transfer
          </NavLink>

          <NavLink
            to="/transactions"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Transactions
          </NavLink>

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