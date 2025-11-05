# Admin Dashboard Features

## Overview
The GoDNS admin dashboard provides a comprehensive interface for managing DNS zones and records with a modern, responsive design.

## Key Features

### 🎨 Theme Support
- **Dark Mode (Default)**: The application starts in dark mode by default
- **Light Mode**: Toggle to light mode using the sun/moon icon in the header
- **Persistent**: Theme preference is saved to localStorage

### 🔍 Animated Search
- **Compact Mode**: 300px search bar in the header
- **Expanded Mode**: Expands to 600px when focused
- **Quick Results**: Shows filtered results as you type
- **Smart Closing**: Closes when clicking outside
- **Animations**: Smooth transitions with slide-down effects

### 📱 Responsive Layout
- **Header**: Sticky header with search and user menu
- **Sidebar**: Collapsible navigation menu
  - Dashboard
  - Zones
  - Records
  - Profile
- **Mobile Optimized**: Icons-only sidebar on small screens

### 👤 User Features
- **Profile Page**: View user information and statistics
- **Avatar Menu**: Quick access to profile and logout
- **User Stats**: Display zones, records, and activity metrics

### 📊 Dashboard
- **Statistics Cards**: Quick overview of zones, records, and activity
- **Recent Activity**: Timeline of recent changes
- **Visual Metrics**: Color-coded status indicators

### 🌐 Zone Management
- **Zone List**: Table view of all DNS zones
- **Zone Types**: Primary/Secondary badge indicators
- **Record Count**: Number of records per zone
- **Status Badges**: Visual status indicators
- **Add Zone**: Quick action button

### 📝 Record Management
- **Record List**: Comprehensive view of all DNS records
- **Type Filtering**: Filter by record type (A, AAAA, CNAME, MX, NS, TXT, etc.)
- **Record Details**: Name, type, value, and TTL displayed
- **Add Record**: Quick action button
- **Sortable**: Easy to scan and manage

### 🔐 Error Pages
- **401 Unauthorized**: Not authenticated
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Page doesn't exist
- **500 Server Error**: Server-side issues

## Components Structure

```
src/
├── components/
│   ├── Header.tsx          # Top navigation bar
│   ├── Header.css
│   ├── Sidebar.tsx         # Side navigation menu
│   ├── Sidebar.css
│   ├── SearchBar.tsx       # Animated search component
│   ├── SearchBar.css
│   ├── Layout.tsx          # Main layout wrapper
│   └── index.ts           # Component exports
├── contexts/
│   └── ThemeContext.tsx    # Theme management
├── pages/
│   ├── DashboardPage.tsx   # Main dashboard
│   ├── ProfilePage.tsx     # User profile
│   ├── ZonesPage.tsx       # DNS zones
│   ├── RecordsPage.tsx     # DNS records
│   ├── LoginPage.tsx       # Login page
│   └── [Error pages...]
└── App.tsx                 # App router & theme wrapper
```

## Navigation

### Main Routes
- `/` - Dashboard (home page)
- `/zones` - DNS Zones management
- `/records` - DNS Records management
- `/profile` - User profile

### Public Routes
- `/login` - Login page

### Error Routes
- `/unauthorized` - 401 error
- `/forbidden` - 403 error
- `/server-error` - 500 error
- `*` - 404 not found (catch-all)

## Styling

### CSS Variables
The app uses Radix UI's CSS variable system:
- `--gray-a*`: Gray scale with alpha
- `--accent-*`: Accent color variations
- `--color-panel`: Background panels

### Responsive Breakpoints
- Mobile: < 768px (icons-only sidebar)
- Tablet: 768px - 1024px
- Desktop: > 1024px

## Next Steps

### Authentication
- [ ] Implement actual login logic
- [ ] Add authentication state management
- [ ] Create protected route wrapper
- [ ] Add token storage and refresh

### API Integration
- [ ] Connect to GoDNS API endpoints
- [ ] Implement zone CRUD operations
- [ ] Implement record CRUD operations
- [ ] Add real-time search functionality

### Features
- [ ] Add pagination to tables
- [ ] Implement sorting on table columns
- [ ] Add bulk actions for zones/records
- [ ] Create zone/record detail pages
- [ ] Add form validation
- [ ] Implement notifications/toasts

### UI Enhancements
- [ ] Add loading states
- [ ] Implement skeleton screens
- [ ] Add confirmation dialogs
- [ ] Create modal forms for add/edit
- [ ] Add keyboard shortcuts
- [ ] Improve mobile experience

## Theme Toggle Usage

```tsx
import { useTheme } from './contexts/ThemeContext';

function MyComponent() {
  const { theme, toggleTheme } = useTheme();
  
  return (
    <button onClick={toggleTheme}>
      Current theme: {theme}
    </button>
  );
}
```

## Search Integration

The search bar in the Header component can be extended with actual API calls:

```tsx
// In SearchBar.tsx
const handleSearch = async (query: string) => {
  const results = await api.search(query);
  setQuickResults(results);
};
```
