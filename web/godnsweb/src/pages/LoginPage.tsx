import { Button, Flex, Card, Heading } from '@radix-ui/themes';
import { useAuth } from '../contexts/useAuth';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function LoginPage() {
  const { login, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    // If already authenticated, redirect to dashboard
    if (isAuthenticated) {
      navigate('/', { replace: true });
    }
  }, [isAuthenticated, navigate]);

  return (
    <Flex
      direction="column"
      align="center"
      justify="center"
      style={{ minHeight: '100vh', padding: '1rem' }}
    >
      <Card size="4" style={{ maxWidth: '400px', width: '100%' }}>
        <Flex direction="column" gap="4" align="center">
          <Heading size="6">GoDNS Login</Heading>
          <Button size="3" onClick={login} style={{ width: '100%' }}>
            Login
          </Button>
        </Flex>
      </Card>
    </Flex>
  );
}
