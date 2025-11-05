import { Flex, IconButton, Text, Avatar, DropdownMenu } from '@radix-ui/themes';
import { SunIcon, MoonIcon, PersonIcon, ExitIcon } from '@radix-ui/react-icons';
import { useTheme } from '../contexts';
import { useAuth } from '../contexts/useAuth';
import SearchBar from './SearchBar';
import './Header.css';

export default function Header() {
  const { theme, toggleTheme } = useTheme();
  const { user, logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  // Get initials from user name or username
  const getInitials = () => {
    if (user?.name) {
      return user.name
        .split(' ')
        .map(n => n[0])
        .join('')
        .toUpperCase()
        .substring(0, 2);
    }
    if (user?.preferred_username) {
      return user.preferred_username.substring(0, 2).toUpperCase();
    }
    return 'U';
  };

  return (
    <header className="header">
      <Flex justify="between" align="center" style={{ height: '100%', padding: '0 1.5rem' }}>
        <Flex align="center" gap="4">
          <a
            href="/"
            style={{
              textDecoration: 'none',
              color: 'inherit',
              cursor: 'pointer',
            }}
            onClick={e => {
              e.preventDefault();
              window.location.href = '/';
            }}
          >
            <Text size="5" weight="bold">
              GoDNS
            </Text>
          </a>
        </Flex>

        <Flex align="center" gap="4" style={{ flex: 1, justifyContent: 'center' }}>
          <SearchBar />
        </Flex>

        <Flex align="center" gap="3">
          <IconButton
            variant="ghost"
            onClick={toggleTheme}
            title={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
          >
            {theme === 'light' ? <MoonIcon /> : <SunIcon />}
          </IconButton>

          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <IconButton variant="soft" radius="full">
                <Avatar size="1" fallback={getInitials()} radius="full" />
              </IconButton>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              <DropdownMenu.Label>
                <Flex direction="column" gap="1">
                  <Text size="2" weight="medium">
                    {user?.name || user?.preferred_username}
                  </Text>
                  {user?.email && (
                    <Text size="1" color="gray">
                      {user.email}
                    </Text>
                  )}
                </Flex>
              </DropdownMenu.Label>
              <DropdownMenu.Separator />
              <DropdownMenu.Item onClick={() => (window.location.href = '/profile')}>
                <PersonIcon /> Profile
              </DropdownMenu.Item>
              <DropdownMenu.Separator />
              <DropdownMenu.Item color="red" onClick={handleLogout}>
                <ExitIcon /> Logout
              </DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
      </Flex>
    </header>
  );
}
