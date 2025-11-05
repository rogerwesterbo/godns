# Theme Implementation Guide

## Overview
The GoDNS web application uses a custom theme system built on top of Radix UI's theming capabilities. Dark mode is the default theme.

## Theme Context

### Location
`src/contexts/ThemeContext.tsx`

### Features
- Persistent theme storage (localStorage)
- Default to dark mode
- Easy theme toggling
- Type-safe theme values

### Usage

```tsx
import { useTheme } from './contexts/ThemeContext';

function MyComponent() {
  const { theme, toggleTheme } = useTheme();
  
  return (
    <div>
      <p>Current theme: {theme}</p>
      <button onClick={toggleTheme}>
        Toggle to {theme === 'light' ? 'dark' : 'light'} mode
      </button>
    </div>
  );
}
```

## App Integration

The theme is integrated at the root level in `App.tsx`:

```tsx
function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
}

function AppContent() {
  const { theme } = useTheme();
  
  return (
    <Theme appearance={theme}>
      {/* App content */}
    </Theme>
  );
}
```

## Theme Toggle Button

The theme toggle is in the Header component:

```tsx
<IconButton onClick={toggleTheme}>
  {theme === 'light' ? <MoonIcon /> : <SunIcon />}
</IconButton>
```

- **Light mode**: Shows moon icon (click to go dark)
- **Dark mode**: Shows sun icon (click to go light)

## CSS Variables

Radix UI automatically manages CSS variables based on the theme:

### Light Mode
- Light backgrounds
- Dark text
- Subtle shadows

### Dark Mode (Default)
- Dark backgrounds
- Light text
- Elevated surfaces

## Accessing Theme in CSS

Use Radix UI's CSS variables:

```css
.my-component {
  background-color: var(--color-panel);
  color: var(--gray-12);
  border: 1px solid var(--gray-a5);
}

.my-accent {
  background-color: var(--accent-9);
  color: var(--accent-contrast);
}
```

## Common CSS Variables

### Backgrounds
- `--color-background` - App background
- `--color-panel` - Panel/card background
- `--color-surface` - Surface background

### Text
- `--gray-12` - Primary text
- `--gray-11` - Secondary text
- `--gray-9` - Disabled text

### Borders
- `--gray-a5` - Subtle borders
- `--gray-a6` - Medium borders
- `--gray-a7` - Strong borders

### Accent Colors
- `--accent-9` - Solid accent
- `--accent-a4` - Subtle accent background
- `--accent-11` - Accent text
- `--accent-12` - High contrast accent text

## Customizing Theme

To change the accent color, update in `App.tsx`:

```tsx
<Theme 
  accentColor="blue"    // blue, green, red, purple, etc.
  grayColor="slate"     // slate, gray, mauve, etc.
  radius="medium"       // none, small, medium, large, full
  appearance={theme}
>
```

## Testing Themes

1. Start the dev server: `npm run dev`
2. Open http://localhost:14200
3. Click the sun/moon icon in the header
4. Theme should toggle and persist on reload

## LocalStorage

Theme preference is stored as:
```
Key: 'theme'
Values: 'light' | 'dark'
Default: 'dark'
```

## Adding Theme-Aware Components

When creating new components that need theme awareness:

```tsx
import { useTheme } from '../contexts/ThemeContext';

function ThemeAwareComponent() {
  const { theme } = useTheme();
  
  return (
    <div className={`component ${theme}`}>
      {/* Content */}
    </div>
  );
}
```

## Best Practices

1. **Use CSS Variables**: Leverage Radix UI's variables for consistency
2. **Avoid Hard-Coded Colors**: Don't use `#fff` or `#000` directly
3. **Test Both Themes**: Ensure components look good in light and dark
4. **Semantic Variables**: Use `--gray-12` not specific hex codes
5. **Respect User Preference**: Theme persists across sessions

## Troubleshooting

### Theme not persisting
- Check localStorage in DevTools
- Ensure ThemeProvider wraps entire app
- Verify no errors in console

### Colors look wrong
- Check you're using Radix CSS variables
- Verify Theme component has correct appearance prop
- Ensure @radix-ui/themes/styles.css is imported

### Toggle not working
- Verify useTheme hook is inside ThemeProvider
- Check toggleTheme function is called correctly
- Look for console errors
