import { useState, useEffect } from 'react';
import { Dialog, Flex, TextField, Button, Text, Select } from '@radix-ui/themes';
import * as api from '../services/api';

interface RecordDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
  zoneDomain: string;
  record?: api.DNSRecord;
  mode: 'create' | 'edit';
  // Optional props for zone selection in create mode
  availableZones?: string[];
  onZoneChange?: (zone: string) => void;
}

const RECORD_TYPES = ['A', 'AAAA', 'CNAME', 'MX', 'NS', 'TXT', 'PTR', 'SRV', 'SOA', 'CAA'];

export function RecordDialog({
  open,
  onOpenChange,
  onSuccess,
  zoneDomain,
  record,
  mode,
  availableZones,
  onZoneChange,
}: RecordDialogProps) {
  const [name, setName] = useState('');
  const [type, setType] = useState('A');
  const [value, setValue] = useState('');
  const [ttl, setTtl] = useState('300');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (record && mode === 'edit') {
      setName(record.name);
      setType(record.type);
      setValue(record.value);
      setTtl(String(record.ttl));
    } else {
      setName('');
      setType('A');
      setValue('');
      setTtl('300');
    }
    setError(null);
  }, [record, mode, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim() || !value.trim()) {
      setError('Name and value are required');
      return;
    }

    const ttlNum = parseInt(ttl);
    if (isNaN(ttlNum) || ttlNum < 0) {
      setError('TTL must be a positive number');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const recordData: api.DNSRecord = {
        name: name.trim(),
        type,
        value: value.trim(),
        ttl: ttlNum,
      };

      if (mode === 'create') {
        await api.createRecord(zoneDomain, recordData);
      } else if (record) {
        await api.updateRecord(zoneDomain, record.name, record.type, recordData);
      }

      onSuccess();
      onOpenChange(false);
    } catch (err) {
      console.error(`Failed to ${mode} record:`, err);
      setError(err instanceof Error ? err.message : `Failed to ${mode} record`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Content style={{ maxWidth: 500 }}>
        <Dialog.Title>{mode === 'create' ? 'Create New Record' : 'Edit Record'}</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          {mode === 'create'
            ? `Add a new DNS record to ${zoneDomain}`
            : `Update DNS record in ${zoneDomain}`}
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="3">
            {mode === 'create' && availableZones && availableZones.length > 0 && onZoneChange && (
              <label>
                <Text as="div" size="2" mb="1" weight="bold">
                  Zone
                </Text>
                <Select.Root
                  value={zoneDomain}
                  onValueChange={onZoneChange}
                  disabled={isSubmitting}
                >
                  <Select.Trigger placeholder="Select zone" />
                  <Select.Content>
                    {availableZones.map(zone => (
                      <Select.Item key={zone} value={zone}>
                        {zone}
                      </Select.Item>
                    ))}
                  </Select.Content>
                </Select.Root>
                <Text as="div" size="1" color="gray" mt="1">
                  Select the zone where this record will be created
                </Text>
              </label>
            )}

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Name
              </Text>
              <TextField.Root
                placeholder="www.example.com."
                value={name}
                onChange={e => setName(e.target.value)}
                disabled={isSubmitting || mode === 'edit'}
              />
              <Text as="div" size="1" color="gray" mt="1">
                Fully qualified domain name (FQDN)
              </Text>
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Type
              </Text>
              <Select.Root
                value={type}
                onValueChange={setType}
                disabled={isSubmitting || mode === 'edit'}
              >
                <Select.Trigger placeholder="Select record type" />
                <Select.Content>
                  {RECORD_TYPES.map(t => (
                    <Select.Item key={t} value={t}>
                      {t}
                    </Select.Item>
                  ))}
                </Select.Content>
              </Select.Root>
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Value
              </Text>
              <TextField.Root
                placeholder={
                  type === 'A'
                    ? '192.168.1.1'
                    : type === 'AAAA'
                      ? '2001:db8::1'
                      : type === 'CNAME'
                        ? 'target.example.com.'
                        : type === 'MX'
                          ? '10 mail.example.com.'
                          : type === 'TXT'
                            ? '"v=spf1 include:_spf.example.com ~all"'
                            : 'Record value'
                }
                value={value}
                onChange={e => setValue(e.target.value)}
                disabled={isSubmitting}
              />
              <Text as="div" size="1" color="gray" mt="1">
                {type === 'MX' && 'Format: priority hostname (e.g., 10 mail.example.com.)'}
                {type === 'TXT' && 'Enclose in quotes if contains spaces'}
                {(type === 'A' || type === 'AAAA') && 'IP address'}
                {type === 'CNAME' && 'Target domain name'}
              </Text>
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                TTL (seconds)
              </Text>
              <TextField.Root
                type="number"
                placeholder="300"
                value={ttl}
                onChange={e => setTtl(e.target.value)}
                disabled={isSubmitting}
                min="0"
              />
              <Text as="div" size="1" color="gray" mt="1">
                Time to live - how long the record is cached (default: 300)
              </Text>
            </label>

            {error && (
              <Text color="red" size="2">
                {error}
              </Text>
            )}

            <Flex gap="3" mt="4" justify="end">
              <Dialog.Close>
                <Button variant="soft" color="gray" disabled={isSubmitting}>
                  Cancel
                </Button>
              </Dialog.Close>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting
                  ? mode === 'create'
                    ? 'Creating...'
                    : 'Updating...'
                  : mode === 'create'
                    ? 'Create Record'
                    : 'Update Record'}
              </Button>
            </Flex>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  );
}
