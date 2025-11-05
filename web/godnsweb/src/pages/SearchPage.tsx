import { useEffect, useState } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import {
  Box,
  Heading,
  TextField,
  Card,
  Text,
  Badge,
  Flex,
  Button,
  Spinner,
} from '@radix-ui/themes';
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';
import * as api from '../services/api';

export function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [query, setQuery] = useState(searchParams.get('q') || '');
  const [results, setResults] = useState<api.SearchResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const q = searchParams.get('q');
    if (q) {
      setQuery(q);
      performSearch(q);
    }
  }, [searchParams]);

  const performSearch = async (searchQuery: string) => {
    if (!searchQuery.trim()) {
      setResults(null);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const response = await api.search(searchQuery.trim());
      setResults(response);
    } catch (err) {
      console.error('Search failed:', err);
      setError(err instanceof Error ? err.message : 'Search failed');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      setSearchParams({ q: query.trim() });
    }
  };

  const getRecordTypeBadgeColor = (type: string) => {
    const colors: Record<string, 'blue' | 'green' | 'orange' | 'purple' | 'red'> = {
      A: 'blue',
      AAAA: 'blue',
      CNAME: 'green',
      MX: 'orange',
      TXT: 'purple',
      NS: 'red',
      SOA: 'red',
    };
    return colors[type] || 'gray';
  };

  return (
    <Box p="6">
      <Heading size="8" mb="6">
        Search
      </Heading>

      <form onSubmit={handleSearch}>
        <Flex gap="3" mb="6">
          <Box style={{ flex: 1 }}>
            <TextField.Root
              size="3"
              placeholder="Search zones and records..."
              value={query}
              onChange={e => setQuery(e.target.value)}
            >
              <TextField.Slot>
                <MagnifyingGlassIcon height="16" width="16" />
              </TextField.Slot>
            </TextField.Root>
          </Box>
          <Button size="3" type="submit">
            Search
          </Button>
        </Flex>
      </form>

      {isLoading && (
        <Flex justify="center" py="8">
          <Spinner size="3" />
        </Flex>
      )}

      {error && (
        <Card>
          <Text color="red" size="3">
            {error}
          </Text>
        </Card>
      )}

      {!isLoading && results && (
        <>
          <Text size="2" color="gray" mb="4">
            Found {results.total} result{results.total !== 1 ? 's' : ''} for "{results.query}"
          </Text>

          {results.results.length === 0 ? (
            <Card>
              <Text color="gray">No results found. Try a different search query.</Text>
            </Card>
          ) : (
            <Flex direction="column" gap="3">
              {results.results.map((result, index) => (
                <Card key={index}>
                  <Flex direction="column" gap="2">
                    <Flex justify="between" align="center">
                      <Badge color={result.type === 'zone' ? 'blue' : 'green'} size="1">
                        {result.type}
                      </Badge>
                    </Flex>

                    {result.type === 'zone' && result.zone && (
                      <Box>
                        <Link to={`/zones/${encodeURIComponent(result.zone.domain)}`}>
                          <Text size="5" weight="bold" style={{ cursor: 'pointer' }}>
                            {result.zone.domain}
                          </Text>
                        </Link>
                        <Text size="2" color="gray" mt="1">
                          {result.zone.records?.length || 0} record
                          {result.zone.records?.length !== 1 ? 's' : ''}
                        </Text>
                      </Box>
                    )}

                    {result.type === 'record' && result.record && (
                      <Box>
                        <Flex gap="2" align="center" mb="1">
                          <Text size="4" weight="bold">
                            {result.record.name}
                          </Text>
                          <Badge color={getRecordTypeBadgeColor(result.record.type)} size="1">
                            {result.record.type}
                          </Badge>
                        </Flex>
                        <Text size="2" color="gray">
                          Zone:{' '}
                          <Link to={`/zones/${encodeURIComponent(result.record.zone)}`}>
                            {result.record.zone}
                          </Link>
                        </Text>
                        <Flex gap="4" mt="2">
                          <Text size="2">
                            <Text color="gray">Value:</Text> {result.record.value}
                          </Text>
                          <Text size="2">
                            <Text color="gray">TTL:</Text> {result.record.ttl}s
                          </Text>
                        </Flex>
                      </Box>
                    )}

                    {result.highlight && (
                      <Text size="2" color="gray" style={{ fontStyle: 'italic' }}>
                        {result.highlight}
                      </Text>
                    )}
                  </Flex>
                </Card>
              ))}
            </Flex>
          )}
        </>
      )}

      {!isLoading && !results && !error && query && (
        <Card>
          <Text color="gray">Enter a search query to find zones and records.</Text>
        </Card>
      )}
    </Box>
  );
}
