import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { CallbackPage } from "./pages/Callback";
import { HomePage } from "./pages/Home";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* OAuth callback route - no authentication required */}
        <Route path="/callback" element={<CallbackPage />} />

        {/* Protected routes - authentication required */}
        <Route
          path="/*"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
