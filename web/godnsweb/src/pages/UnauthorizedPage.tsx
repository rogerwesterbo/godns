import { Flex, Card, Heading, Text, Button } from '@radix-ui/themes';
import { useNavigate } from 'react-router-dom';
import { ExclamationTriangleIcon } from '@radix-ui/react-icons';

export default function UnauthorizedPage() {
  const navigate = useNavigate();

  return (
    <Flex
      direction="column"
      align="center"
      justify="center"
      style={{ minHeight: '100vh', padding: '1rem' }}
    >
      <Card size="4" style={{ maxWidth: '500px', width: '100%' }}>
        <Flex direction="column" gap="4" align="center">
          <ExclamationTriangleIcon width="48" height="48" color="orange" />
          <Heading size="6">401 - Unauthorized</Heading>
          <Text align="center" color="gray">
            You are not authorized to access this resource. Please log in to continue.
          </Text>
          <Button size="3" onClick={() => navigate('/login')} style={{ width: '100%' }}>
            Go to Login
          </Button>
        </Flex>
      </Card>
    </Flex>
  );
}
