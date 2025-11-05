import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { Flex, Spinner, Text } from '@radix-ui/themes';
import { useAuth } from '../contexts/useAuth';

interface ProtectedRouteProps {
  children: ReactNode;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <Flex
        direction="column"
        align="center"
        justify="center"
        style={{ minHeight: '100vh' }}
        gap="4"
      >
        <Spinner size="3" />
        <Text size="2" color="gray">
          Loading...
        </Text>
      </Flex>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}
