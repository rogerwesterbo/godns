import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Theme } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import { ThemeProvider, useTheme, AuthProvider } from './contexts';
import { Layout, ProtectedRoute, ScrollToTop } from './components';

// Pages
import {
  LoginPage,
  CallbackPage,
  DashboardPage,
  ProfilePage,
  ZonesPage,
  ZoneDetailPage,
  RecordsPage,
  SearchPage,
  UnauthorizedPage,
  ForbiddenPage,
  NotFoundPage,
  ServerErrorPage,
} from './pages';

function AppContent() {
  const { theme } = useTheme();

  return (
    <Theme accentColor="blue" grayColor="slate" radius="medium" appearance={theme}>
      <AuthProvider>
        <Router>
          <ScrollToTop />
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/callback" element={<CallbackPage />} />

            {/* Protected routes with layout */}
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Layout>
                    <DashboardPage />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <Layout>
                    <ProfilePage />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/zones"
              element={
                <ProtectedRoute>
                  <Layout>
                    <ZonesPage />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/zones/:domain"
              element={
                <ProtectedRoute>
                  <Layout>
                    <ZoneDetailPage />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/records"
              element={
                <ProtectedRoute>
                  <Layout>
                    <RecordsPage />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/search"
              element={
                <ProtectedRoute>
                  <Layout>
                    <SearchPage />
                  </Layout>
                </ProtectedRoute>
              }
            />

            {/* Error routes */}
            <Route path="/unauthorized" element={<UnauthorizedPage />} />
            <Route path="/forbidden" element={<ForbiddenPage />} />
            <Route path="/server-error" element={<ServerErrorPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </Router>
      </AuthProvider>
    </Theme>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
}

export default App;
