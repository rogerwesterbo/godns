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

const RECORD_TYPES = ['A', 'AAAA', 'CNAME', 'ALIAS', 'MX', 'NS', 'TXT', 'PTR', 'SRV', 'SOA', 'CAA'];

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

  // MX fields
  const [mxPriority, setMxPriority] = useState('10');
  const [mxHost, setMxHost] = useState('');

  // SRV fields
  const [srvPriority, setSrvPriority] = useState('10');
  const [srvWeight, setSrvWeight] = useState('60');
  const [srvPort, setSrvPort] = useState('');
  const [srvTarget, setSrvTarget] = useState('');

  // SOA fields
  const [soaMname, setSoaMname] = useState('');
  const [soaRname, setSoaRname] = useState('');
  const [soaSerial, setSoaSerial] = useState('2024110601');
  const [soaRefresh, setSoaRefresh] = useState('3600');
  const [soaRetry, setSoaRetry] = useState('1800');
  const [soaExpire, setSoaExpire] = useState('604800');
  const [soaMinimum, setSoaMinimum] = useState('300');

  // CAA fields
  const [caaFlags, setCaaFlags] = useState('0');
  const [caaTag, setCaaTag] = useState('issue');
  const [caaValue, setCaaValue] = useState('');

  useEffect(() => {
    if (record && mode === 'edit') {
      setName(record.name);
      setType(record.type);
      setValue(record.value || '');
      setTtl(String(record.ttl));

      // MX fields
      setMxPriority(String(record.mx_priority || 10));
      setMxHost(record.mx_host || '');

      // SRV fields
      setSrvPriority(String(record.srv_priority || 10));
      setSrvWeight(String(record.srv_weight || 60));
      setSrvPort(String(record.srv_port || ''));
      setSrvTarget(record.srv_target || '');

      // SOA fields
      setSoaMname(record.soa_mname || '');
      setSoaRname(record.soa_rname || '');
      setSoaSerial(String(record.soa_serial || 2024110601));
      setSoaRefresh(String(record.soa_refresh || 3600));
      setSoaRetry(String(record.soa_retry || 1800));
      setSoaExpire(String(record.soa_expire || 604800));
      setSoaMinimum(String(record.soa_minimum || 300));

      // CAA fields
      setCaaFlags(String(record.caa_flags || 0));
      setCaaTag(record.caa_tag || 'issue');
      setCaaValue(record.caa_value || '');
    } else {
      setName('');
      setType('A');
      setValue('');
      setTtl('300');

      // Reset MX fields
      setMxPriority('10');
      setMxHost('');

      // Reset SRV fields
      setSrvPriority('10');
      setSrvWeight('60');
      setSrvPort('');
      setSrvTarget('');

      // Reset SOA fields
      setSoaMname('');
      setSoaRname('');
      setSoaSerial('2024110601');
      setSoaRefresh('3600');
      setSoaRetry('1800');
      setSoaExpire('604800');
      setSoaMinimum('300');

      // Reset CAA fields
      setCaaFlags('0');
      setCaaTag('issue');
      setCaaValue('');
    }
    setError(null);
  }, [record, mode, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      setError('Name is required');
      return;
    }

    // Validate based on record type
    if (type === 'MX') {
      if (!mxHost.trim()) {
        setError('MX host is required');
        return;
      }
      const priority = parseInt(mxPriority);
      if (isNaN(priority) || priority < 0 || priority > 65535) {
        setError('MX priority must be between 0 and 65535');
        return;
      }
    } else if (type === 'SRV') {
      if (!srvTarget.trim() || !srvPort.trim()) {
        setError('SRV target and port are required');
        return;
      }
      const priority = parseInt(srvPriority);
      const weight = parseInt(srvWeight);
      const port = parseInt(srvPort);
      if (isNaN(priority) || priority < 0 || priority > 65535) {
        setError('SRV priority must be between 0 and 65535');
        return;
      }
      if (isNaN(weight) || weight < 0 || weight > 65535) {
        setError('SRV weight must be between 0 and 65535');
        return;
      }
      if (isNaN(port) || port < 0 || port > 65535) {
        setError('SRV port must be between 0 and 65535');
        return;
      }
    } else if (type === 'SOA') {
      if (!soaMname.trim() || !soaRname.trim()) {
        setError('SOA mname and rname are required');
        return;
      }
      const serial = parseInt(soaSerial);
      const refresh = parseInt(soaRefresh);
      const retry = parseInt(soaRetry);
      const expire = parseInt(soaExpire);
      const minimum = parseInt(soaMinimum);
      if (isNaN(serial) || isNaN(refresh) || isNaN(retry) || isNaN(expire) || isNaN(minimum)) {
        setError('All SOA numeric fields must be valid numbers');
        return;
      }
    } else if (type === 'CAA') {
      if (!caaValue.trim()) {
        setError('CAA value is required');
        return;
      }
      const flags = parseInt(caaFlags);
      if (isNaN(flags) || (flags !== 0 && flags !== 128)) {
        setError('CAA flags must be 0 or 128');
        return;
      }
    } else {
      // Simple record types require value
      if (!value.trim()) {
        setError('Value is required');
        return;
      }
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
        ttl: ttlNum,
      };

      // Add type-specific fields
      if (type === 'MX') {
        recordData.mx_priority = parseInt(mxPriority);
        recordData.mx_host = mxHost.trim();
      } else if (type === 'SRV') {
        recordData.srv_priority = parseInt(srvPriority);
        recordData.srv_weight = parseInt(srvWeight);
        recordData.srv_port = parseInt(srvPort);
        recordData.srv_target = srvTarget.trim();
      } else if (type === 'SOA') {
        recordData.soa_mname = soaMname.trim();
        recordData.soa_rname = soaRname.trim();
        recordData.soa_serial = parseInt(soaSerial);
        recordData.soa_refresh = parseInt(soaRefresh);
        recordData.soa_retry = parseInt(soaRetry);
        recordData.soa_expire = parseInt(soaExpire);
        recordData.soa_minimum = parseInt(soaMinimum);
      } else if (type === 'CAA') {
        recordData.caa_flags = parseInt(caaFlags);
        recordData.caa_tag = caaTag;
        recordData.caa_value = caaValue.trim();
      } else {
        recordData.value = value.trim();
      }

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

            {/* Simple record types: A, AAAA, CNAME, ALIAS, NS, TXT, PTR */}
            {['A', 'AAAA', 'CNAME', 'ALIAS', 'NS', 'TXT', 'PTR'].includes(type) && (
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
                        : type === 'CNAME' || type === 'ALIAS'
                          ? 'target.example.com.'
                          : type === 'TXT'
                            ? '"v=spf1 include:_spf.example.com ~all"'
                            : type === 'PTR'
                              ? 'hostname.example.com.'
                              : 'Record value'
                  }
                  value={value}
                  onChange={e => setValue(e.target.value)}
                  disabled={isSubmitting}
                />
                <Text as="div" size="1" color="gray" mt="1">
                  {(type === 'A' || type === 'AAAA') && 'IP address'}
                  {(type === 'CNAME' || type === 'ALIAS') && 'Target domain name'}
                  {type === 'NS' && 'Nameserver hostname'}
                  {type === 'TXT' && 'Text value (enclose in quotes if contains spaces)'}
                  {type === 'PTR' && 'Reverse DNS hostname'}
                </Text>
              </label>
            )}

            {/* MX Record Fields */}
            {type === 'MX' && (
              <>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Priority
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="10"
                    value={mxPriority}
                    onChange={e => setMxPriority(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                    max="65535"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Mail server priority (0-65535, lower is higher priority)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Mail Host
                  </Text>
                  <TextField.Root
                    placeholder="mail.example.com."
                    value={mxHost}
                    onChange={e => setMxHost(e.target.value)}
                    disabled={isSubmitting}
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Fully qualified domain name of the mail server
                  </Text>
                </label>
              </>
            )}

            {/* SRV Record Fields */}
            {type === 'SRV' && (
              <>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Priority
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="10"
                    value={srvPriority}
                    onChange={e => setSrvPriority(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                    max="65535"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Service priority (0-65535, lower is higher priority)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Weight
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="60"
                    value={srvWeight}
                    onChange={e => setSrvWeight(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                    max="65535"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Load balancing weight (0-65535)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Port
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="80"
                    value={srvPort}
                    onChange={e => setSrvPort(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                    max="65535"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Service port number (0-65535)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Target
                  </Text>
                  <TextField.Root
                    placeholder="service.example.com."
                    value={srvTarget}
                    onChange={e => setSrvTarget(e.target.value)}
                    disabled={isSubmitting}
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Target hostname providing the service
                  </Text>
                </label>
              </>
            )}

            {/* SOA Record Fields */}
            {type === 'SOA' && (
              <>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Primary Nameserver (MNAME)
                  </Text>
                  <TextField.Root
                    placeholder="ns1.example.com."
                    value={soaMname}
                    onChange={e => setSoaMname(e.target.value)}
                    disabled={isSubmitting}
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Primary nameserver for this zone
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Admin Email (RNAME)
                  </Text>
                  <TextField.Root
                    placeholder="hostmaster.example.com."
                    value={soaRname}
                    onChange={e => setSoaRname(e.target.value)}
                    disabled={isSubmitting}
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Email of zone administrator (@ replaced with .)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Serial Number
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="2024110601"
                    value={soaSerial}
                    onChange={e => setSoaSerial(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Zone version (format: YYYYMMDDnn)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Refresh Interval
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="3600"
                    value={soaRefresh}
                    onChange={e => setSoaRefresh(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Secondary nameserver refresh interval (seconds)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Retry Interval
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="1800"
                    value={soaRetry}
                    onChange={e => setSoaRetry(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Retry interval if refresh fails (seconds)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Expire Time
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="604800"
                    value={soaExpire}
                    onChange={e => setSoaExpire(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    When zone data expires (seconds)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Minimum TTL
                  </Text>
                  <TextField.Root
                    type="number"
                    placeholder="300"
                    value={soaMinimum}
                    onChange={e => setSoaMinimum(e.target.value)}
                    disabled={isSubmitting}
                    min="0"
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    Minimum TTL for negative caching (seconds)
                  </Text>
                </label>
              </>
            )}

            {/* CAA Record Fields */}
            {type === 'CAA' && (
              <>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Flags
                  </Text>
                  <Select.Root value={caaFlags} onValueChange={setCaaFlags} disabled={isSubmitting}>
                    <Select.Trigger placeholder="Select flags" />
                    <Select.Content>
                      <Select.Item value="0">0 (Non-critical)</Select.Item>
                      <Select.Item value="128">128 (Critical)</Select.Item>
                    </Select.Content>
                  </Select.Root>
                  <Text as="div" size="1" color="gray" mt="1">
                    CAA flags (0 = non-critical, 128 = critical)
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Tag
                  </Text>
                  <Select.Root value={caaTag} onValueChange={setCaaTag} disabled={isSubmitting}>
                    <Select.Trigger placeholder="Select tag" />
                    <Select.Content>
                      <Select.Item value="issue">issue (Authorize CA)</Select.Item>
                      <Select.Item value="issuewild">issuewild (Wildcard certs)</Select.Item>
                      <Select.Item value="iodef">iodef (Incident reporting)</Select.Item>
                    </Select.Content>
                  </Select.Root>
                  <Text as="div" size="1" color="gray" mt="1">
                    CAA property tag
                  </Text>
                </label>
                <label>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Value
                  </Text>
                  <TextField.Root
                    placeholder={
                      caaTag === 'iodef' ? 'mailto:security@example.com' : 'letsencrypt.org'
                    }
                    value={caaValue}
                    onChange={e => setCaaValue(e.target.value)}
                    disabled={isSubmitting}
                  />
                  <Text as="div" size="1" color="gray" mt="1">
                    {caaTag === 'iodef'
                      ? 'URL or email for incident reporting'
                      : 'Certificate authority domain name'}
                  </Text>
                </label>
              </>
            )}

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
