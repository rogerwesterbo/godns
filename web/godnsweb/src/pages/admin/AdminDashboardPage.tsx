import { useEffect, useState } from 'react';
import { Flex, Card, Heading, Text, Grid, Spinner, Badge, Button } from '@radix-ui/themes';
import { ReloadIcon, RocketIcon, ActivityLogIcon, LightningBoltIcon } from '@radix-ui/react-icons';
import * as adminApi from '../../services/admin-api';

export default function AdminDashboardPage() {
  const [systemStats, setSystemStats] = useState<adminApi.SystemStats | null>(null);
  const [cacheStats, setCacheStats] = useState<adminApi.CacheStatsDetailed | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const [system, cache] = await Promise.all([
        adminApi.getSystemStats(),
        adminApi.getCacheStats(),
      ]);
      console.log('System stats cache:', system.cache);
      console.log('Detailed cache stats:', cache);
      setSystemStats(system);
      setCacheStats(cache);
    } catch (err) {
      console.error('Failed to load admin stats:', err);
      setError(err instanceof Error ? err.message : 'Failed to load statistics');
    } finally {
      setIsLoading(false);
    }
  };

  const formatPercentage = (value: number | string): string => {
    // If already a string with %, return it
    if (typeof value === 'string') {
      return value;
    }
    // If a number, convert to percentage
    return `${(value * 100).toFixed(1)}%`;
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Admin Dashboard</Heading>
        <Button size="3" variant="soft" onClick={loadStats} disabled={isLoading}>
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
          {/* System Stats */}
          <Grid columns={{ initial: '1', sm: '2', md: '3' }} gap="4">
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
                  <ActivityLogIcon width="20" height="20" />
                  <Heading size="4">Cached</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold" color="green">
                    {systemStats?.query_log?.cached_queries?.toLocaleString() || '0'}
                  </Text>
                  <Text size="2" color="gray">
                    Cached queries
                  </Text>
                </Flex>
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="3">
                <Flex align="center" gap="2">
                  <LightningBoltIcon width="20" height="20" />
                  <Heading size="4">Query Cache Rate</Heading>
                </Flex>
                <Flex direction="column" gap="1">
                  <Text size="7" weight="bold" color="blue">
                    {systemStats?.query_log?.cache_hit_rate
                      ? formatPercentage(systemStats.query_log.cache_hit_rate)
                      : '0%'}
                  </Text>
                  <Text size="2" color="gray">
                    Cache hit rate
                  </Text>
                </Flex>
              </Flex>
            </Card>
          </Grid>

          {/* Cache Stats */}
          <Grid columns={{ initial: '1', md: '2' }} gap="4">
            <Card>
              <Flex direction="column" gap="4">
                <Flex justify="between" align="center">
                  <Heading size="5">DNS Cache</Heading>
                  <Badge color={cacheStats?.enabled ? 'green' : 'gray'}>
                    {cacheStats?.enabled ? 'Enabled' : 'Disabled'}
                  </Badge>
                </Flex>
                {cacheStats?.enabled ? (
                  <Grid columns="2" gap="4">
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Cache Size
                      </Text>
                      <Text size="6" weight="bold">
                        {cacheStats.size} / {cacheStats.capacity}
                      </Text>
                    </Flex>
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Hit Rate
                      </Text>
                      <Text size="6" weight="bold" color="green">
                        {typeof cacheStats.hit_rate === 'string'
                          ? cacheStats.hit_rate
                          : `${(cacheStats.hit_rate * 100).toFixed(1)}%`}
                      </Text>
                    </Flex>
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Cache Hits
                      </Text>
                      <Text size="5" weight="medium">
                        {cacheStats.hits?.toLocaleString() || '0'}
                      </Text>
                    </Flex>
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Cache Misses
                      </Text>
                      <Text size="5" weight="medium">
                        {cacheStats.misses?.toLocaleString() || '0'}
                      </Text>
                    </Flex>
                  </Grid>
                ) : (
                  <Text size="2" color="gray">
                    DNS caching is currently disabled
                  </Text>
                )}
              </Flex>
            </Card>

            <Card>
              <Flex direction="column" gap="4">
                <Flex justify="between" align="center">
                  <Heading size="5">Rate Limiter</Heading>
                  <Badge color={systemStats?.rate_limiter?.enabled ? 'green' : 'gray'}>
                    {systemStats?.rate_limiter?.enabled ? 'Enabled' : 'Disabled'}
                  </Badge>
                </Flex>
                {systemStats?.rate_limiter?.enabled ? (
                  <Grid columns="2" gap="4">
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Rate Limit
                      </Text>
                      <Text size="6" weight="bold">
                        {systemStats.rate_limiter.qps} QPS
                      </Text>
                    </Flex>
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Burst
                      </Text>
                      <Text size="6" weight="bold">
                        {systemStats.rate_limiter.burst}
                      </Text>
                    </Flex>
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Active Limiters
                      </Text>
                      <Text size="5" weight="medium">
                        {systemStats.rate_limiter.active_limiters || 0}
                      </Text>
                    </Flex>
                    <Flex direction="column" gap="1">
                      <Text size="2" color="gray">
                        Blocked Queries
                      </Text>
                      <Text size="5" weight="medium" color="red">
                        {systemStats.rate_limiter.total_blocked?.toLocaleString() || '0'}
                      </Text>
                    </Flex>
                  </Grid>
                ) : (
                  <Text size="2" color="gray">
                    Rate limiting is currently disabled
                  </Text>
                )}
              </Flex>
            </Card>
          </Grid>

          {/* Query Stats - Detailed breakdown */}
          {systemStats?.query_log?.enabled && (
            <Card>
              <Flex direction="column" gap="4">
                <Flex justify="between" align="center">
                  <Heading size="5">Detailed Query Statistics</Heading>
                  <Badge color="blue">
                    <LightningBoltIcon />
                    Active
                  </Badge>
                </Flex>
                <Grid columns={{ initial: '1', sm: '2', md: '4' }} gap="4">
                  <Flex direction="column" gap="1">
                    <Text size="2" color="gray">
                      Total Queries
                    </Text>
                    <Text size="6" weight="bold">
                      {systemStats.query_log.total_queries?.toLocaleString() || '0'}
                    </Text>
                  </Flex>
                  <Flex direction="column" gap="1">
                    <Text size="2" color="gray">
                      Cached Queries
                    </Text>
                    <Text size="6" weight="bold" color="green">
                      {systemStats.query_log.cached_queries?.toLocaleString() || '0'}
                    </Text>
                  </Flex>
                  <Flex direction="column" gap="1">
                    <Text size="2" color="gray">
                      Blocked Queries
                    </Text>
                    <Text size="6" weight="bold" color="red">
                      {systemStats.query_log.blocked_queries?.toLocaleString() || '0'}
                    </Text>
                  </Flex>
                  <Flex direction="column" gap="1">
                    <Text size="2" color="gray">
                      Cache Hit Rate
                    </Text>
                    <Text size="6" weight="bold" color="blue">
                      {formatPercentage(systemStats.query_log.cache_hit_rate)}
                    </Text>
                  </Flex>
                </Grid>
              </Flex>
            </Card>
          )}
        </>
      )}
    </Flex>
  );
}
