import React from 'react';
import { useNavigate } from 'react-router-dom';

export function Logout() {
  const navigate = useNavigate();

  const handleLogout = () => {
    // Remove the auth token from local storage
    localStorage.removeItem('authToken');
    // Redirect to the sign-in page
    navigate('/signin');
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded shadow-md text-center">
        <h1 className="text-2xl font-bold mb-4">Log Out</h1>
        <p className="mb-4">Click the button below to log out and return to the sign-in page.</p>
        <button onClick={handleLogout} className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600">
          Sign In
        </button>
      </div>
    </div>
  );
}
