import React, { useState } from 'react';
import { useMutation } from '@apollo/client';
import { useNavigate } from 'react-router-dom';
import { CREATE_WEBSITE } from '../../lib/graphql/websites';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Switch } from "@/components/ui/switch";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { TemplateMessage } from '@/components/templates/TemplateMessage';
import { Textarea } from "@/components/ui/textarea";
import { CRON_PRESETS } from '@/lib/constants';

export function NewWebsitePage() {
  const navigate = useNavigate();
  const [createWebsite] = useMutation(CREATE_WEBSITE);
  const [cronType, setCronType] = useState<keyof typeof CRON_PRESETS>('30min');
  const [formData, setFormData] = useState({
    name: '',
    url: '',
    mode: 'plain',
    cron: CRON_PRESETS['30min'],
    enabled: true,
    setting: {
      user_agent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36',
      referer: '',
      method: 'GET',
      deduplication: true,
      trim: true,
      sort: true,
      selectors: [],
      json_path: [] as string[],
      template: undefined
    }
  });

  const handleCronChange = (value: keyof typeof CRON_PRESETS) => {
    setCronType(value);
    if (value !== 'custom') {
      setFormData(prev => ({
        ...prev,
        cron: CRON_PRESETS[value]
      }));
    }
  };

  const handleCustomCronChange = (value: string) => {
    setFormData(prev => ({
      ...prev,
      cron: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const { data } = await createWebsite({
        variables: {
          input: formData
        }
      });
      
      if (data?.createWebsite) {
        navigate(`/websites/${data.createWebsite.id}`);
      }
    } catch (error) {
      console.error('Error creating website:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="container mx-auto py-6">
      <Card>
        <CardHeader>
          <CardTitle>New Website</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="My Website"
              required
            />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="url">URL</Label>
            <Input
              id="url"
              value={formData.url}
              onChange={(e) => setFormData({ ...formData, url: e.target.value })}
              placeholder="https://example.com"
              type="url"
              required
            />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="mode">Mode</Label>
            <Select
              value={formData.mode}
              onValueChange={(value) => setFormData({ ...formData, mode: value })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select mode" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="plain">HTML</SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <div className="space-y-2">
            <Label>Check Frequency</Label>
            <RadioGroup
              value={cronType}
              onValueChange={handleCronChange}
              className="flex flex-row flex-wrap gap-4"
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="30min" id="30min" />
                <Label htmlFor="30min" className="font-normal">Every 30 minutes</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="hourly" id="hourly" />
                <Label htmlFor="hourly" className="font-normal">Every hour</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="daily" id="daily" />
                <Label htmlFor="daily" className="font-normal">Once a day</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="custom" id="custom" />
                <Label htmlFor="custom" className="font-normal">Custom schedule</Label>
              </div>
            </RadioGroup>
            
            {cronType === 'custom' && (
              <div className="mt-4 space-y-2">
                <Label htmlFor="custom-cron">Custom Cron Expression</Label>
                <Input
                  id="custom-cron"
                  value={formData.cron}
                  onChange={(e) => handleCustomCronChange(e.target.value)}
                  placeholder="*/5 * * * *"
                  required
                />
                <p className="text-sm text-muted-foreground">
                  Enter custom cron expression (e.g., */5 * * * *)
                </p>
              </div>
            )}
          </div>

          <div className="flex items-center space-x-2">
            <Switch
              id="enabled"
              checked={formData.enabled}
              onCheckedChange={(checked) => setFormData({ ...formData, enabled: checked })}
            />
            <Label htmlFor="enabled">Enabled</Label>
          </div>

          <Accordion type="single" collapsible className="w-full">
            <AccordionItem value="settings">
              <AccordionTrigger>Advanced Settings</AccordionTrigger>
              <AccordionContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="user-agent">User Agent</Label>
                  <Input
                    id="user-agent"
                    value={formData.setting.user_agent}
                    onChange={(e) => setFormData({
                      ...formData,
                      setting: { ...formData.setting, user_agent: e.target.value }
                    })}
                    placeholder="User Agent String"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="referer">Referer</Label>
                  <Input
                    id="referer"
                    value={formData.setting.referer}
                    onChange={(e) => setFormData({
                      ...formData,
                      setting: { ...formData.setting, referer: e.target.value }
                    })}
                    placeholder="Referer URL"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="method">HTTP Method</Label>
                  <Select
                    value={formData.setting.method}
                    onValueChange={(value) => setFormData({
                      ...formData,
                      setting: { ...formData.setting, method: value }
                    })}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select method" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="GET">GET</SelectItem>
                      <SelectItem value="POST">POST</SelectItem>
                      <SelectItem value="PUT">PUT</SelectItem>
                      <SelectItem value="DELETE">DELETE</SelectItem>
                      <SelectItem value="PATCH">PATCH</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-4">
                  <div className="flex items-center space-x-2">
                    <Switch
                      id="deduplication"
                      checked={formData.setting.deduplication}
                      onCheckedChange={(checked) => setFormData({
                        ...formData,
                        setting: { ...formData.setting, deduplication: checked }
                      })}
                    />
                    <Label htmlFor="deduplication">Enable Deduplication</Label>
                  </div>

                  <div className="flex items-center space-x-2">
                    <Switch
                      id="trim"
                      checked={formData.setting.trim}
                      onCheckedChange={(checked) => setFormData({
                        ...formData,
                        setting: { ...formData.setting, trim: checked }
                      })}
                    />
                    <Label htmlFor="trim">Trim Values</Label>
                  </div>

                  <div className="flex items-center space-x-2">
                    <Switch
                      id="sort"
                      checked={formData.setting.sort}
                      onCheckedChange={(checked) => setFormData({
                        ...formData,
                        setting: { ...formData.setting, sort: checked }
                      })}
                    />
                    <Label htmlFor="sort">Sort Results</Label>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="selectors">CSS Selectors</Label>
                  <Input
                    id="selectors"
                    value={formData.setting.selectors.join(', ')}
                    onChange={(e) => setFormData({
                      ...formData,
                      setting: {
                        ...formData.setting,
                        selectors: e.target.value.split(',').map(s => s.trim()).filter(Boolean)
                      }
                    })}
                    placeholder="Enter CSS selectors separated by commas"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="json-path">JSON Path</Label>
                  <Input
                    id="json-path"
                    value={formData.setting.json_path.join(', ')}
                    onChange={(e) => setFormData({
                      ...formData,
                      setting: {
                        ...formData.setting,
                        json_path: e.target.value.split(',').map(s => s.trim()).filter(Boolean)
                      }
                    })}
                    placeholder="Enter JSON paths separated by commas"
                  />
                </div>

                <div className="space-y-4">
                  <Label htmlFor="template">Template Message</Label>
                  <div className="bg-gray-50 rounded-lg p-4 mb-4">
                    <TemplateMessage />
                  </div>
                  <Textarea
                    id="template"
                    value={formData.setting.template || ''}
                    onChange={(e) => setFormData({ ...formData, setting: { ...formData.setting, template: e.target.value } })}
                    placeholder="Enter template text"
                    className="min-h-[200px]"
                  />
                </div>
              </AccordionContent>
            </AccordionItem>
          </Accordion>

          <div className="flex justify-end space-x-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate('/websites')}
            >
              Cancel
            </Button>
            <Button type="submit">
              Create Website
            </Button>
          </div>
        </CardContent>
      </Card>
    </form>
  );
}
