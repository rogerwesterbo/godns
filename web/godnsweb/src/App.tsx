import { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Theme, Flex, Spinner } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import { ThemeProvider, useTheme, AuthProvider } from './contexts';
import { Layout, ProtectedRoute, ScrollToTop } from './components';

// Route-level code splitting — each page is its own chunk.
const LoginPage = lazy(() => import('./pages/LoginPage'));
const CallbackPage = lazy(() => import('./pages/CallbackPage'));
const DashboardPage = lazy(() => import('./pages/DashboardPage'));
const ProfilePage = lazy(() => import('./pages/ProfilePage'));
const ZonesPage = lazy(() => import('./pages/ZonesPage'));
const ZoneDetailPage = lazy(() => import('./pages/ZoneDetailPage'));
const RecordsPage = lazy(() => import('./pages/RecordsPage'));
const SearchPage = lazy(() => import('./pages/SearchPage').then(m => ({ default: m.SearchPage })));
const ExportPage = lazy(() => import('./pages/ExportPage'));
const UnauthorizedPage = lazy(() => import('./pages/UnauthorizedPage'));
const ForbiddenPage = lazy(() => import('./pages/ForbiddenPage'));
const NotFoundPage = lazy(() => import('./pages/NotFoundPage'));
const ServerErrorPage = lazy(() => import('./pages/ServerErrorPage'));

const AdminDashboardPage = lazy(() => import('./pages/admin/AdminDashboardPage'));
const CachePage = lazy(() => import('./pages/admin/CachePage'));

function RouteFallback() {
  return (
    <Flex align="center" justify="center" style={{ minHeight: '100vh' }}>
      <Spinner size="3" />
    </Flex>
  );
}

function AppContent() {
  const { theme } = useTheme();

  return (
    <Theme accentColor="blue" grayColor="slate" radius="medium" appearance={theme}>
      <AuthProvider>
        <Router>
          <ScrollToTop />
          <Suspense fallback={<RouteFallback />}>
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
              <Route
                path="/export"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <ExportPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />

              {/* Admin routes */}
              <Route
                path="/admin"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <AdminDashboardPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/admin/cache"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <CachePage />
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
          </Suspense>
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
