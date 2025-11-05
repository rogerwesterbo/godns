import { Flex, Text, Box } from '@radix-ui/themes';
import { DashboardIcon, GlobeIcon, FileTextIcon, PersonIcon } from '@radix-ui/react-icons';
import { NavLink } from 'react-router-dom';
import './Sidebar.css';

interface NavItem {
  to: string;
  icon: React.ReactNode;
  label: string;
}

const navItems: NavItem[] = [
  { to: '/', icon: <DashboardIcon />, label: 'Dashboard' },
  { to: '/zones', icon: <GlobeIcon />, label: 'Zones' },
  { to: '/records', icon: <FileTextIcon />, label: 'Records' },
  { to: '/profile', icon: <PersonIcon />, label: 'Profile' },
];

export default function Sidebar() {
  return (
    <aside className="sidebar">
      <Flex direction="column" gap="2" style={{ padding: '1rem' }}>
        {navItems.map(item => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            end={item.to === '/'}
          >
            <Flex align="center" gap="3">
              <Box className="nav-icon">{item.icon}</Box>
              <Text size="2" weight="medium">
                {item.label}
              </Text>
            </Flex>
          </NavLink>
        ))}
      </Flex>
    </aside>
  );
}
