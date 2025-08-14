import React from 'react';
import { Check, X } from 'lucide-react';

interface MessageAlertProps {
  message: { type: 'success' | 'error'; text: string } | null;
}

export const MessageAlert: React.FC<MessageAlertProps> = ({ message }) => {
  if (!message) return null;

  return (
    <div className={`p-4 rounded-lg border ${
      message.type === 'success'
        ? 'bg-green-400/10 border-green-400/20 text-green-400'
        : 'bg-red-400/10 border-red-400/20 text-red-400'
    }`}>
      <div className="flex items-center space-x-2">
        {message.type === 'success' ? (
          <Check className="w-4 h-4" />
        ) : (
          <X className="w-4 h-4" />
        )}
        <span className="text-sm">{message.text}</span>
      </div>
    </div>
  );
};