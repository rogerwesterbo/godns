import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
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
  IconButton,
  AlertDialog,
} from '@radix-ui/themes';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  TrashIcon,
  Pencil1Icon,
  ArrowLeftIcon,
  LockClosedIcon,
  LockOpen1Icon,
  ReloadIcon,
} from '@radix-ui/react-icons';
import * as api from '../services/api';
import { RecordDialog, SortableColumnHeader } from '../components';
import { formatRecordValue } from '../utils/recordFormatting';
import { useSortableData } from '../hooks';

export default function ZoneDetailPage() {
  const { domain } = useParams<{ domain: string }>();
  const navigate = useNavigate();
  const [zone, setZone] = useState<api.DNSZone | null>(null);
  const [filteredRecords, setFilteredRecords] = useState<api.DNSRecord[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState('');
  const [typeFilter, setTypeFilter] = useState('All');
  const [currentPage, setCurrentPage] = useState(1);
  const [showRecordDialog, setShowRecordDialog] = useState(false);
  const [editingRecord, setEditingRecord] = useState<api.DNSRecord | undefined>();
  const [recordDialogMode, setRecordDialogMode] = useState<'create' | 'edit'>('create');
  const [deletingRecord, setDeletingRecord] = useState<api.DNSRecord | null>(null);
  const [togglingRecord, setTogglingRecord] = useState<api.DNSRecord | null>(null);
  const itemsPerPage = 15;

  useEffect(() => {
    if (domain) {
      loadZone(decodeURIComponent(domain));
    }
  }, [domain]);

  useEffect(() => {
    if (!zone) return;

    let records = zone.records;

    // Apply type filter
    if (typeFilter !== 'All') {
      records = records.filter(r => r.type === typeFilter);
    }

    // Apply text filter
    if (filter.trim()) {
      records = records.filter(
        r =>
          r.name.toLowerCase().includes(filter.toLowerCase()) ||
          (r.value && r.value.toLowerCase().includes(filter.toLowerCase())) ||
          r.type.toLowerCase().includes(filter.toLowerCase()) ||
          formatRecordValue(r).toLowerCase().includes(filter.toLowerCase())
      );
    }

    setFilteredRecords(records);
    setCurrentPage(1);
  }, [zone, filter, typeFilter]);

  // Sortable data
  const {
    items: sortedRecords,
    requestSort,
    sortConfig,
  } = useSortableData<api.DNSRecord>(filteredRecords, 'name');

  // Pagination
  const totalPages = Math.ceil(sortedRecords.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentRecords = sortedRecords.slice(startIndex, endIndex);

  const loadZone = async (zoneDomain: string) => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.getZone(zoneDomain);
      setZone(data);
      setFilteredRecords(data.records);
    } catch (err) {
      console.error('Failed to load zone:', err);
      setError(err instanceof Error ? err.message : 'Failed to load zone');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteZone = async () => {
    if (!zone) return;

    try {
      await api.deleteZone(zone.domain);
      navigate('/zones');
    } catch (err) {
      console.error('Failed to delete zone:', err);
      setError(err instanceof Error ? err.message : 'Failed to delete zone');
    }
  };

  const handleToggleZoneStatus = async (enabled: boolean) => {
    if (!zone) return;

    try {
      await api.setZoneStatus(zone.domain, enabled);
      // Update local state
      setZone(prevZone => (prevZone ? { ...prevZone, enabled } : null));
    } catch (err) {
      console.error('Failed to update zone status:', err);
      setError(err instanceof Error ? err.message : 'Failed to update zone status');
      // Reload zone to get the correct state
      if (domain) {
        loadZone(domain);
      }
    }
  };

  const handleCreateRecord = () => {
    setRecordDialogMode('create');
    setEditingRecord(undefined);
    setShowRecordDialog(true);
  };

  const handleEditRecord = (record: api.DNSRecord) => {
    setRecordDialogMode('edit');
    setEditingRecord(record);
    setShowRecordDialog(true);
  };

  const handleDeleteRecord = async () => {
    if (!zone || !deletingRecord) return;

    try {
      await api.deleteRecord(zone.domain, deletingRecord.name, deletingRecord.type);
      await loadZone(zone.domain);
      setDeletingRecord(null);
    } catch (err) {
      console.error('Failed to delete record:', err);
      setError(err instanceof Error ? err.message : 'Failed to delete record');
    }
  };

  const handleToggleRecord = async () => {
    if (!zone || !togglingRecord) return;

    try {
      await api.setRecordStatus(
        zone.domain,
        togglingRecord.name,
        togglingRecord.type,
        !!togglingRecord.disabled
      );
      await loadZone(zone.domain);
      setTogglingRecord(null);
    } catch (err) {
      console.error('Failed to toggle record:', err);
      setError(err instanceof Error ? err.message : 'Failed to toggle record status');
    }
  };

  const handleRecordSuccess = async () => {
    if (zone) {
      await loadZone(zone.domain);
    }
  };

  // Get unique record types from zone
  const recordTypes = zone
    ? ['All', ...Array.from(new Set(zone.records.map(r => r.type))).sort()]
    : ['All'];

  const getRecordTypeBadgeColor = (
    type: string
  ): 'blue' | 'green' | 'orange' | 'purple' | 'red' | 'gray' => {
    const colors: Record<string, 'blue' | 'green' | 'orange' | 'purple' | 'red' | 'gray'> = {
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

  if (isLoading) {
    return (
      <Flex direction="column" gap="6">
        <Flex justify="center" py="8">
          <Spinner size="3" />
        </Flex>
      </Flex>
    );
  }

  if (error || !zone) {
    return (
      <Flex direction="column" gap="6">
        <Flex align="center" gap="3">
          <IconButton variant="ghost" onClick={() => navigate('/zones')}>
            <ArrowLeftIcon />
          </IconButton>
          <Heading size="8">Zone Not Found</Heading>
        </Flex>
        <Card>
          <Text color="red" size="3">
            {error || 'Zone not found'}
          </Text>
        </Card>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Flex align="center" gap="3">
          <IconButton variant="ghost" onClick={() => navigate('/zones')}>
            <ArrowLeftIcon />
          </IconButton>
          <Box>
            <Flex align="center" gap="2">
              <Heading size="8">{zone.domain}</Heading>
              <Badge color={(zone.enabled ?? true) ? 'green' : 'red'}>
                {(zone.enabled ?? true) ? 'Active' : 'Disabled'}
              </Badge>
            </Flex>
            <Text size="2" color="gray">
              {zone.records.length} record{zone.records.length !== 1 ? 's' : ''}
            </Text>
          </Box>
        </Flex>
        <Flex gap="2">
          <Button size="3" variant="soft" onClick={() => domain && loadZone(domain)}>
            <ReloadIcon /> Refresh
          </Button>
          <Button size="3" variant="soft" onClick={handleCreateRecord}>
            <PlusIcon /> Add Record
          </Button>
          <AlertDialog.Root>
            <AlertDialog.Trigger>
              <Button size="3" variant="soft" color={(zone.enabled ?? true) ? 'orange' : 'green'}>
                {(zone.enabled ?? true) ? <LockClosedIcon /> : <LockOpen1Icon />}
                {(zone.enabled ?? true) ? 'Disable Zone' : 'Enable Zone'}
              </Button>
            </AlertDialog.Trigger>
            <AlertDialog.Content>
              <AlertDialog.Title>
                {(zone.enabled ?? true) ? 'Disable' : 'Enable'} Zone
              </AlertDialog.Title>
              <AlertDialog.Description>
                {(zone.enabled ?? true) ? (
                  <>
                    Are you sure you want to <strong>disable</strong> the zone{' '}
                    <strong>{zone.domain}</strong>?
                    <br />
                    <br />
                    This will prevent the DNS server from responding to queries for this zone and
                    all its {zone.records.length} record{zone.records.length !== 1 ? 's' : ''}. The
                    zone data will be preserved and can be re-enabled later.
                  </>
                ) : (
                  <>
                    Are you sure you want to <strong>enable</strong> the zone{' '}
                    <strong>{zone.domain}</strong>?
                    <br />
                    <br />
                    This will allow the DNS server to respond to queries for this zone and all its{' '}
                    {zone.records.length} record{zone.records.length !== 1 ? 's' : ''}.
                  </>
                )}
              </AlertDialog.Description>
              <Flex gap="3" mt="4" justify="end">
                <AlertDialog.Cancel>
                  <Button variant="soft" color="gray">
                    Cancel
                  </Button>
                </AlertDialog.Cancel>
                <AlertDialog.Action>
                  <Button
                    color={(zone.enabled ?? true) ? 'orange' : 'green'}
                    onClick={() => handleToggleZoneStatus(!(zone.enabled ?? true))}
                  >
                    {(zone.enabled ?? true) ? 'Disable Zone' : 'Enable Zone'}
                  </Button>
                </AlertDialog.Action>
              </Flex>
            </AlertDialog.Content>
          </AlertDialog.Root>
          <AlertDialog.Root>
            <AlertDialog.Trigger>
              <Button size="3" color="red" variant="soft">
                <TrashIcon /> Delete Zone
              </Button>
            </AlertDialog.Trigger>
            <AlertDialog.Content>
              <AlertDialog.Title>Delete Zone</AlertDialog.Title>
              <AlertDialog.Description>
                Are you sure you want to delete the zone <strong>{zone.domain}</strong>? This will
                remove all {zone.records.length} record{zone.records.length !== 1 ? 's' : ''} and
                cannot be undone.
              </AlertDialog.Description>
              <Flex gap="3" mt="4" justify="end">
                <AlertDialog.Cancel>
                  <Button variant="soft" color="gray">
                    Cancel
                  </Button>
                </AlertDialog.Cancel>
                <AlertDialog.Action>
                  <Button color="red" onClick={handleDeleteZone}>
                    Delete Zone
                  </Button>
                </AlertDialog.Action>
              </Flex>
            </AlertDialog.Content>
          </AlertDialog.Root>
        </Flex>
      </Flex>

      <Card>
        <Flex direction="column" gap="4">
          <Flex justify="between" align="center" gap="3">
            <Text size="2" color="gray">
              DNS records for this zone
            </Text>
            <Flex gap="2" align="center">
              <Box style={{ width: '250px' }}>
                <TextField.Root
                  size="2"
                  placeholder="Filter records..."
                  value={filter}
                  onChange={e => setFilter(e.target.value)}
                >
                  <TextField.Slot>
                    <MagnifyingGlassIcon height="14" width="14" />
                  </TextField.Slot>
                </TextField.Root>
              </Box>
              <Flex gap="1">
                {recordTypes.map(type => (
                  <Button
                    key={type}
                    size="1"
                    variant={typeFilter === type ? 'solid' : 'soft'}
                    onClick={() => setTypeFilter(type)}
                  >
                    {type}
                  </Button>
                ))}
              </Flex>
            </Flex>
          </Flex>

          {currentRecords.length === 0 ? (
            <Box py="8">
              <Text color="gray" align="center">
                {zone.records.length === 0
                  ? 'No records in this zone yet.'
                  : 'No records match your filters.'}
              </Text>
            </Box>
          ) : (
            <>
              <Text size="2" color="gray">
                Showing {startIndex + 1}-{Math.min(endIndex, sortedRecords.length)} of{' '}
                {sortedRecords.length} record{sortedRecords.length !== 1 ? 's' : ''}
                {(filter || typeFilter !== 'All') && ' (filtered)'}
              </Text>

              <Table.Root variant="surface">
                <Table.Header>
                  <Table.Row>
                    <SortableColumnHeader<api.DNSRecord>
                      column="name"
                      currentSortKey={sortConfig.key as keyof api.DNSRecord | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Name
                    </SortableColumnHeader>
                    <SortableColumnHeader<api.DNSRecord>
                      column="type"
                      currentSortKey={sortConfig.key as keyof api.DNSRecord | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Type
                    </SortableColumnHeader>
                    <Table.ColumnHeaderCell>Value</Table.ColumnHeaderCell>
                    <SortableColumnHeader<api.DNSRecord>
                      column="ttl"
                      currentSortKey={sortConfig.key as keyof api.DNSRecord | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      TTL
                    </SortableColumnHeader>
                    <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                  </Table.Row>
                </Table.Header>

                <Table.Body>
                  {currentRecords.map((record, index) => (
                    <Table.Row key={`${record.name}-${record.type}-${index}`}>
                      <Table.Cell>
                        <Text weight="medium">{record.name}</Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Badge color={getRecordTypeBadgeColor(record.type)}>{record.type}</Badge>
                      </Table.Cell>
                      <Table.Cell>
                        <Text
                          size="2"
                          style={{
                            maxWidth: '400px',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            display: 'block',
                          }}
                        >
                          {formatRecordValue(record)}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="2" color="gray">
                          {record.ttl}s
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Badge color={record.disabled ? 'gray' : 'green'}>
                          {record.disabled ? 'Disabled' : 'Active'}
                        </Badge>
                      </Table.Cell>
                      <Table.Cell>
                        <Flex gap="2">
                          <IconButton
                            size="1"
                            variant="ghost"
                            color={record.disabled ? 'green' : 'orange'}
                            onClick={() => setTogglingRecord(record)}
                            title={record.disabled ? 'Enable record' : 'Disable record'}
                          >
                            {record.disabled ? <LockOpen1Icon /> : <LockClosedIcon />}
                          </IconButton>
                          <IconButton
                            size="1"
                            variant="ghost"
                            onClick={() => handleEditRecord(record)}
                          >
                            <Pencil1Icon />
                          </IconButton>
                          <IconButton
                            size="1"
                            variant="ghost"
                            color="red"
                            onClick={() => setDeletingRecord(record)}
                          >
                            <TrashIcon />
                          </IconButton>
                        </Flex>
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

      {zone && (
        <>
          <RecordDialog
            open={showRecordDialog}
            onOpenChange={setShowRecordDialog}
            onSuccess={handleRecordSuccess}
            zoneDomain={zone.domain}
            record={editingRecord}
            mode={recordDialogMode}
          />

          <AlertDialog.Root
            open={!!deletingRecord}
            onOpenChange={open => !open && setDeletingRecord(null)}
          >
            <AlertDialog.Content>
              <AlertDialog.Title>Delete Record</AlertDialog.Title>
              <AlertDialog.Description>
                Are you sure you want to delete the record <strong>{deletingRecord?.name}</strong> (
                {deletingRecord?.type})? This action cannot be undone.
              </AlertDialog.Description>
              <Flex gap="3" mt="4" justify="end">
                <AlertDialog.Cancel>
                  <Button variant="soft" color="gray">
                    Cancel
                  </Button>
                </AlertDialog.Cancel>
                <AlertDialog.Action>
                  <Button color="red" onClick={handleDeleteRecord}>
                    Delete Record
                  </Button>
                </AlertDialog.Action>
              </Flex>
            </AlertDialog.Content>
          </AlertDialog.Root>

          {/* Toggle Record Status Dialog */}
          <AlertDialog.Root
            open={!!togglingRecord}
            onOpenChange={open => !open && setTogglingRecord(null)}
          >
            <AlertDialog.Content>
              <AlertDialog.Title>
                {togglingRecord?.disabled ? 'Enable' : 'Disable'} Record
              </AlertDialog.Title>
              <AlertDialog.Description>
                Are you sure you want to {togglingRecord?.disabled ? 'enable' : 'disable'} the
                record <strong>{togglingRecord?.name}</strong> ({togglingRecord?.type})?
              </AlertDialog.Description>
              {!togglingRecord?.disabled && (
                <Text as="p" mt="2" color="orange">
                  Disabled records will not be served by the DNS server.
                </Text>
              )}
              <Flex gap="3" mt="4" justify="end">
                <AlertDialog.Cancel>
                  <Button variant="soft" color="gray">
                    Cancel
                  </Button>
                </AlertDialog.Cancel>
                <AlertDialog.Action>
                  <Button
                    color={togglingRecord?.disabled ? 'green' : 'orange'}
                    onClick={handleToggleRecord}
                  >
                    {togglingRecord?.disabled ? 'Enable' : 'Disable'} Record
                  </Button>
                </AlertDialog.Action>
              </Flex>
            </AlertDialog.Content>
          </AlertDialog.Root>
        </>
      )}
    </Flex>
  );
}
