import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation } from '@apollo/client';
import { GET_NOTIFICATION, UPDATE_NOTIFICATION } from '../../lib/graphql/notifications';
import { getIconForType } from '../../lib/getIconForType';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Send } from 'lucide-react';
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { SnackbarProvider, useSnackbar } from 'notistack';

export function EditNotificationPage() {
  const { id } = useParams();
  const { data, loading, error } = useQuery(GET_NOTIFICATION, { variables: { id } });
  const [updateNotification] = useMutation(UPDATE_NOTIFICATION);
  const [formData, setFormData] = useState({
    name: '',
    type: '',
    token: undefined,
    destination: undefined,
    websiteId: '',
  });
  const { enqueueSnackbar } = useSnackbar();

  useEffect(() => {
    if (data) {
      setFormData(data.getNotification);
    }
  }, [data]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const { createdAt, updatedAt, ...rest } = formData;
      await updateNotification({ variables: { input: { id, ...rest } } });
      enqueueSnackbar('Notification updated successfully!', { variant: 'success' });
    } catch (error) {
      console.error('Error updating notification:', error);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading notification details</div>;

  return (
    <SnackbarProvider maxSnack={3}>
      <div className="container py-4">
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Notification Name"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">Type</Label>
              <div className="flex items-center">
                <Select
                  name="type"
                  value={formData.type}
                  onValueChange={(value) =>
                    handleChange({
                      target: { name: 'type', value }
                    } as any)
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select notification type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="telegram">
                      <div className="flex items-center">
                        {getIconForType('telegram')}
                        Telegram
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="mt-4">
          <CardContent className="space-y-4">
            <Accordion type="single" collapsible className="w-full">
              <AccordionItem value="item-1">
                <AccordionTrigger>Advanced Settings</AccordionTrigger>
                <AccordionContent>
                  <div className="space-y-2">
                    <Label htmlFor="token">Token</Label>
                    <Input
                      id="token"
                      name="token"
                      value={formData.token}
                      onChange={handleChange}
                      type="password"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="destination">Destination</Label>
                    <Input
                      id="destination"
                      name="destination"
                      value={formData.destination}
                      onChange={handleChange}
                    />
                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </CardContent>
        </Card>

        <div className="mt-4"></div>

        <Button type="submit" onClick={handleSubmit} className="w-full">
          Save Changes
        </Button>
      </div>
    </SnackbarProvider>
  );
}
