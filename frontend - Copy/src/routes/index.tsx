import { createBrowserRouter } from 'react-router-dom'
import AppLayout from '../layouts/AppLayout'
import DashboardPage from '../pages/DashboardPage'
import LoginPage from '../pages/LoginPage'
import ProtectedRoute from './ProtectedRoute'
import RegisterPage from '../pages/RegisterPage'
import ProfilePage from '../pages/ProfilePage'
import TopUpPage from '../pages/TopUpPage'
import TransferPage from '../pages/TransferPage'
import TransactionsPage from '../pages/TransactionsPage'

const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <AppLayout>
          <DashboardPage />
        </AppLayout>
      </ProtectedRoute>
    ),
  },
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/profile',
    element: (
      <ProtectedRoute>
        <AppLayout>
          <ProfilePage />
        </AppLayout>
      </ProtectedRoute>
    ),
  },
  {
    path: '/top-up',
    element: (
      <ProtectedRoute>
        <AppLayout>
          <TopUpPage />
        </AppLayout>
      </ProtectedRoute>
    ),
  },

  {
    path: '/transfer',
    element: (
      <ProtectedRoute>
        <AppLayout>
          <TransferPage />
        </AppLayout>
      </ProtectedRoute>
    ),
  },

  {
  path: '/transactions',
    element: (
      <ProtectedRoute>
        <AppLayout>
          <TransactionsPage />
        </AppLayout>
      </ProtectedRoute>
    ),
  },
])

export default router