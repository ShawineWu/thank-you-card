import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import LoginPage from './pages/Login/LoginPage'
import FeedPage from './pages/Feed/FeedPage'
import CreateCardPage from './pages/CreateCard/CreateCardPage'
import StatsPage from './pages/Stats/StatsPage'
import TopEmployeesPage from './pages/TopEmployees/TopEmployeesPage'
import AnalyticsLayout from './pages/Analytics/AnalyticsLayout'
import DashboardPage from './pages/Analytics/DashboardPage'
import RecognizersPage from './pages/Analytics/RecognizersPage'
import TeamPatternsPage from './pages/Analytics/TeamPatternsPage'
import ValuesDistributionPage from './pages/Analytics/ValuesDistributionPage'
import CompanyValuesPage from './pages/CompanyValues/CompanyValuesPage'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/" element={<Layout />}>
            <Route index element={<ProtectedRoute><FeedPage /></ProtectedRoute>} />
            <Route path="create" element={<ProtectedRoute><CreateCardPage /></ProtectedRoute>} />
            <Route path="my-cards" element={<Navigate to="/stats" replace />} />
            <Route path="stats" element={<ProtectedRoute><StatsPage /></ProtectedRoute>} />
            <Route path="top-employees" element={<ProtectedRoute><TopEmployeesPage /></ProtectedRoute>} />
            <Route path="company-values" element={<ProtectedRoute><CompanyValuesPage /></ProtectedRoute>} />
            <Route path="analytics" element={<ProtectedRoute requireHR><AnalyticsLayout /></ProtectedRoute>}>
              <Route index element={<Navigate to="dashboard" replace />} />
              <Route path="dashboard" element={<DashboardPage />} />
              <Route path="recognizers" element={<RecognizersPage />} />
              <Route path="teams" element={<TeamPatternsPage />} />
              <Route path="values" element={<ValuesDistributionPage />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
