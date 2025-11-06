import { useState, useEffect } from 'react';
import {
  Flex,
  Card,
  Heading,
  Text,
  Button,
  Select,
  TextArea,
  Spinner,
  Callout,
  Table,
  Box,
} from '@radix-ui/themes';
import { DownloadIcon, CodeIcon, InfoCircledIcon } from '@radix-ui/react-icons';
import * as api from '../services/api';

export default function ExportPage() {
  const [zones, setZones] = useState<api.DNSZone[]>([]);
  const [selectedZone, setSelectedZone] = useState<string>('all');
  const [format, setFormat] = useState<string>('bind');
  const [exportedData, setExportedData] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const [isExporting, setIsExporting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadZones();
  }, []);

  const loadZones = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.listZones();
      // Filter to only show enabled zones
      const enabledZones = data.filter(zone => zone.enabled);
      setZones(enabledZones);
    } catch (err) {
      console.error('Failed to load zones:', err);
      setError(err instanceof Error ? err.message : 'Failed to load zones');
    } finally {
      setIsLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      setIsExporting(true);
      setError(null);
      setExportedData('');

      let data: string;
      if (selectedZone === 'all') {
        data = await api.exportAllZones(format);
      } else {
        data = await api.exportZone(selectedZone, format);
      }

      setExportedData(data);
    } catch (err) {
      console.error('Failed to export:', err);
      setError(err instanceof Error ? err.message : 'Failed to export zones');
    } finally {
      setIsExporting(false);
    }
  };

  const handleDownload = () => {
    if (!exportedData) return;

    const zoneName = selectedZone === 'all' ? 'all-zones' : selectedZone.replace('.', '-');
    const filename = `${zoneName}-${format}.txt`;
    const blob = new Blob([exportedData], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const handleCopy = async () => {
    if (!exportedData) return;
    
    try {
      await navigator.clipboard.writeText(exportedData);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  };

  const formatDescriptions = {
    bind: 'BIND zone file format - Standard DNS zone file format used by BIND and other DNS servers',
    coredns: 'CoreDNS Corefile format - Configuration format for CoreDNS',
    powerdns: 'PowerDNS zone format - Format compatible with PowerDNS',
    zonefile: 'Generic zone file format - Standard RFC 1035 compliant zone file',
  };

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="center">
        <Heading size="8">Export Zones</Heading>
      </Flex>

      <Callout.Root color="blue">
        <Callout.Icon>
          <InfoCircledIcon />
        </Callout.Icon>
        <Callout.Text>
          Export your DNS zones in various formats compatible with different DNS servers.
          Only <strong>enabled zones</strong> can be exported. Disabled zones are automatically excluded.
        </Callout.Text>
      </Callout.Root>

      <Card>
        <Flex direction="column" gap="4">
          <Heading size="5">Export Configuration</Heading>

          {isLoading ? (
            <Flex justify="center" py="4">
              <Spinner size="3" />
            </Flex>
          ) : error && zones.length === 0 ? (
            <Text color="red" size="3">
              {error}
            </Text>
          ) : (
            <Flex direction="column" gap="4">
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">
                  Select Zone
                </Text>
                <Select.Root value={selectedZone} onValueChange={setSelectedZone}>
                  <Select.Trigger placeholder="Choose a zone..." />
                  <Select.Content>
                    <Select.Item value="all">All Zones</Select.Item>
                    <Select.Separator />
                    {zones.map(zone => (
                      <Select.Item key={zone.domain} value={zone.domain}>
                        {zone.domain} ({zone.records.length} records)
                      </Select.Item>
                    ))}
                  </Select.Content>
                </Select.Root>
              </Flex>

              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">
                  Export Format
                </Text>
                <Select.Root value={format} onValueChange={setFormat}>
                  <Select.Trigger />
                  <Select.Content>
                    <Select.Item value="bind">BIND (Zone File)</Select.Item>
                    <Select.Item value="zonefile">Generic Zone File</Select.Item>
                    <Select.Item value="coredns">CoreDNS</Select.Item>
                    <Select.Item value="powerdns">PowerDNS</Select.Item>
                  </Select.Content>
                </Select.Root>
                <Text size="1" color="gray">
                  {formatDescriptions[format as keyof typeof formatDescriptions]}
                </Text>
              </Flex>

              <Button size="3" onClick={handleExport} disabled={isExporting || !selectedZone}>
                <CodeIcon />
                {isExporting ? 'Exporting...' : 'Export'}
              </Button>
            </Flex>
          )}
        </Flex>
      </Card>

      {error && exportedData === '' && (
        <Callout.Root color="red">
          <Callout.Text>{error}</Callout.Text>
        </Callout.Root>
      )}

      {exportedData && (
        <Card>
          <Flex direction="column" gap="4">
            <Flex justify="between" align="center">
              <Heading size="5">Exported Configuration</Heading>
              <Flex gap="2">
                <Button size="2" variant="soft" onClick={handleCopy}>
                  Copy to Clipboard
                </Button>
                <Button size="2" onClick={handleDownload}>
                  <DownloadIcon /> Download
                </Button>
              </Flex>
            </Flex>

            <Box
              style={{
                border: '1px solid var(--gray-a6)',
                borderRadius: 'var(--radius-3)',
                overflow: 'hidden',
              }}
            >
              <TextArea
                value={exportedData}
                readOnly
                style={{
                  minHeight: '400px',
                  fontFamily: 'monospace',
                  fontSize: '13px',
                  resize: 'vertical',
                }}
              />
            </Box>

            <Text size="2" color="gray">
              {exportedData.split('\n').length} lines •{' '}
              {new Blob([exportedData]).size.toLocaleString()} bytes
            </Text>
          </Flex>
        </Card>
      )}

      <Card>
        <Flex direction="column" gap="4">
          <Heading size="5">Format Reference</Heading>
          <Table.Root variant="surface">
            <Table.Header>
              <Table.Row>
                <Table.ColumnHeaderCell>Format</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Use Case</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Description</Table.ColumnHeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">BIND</Text>
                </Table.Cell>
                <Table.Cell>BIND, Named</Table.Cell>
                <Table.Cell>Standard zone file format for BIND DNS server</Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">Zone File</Text>
                </Table.Cell>
                <Table.Cell>Universal</Table.Cell>
                <Table.Cell>RFC 1035 compliant generic zone file format</Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">CoreDNS</Text>
                </Table.Cell>
                <Table.Cell>CoreDNS, Kubernetes</Table.Cell>
                <Table.Cell>Corefile format for CoreDNS configuration</Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>
                  <Text weight="bold">PowerDNS</Text>
                </Table.Cell>
                <Table.Cell>PowerDNS</Table.Cell>
                <Table.Cell>Zone format compatible with PowerDNS server</Table.Cell>
              </Table.Row>
            </Table.Body>
          </Table.Root>
        </Flex>
      </Card>
    </Flex>
  );
}
