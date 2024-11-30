export interface Website {
  id: string;
  url: string;
  name: string;
  enabled: boolean;
  mode: 'plain';
  next_check_at: string;
  cron: string;
  setting: {
    user_agent?: string;
    referer?: string;
    method: string;
    deduplication: boolean;
    trim: boolean;
    sort: boolean;
    selectors: string[];
    json_path: string[];
  };
}

export interface AuthInput {
  email: string;
  password: string;
}

export interface SignUpInput extends AuthInput {
  firstName: string;
  lastName: string;
}