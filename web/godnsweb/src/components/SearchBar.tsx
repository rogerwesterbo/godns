import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { TextField, Box } from '@radix-ui/themes';
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';
import './SearchBar.css';

export default function SearchBar() {
  const [isExpanded, setIsExpanded] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsExpanded(false);
      }
    };

    if (isExpanded) {
      document.addEventListener('mousedown', handleClickOutside);
      inputRef.current?.focus();
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isExpanded]);

  const handleFocus = () => {
    setIsExpanded(true);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (searchQuery.trim()) {
        navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
        setIsExpanded(false);
        setSearchQuery('');
      }
    } else if (e.key === 'Escape') {
      setIsExpanded(false);
    }
  };

  return (
    <Box ref={containerRef} className={`search-container ${isExpanded ? 'expanded' : ''}`}>
      <TextField.Root
        ref={inputRef}
        placeholder="Search zones, records..."
        value={searchQuery}
        onChange={e => setSearchQuery(e.target.value)}
        onFocus={handleFocus}
        onKeyDown={handleKeyDown}
        size="2"
      >
        <TextField.Slot>
          <MagnifyingGlassIcon height="16" width="16" />
        </TextField.Slot>
      </TextField.Root>
    </Box>
  );
}
