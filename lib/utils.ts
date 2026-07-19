import { type ClassValue, clsx } from "clsx";

/**
 * Merge Tailwind class names, resolving conflicts via clsx/twMerge pattern.
 * Lightweight — no dependency on `tailwind-merge` for now.
 * If class conflicts become an issue, swap `clsx` for `tailwind-merge`.
 */
export function cn(...inputs: ClassValue[]): string {
  return clsx(inputs);
}

/**
 * Format a date for display.
 */
export function formatDate(date: Date): string {
  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(date);
}
