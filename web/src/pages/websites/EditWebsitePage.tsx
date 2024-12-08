import React, {useState, useEffect} from 'react';
import {useParams, useNavigate} from 'react-router-dom';
import {useQuery, useMutation} from '@apollo/client';
import {GET_WEBSITE_BY_ID, UPDATE_WEBSITE} from '../../lib/graphql/websites';
import {Button} from "@/components/ui/button";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Input} from "@/components/ui/input";
import {Label} from "@/components/ui/label";
import {Switch} from "@/components/ui/switch";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {RadioGroup, RadioGroupItem} from "@/components/ui/radio-group";
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion";
import {TemplateMessage} from '@/components/templates/TemplateMessage';
import {CRON_PRESETS} from '@/lib/constants';
import {Textarea} from "@/components/ui/textarea";

export function EditWebsitePage() {
    const navigate = useNavigate();
    const {id} = useParams();
    const {data, loading, error} = useQuery(GET_WEBSITE_BY_ID, {variables: {id}});
    const [updateWebsite] = useMutation(UPDATE_WEBSITE);
    const [cronType, setCronType] = useState<keyof typeof CRON_PRESETS>('custom');
    const [formData, setFormData] = useState({
        url: '',
        name: '',
        enabled: false,
        mode: 'plain',
        cron: '',
        setting: {
            user_agent: '',
            referer: '',
            method: 'GET',
            deduplication: false,
            trim: false,
            sort: false,
            selectors: [],
            json_path: [] as string[],
            template: ''
        },
    });

    useEffect(() => {
        if (data) {
            setFormData(data.getWebsiteByID);
            // Determine the cron type based on the loaded data
            const cronValue = data.getWebsiteByID.cron;
            const matchingPreset = Object.entries(CRON_PRESETS).find(
                ([key, value]) => value === cronValue && key !== 'custom'
            );
            setCronType(matchingPreset ? matchingPreset[0] as keyof typeof CRON_PRESETS : 'custom');
        }
    }, [data]);

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
        const {name, url, mode, cron, enabled, setting} = formData;
        try {
            await updateWebsite({
                variables: {
                    input: {
                        id,
                        name,
                        url,
                        mode,
                        cron,
                        enabled,
                        setting: {
                            ...setting,
                            json_path: Array.isArray(setting.json_path) ? setting.json_path.join(',') : ''
                        }
                    }
                }
            });
            navigate('/websites');
        } catch (error) {
            console.error('Error updating website:', error);
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading website details</div>;

    return (
        <form onSubmit={handleSubmit} className="container mx-auto py-6">
            <Card>
                <CardHeader>
                    <CardTitle>Edit Website</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="name">Name</Label>
                        <Input
                            id="name"
                            value={formData.name}
                            onChange={(e) => setFormData({...formData, name: e.target.value})}
                            placeholder="My Website"
                            required
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="url">URL</Label>
                        <Input
                            id="url"
                            value={formData.url}
                            onChange={(e) => setFormData({...formData, url: e.target.value})}
                            placeholder="https://example.com"
                            type="url"
                            required
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="mode">Mode</Label>
                        <Select
                            value={formData.mode}
                            onValueChange={(value) => setFormData({...formData, mode: value})}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Select mode"/>
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="plain">HTML</SelectItem>
                                <SelectItem value="renderer">Google Chrome Renderer</SelectItem>
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
                                <RadioGroupItem value="30min" id="30min"/>
                                <Label htmlFor="30min" className="font-normal">Every 30 minutes</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="hourly" id="hourly"/>
                                <Label htmlFor="hourly" className="font-normal">Every hour</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="daily" id="daily"/>
                                <Label htmlFor="daily" className="font-normal">Once a day</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="custom" id="custom"/>
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
                            onCheckedChange={(checked) => setFormData({...formData, enabled: checked})}
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
                                            setting: {...formData.setting, user_agent: e.target.value}
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
                                            setting: {...formData.setting, referer: e.target.value}
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
                                            setting: {...formData.setting, method: value}
                                        })}
                                    >
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select method"/>
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
                                                setting: {...formData.setting, deduplication: checked}
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
                                                setting: {...formData.setting, trim: checked}
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
                                                setting: {...formData.setting, sort: checked}
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
                                        value={Array.isArray(formData.setting.json_path) ? formData.setting.json_path.join(', ') : ''}
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
                                        <TemplateMessage/>
                                    </div>
                                    <Textarea
                                        id="template"
                                        value={formData.setting.template || ''}
                                        onChange={(e) => setFormData({
                                            ...formData,
                                            setting: {...formData.setting, template: e.target.value}
                                        })}
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
                            Save Changes
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </form>
    );
}
