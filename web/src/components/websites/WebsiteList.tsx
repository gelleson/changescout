import React from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { Link } from 'react-router-dom';
import { Power, Edit, Plus, Trash2 } from 'lucide-react';
import { GET_WEBSITES, CHANGE_WEBSITE_STATUS, DELETE_WEBSITE } from '../../lib/graphql/websites';
import type { Website } from '../../types';
import { EmptyState } from '../EmptyState';

export function WebsiteList() {
  const { data, loading, error } = useQuery(GET_WEBSITES);
  const [changeStatus] = useMutation(CHANGE_WEBSITE_STATUS);
  const [deleteWebsite] = useMutation(DELETE_WEBSITE);

  const handleStatusChange = async (id: string, enabled: boolean) => {
    try {
      await changeStatus({
        variables: { id, enabled: !enabled },
        refetchQueries: [{ query: GET_WEBSITES }]
      });
    } catch (error) {
      console.error('Error changing website status:', error);
    }
  };

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this website?')) {
      try {
        await deleteWebsite({
          variables: { id },
          refetchQueries: [{ query: GET_WEBSITES }]
        });
      } catch (error) {
        console.error('Error deleting website:', error);
      }
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading websites</div>;
  if (!data?.getWebsites.length) return <EmptyState message="No websites found. Add a new website to get started!" buttonText="Add Website" buttonLink="/websites/new" />;

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold text-gray-900">Websites</h1>
        <Link
          to="/websites/new"
          className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
        >
          <Plus className="h-4 w-4 mr-2" />
          Add Website
        </Link>
      </div>
      
      <div className="space-y-6">
        {data?.getWebsites.map((website: Website) => (
          <div key={website.id} className="bg-white shadow overflow-hidden sm:rounded-md">
            <div className="px-4 py-4 flex items-center justify-between sm:px-6">
              <div className="flex-1 min-w-0">
                <Link to={`/websites/${website.id}`} className="hover:text-indigo-600">
                  <h3 className="text-lg font-medium text-gray-900 truncate">
                    {website.name}
                  </h3>
                  <p className="mt-1 text-sm text-gray-500">{website.url}</p>
                  <p className="mt-1 text-xs text-gray-400">
                    Last check: {website.last_check_at ? new Date(website.last_check_at).toLocaleString() : 'N/A'}
                  </p>
                  <p className="mt-1 text-xs text-gray-400">
                    Next check: {new Date(website.next_check_at).toLocaleString()}
                  </p>
                </Link>
              </div>
              <div className="flex items-center space-x-4">
                <button
                  onClick={() => handleStatusChange(website.id, website.enabled)}
                  className={`p-2 rounded-full ${
                    website.enabled ? 'text-green-600' : 'text-gray-400'
                  }`}
                >
                  <Power className="h-5 w-5" />
                </button>
                <Link
                  to={`/websites/${website.id}/edit`}
                  className="p-2 rounded-full text-gray-400 hover:text-gray-500"
                >
                  <Edit className="h-5 w-5" />
                </Link>
                <button
                  onClick={() => handleDelete(website.id)}
                  className="p-2 rounded-full text-red-400 hover:text-red-500"
                >
                  <Trash2 className="h-5 w-5" />
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}