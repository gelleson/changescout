import React from 'react';

interface EmptyStateProps {
  message: string;
  buttonText?: string;
  buttonLink?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({ message, buttonText, buttonLink }) => {
  return (
    <div className="flex flex-col items-center justify-center h-full text-center space-y-4">
      <svg
        className="w-16 h-16 mb-4 text-gray-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M3 15a4 4 0 004 4h10a4 4 0 004-4M7 10l5-5m0 0l5 5m-5-5v12"
        />
      </svg>
      <p className="text-lg font-medium text-gray-600 mb-4">{message}</p>
      {buttonText && buttonLink && (
        <a
          href={buttonLink}
          className="mt-4 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
        >
          {buttonText}
        </a>
      )}
    </div>
  );
};
