import { useState } from 'react';
import { Dialog, Flex, TextField, Button, Text } from '@radix-ui/themes';
import * as api from '../services/api';

interface CreateZoneDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export function CreateZoneDialog({ open, onOpenChange, onSuccess }: CreateZoneDialogProps) {
  const [domain, setDomain] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!domain.trim()) {
      setError('Domain is required');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await api.createZone({
        domain: domain.trim(),
        records: [],
        enabled: true,
      });
      setDomain('');
      onSuccess();
      onOpenChange(false);
    } catch (err) {
      console.error('Failed to create zone:', err);
      setError(err instanceof Error ? err.message : 'Failed to create zone');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Content style={{ maxWidth: 450 }}>
        <Dialog.Title>Create New Zone</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Create a new DNS zone. Records can be added after creation.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="3">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Domain Name
              </Text>
              <TextField.Root
                placeholder="example.com"
                value={domain}
                onChange={e => setDomain(e.target.value)}
                disabled={isSubmitting}
              />
              <Text as="div" size="1" color="gray" mt="1">
                The domain will automatically get a trailing dot if needed
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
                {isSubmitting ? 'Creating...' : 'Create Zone'}
              </Button>
            </Flex>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  );
}
