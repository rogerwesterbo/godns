import { Flex, Card, Heading, Text, Button } from '@radix-ui/themes';
import { useNavigate } from 'react-router-dom';
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';

export default function NotFoundPage() {
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
          <MagnifyingGlassIcon width="48" height="48" color="gray" />
          <Heading size="6">404 - Not Found</Heading>
          <Text align="center" color="gray">
            The page you're looking for doesn't exist.
          </Text>
          <Button size="3" onClick={() => navigate('/')} style={{ width: '100%' }}>
            Go Home
          </Button>
        </Flex>
      </Card>
    </Flex>
  );
}
