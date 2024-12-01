import { useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery, useMutation } from '@apollo/client';
import { GET_WEBSITE_BY_ID, DELETE_WEBSITE, CREATE_PREVIEW_WEBSITE } from '../../lib/graphql/websites';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Edit, ArrowLeft, Globe, Clock, Settings, Trash2, Menu } from 'lucide-react';
import { formatDate } from '../../lib/utils';
import { Sheet} from 'react-modal-sheet';

export function WebsiteDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { data, loading, error } = useQuery(GET_WEBSITE_BY_ID, { variables: { id } });
  const [deleteWebsite] = useMutation(DELETE_WEBSITE);
  const [createPreviewWebsite] = useMutation(CREATE_PREVIEW_WEBSITE);
  const [isPreviewOpen, setPreviewOpen] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this website?')) {
      try {
        await deleteWebsite({
          variables: { id },
        });
        navigate('/websites');
      } catch (error) {
        console.error('Error deleting website:', error);
      }
    }
  };

  const handleCreatePreview = async () => {
    try {
      const { data } = await createPreviewWebsite({ variables: { url: id } });
      setPreview(data.createPreviewWebsite.result);
      setPreviewOpen(true);
    } catch (error) {
      console.error('Error creating preview:', error);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading website details</div>;

  const website = data?.getWebsiteByID;

  if (!website) return <div>No website data available.</div>;

  return (
    <div className="container mx-auto py-6">
      <div className="mb-6">
        <Link
          to="/websites"
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back to Websites
        </Link>
        
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold text-gray-900">{website.name}</h1>
          <div className="flex items-center space-x-2">
            <Link
              to={`/websites/${id}/edit`}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
            >
              <Edit className="h-4 w-4 mr-2" />
              Edit Website
            </Link>
            <button
              onClick={handleDelete}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete Website
            </button>
            <button
              onClick={() => handleCreatePreview()}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
            >
              Preview
            </button>
          </div>
        </div>
      </div>

      <Sheet isOpen={isPreviewOpen} onClose={() => setPreviewOpen(false)}>
        <Sheet.Container>
          <Sheet.Header>
            <div className="flex justify-between items-center p-4 border-b">
              <span className="text-xl font-semibold">Website Preview</span>
              <button className="p-2 rounded-full text-gray-500 hover:text-gray-700" onClick={() => setPreviewOpen(false)}>
                <Menu className="h-6 w-6" />
              </button>
            </div>
          </Sheet.Header>
          <Sheet.Content>
            <div className="p-4">
              <p>{preview}</p>
            </div>
          </Sheet.Content>
        </Sheet.Container>
        <Sheet.Backdrop />
      </Sheet>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Basic Info Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center text-lg">
              <Globe className="h-5 w-5 mr-2 text-gray-500" />
              Basic Information
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div>
              <label className="text-sm font-medium text-gray-500">URL</label>
              <p className="text-gray-900 break-all">{website.url}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-500">Mode</label>
              <p className="text-gray-900 capitalize">{website.mode}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-500">Status</label>
              <p className={`text-${website.enabled ? 'green' : 'red'}-600`}>
                {website.enabled ? 'Active' : 'Inactive'}
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Schedule Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center text-lg">
              <Clock className="h-5 w-5 mr-2 text-gray-500" />
              Schedule
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div>
              <label className="text-sm font-medium text-gray-500">Cron Expression</label>
              <p className="text-gray-900 font-mono">{website.cron}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-500">Next Check</label>
              <p className="text-gray-900">{formatDate(new Date(website.next_check_at))}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-500">Last Check</label>
              <p className="text-gray-900">{website.last_check_at ? formatDate(new Date(website.last_check_at)) : 'N/A'}</p>
            </div>
          </CardContent>
        </Card>

        {/* Advanced Settings Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center text-lg">
              <Settings className="h-5 w-5 mr-2 text-gray-500" />
              Advanced Settings
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div>
              <label className="text-sm font-medium text-gray-500">HTTP Method</label>
              <p className="text-gray-900">{website.setting.method}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-500">User Agent</label>
              <p className="text-gray-900 break-all text-sm">
                {website.setting.user_agent || 'Default'}
              </p>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-500">Features</label>
              <div className="flex flex-wrap gap-2 mt-1">
                {website.setting.deduplication && (
                  <span className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded">
                    Deduplication
                  </span>
                )}
                {website.setting.trim && (
                  <span className="px-2 py-1 text-xs bg-green-100 text-green-800 rounded">
                    Trim
                  </span>
                )}
                {website.setting.sort && (
                  <span className="px-2 py-1 text-xs bg-purple-100 text-purple-800 rounded">
                    Sort
                  </span>
                )}
              </div>
            </div>
            {website.setting.selectors.length > 0 && (
              <div>
                <label className="text-sm font-medium text-gray-500">CSS Selectors</label>
                <div className="mt-1 space-y-1">
                  {website.setting.selectors.map((selector: string, index: number) => (
                    <p key={index} className="text-sm font-mono text-gray-900">
                      {selector}
                    </p>
                  ))}
                </div>
              </div>
            )}
            {website.setting.json_path && website.setting.json_path.length > 0 && (
              <div>
                <label className="text-sm font-medium text-gray-500">JSON Paths</label>
                <div className="mt-1 space-y-1">
                  {website.setting.json_path.map((path: string, index: number) => (
                    <p key={index} className="text-sm font-mono text-gray-900">
                      {path}
                    </p>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
