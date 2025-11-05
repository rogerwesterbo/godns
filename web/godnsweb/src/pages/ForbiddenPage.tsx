import { Flex, Card, Heading, Text, Button } from '@radix-ui/themes';
import { useNavigate } from 'react-router-dom';
import { LockClosedIcon } from '@radix-ui/react-icons';

export default function ForbiddenPage() {
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
          <LockClosedIcon width="48" height="48" color="red" />
          <Heading size="6">403 - Forbidden</Heading>
          <Text align="center" color="gray">
            You don't have permission to access this resource.
          </Text>
          <Button size="3" onClick={() => navigate('/')} style={{ width: '100%' }}>
            Go Home
          </Button>
        </Flex>
      </Card>
    </Flex>
  );
}
