import React, { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useQuery, gql } from '@apollo/client';

const TEST_AUTH = gql`
  query GetWebsites {
    getWebsites {
      id
    }
  }
`;

interface AuthGuardProps {
  children: React.ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { loading, error } = useQuery(TEST_AUTH, {
    onError: (error) => {
      if (error.message.includes('not authorized')) {
        navigate('/signin', { state: { from: location.pathname } });
      }
    }
  });

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  return <>{children}</>;
}