
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { register } from '../services/authService'
import { getApiErrorMessage } from '../utils/apiError'
import ErrorState from '../components/ErrorState'
import FormField from '../components/FormField'
import AuthContainer from '../components/AuthContainer'
import FormActions from '../components/FormActions'
import Button from '../components/Button'
import Input from '../components/Input'
import '../assets/css/LoginPage.css'
import '../assets/css/RegisterPage.css'
import logo from '../assets/images/1.png'

function RegisterPage() {
  const navigate = useNavigate()

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [nameError, setNameError] = useState('')
  const [emailError, setEmailError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()

    setNameError('')
    setEmailError('')
    setPasswordError('')

    if (loading) return

    if (name.trim().length < 3) {
      setNameError('Name minimal 3 karakter.')
    }

    if (!email.trim()) {
      setEmailError('Email wajib diisi.')
    }

    if (password.length < 8) {
      setPasswordError('Password minimal 8 karakter.')
    }

    if (
      name.trim().length < 3 ||
      !email.trim() ||
      password.length < 8
    ) {
      return
    }

    setError('')
    setLoading(true)

    try {
      await register({
        name,
        email,
        password,
      })

      navigate('/login')
    } catch (error) {
      setError(
        getApiErrorMessage(
          error,
          'Register gagal. Periksa data yang dimasukkan.',
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

      {/* Register Form Panel */}
      <section className="login-form-panel">
        {/* <div className="mobile-brand-logo">
          <img src={logo} alt="DANQIVA Logo" />
          <span>DANQIVA</span>
        </div> */}

        <AuthContainer>
          <section className="login-page register-page">
            <div className="login-heading">
              <h2>Create your account.</h2>

              <p className="auth-description">
                Join DANQIVA and start managing your finances
                in one place.
              </p>
            </div>

            <form onSubmit={handleSubmit}>
              <FormField
                label="Full Name"
                htmlFor="name"
                error={nameError}
                required
              >
                <Input
                  id="name"
                  type="text"
                  value={name}
                  onChange={(event) => {
                    setName(event.target.value)
                    setNameError('')
                  }}
                  maxLength={100}
                  placeholder="Enter your full name"
                  autoComplete="name"
                  aria-invalid={Boolean(nameError)}
                  aria-describedby="name-error"
                  required
                />
              </FormField>

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
                  placeholder="Create your password"
                  autoComplete="new-password"
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
                  {loading ? 'Creating account...' : 'Create Account'}
                </Button>
              </FormActions>
            </form>

            <div className="login-divider">
              <span />
              <span>SECURE REGISTRATION</span>
              <span />
            </div>

            <p className="auth-link">
              Already have an account?{' '}
              <Link to="/login">Sign in</Link>
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

export default RegisterPage
