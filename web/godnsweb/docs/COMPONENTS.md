# Component Documentation

## Layout Components

### Layout
**Location**: `src/components/Layout.tsx`

The main layout wrapper that includes Header and Sidebar.

**Props**:
- `children: ReactNode` - Page content to display

**Usage**:
```tsx
<Layout>
  <DashboardPage />
</Layout>
```

**Features**:
- Sticky header (60px height)
- Fixed sidebar (240px width, 60px on mobile)
- Scrollable main content area
- Responsive flex layout

---

### Header
**Location**: `src/components/Header.tsx`

Top navigation bar with search, theme toggle, and user menu.

**Props**: None (uses context)

**Features**:
- **Branding**: "GoDNS" text logo on left
- **Search**: Animated SearchBar component in center
- **Theme Toggle**: Sun/Moon icon button
- **User Menu**: Avatar dropdown with profile and logout

**Styling**:
- Height: 60px
- Sticky position at top
- Backdrop blur effect
- Border bottom

---

### Sidebar
**Location**: `src/components/Sidebar.tsx`

Side navigation menu with route links.

**Props**: None

**Navigation Items**:
1. Dashboard (/) - DashboardIcon
2. Zones (/zones) - GlobeIcon  
3. Records (/records) - FileTextIcon
4. Profile (/profile) - PersonIcon

**Features**:
- Active route highlighting
- Hover effects
- Icon and text labels
- Responsive (icon-only on mobile)

**Styling**:
- Width: 240px (desktop), 60px (mobile)
- Sticky position below header
- Scrollable overflow

---

### SearchBar
**Location**: `src/components/SearchBar.tsx`

Animated search input with quick results dropdown.

**Props**: None

**State**:
- `isExpanded: boolean` - Expansion state
- `searchQuery: string` - Current search text

**Features**:
- **Compact**: 300px width by default
- **Expanded**: 600px width when focused
- **Quick Results**: Filtered list below search
- **Auto-close**: Closes on outside click
- **Animations**: Smooth CSS transitions

**Mock Data**:
Currently shows mock zones and records. Replace with API calls:

```tsx
const quickResults = await api.search(searchQuery);
```

**Styling**:
- Transition duration: 0.3s
- Animation: slideDown on results
- Max results height: 400px

---

## Theme System

### ThemeContext
**Location**: `src/contexts/ThemeContext.tsx`

Manages application theme state.

**Exports**:
- `ThemeProvider` - Context provider component
- `useTheme` - Hook to access theme

**API**:
```tsx
const { theme, toggleTheme } = useTheme();
// theme: 'light' | 'dark'
// toggleTheme: () => void
```

**Storage**:
- localStorage key: 'theme'
- Default: 'dark'

---

## Page Components

### DashboardPage
**Location**: `src/pages/DashboardPage.tsx`

Main admin dashboard with statistics and activity.

**Features**:
- Statistics cards (Zones, Records, Activity)
- Recent activity timeline
- Mock data (replace with API)

---

### ZonesPage
**Location**: `src/pages/ZonesPage.tsx`

DNS zones management interface.

**Features**:
- Zone list table
- Type badges (Primary/Secondary)
- Record count per zone
- Add zone button
- Mock data (replace with API)

---

### RecordsPage
**Location**: `src/pages/RecordsPage.tsx`

DNS records management interface.

**Features**:
- Records list table
- Type filter dropdown
- Record details (name, type, value, TTL)
- Add record button
- Mock data (replace with API)

---

### ProfilePage
**Location**: `src/pages/ProfilePage.tsx`

User profile and account information.

**Features**:
- User avatar and basic info
- Role badge
- Account statistics grid
- Member since date
- Mock data (replace with API)

---

### LoginPage
**Location**: `src/pages/LoginPage.tsx`

Authentication page.

**Features**:
- Centered login card
- Single "Login" button
- TODO: Implement actual auth

---

### Error Pages
**Locations**: `src/pages/*Page.tsx`

Standardized error pages.

**Pages**:
- `UnauthorizedPage` - 401 (orange warning icon)
- `ForbiddenPage` - 403 (red lock icon)
- `NotFoundPage` - 404 (gray magnifying glass)
- `ServerErrorPage` - 500 (red cross icon)

**Features**:
- Centered error cards
- Icon, title, description
- Navigation button back to home/login
- Consistent styling

---

## Styling Guidelines

### Component Styles
Each component with custom CSS has a corresponding `.css` file:
- `Header.css`
- `Sidebar.css`
- `SearchBar.css`

### CSS Variables
Use Radix UI variables:

```css
/* Colors */
var(--gray-12)      /* Primary text */
var(--gray-11)      /* Secondary text */
var(--gray-a5)      /* Borders */
var(--accent-9)     /* Accent solid */
var(--accent-a4)    /* Accent subtle */

/* Backgrounds */
var(--color-panel)      /* Panel background */
var(--color-background) /* App background */
```

### Responsive Design
Mobile breakpoint: 768px

```css
@media (max-width: 768px) {
  /* Mobile styles */
}
```

---

## Icons

### Radix Icons
All icons from `@radix-ui/react-icons`:

**Navigation**:
- `DashboardIcon` - Dashboard
- `GlobeIcon` - Zones
- `FileTextIcon` - Records
- `PersonIcon` - Profile

**Actions**:
- `PlusIcon` - Add button
- `MagnifyingGlassIcon` - Search
- `SunIcon` - Light mode
- `MoonIcon` - Dark mode
- `ExitIcon` - Logout

**Errors**:
- `ExclamationTriangleIcon` - Warning (401)
- `LockClosedIcon` - Forbidden (403)
- `MagnifyingGlassIcon` - Not found (404)
- `CrossCircledIcon` - Error (500)

**Stats**:
- `ActivityLogIcon` - Activity
- `CalendarIcon` - Date
- `EnvelopeClosedIcon` - Email

**Usage**:
```tsx
import { IconName } from '@radix-ui/react-icons';

<IconName width="20" height="20" />
```

---

## Adding New Components

1. Create component file in `src/components/`
2. Add CSS file if needed
3. Export from `src/components/index.ts`
4. Import in parent component

Example:
```tsx
// src/components/NewComponent.tsx
export default function NewComponent() {
  return <div>New Component</div>;
}

// src/components/index.ts
export { default as NewComponent } from './NewComponent';

// Usage
import { NewComponent } from './components';
```

---

## Best Practices

1. **Use Radix UI components**: Leverage Button, Card, Flex, etc.
2. **Consistent spacing**: Use Radix `gap` props (1-9 scale)
3. **Type safety**: Define prop interfaces
4. **Responsive**: Test mobile and desktop
5. **Accessibility**: Use semantic HTML and ARIA labels
6. **CSS variables**: Avoid hard-coded colors
7. **Component exports**: Use index.ts for clean imports
