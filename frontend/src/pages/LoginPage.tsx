
import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login, getMe } from '../services/authService'
import { setTokens } from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import { getApiErrorMessage } from '../utils/apiError'
import ErrorState from '../components/ErrorState'
import FormField from '../components/FormField'
import AuthContainer from '../components/AuthContainer'
import FormActions from '../components/FormActions'
import Button from '../components/Button'
import Input from '../components/Input'
import '../assets/css/LoginPage.css'
import logo from "../assets/images/1.png";

function LoginPage() {
  const navigate = useNavigate()
  const setUser = useAuthStore((state) => state.setUser)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [emailError, setEmailError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()
    setEmailError('')
    setPasswordError('')

    if (loading) {
      return
    }

    if (!email.trim()) {
      setEmailError('Email wajib diisi.')
    }

    if (!password) {
      setPasswordError('Password wajib diisi.')
    }

    if (!email.trim() || !password) {
      return
    }

    setError('')
    setLoading(true)

    try {
      const tokenResponse = await login({
        email,
        password,
      })

      setTokens(
        tokenResponse.data.access_token,
        tokenResponse.data.refresh_token,
      )

      const user = await getMe()

      setUser(user)
      navigate('/')
    } catch (error) {
      setError(
        getApiErrorMessage(
          error,
          'Login gagal. Periksa email dan password.',
        ),
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="login-screen">
      {/* Branding Panel */}
      <aside className="login-brand-panel">
        <div className="brand-content">
          <Link to="/" className="brand-logo">
            <img src={logo} alt="DANQIVA Logo" />
            <span>DANQIVA</span>
          </Link>

          <div className="brand-message">
            <span className="brand-eyebrow">
              YOUR FINANCIAL JOURNEY
            </span>

            <h1>
              Take control
              <br />
              of your <span>financial future.</span>
            </h1>

            <p>
              Manage your money, track your transactions,
              and stay on top of your finances.
              Everything you need, in one place.
            </p>
          </div>

          <div className="brand-visual">
            <div className="visual-glow" />

            <div className="visual-card">
              <div className="visual-card-header">
                <span>YOUR WALLET</span>
                <span className="visual-card-dot" />
              </div>

              <div className="visual-card-label">
                Total Balance
              </div>

              <div className="visual-card-balance">
                <span>Rp</span> 12.850.000
              </div>

              <div className="visual-card-footer">
                <div>
                  <span className="visual-small-label">
                    ACCOUNT STATUS
                  </span>
                  <strong>
                    <span className="status-dot" />
                    Active
                  </strong>
                </div>

                <div className="visual-card-symbol">
                  D
                </div>
              </div>
            </div>

            <div className="floating-badge">
              <span className="floating-badge-icon">↗</span>
              <div>
                <strong>Smart Finance</strong>
                <span>Made simple</span>
              </div>
            </div>

            <div className="visual-decoration decoration-one" />
            <div className="visual-decoration decoration-two" />
          </div>
        </div>

        <div className="brand-footer">
          <span>© {new Date().getFullYear()} DANQIVA</span>
          <span>Manage smarter. Live better.</span>
        </div>
      </aside>

      {/* Login Form Panel */}
      <section className="login-form-panel">
        <div className="mobile-brand-logo">
          <span className="brand-logo-icon">D</span>
          <span>DANQIVA</span>
        </div>

        <AuthContainer>
          <section className="login-page">
            <div className="login-heading">
              <span className="login-eyebrow">
                WELCOME BACK
              </span>

              <h2>Welcome back.</h2>

              <p className="auth-description">
                Sign in to access your wallet and manage
                your finances.
              </p>
            </div>

            <form onSubmit={handleSubmit}>
              <FormField
                label="Email Address"
                htmlFor="email"
                error={emailError}
                required
              >
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(event) => {
                    setEmail(event.target.value)
                    setEmailError('')
                  }}
                  maxLength={150}
                  placeholder="Enter your email"
                  autoComplete="email"
                  aria-invalid={Boolean(emailError)}
                  aria-describedby="email-error"
                  required
                />
              </FormField>

              <FormField
                label="Password"
                htmlFor="password"
                error={passwordError}
                required
              >
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(event) => {
                    setPassword(event.target.value)
                    setPasswordError('')
                  }}
                  maxLength={72}
                  placeholder="Enter your password"
                  autoComplete="current-password"
                  aria-invalid={Boolean(passwordError)}
                  aria-describedby="password-error"
                  required
                />
              </FormField>

              {error && <ErrorState message={error} />}

              <FormActions>
                <Button
                  type="submit"
                  variant="primary"
                  loading={loading}
                >
                  {loading ? 'Logging in...' : 'Sign In'}
                </Button>
              </FormActions>
            </form>

            <div className="login-divider">
              <span />
              <span>SECURE ACCESS</span>
              <span />
            </div>

            <p className="auth-link">
              Don't have an account?{' '}
              <Link to="/register">Create an account</Link>
            </p>

            <div className="login-security-note">
              <span className="security-icon">✓</span>
              Your account is protected with secure authentication.
            </div>
          </section>
        </AuthContainer>

        <div className="login-panel-footer">
          <span>Powered by DANQIVA</span>
          <span>Secure Financial Management</span>
        </div>
      </section>
    </main>
  )
}

export default LoginPage
