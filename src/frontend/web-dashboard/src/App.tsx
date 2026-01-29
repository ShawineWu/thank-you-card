import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { CallbackPage } from "./pages/Callback";
import { Dashboard } from "./pages/Dashboard";
import { Feed } from "./pages/Feed";
import { MainLayout } from "./components/Layout/MainLayout";
import { Profile } from "./pages/Profile";
import { Settings } from "./pages/Settings";
import { HRAnalytics } from "./pages/HRAnalytics";
import { CompanyValues } from "./pages/CompanyValues";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* OAuth callback route - no authentication required */}
        <Route path="/callback" element={<CallbackPage />} />

        {/* Protected routes wrapped in MainLayout */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <MainLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="feed" element={<Feed />} />
          <Route path="profile" element={<Profile />} />
          <Route path="values" element={<CompanyValues />} />
          <Route path="analytics" element={<HRAnalytics />} />
          <Route path="settings" element={<Settings />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
