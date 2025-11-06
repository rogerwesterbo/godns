import { useEffect, useState } from 'react';
import { Flex, Card, Heading, Text, Grid, Spinner, Badge, Button, Table } from '@radix-ui/themes';
import { ReloadIcon, CheckCircledIcon, CrossCircledIcon } from '@radix-ui/react-icons';
import * as adminApi from '../../services/admin-api';

export default function LoadBalancerPage() {
  const [stats, setStats] = useState<adminApi.LoadBalancerStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await adminApi.getLoadBalancerStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to load load balancer stats:', err);
      setError(err instanceof Error ? err.message : 'Failed to load load balancer statistics');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Load Balancer</Heading>
        <Button onClick={loadStats} disabled={isLoading}>
          <ReloadIcon />
          Refresh
        </Button>
      </Flex>

      {isLoading ? (
        <Flex justify="center" py="8">
          <Spinner size="3" />
        </Flex>
      ) : error ? (
        <Card>
          <Text color="red" size="3">
            {error}
          </Text>
        </Card>
      ) : (
        <>
          <Grid columns={{ initial: '1', sm: '2', md: '3' }} gap="4">
            <Card>
              <Flex direction="column" gap="3">
                <Heading size="4">Status</Heading>
                <Badge color={stats?.enabled ? 'green' : 'gray'} size="2">
                  {stats?.enabled ? 'Enabled' : 'Disabled'}
                </Badge>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Heading size="4">Strategy</Heading>
                <Text size="5" weight="bold">
                  {stats?.strategy || 'N/A'}
                </Text>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Heading size="4">Health</Heading>
                <Text size="5" weight="bold">
                  {stats?.healthy_backends || 0} / {stats?.total_backends || 0}
                </Text>
                <Text size="2" color="gray">
                  Healthy backends
                </Text>
              </Flex>
            </Card>
          </Grid>

          {stats?.enabled && stats.backends && stats.backends.length > 0 ? (
            <Card>
              <Flex direction="column" gap="4">
                <Heading size="5">Backend Servers</Heading>
                <Table.Root>
                  <Table.Header>
                    <Table.Row>
                      <Table.ColumnHeaderCell>Address</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Response Time</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Last Check</Table.ColumnHeaderCell>
                    </Table.Row>
                  </Table.Header>
                  <Table.Body>
                    {stats.backends.map((backend, index) => (
                      <Table.Row key={index}>
                        <Table.Cell>
                          <Text size="2" weight="medium">
                            {backend.address}
                          </Text>
                        </Table.Cell>
                        <Table.Cell>
                          <Flex align="center" gap="2">
                            {backend.healthy ? (
                              <CheckCircledIcon color="green" />
                            ) : (
                              <CrossCircledIcon color="red" />
                            )}
                            <Badge color={backend.healthy ? 'green' : 'red'} size="1">
                              {backend.healthy ? 'Healthy' : 'Unhealthy'}
                            </Badge>
                          </Flex>
                        </Table.Cell>
                        <Table.Cell>
                          <Text
                            size="2"
                            color={backend.response_time_ms < 100 ? 'green' : 'orange'}
                          >
                            {backend.response_time_ms}ms
                          </Text>
                        </Table.Cell>
                        <Table.Cell>
                          <Text size="2" color="gray">
                            {new Date(backend.last_check).toLocaleString()}
                          </Text>
                        </Table.Cell>
                      </Table.Row>
                    ))}
                  </Table.Body>
                </Table.Root>
              </Flex>
            </Card>
          ) : (
            <Card>
              <Text size="2" color="gray">
                {stats?.enabled
                  ? 'No backend servers configured'
                  : 'Load balancer is currently disabled'}
              </Text>
            </Card>
          )}
        </>
      )}
    </Flex>
  );
}
