import React from "react";
import {
  ApolloClient,
  InMemoryCache,
  gql,
  ApolloProvider,
  useQuery,
} from "@apollo/client";

import "./App.css";

// Define GraphQL query
const USER_QUERY = gql`
  query GetUser($id: String!) {
    user(id: $id) {
      id
      name
    }
  }
`;

// Server URI
const GRAPHQL_SERVER_URI = "http://localhost:8080/graphql";

// Initialize Apollo Client
const client = new ApolloClient({
  uri: GRAPHQL_SERVER_URI,
  cache: new InMemoryCache(),
});

// Define the User type, matching the schema on the server
interface User {
  id: string;
  name: string;
}

// Define the shape of the query response data
interface QueryResponse {
  user: User;
}

interface UserProps {
  id: string;
}

const UserComponent: React.FC<UserProps> = ({ id }) => {
  const { loading, error, data } = useQuery<QueryResponse>(USER_QUERY, {
    variables: { id },
  });

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error :(</p>;

  return data?.user ? (
    <div>
      <p>User ID: {data.user.id}</p>
      <p>User Name: {data.user.name}</p>
    </div>
  ) : (
    <p>No user found</p>
  );
};

const App: React.FC = () => {
  return (
    <ApolloProvider client={client}>
      <div>
        <h2>GraphQL sample 🚀</h2>
        <UserComponent id="1" />
      </div>
    </ApolloProvider>
  );
};

export default App;
