import { useEffect, useState } from 'react';
import { Flex, Card, Heading, Text, Grid, Spinner, Badge, Button, Callout } from '@radix-ui/themes';
import { ReloadIcon, TrashIcon, InfoCircledIcon } from '@radix-ui/react-icons';
import * as adminApi from '../../services/admin-api';

export default function CachePage() {
  const [stats, setStats] = useState<adminApi.CacheStatsDetailed | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isClearing, setIsClearing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [clearSuccess, setClearSuccess] = useState(false);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await adminApi.getCacheStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to load cache stats:', err);
      setError(err instanceof Error ? err.message : 'Failed to load cache statistics');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClearCache = async () => {
    if (!window.confirm('Are you sure you want to clear the DNS cache?')) {
      return;
    }

    try {
      setIsClearing(true);
      setClearSuccess(false);
      await adminApi.clearCache();
      setClearSuccess(true);
      await loadStats();
      setTimeout(() => setClearSuccess(false), 3000);
    } catch (err) {
      console.error('Failed to clear cache:', err);
      setError(err instanceof Error ? err.message : 'Failed to clear cache');
    } finally {
      setIsClearing(false);
    }
  };

  const formatPercentage = (value: number | string): string => {
    if (typeof value === 'string') {
      return value; // Already formatted as string like "93.10%"
    }
    return `${(value * 100).toFixed(1)}%`;
  };

  const usagePercentage = stats ? (stats.size / stats.capacity) * 100 : 0;

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">DNS Cache</Heading>
        <Flex gap="2">
          <Button size="3" onClick={loadStats} disabled={isLoading} variant="soft">
            <ReloadIcon />
            Refresh
          </Button>
          {stats?.enabled && (
            <Button
              size="3"
              onClick={handleClearCache}
              disabled={isClearing || !stats?.enabled}
              color="red"
              variant="soft"
            >
              <TrashIcon />
              Clear Cache
            </Button>
          )}
        </Flex>
      </Flex>

      {clearSuccess && (
        <Callout.Root color="green">
          <Callout.Icon>
            <InfoCircledIcon />
          </Callout.Icon>
          <Callout.Text>Cache cleared successfully!</Callout.Text>
        </Callout.Root>
      )}

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
          <Card>
            <Flex direction="column" gap="4">
              <Flex justify="between" align="center">
                <Heading size="5">Cache Status</Heading>
                <Badge color={stats?.enabled ? 'green' : 'gray'} size="2">
                  {stats?.enabled ? 'Enabled' : 'Disabled'}
                </Badge>
              </Flex>

              {stats?.enabled ? (
                <>
                  <Grid columns={{ initial: '1', sm: '2', md: '4' }} gap="4">
                    <Flex direction="column" gap="2">
                      <Text size="2" color="gray">
                        Cache Size
                      </Text>
                      <Text size="7" weight="bold">
                        {stats.size}
                      </Text>
                      <Text size="2" color="gray">
                        of {stats.capacity} entries
                      </Text>
                    </Flex>

                    <Flex direction="column" gap="2">
                      <Text size="2" color="gray">
                        Hit Rate
                      </Text>
                      <Text size="7" weight="bold" color="green">
                        {formatPercentage(stats.hit_rate)}
                      </Text>
                      <Text size="2" color="gray">
                        Cache efficiency
                      </Text>
                    </Flex>

                    <Flex direction="column" gap="2">
                      <Text size="2" color="gray">
                        Cache Hits
                      </Text>
                      <Text size="7" weight="bold" color="green">
                        {stats.hits?.toLocaleString() || '0'}
                      </Text>
                      <Text size="2" color="gray">
                        Successful lookups
                      </Text>
                    </Flex>

                    <Flex direction="column" gap="2">
                      <Text size="2" color="gray">
                        Cache Misses
                      </Text>
                      <Text size="7" weight="bold" color="orange">
                        {stats.misses?.toLocaleString() || '0'}
                      </Text>
                      <Text size="2" color="gray">
                        Lookups missed
                      </Text>
                    </Flex>
                  </Grid>

                  <Flex direction="column" gap="2">
                    <Flex justify="between">
                      <Text size="2" color="gray">
                        Cache Usage
                      </Text>
                      <Text size="2" weight="medium">
                        {usagePercentage.toFixed(1)}%
                      </Text>
                    </Flex>
                    <div
                      style={{
                        width: '100%',
                        height: '12px',
                        backgroundColor: 'var(--gray-a3)',
                        borderRadius: '6px',
                        overflow: 'hidden',
                      }}
                    >
                      <div
                        style={{
                          width: `${usagePercentage}%`,
                          height: '100%',
                          backgroundColor:
                            usagePercentage > 90
                              ? 'var(--red-9)'
                              : usagePercentage > 70
                                ? 'var(--orange-9)'
                                : 'var(--green-9)',
                          transition: 'width 0.3s ease',
                        }}
                      />
                    </div>
                  </Flex>

                  {stats.evictions > 0 && (
                    <Callout.Root color="orange">
                      <Callout.Icon>
                        <InfoCircledIcon />
                      </Callout.Icon>
                      <Callout.Text>
                        {stats.evictions?.toLocaleString() || '0'} entries have been evicted due to
                        capacity limits
                      </Callout.Text>
                    </Callout.Root>
                  )}
                </>
              ) : (
                <Text size="2" color="gray">
                  DNS caching is currently disabled. Enable it in the server configuration to
                  improve query performance.
                </Text>
              )}
            </Flex>
          </Card>
        </>
      )}
    </Flex>
  );
}
