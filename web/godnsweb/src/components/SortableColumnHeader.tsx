import { Table, Flex } from '@radix-ui/themes';
import { ChevronUpIcon, ChevronDownIcon } from '@radix-ui/react-icons';
import type { SortDirection } from '../hooks';

interface SortableColumnHeaderProps<T> {
  column: keyof T | string;
  currentSortKey: keyof T | string | null;
  currentSortDirection: SortDirection;
  onSort: (column: keyof T | string) => void;
  children: React.ReactNode;
}

export function SortableColumnHeader<T>({
  column,
  currentSortKey,
  currentSortDirection,
  onSort,
  children,
}: SortableColumnHeaderProps<T>) {
  const isActive = currentSortKey === column;
  
  return (
    <Table.ColumnHeaderCell
      style={{ cursor: 'pointer', userSelect: 'none' }}
      onClick={() => onSort(column)}
    >
      <Flex align="center" gap="1">
        {children}
        {isActive && currentSortDirection === 'asc' && <ChevronUpIcon />}
        {isActive && currentSortDirection === 'desc' && <ChevronDownIcon />}
      </Flex>
    </Table.ColumnHeaderCell>
  );
}
