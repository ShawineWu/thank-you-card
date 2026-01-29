import { useState, useEffect } from 'react';

/**
 * Custom hook for debouncing values
 * 
 * Delays updating the returned value until the specified delay has passed
 * without the input value changing. Useful for optimizing expensive operations
 * like API calls or preview updates.
 * 
 * @param value - The value to debounce
 * @param delay - Delay in milliseconds before updating the debounced value
 * @returns The debounced value
 * 
 * @example
 * ```tsx
 * const [searchTerm, setSearchTerm] = useState('');
 * const debouncedSearchTerm = useDebounce(searchTerm, 500);
 * 
 * useEffect(() => {
 *   // This will only run 500ms after the user stops typing
 *   fetchSearchResults(debouncedSearchTerm);
 * }, [debouncedSearchTerm]);
 * ```
 */
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    // Set up a timer to update the debounced value after the delay
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    // Clean up the timer if value changes before delay expires
    // or if the component unmounts
    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}
