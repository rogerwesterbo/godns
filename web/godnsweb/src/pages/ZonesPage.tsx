import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Table,
  Badge,
  Spinner,
  TextField,
  Box,
} from '@radix-ui/themes';
import { PlusIcon, MagnifyingGlassIcon, ReloadIcon } from '@radix-ui/react-icons';
import * as api from '../services/api';
import { CreateZoneDialog, SortableColumnHeader } from '../components';
import { useSortableData } from '../hooks';

export default function ZonesPage() {
  const [zones, setZones] = useState<api.DNSZone[]>([]);
  const [filteredZones, setFilteredZones] = useState<api.DNSZone[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const itemsPerPage = 10;

  // Sortable data hook with custom comparator for records count
  const {
    items: sortedZones,
    requestSort,
    sortConfig,
  } = useSortableData(filteredZones, 'domain', 'asc', {
    recordCount: (a, b) => a.records.length - b.records.length,
  });

  useEffect(() => {
    loadZones();
  }, []);

  useEffect(() => {
    if (filter.trim()) {
      const filtered = zones.filter(zone =>
        zone.domain.toLowerCase().includes(filter.toLowerCase())
      );
      setFilteredZones(filtered);
    } else {
      setFilteredZones(zones);
    }
    setCurrentPage(1);
  }, [filter, zones]);

  const loadZones = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.listZones();
      setZones(data);
      setFilteredZones(data);
    } catch (err) {
      console.error('Failed to load zones:', err);
      setError(err instanceof Error ? err.message : 'Failed to load zones');
    } finally {
      setIsLoading(false);
    }
  };

  // Pagination
  const totalPages = Math.ceil(sortedZones.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentZones = sortedZones.slice(startIndex, endIndex);

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">DNS Zones</Heading>
        <Flex gap="2">
          <Button size="3" variant="soft" onClick={loadZones}>
            <ReloadIcon /> Refresh
          </Button>
          <Button size="3" onClick={() => setShowCreateDialog(true)}>
            <PlusIcon /> Add Zone
          </Button>
        </Flex>
      </Flex>

      <CreateZoneDialog
        open={showCreateDialog}
        onOpenChange={setShowCreateDialog}
        onSuccess={loadZones}
      />

      <Card>
        <Flex direction="column" gap="4">
          <Flex justify="between" align="center">
            <Text size="2" color="gray">
              Manage your DNS zones and their configurations
            </Text>
            <Box style={{ width: '300px' }}>
              <TextField.Root
                size="2"
                placeholder="Filter zones..."
                value={filter}
                onChange={e => setFilter(e.target.value)}
              >
                <TextField.Slot>
                  <MagnifyingGlassIcon height="14" width="14" />
                </TextField.Slot>
              </TextField.Root>
            </Box>
          </Flex>

          {isLoading && (
            <Flex justify="center" py="8">
              <Spinner size="3" />
            </Flex>
          )}

          {error && (
            <Text color="red" size="3">
              {error}
            </Text>
          )}

          {!isLoading && !error && currentZones.length === 0 && (
            <Box py="8">
              <Text color="gray" align="center">
                {filter ? 'No zones match your filter.' : 'No zones configured yet.'}
              </Text>
            </Box>
          )}

          {!isLoading && !error && currentZones.length > 0 && (
            <>
              <Text size="2" color="gray">
                Showing {startIndex + 1}-{Math.min(endIndex, sortedZones.length)} of{' '}
                {sortedZones.length} zone{sortedZones.length !== 1 ? 's' : ''}
                {filter && ` matching "${filter}"`}
              </Text>

              <Table.Root variant="surface">
                <Table.Header>
                  <Table.Row>
                    <SortableColumnHeader
                      column={'domain' as keyof api.DNSZone}
                      currentSortKey={sortConfig.key as keyof api.DNSZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Zone Name
                    </SortableColumnHeader>
                    <SortableColumnHeader
                      column={'recordCount' as keyof api.DNSZone}
                      currentSortKey={sortConfig.key as keyof api.DNSZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Records
                    </SortableColumnHeader>
                    <SortableColumnHeader
                      column={'enabled' as keyof api.DNSZone}
                      currentSortKey={sortConfig.key as keyof api.DNSZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Status
                    </SortableColumnHeader>
                  </Table.Row>
                </Table.Header>

                <Table.Body>
                  {currentZones.map(zone => (
                    <Table.Row key={zone.domain}>
                      <Table.Cell>
                        <Link
                          to={`/zones/${encodeURIComponent(zone.domain)}`}
                          style={{ textDecoration: 'none' }}
                        >
                          <Text weight="medium">{zone.domain}</Text>
                        </Link>
                      </Table.Cell>
                      <Table.Cell>{zone.records.length}</Table.Cell>
                      <Table.Cell>
                        <Badge color={(zone.enabled ?? true) ? 'green' : 'red'}>
                          {(zone.enabled ?? true) ? 'Active' : 'Disabled'}
                        </Badge>
                      </Table.Cell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table.Root>

              {totalPages > 1 && (
                <Flex justify="center" gap="2" mt="2">
                  <Button
                    size="2"
                    variant="soft"
                    disabled={currentPage === 1}
                    onClick={() => setCurrentPage(p => p - 1)}
                  >
                    Previous
                  </Button>
                  <Flex align="center" gap="1">
                    {Array.from({ length: Math.min(totalPages, 7) }, (_, i) => {
                      let page;
                      if (totalPages <= 7) {
                        page = i + 1;
                      } else if (currentPage <= 4) {
                        page = i + 1;
                      } else if (currentPage >= totalPages - 3) {
                        page = totalPages - 6 + i;
                      } else {
                        page = currentPage - 3 + i;
                      }
                      return (
                        <Button
                          key={page}
                          size="2"
                          variant={page === currentPage ? 'solid' : 'soft'}
                          onClick={() => setCurrentPage(page)}
                        >
                          {page}
                        </Button>
                      );
                    })}
                  </Flex>
                  <Button
                    size="2"
                    variant="soft"
                    disabled={currentPage === totalPages}
                    onClick={() => setCurrentPage(p => p + 1)}
                  >
                    Next
                  </Button>
                </Flex>
              )}
            </>
          )}
        </Flex>
      </Card>
    </Flex>
  );
}
