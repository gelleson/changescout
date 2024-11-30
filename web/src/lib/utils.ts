import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export const PROJECT_NAME = "Change Scout";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
