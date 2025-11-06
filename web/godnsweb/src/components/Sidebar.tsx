import { useState } from 'react';
import { Flex, Text, Box } from '@radix-ui/themes';
import {
  DashboardIcon,
  GlobeIcon,
  FileTextIcon,
  PersonIcon,
  GearIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  RocketIcon,
  BarChartIcon,
  DownloadIcon,
} from '@radix-ui/react-icons';
import { NavLink, useLocation } from 'react-router-dom';
import './Sidebar.css';

interface NavItem {
  to: string;
  icon: React.ReactNode;
  label: string;
  children?: NavItem[];
}

const navItems: NavItem[] = [
  { to: '/', icon: <DashboardIcon />, label: 'Dashboard' },
  { to: '/zones', icon: <GlobeIcon />, label: 'Zones' },
  { to: '/records', icon: <FileTextIcon />, label: 'Records' },
  { to: '/export', icon: <DownloadIcon />, label: 'Export' },
  {
    to: '/admin',
    icon: <GearIcon />,
    label: 'Admin',
    children: [
      { to: '/admin', icon: <BarChartIcon />, label: 'Overview' },
      { to: '/admin/cache', icon: <RocketIcon />, label: 'Cache' },
    ],
  },
  { to: '/profile', icon: <PersonIcon />, label: 'Profile' },
];

export default function Sidebar() {
  const location = useLocation();
  const [expandedMenus, setExpandedMenus] = useState<Set<string>>(
    new Set(
      navItems
        .filter(item => item.children && location.pathname.startsWith(item.to))
        .map(item => item.to)
    )
  );

  const toggleMenu = (path: string) => {
    setExpandedMenus(prev => {
      const next = new Set(prev);
      if (next.has(path)) {
        next.delete(path);
      } else {
        next.add(path);
      }
      return next;
    });
  };

  const renderNavItem = (item: NavItem, isChild = false) => {
    const hasChildren = item.children && item.children.length > 0;
    const isExpanded = expandedMenus.has(item.to);
    const isActive =
      location.pathname === item.to || (hasChildren && location.pathname.startsWith(item.to));

    if (hasChildren) {
      return (
        <div key={item.to}>
          <div
            className={`nav-item ${isActive ? 'active' : ''}`}
            onClick={() => toggleMenu(item.to)}
            style={{ cursor: 'pointer' }}
          >
            <Flex align="center" gap="3" justify="between" style={{ width: '100%' }}>
              <Flex align="center" gap="3">
                <Box className="nav-icon">{item.icon}</Box>
                <Text size="2" weight="medium">
                  {item.label}
                </Text>
              </Flex>
              <Box className="nav-icon">
                {isExpanded ? <ChevronDownIcon /> : <ChevronRightIcon />}
              </Box>
            </Flex>
          </div>
          {isExpanded && item.children && (
            <Flex direction="column" gap="1" style={{ paddingLeft: '2rem', marginTop: '0.25rem' }}>
              {item.children.map(child => renderNavItem(child, true))}
            </Flex>
          )}
        </div>
      );
    }

    return (
      <NavLink
        key={item.to}
        to={item.to}
        className={({ isActive }) =>
          `nav-item ${isChild ? 'nav-item-child' : ''} ${isActive ? 'active' : ''}`
        }
        end={item.to === '/' || isChild}
      >
        <Flex align="center" gap="3">
          <Box className="nav-icon">{item.icon}</Box>
          <Text size="2" weight="medium">
            {item.label}
          </Text>
        </Flex>
      </NavLink>
    );
  };

  return (
    <aside className="sidebar">
      <Flex direction="column" gap="2" style={{ padding: '1rem' }}>
        {navItems.map(item => renderNavItem(item))}
      </Flex>
    </aside>
  );
}
