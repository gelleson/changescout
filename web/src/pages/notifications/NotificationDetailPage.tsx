import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@apollo/client';
import { GET_NOTIFICATION } from '../../lib/graphql/notifications';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Edit, ArrowLeft, Bell } from 'lucide-react';

export function NotificationDetailPage() {
  const { id } = useParams();
  const { data, loading, error } = useQuery(GET_NOTIFICATION, { variables: { id } });

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading notification details</div>;

  const notification = data?.getNotification;

  return (
    <div className="container mx-auto py-6">
      <div className="mb-6">
        <Link
          to="/notifications"
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back to Notifications
        </Link>
        
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold text-gray-900">{notification.name}</h1>
          <Link
            to={`/notifications/${id}/edit`}
            className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
          >
            <Edit className="h-4 w-4 mr-2" />
            Edit Notification
          </Link>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center text-lg">
            <Bell className="h-5 w-5 mr-2 text-gray-500" />
            Notification Details
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <div>
            <label className="text-sm font-medium text-gray-500">Type</label>
            <p className="text-gray-900">{notification.type}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Destination</label>
            <p className="text-gray-900">{notification.destination}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Created At</label>
            <p className="text-gray-900">
              {new Date(notification.createdAt).toLocaleString()}
            </p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Updated At</label>
            <p className="text-gray-900">
              {new Date(notification.updatedAt).toLocaleString()}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
