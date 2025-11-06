import { useEffect, useState } from 'react';
import { Flex, Card, Heading, Text, Grid, Spinner, Badge, Button, Table } from '@radix-ui/themes';
import { ReloadIcon, CheckCircledIcon, CrossCircledIcon } from '@radix-ui/react-icons';
import * as adminApi from '../../services/admin-api';

export default function HealthCheckPage() {
  const [stats, setStats] = useState<adminApi.HealthCheckStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await adminApi.getHealthCheckStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to load health check stats:', err);
      setError(err instanceof Error ? err.message : 'Failed to load health check statistics');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Health Checks</Heading>
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
                <Heading size="4">Total Targets</Heading>
                <Text size="7" weight="bold">
                  {stats?.total_targets || 0}
                </Text>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Heading size="4">Healthy Targets</Heading>
                <Text size="7" weight="bold" color="green">
                  {stats?.healthy_targets || 0}
                </Text>
              </Flex>
            </Card>
          </Grid>

          {stats?.enabled && stats.targets && stats.targets.length > 0 ? (
            <Card>
              <Flex direction="column" gap="4">
                <Heading size="5">Health Check Targets</Heading>
                <Table.Root>
                  <Table.Header>
                    <Table.Row>
                      <Table.ColumnHeaderCell>Target</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Protocol</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Last Check</Table.ColumnHeaderCell>
                      <Table.ColumnHeaderCell>Error</Table.ColumnHeaderCell>
                    </Table.Row>
                  </Table.Header>
                  <Table.Body>
                    {stats.targets.map((target, index) => (
                      <Table.Row key={index}>
                        <Table.Cell>
                          <Text size="2" weight="medium">
                            {target.host}:{target.port}
                          </Text>
                        </Table.Cell>
                        <Table.Cell>
                          <Badge variant="soft" size="1">
                            {target.protocol.toUpperCase()}
                          </Badge>
                        </Table.Cell>
                        <Table.Cell>
                          <Flex align="center" gap="2">
                            {target.healthy ? (
                              <CheckCircledIcon color="green" />
                            ) : (
                              <CrossCircledIcon color="red" />
                            )}
                            <Badge color={target.healthy ? 'green' : 'red'} size="1">
                              {target.healthy ? 'Healthy' : 'Unhealthy'}
                            </Badge>
                          </Flex>
                        </Table.Cell>
                        <Table.Cell>
                          <Text size="2" color="gray">
                            {target.last_check
                              ? new Date(target.last_check).toLocaleString()
                              : 'Never'}
                          </Text>
                        </Table.Cell>
                        <Table.Cell>
                          {target.last_error ? (
                            <Text size="2" color="red">
                              {target.last_error}
                            </Text>
                          ) : (
                            <Text size="2" color="gray">
                              -
                            </Text>
                          )}
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
                  ? 'No health check targets configured'
                  : 'Health checks are currently disabled'}
              </Text>
            </Card>
          )}
        </>
      )}
    </Flex>
  );
}
