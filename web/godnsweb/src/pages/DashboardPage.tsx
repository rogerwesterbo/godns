import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Flex, Card, Heading, Text, Grid, Box, Spinner, Badge, Button } from '@radix-ui/themes';
import {
  GlobeIcon,
  FileTextIcon,
  RocketIcon,
  LightningBoltIcon,
  ReloadIcon,
} from '@radix-ui/react-icons';
import * as api from '../services/api';
import * as adminApi from '../services/admin-api';

export default function DashboardPage() {
  const [zones, setZones] = useState<api.DNSZone[]>([]);
  const [systemStats, setSystemStats] = useState<adminApi.SystemStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const [zonesData, system] = await Promise.all([
        api.listZones(),
        adminApi.getSystemStats().catch(() => null),
      ]);
      console.log('DashboardPage - System stats:', system);
      setZones(zonesData);
      setSystemStats(system);
    } catch (err) {
      console.error('Failed to load data:', err);
      setError(err instanceof Error ? err.message : 'Failed to load dashboard data');
    } finally {
      setIsLoading(false);
    }
  };

  // Calculate metrics
  const totalZones = zones.length;
  const totalRecords = zones.reduce((sum, zone) => sum + zone.records.length, 0);

  // Record type distribution
  const recordTypes = zones
    .flatMap(z => z.records)
    .reduce(
      (acc, record) => {
        acc[record.type] = (acc[record.type] || 0) + 1;
        return acc;
      },
      {} as Record<string, number>
    );

  const topRecordTypes = Object.entries(recordTypes)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 5);

  // Recent zones (last 5)
  const recentZones = [...zones].sort((a, b) => b.records.length - a.records.length).slice(0, 5);

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Dashboard</Heading>
        <Flex gap="3" align="center">
          <Text size="2" color="gray">
            {new Date().toLocaleDateString('en-US', {
              weekday: 'long',
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}
          </Text>
          <Button size="3" variant="soft" onClick={loadData} disabled={isLoading}>
            <ReloadIcon />
            Refresh
          </Button>
        </Flex>
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
          <Grid columns={{ initial: '1', sm: '2', md: '4' }} gap="4">
            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <GlobeIcon width="20" height="20" />
                  <Heading size="4">Zones</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold">
                    {totalZones}
                  </Text>
                  <Text size="2" color="gray">
                    Total DNS zones
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <FileTextIcon width="20" height="20" />
                  <Heading size="4">Records</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold">
                    {totalRecords}
                  </Text>
                  <Text size="2" color="gray">
                    Total DNS records
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <RocketIcon width="20" height="20" />
                  <Heading size="4">Queries</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold">
                    {systemStats?.query_log?.total_queries?.toLocaleString() || '0'}
                  </Text>
                  <Text size="2" color="gray">
                    Total queries
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <LightningBoltIcon width="20" height="20" />
                  <Heading size="4">Cache Hit Rate</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text
                    size="7"
                    weight="bold"
                    color={systemStats?.query_log?.enabled ? 'green' : 'gray'}
                  >
                    {systemStats?.query_log?.enabled
                      ? `${(systemStats.query_log.cache_hit_rate * 100).toFixed(1)}%`
                      : 'Disabled'}
                  </Text>
                  <Text size="2" color="gray">
                    {systemStats?.query_log?.enabled ? 'Cache efficiency' : 'Enable in settings'}
                  </Text>
                </Flex>
              </Flex>
            </Card>
          </Grid>

          <Grid columns={{ initial: '1', md: '2' }} gap="4">
            <Card>
              <Flex direction="column" gap="4">
                <Heading size="5">Record Type Distribution</Heading>
                {topRecordTypes.length === 0 ? (
                  <Text size="2" color="gray">
                    No records found
                  </Text>
                ) : (
                  <Flex direction="column" gap="3">
                    {topRecordTypes.map(([type, count]) => (
                      <Flex key={type} justify="between" align="center">
                        <Flex align="center" gap="2">
                          <Badge
                            color={
                              type === 'A' || type === 'AAAA'
                                ? 'blue'
                                : type === 'CNAME'
                                  ? 'green'
                                  : type === 'MX'
                                    ? 'orange'
                                    : type === 'TXT'
                                      ? 'purple'
                                      : type === 'NS' || type === 'SOA'
                                        ? 'red'
                                        : 'gray'
                            }
                          >
                            {type}
                          </Badge>
                          <Text size="2">
                            {count} record{count !== 1 ? 's' : ''}
                          </Text>
                        </Flex>
                        <Box
                          style={{
                            width: '100px',
                            height: '8px',
                            backgroundColor: 'var(--gray-a3)',
                            borderRadius: '4px',
                            overflow: 'hidden',
                          }}
                        >
                          <Box
                            style={{
                              width: `${(count / totalRecords) * 100}%`,
                              height: '100%',
                              backgroundColor: 'var(--accent-9)',
                            }}
                          />
                        </Box>
                      </Flex>
                    ))}
                  </Flex>
                )}
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="4">
                <Heading size="5">Top Zones by Records</Heading>
                {recentZones.length === 0 ? (
                  <Text size="2" color="gray">
                    No zones found
                  </Text>
                ) : (
                  <Flex direction="column" gap="3">
                    {recentZones.map(zone => (
                      <Box
                        key={zone.domain}
                        style={{
                          padding: '12px',
                          borderRadius: '6px',
                          backgroundColor: 'var(--gray-a2)',
                        }}
                      >
                        <Flex justify="between" align="center">
                          <Flex direction="column" gap="1">
                            <Link
                              to={`/zones/${encodeURIComponent(zone.domain)}`}
                              style={{ textDecoration: 'none' }}
                            >
                              <Text size="2" weight="medium">
                                {zone.domain}
                              </Text>
                            </Link>
                            <Text size="1" color="gray">
                              {zone.records.length} record{zone.records.length !== 1 ? 's' : ''}
                            </Text>
                          </Flex>
                          <Badge color="blue">{zone.records.length}</Badge>
                        </Flex>
                      </Box>
                    ))}
                  </Flex>
                )}
              </Flex>
            </Card>
          </Grid>
        </>
      )}
    </Flex>
  );
}
