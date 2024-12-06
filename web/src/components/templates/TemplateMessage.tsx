import React from 'react';

interface TemplateMessageProps {
  name?: string;
  mode?: string;
  url?: string;
  lastChecked?: string;
  diff?: string;
}

const defaultProps: Required<TemplateMessageProps> = {
  name: '{{.Name}}',
  mode: '{{.Mode}}',
  url: '{{.URL}}',
  lastChecked: '{{.LastChecked}}',
  diff: '{{.Diff}}'
};

const TemplateMessage: React.FC<TemplateMessageProps> = (props) => {
  const { name, mode, url, lastChecked, diff } = { ...defaultProps, ...props };

  return (
    <div className="text-sm text-muted-foreground space-y-4">
      <div>
        <p>Here example of template:</p>
        <ul className="list-none pl-4 space-y-1">
          <li><strong>Website Name</strong>: {name}</li>
          <li><strong>Mode</strong>: {mode}</li>
        </ul>
      </div>
      <div>
        <p><strong>Website Details:</strong></p>
        <ul className="list-none pl-4 space-y-1">
          <li><strong>URL</strong>: {url}</li>
          <li><strong>Last Checked</strong>: {lastChecked}</li>
        </ul>
      </div>
      <div>
        <p><strong>Changes Detected:</strong></p>
        <pre className="bg-gray-100 p-2 rounded-md mt-2">
          <code className="text-sm">{diff}</code>
        </pre>
      </div>
    </div>
  );
};

TemplateMessage.defaultProps = defaultProps;

export { TemplateMessage };
export type { TemplateMessageProps };
