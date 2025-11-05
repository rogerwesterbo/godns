import { Flex, Card, Heading, Text, Button } from '@radix-ui/themes';
import { useNavigate } from 'react-router-dom';
import { CrossCircledIcon } from '@radix-ui/react-icons';

export default function ServerErrorPage() {
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
          <CrossCircledIcon width="48" height="48" color="red" />
          <Heading size="6">500 - Server Error</Heading>
          <Text align="center" color="gray">
            Something went wrong on our end. Please try again later.
          </Text>
          <Button size="3" onClick={() => navigate('/')} style={{ width: '100%' }}>
            Go Home
          </Button>
        </Flex>
      </Card>
    </Flex>
  );
}
