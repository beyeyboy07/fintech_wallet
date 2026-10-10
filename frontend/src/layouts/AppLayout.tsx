
import { useState, type ReactNode } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { logout } from '../services/authService'
import { clearTokens } from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import PageContainer from '../components/PageContainer'
import Button from '../components/Button'
import '../assets/css/AppLayout.css'
import logo from "../assets/images/1.png";

interface AppLayoutProps {
  children: ReactNode
}

function AppLayout({ children }: AppLayoutProps) {
  const navigate = useNavigate()
  const clearAuth = useAuthStore((state) => state.clearAuth)
  const [loggingOut, setLoggingOut] = useState(false)
  const [menuOpen, setMenuOpen] = useState(false)

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

  function closeMenu() {
    setMenuOpen(false)
  }

  return (
    <div className="app-layout">
      <header className="app-header">
        <div className="app-header-inner">
          <NavLink
            to="/"
            className="app-logo"
            onClick={closeMenu}
          >
            <img src={logo} alt="DANQIVA" />
          </NavLink>

          <button
            type="button"
            className="app-menu-toggle"
            aria-label="Toggle navigation menu"
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen((open) => !open)}
          >
            <span />
            <span />
            <span />
          </button>

          <nav
            className={`app-nav ${menuOpen ? 'open' : ''}`}
          >
            <NavLink
              to="/"
              end
              className={({ isActive }) =>
                isActive ? 'active' : ''
              }
              onClick={closeMenu}
            >
              Dashboard
            </NavLink>

            <NavLink
              to="/profile"
              className={({ isActive }) =>
                isActive ? 'active' : ''
              }
              onClick={closeMenu}
            >
              Profile
            </NavLink>

            <NavLink
              to="/top-up"
              className={({ isActive }) =>
                isActive ? 'active' : ''
              }
              onClick={closeMenu}
            >
              Top Up
            </NavLink>

            <NavLink
              to="/transfer"
              className={({ isActive }) =>
                isActive ? 'active' : ''
              }
              onClick={closeMenu}
            >
              Transfer
            </NavLink>

            <NavLink
              to="/transactions"
              className={({ isActive }) =>
                isActive ? 'active' : ''
              }
              onClick={closeMenu}
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
        </div>
      </header>

      <main>
        <PageContainer>{children}</PageContainer>
      </main>
    </div>
  )
}

export default AppLayout
