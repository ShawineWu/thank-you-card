import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import FeedPage from './pages/Feed/FeedPage'
import CreateCardPage from './pages/CreateCard/CreateCardPage'
import MyCardsPage from './pages/MyCards/MyCardsPage'
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
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<FeedPage />} />
          <Route path="create" element={<CreateCardPage />} />
          <Route path="my-cards" element={<MyCardsPage />} />
          <Route path="stats" element={<StatsPage />} />
          <Route path="top-employees" element={<TopEmployeesPage />} />
          <Route path="company-values" element={<CompanyValuesPage />} />
          <Route path="analytics" element={<AnalyticsLayout />}>
            <Route index element={<Navigate to="dashboard" replace />} />
            <Route path="dashboard" element={<DashboardPage />} />
            <Route path="recognizers" element={<RecognizersPage />} />
            <Route path="teams" element={<TeamPatternsPage />} />
            <Route path="values" element={<ValuesDistributionPage />} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
