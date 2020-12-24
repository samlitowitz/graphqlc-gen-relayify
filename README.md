# graphqlc-gen-relayify

This is a code generator designed to work with [graphqlc](https://github.com/samlitowitz/graphqlc).

Generate a GraphQL schema from a GraphQL schema with added Relay Node interfaces and Connection types for the specified types.
See the [examples/](examples/) directory for more... examples.

# Installation
Install [graphqlc](https://github.com/samlitowitz/graphqlc).

`go get -u github.com/samlitowitz/graphqlc-gen-relayify/cmd/graphqlc-gen-relayify`

# Usage
Specify the cursor type.
See the Relay specification (https://facebook.github.io/relay/graphql/connections.htm#sec-Cursor)
```json
{
  "cursorType": {
    "type": "String",
    "nullable": false
  }
}
```

Create PageInfo type if it does not exist and connectify is not empty
Create <TYPE>Connection and <TYPE>Edge types if they do not exist
```json
{
  "connectify": [
    {
      "type": "Todo",
      "fields": [
        {
          "type": "User",
          "field": "todos"
        },
        {
          "type": "MyQuery",
          "field": "todos",
          "overwrite": true
        }
      ]
    }
  ]
}
```

Create Node interface if it does not exist AND nodeify is not empty
Implement the Node interface for each type specified in nodeify
```json
{
  "nodeify": [
    "Todo",
    "User"
  ]
}
```

## Parameters
  * config, required, name of relayify configuration file as defined directly above,
  * suffix, optional, default = .echo.graphql, suffix for output file

`graphqlc --relayify_out=config=relayify.json:. schema.graphql`
  