//go:generate graphqlc --relayify_out=config=relayify.yaml:. schema.graphql
//go:generate graphqlc --relayify_out=config=relayify.yaml,typeSuffix=.test.graphql:. schema.graphql
package todo
