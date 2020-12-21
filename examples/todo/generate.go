//go:generate graphqlc --relayify_out=config=relayify.json:. schema.graphql
//go:generate graphqlc --relayify_out=config=relayify.json,typeSuffix=.test.graphql:. schema.graphql
package todo
