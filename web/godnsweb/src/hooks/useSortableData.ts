import { useState, useMemo } from 'react';

export type SortDirection = 'asc' | 'desc' | null;

export interface SortConfig<T> {
  key: keyof T | string | null;
  direction: SortDirection;
}

export function useSortableData<T>(
  items: T[],
  initialSortKey?: keyof T | string,
  initialSortDirection: SortDirection = 'asc',
  customComparators?: Record<string, (a: T, b: T) => number>
) {
  const [sortConfig, setSortConfig] = useState<SortConfig<T>>({
    key: initialSortKey || null,
    direction: initialSortDirection,
  });

  const sortedItems = useMemo(() => {
    const sortableItems = [...items];
    
    if (sortConfig.key !== null && sortConfig.direction !== null) {
      sortableItems.sort((a, b) => {
        // Check if there's a custom comparator
        if (customComparators && sortConfig.key && customComparators[sortConfig.key as string]) {
          const comparison = customComparators[sortConfig.key as string](a, b);
          return sortConfig.direction === 'asc' ? comparison : -comparison;
        }

        const aValue = a[sortConfig.key as keyof T];
        const bValue = b[sortConfig.key as keyof T];

        if (aValue === bValue) return 0;
        if (aValue === null || aValue === undefined) return 1;
        if (bValue === null || bValue === undefined) return -1;

        let comparison = 0;
        
        // Handle different types
        if (typeof aValue === 'string' && typeof bValue === 'string') {
          comparison = aValue.toLowerCase().localeCompare(bValue.toLowerCase());
        } else if (typeof aValue === 'number' && typeof bValue === 'number') {
          comparison = aValue - bValue;
        } else {
          comparison = String(aValue).localeCompare(String(bValue));
        }

        return sortConfig.direction === 'asc' ? comparison : -comparison;
      });
    }

    return sortableItems;
  }, [items, sortConfig, customComparators]);

  const requestSort = (key: keyof T | string) => {
    let direction: SortDirection = 'asc';
    
    if (sortConfig.key === key) {
      if (sortConfig.direction === 'asc') {
        direction = 'desc';
      } else if (sortConfig.direction === 'desc') {
        direction = null;
      }
    }

    setSortConfig({ key: direction === null ? null : key, direction });
  };

  return { items: sortedItems, requestSort, sortConfig };
}
