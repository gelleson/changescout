import { gql } from 'graphql-tag';

// GraphQL queries and mutations for websites
export const GET_WEBSITE_BY_ID = gql`query GetWebsiteByID($id: ID!) {
  getWebsiteByID(id: $id) {
    id
    url
    name
    enabled
    mode
    next_check_at
    last_check_at
    cron
    setting {
      user_agent
      referer
      method
      deduplication
      trim
      sort
      selectors
      json_path
    }
  }
}`;

export const GET_WEBSITE_BY_URL = gql`query GetWebsiteByURL($url: String!) {
  getWebsiteByURL(url: $url) {
    id
    url
    name
    enabled
    mode
    next_check_at
    last_check_at
    cron
    setting {
      user_agent
      referer
      method
      deduplication
      trim
      sort
      selectors
      json_path
    }
  }
}`;

export const GET_WEBSITES = gql`query GetWebsites {
  getWebsites {
    id
    url
    name
    enabled
    mode
    next_check_at
    last_check_at
    cron
    setting {
      user_agent
      referer
      method
      deduplication
      trim
      sort
      selectors
      json_path
    }
  }
}`;

export const CREATE_WEBSITE = gql`mutation CreateWebsite($input: WebsiteCreateInput!) {
  createWebsite(input: $input) {
    id
    url
    name
    enabled
    mode
    next_check_at
    last_check_at
    cron
    setting {
      user_agent
      referer
      method
      deduplication
      trim
      sort
      selectors
      json_path
    }
  }
}`;

export const UPDATE_WEBSITE = gql`mutation UpdateWebsite($input: WebsiteUpdateInput!) {
  updateWebsite(input: $input) {
    id
    url
    name
    enabled
    mode
    next_check_at
    last_check_at
    cron
    setting {
      user_agent
      referer
      method
      deduplication
      trim
      sort
      selectors
      json_path
    }
  }
}`;

export const CHANGE_WEBSITE_STATUS = gql`mutation ChangeWebsiteStatus($id: ID!, $enabled: Boolean!) {
  changeWebsiteStatus(id: $id, enabled: $enabled) {
    id
    url
    name
    enabled
    mode
    next_check_at
    last_check_at
    cron
    setting {
      user_agent
      referer
      method
      deduplication
      trim
      sort
      selectors
      json_path
    }
  }
}`;

export const DELETE_WEBSITE = gql`mutation DeleteWebsite($id: ID!) {
  deleteWebsite(id: $id)
}`;
