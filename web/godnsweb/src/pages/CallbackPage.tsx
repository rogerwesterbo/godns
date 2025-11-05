import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Flex, Card, Text, Spinner } from '@radix-ui/themes';
import * as authService from '../services/auth';

export default function CallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    handleCallback();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleCallback = async () => {
    try {
      const code = searchParams.get('code');
      const state = searchParams.get('state');
      const errorParam = searchParams.get('error');

      if (errorParam) {
        throw new Error(searchParams.get('error_description') || errorParam);
      }

      if (!code || !state) {
        throw new Error('Missing code or state parameter');
      }

      // Exchange code for tokens
      const tokens = await authService.exchangeCodeForTokens(code, state);

      // Store tokens
      authService.storeTokens(tokens);

      // Redirect to dashboard
      navigate('/', { replace: true });

      // Reload to update auth context
      window.location.reload();
    } catch (err) {
      console.error('Callback error:', err);
      setError(err instanceof Error ? err.message : 'Authentication failed');

      // Redirect to login after error
      setTimeout(() => navigate('/login', { replace: true }), 3000);
    }
  };

  if (error) {
    return (
      <Flex
        direction="column"
        align="center"
        justify="center"
        style={{ minHeight: '100vh', padding: '1rem' }}
      >
        <Card size="4" style={{ maxWidth: '500px', width: '100%' }}>
          <Flex direction="column" gap="4" align="center">
            <Text size="6" weight="bold" color="red">
              Authentication Error
            </Text>
            <Text align="center" color="gray">
              {error}
            </Text>
            <Text size="2" color="gray">
              Redirecting to login...
            </Text>
          </Flex>
        </Card>
      </Flex>
    );
  }

  return (
    <Flex
      direction="column"
      align="center"
      justify="center"
      style={{ minHeight: '100vh', padding: '1rem' }}
    >
      <Card size="4" style={{ maxWidth: '500px', width: '100%' }}>
        <Flex direction="column" gap="4" align="center">
          <Spinner size="3" />
          <Text size="5" weight="medium">
            Completing authentication...
          </Text>
          <Text size="2" color="gray">
            Please wait while we log you in
          </Text>
        </Flex>
      </Card>
    </Flex>
  );
}
