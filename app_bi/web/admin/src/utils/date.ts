import { format } from 'date-fns';

  export function formatTimeToStr(times: Date, pattern?: string): string {
    const defaultPattern = 'yyyy-MM-dd HH:mm:ss';
    const fmt = pattern || defaultPattern;
    return format(times, fmt);
  }