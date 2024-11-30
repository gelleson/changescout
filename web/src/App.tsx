import { ApolloProvider } from '@apollo/client';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/Layout';
import { SignIn } from './pages/auth/SignIn';
import { SignUp } from './pages/auth/SignUp';
import { Logout } from './pages/auth/Logout';
import { WebsitesPage } from './pages/websites/WebsitesPage';
import { NotificationsPage } from './pages/notifications/NotificationsPage';
import { AuthGuard } from './components/AuthGuard';
import { client } from './lib/apollo';
import { EditWebsitePage } from './pages/websites/EditWebsitePage';
import { EditNotificationPage } from './pages/notifications/EditNotificationPage';
import { NewWebsitePage } from './pages/websites/NewWebsitePage';
import { WebsiteDetailPage } from './pages/websites/WebsiteDetailPage';
import { NotificationDetailPage } from './pages/notifications/NotificationDetailPage';
import { AddNotificationPage } from './pages/notifications/AddNotificationPage';
import { SnackbarProvider } from 'notistack';

function App() {

  return (
    <ApolloProvider client={client}>
      <SnackbarProvider maxSnack={3}>
        <BrowserRouter>
          <Routes>
            <Route path="/signin" element={<SignIn />} />
            <Route path="/signup" element={<SignUp />} />
            <Route path="/logout" element={<Logout />} />
            <Route
              path="/"
              element={
                <AuthGuard>
                  <Layout />
                </AuthGuard>
              }
            >
              <Route index element={<Navigate to="/websites" replace />} />
              <Route path="websites" element={<WebsitesPage />} />
              <Route path="websites/new" element={<NewWebsitePage />} />
              <Route path="websites/:id" element={<WebsiteDetailPage />} />
              <Route path="websites/:id/edit" element={<EditWebsitePage />} />
              <Route path="notifications" element={<NotificationsPage />} />
              <Route path="notifications/new" element={<AddNotificationPage />} />
              <Route path="notifications/:id" element={<NotificationDetailPage />} />
              <Route path="notifications/:id/edit" element={<EditNotificationPage />} />
              <Route path="*" element={<Navigate to="/websites" replace />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </SnackbarProvider>
    </ApolloProvider>
  );
}

export default App;