export const CRON_PRESETS = {
  '30min': '*/30 * * * *',
  'hourly': '0 * * * *',
  'daily': '0 0 * * *',
  'custom': 'custom'
} as const;
