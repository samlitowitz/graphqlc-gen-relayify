# graphqlc-gen-relayify


This is a code generator designed to work with [graphqlc](https://github.com/samlitowitz/graphqlc).

Generate a GraphQL schema from a GraphQL schema with added Relay Node interfaces and Connection types for the specified types.

```graphql
# input.graphql
type Query {}
type A {}
```
```graphql
# output.graphql
type Query {
    node(id: ID): Node
}

type A implements Node {
    id: ID!
}

interface Node {
    id: ID!
}
```

# Installation
Install [graphqlc](https://github.com/samlitowitz/graphqlc).

`go get -u github.com/samlitowitz/graphqlc-gen-relayify/cmd/graphqlc-gen-relayify`


# Usage
```yaml
# relayify.yaml
# the file can be named anything, you just have to specify it to the config parameter!

# Create Node interface if it does not exist AND nodeify is not empty
# Implement the Node interface for each type specified in nodeify
nodeify:
  - Todo
  - User
```

`graphqlc --relayify_out=config=relayify.yml:. schema.graphql

# In the works!
  * Create *Connection and *Edge types for specified types
