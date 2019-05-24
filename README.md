# graphqlc-gen-relayify


This is a code generator designed to work with [graphqlc](https://github.com/samlitowitz/graphqlc).

Generate a GraphQL schema from a GraphQL schema with added Relay Node interfaces and Connection types for the specified types.
See the [examples/](https://github.com/samlitowitz/graphqlc-gen-relayify/tree/master/examples) directory for more... examples.

```graphql
# input.graphql
type Query {}
type AType {}
```
```graphql
# output.graphql
type Query {
    node(id: ID): Node
}

type AType implements Node {
    id: ID!
}

type PageInfo {
  hasPreviousPage: Boolean!
  hasNextPage: Boolean!
}

type ATypeEdge {
  node: AType
  cursor: String!
}

type ATypeConnection {
  edge: [ATypeEdge]
  PageInfo: PageInfo!
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

# Create PageInfo type if it does not exist and connectify is not empty
# Create <TYPE>Connection and <TYPE>Edge types if they do not exist
connectify:
  - Todo

# Create Node interface if it does not exist AND nodeify is not empty
# Implement the Node interface for each type specified in nodeify
nodeify:
  - Todo
  - User
```

## Parameters
  * config, required, name of relayify configuration file as defined directly above,
  * typeSuffix, optional, default = .relayified.graphql, suffix for output file

`graphqlc --relayify_out=config=relayify.yml:. schema.graphql`
