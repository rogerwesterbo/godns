import type { ReactNode } from 'react';
import { Flex, Box } from '@radix-ui/themes';
import Header from './Header';
import Sidebar from './Sidebar';

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <Flex direction="column" style={{ minHeight: '100vh' }}>
      <Header />
      <Flex style={{ flex: 1 }}>
        <Sidebar />
        <Box style={{ flex: 1, padding: '2rem', overflowY: 'auto' }}>{children}</Box>
      </Flex>
    </Flex>
  );
}
