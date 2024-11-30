import React from 'react';
import { Send } from 'lucide-react';

export function getIconForType(type: string): React.ReactNode | null {
  switch (type) {
    case 'telegram':
      return <Send className="h-4 w-4 mr-2 text-blue-500" />;
    default:
      return null;
  }
}
