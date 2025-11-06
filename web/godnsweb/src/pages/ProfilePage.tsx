import { Flex, Card, Heading, Text, Avatar, Box, Badge, Grid, Button } from '@radix-ui/themes';
import { PersonIcon, EnvelopeClosedIcon, CalendarIcon, ReloadIcon } from '@radix-ui/react-icons';
import { useAuth } from '../contexts/useAuth';

export default function ProfilePage() {
  const { user } = useAuth();

  const handleRefresh = () => {
    window.location.reload();
  };

  if (!user) {
    return null;
  }

  // Get user initials
  const getInitials = () => {
    if (user.name) {
      return user.name
        .split(' ')
        .map(n => n[0])
        .join('')
        .toUpperCase();
    }
    if (user.preferred_username) {
      return user.preferred_username.substring(0, 2).toUpperCase();
    }
    return 'U';
  };

  // Get primary role
  const getPrimaryRole = () => {
    if (user.roles && user.roles.length > 0) {
      const dnsRole = user.roles.find(r => r.startsWith('dns-'));
      if (dnsRole) {
        return dnsRole.replace('dns-', '').replace('-', ' ').toUpperCase();
      }
      return user.roles[0];
    }
    return 'User';
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Profile</Heading>
        <Button size="3" variant="soft" onClick={handleRefresh}>
          <ReloadIcon /> Refresh
        </Button>
      </Flex>

      <Card size="4">
        <Flex direction="column" gap="6">
          <Flex align="center" gap="4">
            <Avatar size="6" fallback={getInitials()} radius="full" color="blue" />
            <Flex direction="column" gap="2">
              <Flex align="center" gap="2">
                <Heading size="6">{user.name || user.preferred_username}</Heading>
                <Badge color="blue">{getPrimaryRole()}</Badge>
                {user.email_verified && <Badge color="green">Verified</Badge>}
              </Flex>
              {user.email && (
                <Text size="2" color="gray">
                  {user.email}
                </Text>
              )}
            </Flex>
          </Flex>

          <Box style={{ height: '1px', backgroundColor: 'var(--gray-a5)' }} />

          <Grid columns={{ initial: '1', sm: '2' }} gap="4">
            <Card variant="surface">
              <Flex direction="column" gap="2">
                <Flex align="center" gap="2">
                  <PersonIcon width="16" height="16" />
                  <Text size="2" weight="bold">
                    Username
                  </Text>
                </Flex>
                <Text size="2" color="gray">
                  {user.preferred_username || 'N/A'}
                </Text>
              </Flex>
            </Card>

            <Card variant="surface">
              <Flex direction="column" gap="2">
                <Flex align="center" gap="2">
                  <EnvelopeClosedIcon width="16" height="16" />
                  <Text size="2" weight="bold">
                    Email
                  </Text>
                </Flex>
                <Text size="2" color="gray">
                  {user.email || 'N/A'}
                </Text>
              </Flex>
            </Card>

            {user.given_name && (
              <Card variant="surface">
                <Flex direction="column" gap="2">
                  <Flex align="center" gap="2">
                    <PersonIcon width="16" height="16" />
                    <Text size="2" weight="bold">
                      First Name
                    </Text>
                  </Flex>
                  <Text size="2" color="gray">
                    {user.given_name}
                  </Text>
                </Flex>
              </Card>
            )}

            {user.family_name && (
              <Card variant="surface">
                <Flex direction="column" gap="2">
                  <Flex align="center" gap="2">
                    <PersonIcon width="16" height="16" />
                    <Text size="2" weight="bold">
                      Last Name
                    </Text>
                  </Flex>
                  <Text size="2" color="gray">
                    {user.family_name}
                  </Text>
                </Flex>
              </Card>
            )}

            <Card variant="surface">
              <Flex direction="column" gap="2">
                <Flex align="center" gap="2">
                  <CalendarIcon width="16" height="16" />
                  <Text size="2" weight="bold">
                    User ID
                  </Text>
                </Flex>
                <Text
                  size="2"
                  color="gray"
                  style={{
                    fontSize: '11px',
                    wordBreak: 'break-all',
                  }}
                >
                  {user.sub}
                </Text>
              </Flex>
            </Card>
          </Grid>
        </Flex>
      </Card>

      {user.roles && user.roles.length > 0 && (
        <Card>
          <Flex direction="column" gap="4">
            <Heading size="5">Roles & Permissions</Heading>
            <Flex gap="2" wrap="wrap">
              {user.roles.map(role => (
                <Badge
                  key={role}
                  size="2"
                  color={
                    role.includes('admin') ? 'red' : role.includes('write') ? 'orange' : 'blue'
                  }
                >
                  {role}
                </Badge>
              ))}
            </Flex>
          </Flex>
        </Card>
      )}
    </Flex>
  );
}
