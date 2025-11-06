import { useEffect, useState } from 'react';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Table,
  Badge,
  Select,
  Spinner,
  TextField,
  Box,
  AlertDialog,
  IconButton,
} from '@radix-ui/themes';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  Pencil1Icon,
  TrashIcon,
  ReloadIcon,
} from '@radix-ui/react-icons';
import * as api from '../services/api';
import { RecordDialog, SortableColumnHeader } from '../components';
import { formatRecordValue } from '../utils/recordFormatting';
import { useSortableData } from '../hooks';

type RecordWithZone = api.DNSRecord & { zone: string };

export default function RecordsPage() {
  const [zones, setZones] = useState<api.DNSZone[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState('');
  const [typeFilter, setTypeFilter] = useState('All');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 15;

  // CRUD states
  const [showRecordDialog, setShowRecordDialog] = useState(false);
  const [editingRecord, setEditingRecord] = useState<(api.DNSRecord & { zone: string }) | null>(
    null
  );
  const [recordDialogMode, setRecordDialogMode] = useState<'create' | 'edit'>('create');
  const [selectedZoneForCreate, setSelectedZoneForCreate] = useState<string>('');
  const [deletingRecord, setDeletingRecord] = useState<(api.DNSRecord & { zone: string }) | null>(
    null
  );

  useEffect(() => {
    loadZones();
  }, []);

  const loadZones = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.listZones();
      setZones(data);
    } catch (err) {
      console.error('Failed to load zones:', err);
      setError(err instanceof Error ? err.message : 'Failed to load zones');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateRecord = () => {
    if (zones.length === 0) {
      setError('No zones available. Please create a zone first.');
      return;
    }
    setSelectedZoneForCreate(zones[0].domain); // Default to first zone
    setEditingRecord(null);
    setRecordDialogMode('create');
    setShowRecordDialog(true);
  };

  const handleEditRecord = (record: api.DNSRecord & { zone: string }) => {
    setEditingRecord(record);
    setRecordDialogMode('edit');
    setShowRecordDialog(true);
  };

  const handleDeleteRecord = async () => {
    if (!deletingRecord) return;

    try {
      await api.deleteRecord(deletingRecord.zone, deletingRecord.name, deletingRecord.type);
      await loadZones(); // Reload to get updated data
      setDeletingRecord(null);
    } catch (err) {
      console.error('Failed to delete record:', err);
      setError(err instanceof Error ? err.message : 'Failed to delete record');
    }
  };

  const handleRecordSuccess = async () => {
    setShowRecordDialog(false);
    setEditingRecord(null);
    await loadZones(); // Reload to get updated data
  };

  // Flatten all records from all zones
  const allRecords = zones.flatMap(zone =>
    zone.records.map(record => ({
      ...record,
      zone: zone.domain,
    }))
  );

  // Apply filters
  const filteredRecords = allRecords.filter(record => {
    const matchesType = typeFilter === 'All' || record.type === typeFilter;
    const matchesFilter =
      !filter.trim() ||
      record.name.toLowerCase().includes(filter.toLowerCase()) ||
      (record.value && record.value.toLowerCase().includes(filter.toLowerCase())) ||
      record.zone.toLowerCase().includes(filter.toLowerCase()) ||
      formatRecordValue(record).toLowerCase().includes(filter.toLowerCase());
    return matchesType && matchesFilter;
  });

  // Sortable data
  const {
    items: sortedRecords,
    requestSort,
    sortConfig,
  } = useSortableData<RecordWithZone>(filteredRecords, 'name');

  // Pagination
  const totalPages = Math.ceil(sortedRecords.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentRecords = sortedRecords.slice(startIndex, endIndex);

  // Get unique record types
  const recordTypes = ['All', ...Array.from(new Set(allRecords.map(r => r.type))).sort()];

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

  // Reset to page 1 when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [filter, typeFilter]);

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">DNS Records</Heading>
        <Flex gap="2">
          <Button size="3" variant="soft" onClick={loadZones}>
            <ReloadIcon /> Refresh
          </Button>
          <Button size="3" onClick={handleCreateRecord}>
            <PlusIcon /> Add Record
          </Button>
        </Flex>
      </Flex>

      <Card>
        <Flex direction="column" gap="4">
          <Flex justify="between" align="center" gap="3">
            <Text size="2" color="gray">
              Manage DNS records across all zones
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
              <Select.Root value={typeFilter} onValueChange={setTypeFilter}>
                <Select.Trigger placeholder="Filter by type" />
                <Select.Content>
                  {recordTypes.map(type => (
                    <Select.Item key={type} value={type}>
                      {type}
                    </Select.Item>
                  ))}
                </Select.Content>
              </Select.Root>
            </Flex>
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

          {!isLoading && !error && currentRecords.length === 0 && (
            <Box py="8">
              <Text color="gray" align="center">
                {allRecords.length === 0 ? 'No records found.' : 'No records match your filters.'}
              </Text>
            </Box>
          )}

          {!isLoading && !error && currentRecords.length > 0 && (
            <>
              <Text size="2" color="gray">
                Showing {startIndex + 1}-{Math.min(endIndex, sortedRecords.length)} of{' '}
                {sortedRecords.length} record{sortedRecords.length !== 1 ? 's' : ''}
                {(filter || typeFilter !== 'All') && ' (filtered)'}
              </Text>

              <Table.Root variant="surface">
                <Table.Header>
                  <Table.Row>
                    <SortableColumnHeader<RecordWithZone>
                      column="name"
                      currentSortKey={sortConfig.key as keyof RecordWithZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Name
                    </SortableColumnHeader>
                    <SortableColumnHeader<RecordWithZone>
                      column="type"
                      currentSortKey={sortConfig.key as keyof RecordWithZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Type
                    </SortableColumnHeader>
                    <Table.ColumnHeaderCell>Value</Table.ColumnHeaderCell>
                    <SortableColumnHeader<RecordWithZone>
                      column="ttl"
                      currentSortKey={sortConfig.key as keyof RecordWithZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      TTL
                    </SortableColumnHeader>
                    <SortableColumnHeader<RecordWithZone>
                      column="zone"
                      currentSortKey={sortConfig.key as keyof RecordWithZone | null}
                      currentSortDirection={sortConfig.direction}
                      onSort={col => requestSort(col as string)}
                    >
                      Zone
                    </SortableColumnHeader>
                    <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                    <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
                  </Table.Row>
                </Table.Header>

                <Table.Body>
                  {currentRecords.map((record, index) => (
                    <Table.Row key={`${record.zone}-${record.name}-${record.type}-${index}`}>
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
                            maxWidth: '300px',
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
                        <Text size="2" color="gray">
                          {record.zone}
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
                            color="blue"
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

      {/* Record Dialog for Create/Edit */}
      <RecordDialog
        open={showRecordDialog}
        onOpenChange={setShowRecordDialog}
        mode={recordDialogMode}
        record={editingRecord || undefined}
        zoneDomain={
          recordDialogMode === 'create' ? selectedZoneForCreate : editingRecord?.zone || ''
        }
        onSuccess={handleRecordSuccess}
        onZoneChange={recordDialogMode === 'create' ? setSelectedZoneForCreate : undefined}
        availableZones={recordDialogMode === 'create' ? zones.map(z => z.domain) : undefined}
      />

      {/* Delete Record Confirmation */}
      <AlertDialog.Root
        open={!!deletingRecord}
        onOpenChange={open => !open && setDeletingRecord(null)}
      >
        <AlertDialog.Content>
          <AlertDialog.Title>Delete Record</AlertDialog.Title>
          <AlertDialog.Description>
            Are you sure you want to delete the record <strong>{deletingRecord?.name}</strong> (
            {deletingRecord?.type}) from zone <strong>{deletingRecord?.zone}</strong>? This action
            cannot be undone.
          </AlertDialog.Description>
          <Flex gap="3" mt="4" justify="end">
            <AlertDialog.Cancel>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </AlertDialog.Cancel>
            <AlertDialog.Action>
              <Button variant="solid" color="red" onClick={handleDeleteRecord}>
                Delete Record
              </Button>
            </AlertDialog.Action>
          </Flex>
        </AlertDialog.Content>
      </AlertDialog.Root>
    </Flex>
  );
}
