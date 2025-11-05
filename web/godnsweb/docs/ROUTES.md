# GoDNS Web Routes

## Available Routes

| Route | Component | Description | Layout |
|-------|-----------|-------------|--------|
| `/` | DashboardPage | Main dashboard with statistics and activity | Yes |
| `/zones` | ZonesPage | DNS zones management | Yes |
| `/records` | RecordsPage | DNS records management | Yes |
| `/profile` | ProfilePage | User profile and settings | Yes |
| `/login` | LoginPage | Login page with "Login" button | No |
| `/unauthorized` | UnauthorizedPage | 401 - User not authenticated | No |
| `/forbidden` | ForbiddenPage | 403 - User lacks permissions | No |
| `/server-error` | ServerErrorPage | 500 - Server error occurred | No |
| `*` | NotFoundPage | 404 - Page not found (catch-all) | No |

## Layout Structure

Routes with layout include:
- **Header**: Search bar, theme toggle, user menu
- **Sidebar**: Navigation menu with Dashboard, Zones, Records, Profile
- **Main Content**: Page-specific content

## Navigation Examples

```tsx
// Navigate to pages
navigate('/');           // Dashboard
navigate('/zones');      // Zones
navigate('/records');    // Records
navigate('/profile');    // Profile
navigate('/login');      // Login

// Navigate to error pages
navigate('/unauthorized');
navigate('/forbidden');
navigate('/server-error');
```

## Protected Routes

All routes with layout (Dashboard, Zones, Records, Profile) should be protected:

```tsx
// Future implementation
<Route 
  path="/" 
  element={
    <ProtectedRoute>
      <Layout><DashboardPage /></Layout>
    </ProtectedRoute>
  } 
/>
```

## Testing the Routes

Once the dev server is running, you can test the routes by navigating to:

- http://localhost:14200/ - Dashboard (with full layout)
- http://localhost:14200/zones - Zones page
- http://localhost:14200/records - Records page
- http://localhost:14200/profile - Profile page
- http://localhost:14200/login - Login page
- http://localhost:14200/unauthorized - Unauthorized page
- http://localhost:14200/forbidden - Forbidden page
- http://localhost:14200/server-error - Server error page
- http://localhost:14200/any-invalid-path - 404 Not Found page

## Sidebar Navigation

The sidebar includes these menu items:
- 🏠 Dashboard (/)
- 🌐 Zones (/zones)
- 📄 Records (/records)
- 👤 Profile (/profile)

Active routes are highlighted with accent color.
