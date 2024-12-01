import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export const PROJECT_NAME = "Change Scout";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: Date) {
  const now = new Date();
  const options = {
    month: 'long',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  };
  if (date.getFullYear() !== now.getFullYear()) {
    return date.toLocaleString(undefined, {
      year: 'numeric',
      ...options
    });
  }
  return date.toLocaleString(undefined, options);
}
