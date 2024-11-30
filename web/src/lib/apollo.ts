import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';

// Initialize Apollo Client
const httpLink = createHttpLink({
  uri: import.meta.env.VITE_GRAPHQL_API_URL || 'http://localhost:3311/query',
});

const authLink = setContext((_, { headers }) => {
  console.log(import.meta.env.VITE_GRAPHQL_API_URL);
  // Get the authentication token from local storage if it exists
  const token = localStorage.getItem('authToken');
  // Return the headers to the context so httpLink can read them
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
    }
  }
});

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    for (let err of graphQLErrors) {
      if (err.message === 'not authenticated') {
        // Redirect to sign-in page
        window.location.href = '/signin';
      }
    }
  }
  if (networkError) console.log(`[Network error]: ${networkError}`);
});

export const client = new ApolloClient({
  link: errorLink.concat(authLink).concat(httpLink),
  cache: new InMemoryCache(),
});
