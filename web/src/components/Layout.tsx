import React, { useState } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { Monitor, Bell, Settings, LogOut, Menu } from 'lucide-react';
import { PROJECT_NAME } from "../lib/utils";
import { Sheet, SheetContainer, SheetHeader, SheetContent, SheetBackdrop } from 'react-modal-sheet';

export function Layout() {
  const location = useLocation();
  const [isDrawerOpen, setDrawerOpen] = useState(false);

  const toggleDrawer = () => {
    setDrawerOpen(!isDrawerOpen);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <button className="sm:hidden p-2 rounded-full text-gray-500 hover:text-gray-700" onClick={toggleDrawer}>
                <Menu className="h-6 w-6" />
              </button>
              <Link to="/" className="flex items-center px-2 py-2 text-gray-900">
                <Monitor className="h-6 w-6 text-indigo-600" />
                <span className="ml-2 text-xl font-semibold">{PROJECT_NAME}</span>
              </Link>
              <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                <Link
                  to="/websites"
                  className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium ${
                    location.pathname.startsWith('/websites')
                      ? 'border-indigo-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                  }`}
                >
                  Websites
                </Link>
                <Link
                  to="/notifications"
                  className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium ${
                    location.pathname.startsWith('/notifications')
                      ? 'border-indigo-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                  }`}
                >
                  Notifications
                </Link>
              </div>
            </div>
            <div className="flex items-center">
              <button className="p-2 rounded-full text-gray-500 hover:text-gray-700">
                <Bell className="h-5 w-5" />
              </button>
              <button className="ml-3 p-2 rounded-full text-gray-500 hover:text-gray-700">
                <Settings className="h-5 w-5" />
              </button>
              <Link to="/logout" className="ml-3 p-2 rounded-full text-gray-500 hover:text-gray-700">
                <LogOut className="h-5 w-5" />
              </Link>
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <Outlet />
      </main>
      <Sheet isOpen={isDrawerOpen} onClose={() => setDrawerOpen(false)}>
        <Sheet.Container>
          <Sheet.Header>
            <div className="flex justify-between items-center p-4 border-b">
              <span className="text-xl font-semibold">Navigation</span>
              <button className="p-2 rounded-full text-gray-500 hover:text-gray-700" onClick={() => setDrawerOpen(false)}>
                <Menu className="h-6 w-6" />
              </button>
            </div>
          </Sheet.Header>
          <Sheet.Content>
            <div className="flex flex-col space-y-4">
              <Link to="/websites" className="block p-4 bg-white shadow rounded-lg hover:shadow-md transition duration-200" onClick={() => setDrawerOpen(false)}>
                <div className="flex items-center">
                  <Monitor className="h-6 w-6 text-indigo-600" />
                  <span className="ml-3 text-black text-lg font-medium">Websites</span>
                </div>
              </Link>
              <Link to="/notifications" className="block p-4 bg-white shadow rounded-lg hover:shadow-md transition duration-200" onClick={() => setDrawerOpen(false)}>
                <div className="flex items-center">
                  <Bell className="h-6 w-6 text-indigo-600" />
                  <span className="ml-3 text-black text-lg font-medium">Notifications</span>
                </div>
              </Link>
              <Link to="/logout" className="block p-4 bg-white shadow rounded-lg hover:shadow-md transition duration-200" onClick={() => setDrawerOpen(false)}>
                <div className="flex items-center">
                  <LogOut className="h-6 w-6 text-indigo-600" />
                  <span className="ml-3 text-black text-lg font-medium">Logout</span>
                </div>
              </Link>
            </div>
          </Sheet.Content>
        </Sheet.Container>
        <Sheet.Backdrop />
      </Sheet>
    </div>
  );
}