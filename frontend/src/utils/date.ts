import { format, formatDistanceToNow, parseISO } from 'date-fns';

export const formatDate = (dateString: string): string => {
  try {
    const date = parseISO(dateString);
    return format(date, 'MMM d, yyyy');
  } catch {
    return dateString;
  }
};

export const formatDateTime = (dateString: string): string => {
  try {
    const date = parseISO(dateString);
    return format(date, 'MMM d, yyyy HH:mm');
  } catch {
    return dateString;
  }
};

export const formatRelativeTime = (dateString: string): string => {
  try {
    const date = parseISO(dateString);
    return formatDistanceToNow(date, { addSuffix: true });
  } catch {
    return dateString;
  }
};

export const formatDateRange = (from?: string, to?: string): string => {
  if (!from && !to) return 'All time';
  if (!from) return `Until ${formatDate(to!)}`;
  if (!to) return `From ${formatDate(from)}`;
  return `${formatDate(from)} - ${formatDate(to)}`;
};
