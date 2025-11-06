import type { DNSRecord } from '../services/api';

/**
 * Format a DNS record's value for display based on its type.
 * For records with type-specific fields, returns a formatted string.
 * For simple records, returns the value field.
 */
export function formatRecordValue(record: DNSRecord): string {
  switch (record.type) {
    case 'MX':
      if (record.mx_priority !== undefined && record.mx_host) {
        return `${record.mx_priority} → ${record.mx_host}`;
      }
      return record.value || '';

    case 'SRV':
      if (
        record.srv_priority !== undefined &&
        record.srv_weight !== undefined &&
        record.srv_port !== undefined &&
        record.srv_target
      ) {
        return `Pri: ${record.srv_priority}, Wgt: ${record.srv_weight}, Port: ${record.srv_port} → ${record.srv_target}`;
      }
      return record.value || '';

    case 'SOA':
      if (record.soa_mname && record.soa_rname && record.soa_serial !== undefined) {
        return `${record.soa_mname} ${record.soa_rname} (Serial: ${record.soa_serial})`;
      }
      return record.value || '';

    case 'CAA':
      if (record.caa_flags !== undefined && record.caa_tag && record.caa_value) {
        const flagText = record.caa_flags === 128 ? 'Critical' : 'Non-critical';
        return `[${flagText}] ${record.caa_tag}: ${record.caa_value}`;
      }
      return record.value || '';

    default:
      return record.value || '';
  }
}

/**
 * Get a detailed description of a DNS record for tooltips or detailed views.
 */
export function getRecordDetails(record: DNSRecord): string[] {
  const details: string[] = [];

  switch (record.type) {
    case 'MX':
      if (record.mx_priority !== undefined) {
        details.push(`Priority: ${record.mx_priority}`);
      }
      if (record.mx_host) {
        details.push(`Mail Server: ${record.mx_host}`);
      }
      break;

    case 'SRV':
      if (record.srv_priority !== undefined) {
        details.push(`Priority: ${record.srv_priority}`);
      }
      if (record.srv_weight !== undefined) {
        details.push(`Weight: ${record.srv_weight}`);
      }
      if (record.srv_port !== undefined) {
        details.push(`Port: ${record.srv_port}`);
      }
      if (record.srv_target) {
        details.push(`Target: ${record.srv_target}`);
      }
      break;

    case 'SOA':
      if (record.soa_mname) {
        details.push(`Primary NS: ${record.soa_mname}`);
      }
      if (record.soa_rname) {
        details.push(`Admin: ${record.soa_rname}`);
      }
      if (record.soa_serial !== undefined) {
        details.push(`Serial: ${record.soa_serial}`);
      }
      if (record.soa_refresh !== undefined) {
        details.push(`Refresh: ${record.soa_refresh}s`);
      }
      if (record.soa_retry !== undefined) {
        details.push(`Retry: ${record.soa_retry}s`);
      }
      if (record.soa_expire !== undefined) {
        details.push(`Expire: ${record.soa_expire}s`);
      }
      if (record.soa_minimum !== undefined) {
        details.push(`Minimum: ${record.soa_minimum}s`);
      }
      break;

    case 'CAA':
      if (record.caa_flags !== undefined) {
        const flagText = record.caa_flags === 128 ? 'Critical' : 'Non-critical';
        details.push(`Flags: ${record.caa_flags} (${flagText})`);
      }
      if (record.caa_tag) {
        details.push(`Tag: ${record.caa_tag}`);
      }
      if (record.caa_value) {
        details.push(`Value: ${record.caa_value}`);
      }
      break;

    default:
      if (record.value) {
        details.push(record.value);
      }
      break;
  }

  return details;
}
