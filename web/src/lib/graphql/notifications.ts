// GraphQL queries and mutations for notifications
import { gql } from 'graphql-tag';

export const GET_NOTIFICATION = gql`query GetNotification($id: ID!) {
  getNotification(id: $id) {
    id
    name
    type
    websiteId
    createdAt
    updatedAt
  }
}`;

export const GET_NOTIFICATIONS = gql`query GetNotifications($input: NotificationFilters) {
  getNotifications(input: $input) {
    id
    name
    type
    websiteId
    createdAt
    updatedAt
  }
}`;

export const CREATE_NOTIFICATION = gql`mutation CreateNotification($input: NotificationCreateInput!) {
  createNotification(input: $input) {
    id
    name
    type
    websiteId
    createdAt
    updatedAt
  }
}`;

export const UPDATE_NOTIFICATION = gql`mutation UpdateNotification($input: NotificationUpdateInput!) {
  updateNotification(input: $input) {
    id
    name
    type
    websiteId
    createdAt
    updatedAt
  }
}`;

export const DELETE_NOTIFICATION = gql`mutation DeleteNotification($id: ID!) {
  deleteNotification(id: $id) {
    id
    name
    type
    websiteId
    createdAt
    updatedAt
  }
}`;
