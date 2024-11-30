import React from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { Link } from 'react-router-dom';
import { Trash2, Edit, Plus, Bell } from 'lucide-react';
import { GET_NOTIFICATIONS, DELETE_NOTIFICATION } from '../../lib/graphql/notifications';
import { EmptyState } from '../EmptyState';

export function NotificationList() {
  const { data, loading, error } = useQuery(GET_NOTIFICATIONS, {
    variables: { input: {} }
  });
  const [deleteNotification] = useMutation(DELETE_NOTIFICATION);

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this notification?')) {
      try {
        await deleteNotification({
          variables: { id },
          refetchQueries: [{ query: GET_NOTIFICATIONS, variables: { input: {} } }]
        });
      } catch (error) {
        console.error('Error deleting notification:', error);
      }
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading notifications</div>;
  if (!data?.getNotifications.length) return <EmptyState message="No notifications found. Add a new notification to get started!" buttonText="Add Notification" buttonLink="/notifications/new" />;

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold text-gray-900">Notifications</h1>
        <Link
          to="/notifications/new"
          className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
        >
          <Plus className="h-4 w-4 mr-2" />
          Add Notification
        </Link>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {data?.getNotifications.map((notification: any) => (
            <li key={notification.id}>
              <div className="px-4 py-4 flex items-center justify-between sm:px-6">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center">
                    <Bell className="h-5 w-5 text-indigo-600 mr-2" />
                    <Link to={`/notifications/${notification.id}`} className="hover:text-indigo-600">
                      <h3 className="text-lg font-medium text-gray-900 truncate">
                        {notification.name}
                      </h3>
                      <span className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded-full mt-1 inline-block">
                        {notification.type}
                      </span>
                    </Link>
                  </div>

                </div>
                <div className="flex items-center space-x-4">
                  <Link
                    to={`/notifications/${notification.id}/edit`}
                    className="p-2 rounded-full text-gray-400 hover:text-gray-500"
                  >
                    <Edit className="h-5 w-5" />
                  </Link>
                  <button
                    onClick={() => handleDelete(notification.id)}
                    className="p-2 rounded-full text-red-400 hover:text-red-500"
                  >
                    <Trash2 className="h-5 w-5" />
                  </button>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}