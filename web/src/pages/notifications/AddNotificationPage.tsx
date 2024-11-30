import React, { useState } from 'react';
import { useMutation } from '@apollo/client';
import { CREATE_NOTIFICATION } from '../../lib/graphql/notifications';
import { getIconForType } from '../../lib/getIconForType';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useSnackbar } from 'notistack';

export function AddNotificationPage() {
  const [formData, setFormData] = useState({
    name: '',
    type: '',
    token: '',
    destination: '',
  });
  const [createNotification] = useMutation(CREATE_NOTIFICATION);
  const { enqueueSnackbar } = useSnackbar();

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
      await createNotification({ variables: { input: formData } });
      enqueueSnackbar('Notification created successfully!', { variant: 'success' });
    } catch (error) {
      console.error('Error creating notification:', error);
      enqueueSnackbar('Failed to create notification.', { variant: 'error' });
    }
  };

  return (
    <div className="container py-4">
      <Card>
        <CardHeader>
          <CardTitle>Add New Notification</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Notification Name"
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="type">Type</Label>
            <Select
              name="type"
              value={formData.type}
              onValueChange={(value) =>
                handleChange({ target: { name: 'type', value } } as any)
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
          <Button type="submit" onClick={handleSubmit} className="w-full">
            Create Notification
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
