package generator

import (
	"path"
	"strings"

	"github.com/samlitowitz/graphqlc/pkg/graphqlc"

	"github.com/samlitowitz/graphqlc-gen-relayify/pkg/graphqlc-gen-relayify/generator"
	echo "github.com/samlitowitz/graphqlc/pkg/echo/generator"
)

type Generator struct {
	*echo.Generator

	Param map[string]string // Command-line parameters

	config          *generator.Config
	typeSuffix      string // Append type suffix for rename
	genFileNames    map[string]bool
	connectifyTypes map[string]bool
	nodeifyTypes    map[string]bool
}

func New() *Generator {
	g := new(Generator)
	g.Generator = echo.New()
	g.Generator.FnRenameFile = g.graphqlRelayifyFileName
	return g
}

func (g *Generator) CommandLineArguments(parameter string) {
	g.Param = make(map[string]string)
	for _, p := range strings.Split(parameter, ",") {
		if i := strings.Index(p, "="); i < 0 {
			g.Param[p] = ""
		} else {
			g.Param[p[0:i]] = p[i+1:]
		}
	}

	for k, v := range g.Param {
		switch k {
		case "config":
			config, err := generator.LoadConfig(v)
			if err != nil {
				g.Error(err)
			}
			g.config = config
		case "typeSuffix":
			if v != "" {
				g.typeSuffix = v
			}
		}
	}
	if g.config == nil {
		g.Fail("a configuration must be provided")
	}
	if g.typeSuffix == "" {
		g.typeSuffix = ".relayified.graphql"
	}
}

func (g *Generator) BuildSchemas() {
	g.genFileNames = make(map[string]bool)
	for _, n := range g.Request.FileToGenerate {
		g.genFileNames[n] = true
	}

	g.connectifyTypes = make(map[string]bool)
	g.nodeifyTypes = make(map[string]bool)

	for _, typ := range g.config.Connectify {
		g.connectifyTypes[typ] = true
	}

	for _, typ := range g.config.Nodeify {
		g.nodeifyTypes[typ] = true
	}
	if len(g.connectifyTypes) > 0 {
		g.connectify()
	}

	if len(g.nodeifyTypes) > 0 {
		g.nodeify()
	}
}

func (g *Generator) connectify() {
	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFileNames[fd.Name]; !ok {
			continue
		}
		// Create PageInfo type if it does not exist
		pageInfo := getObjectType(fd.Objects, "PageInfo")
		if pageInfo == nil {
			pageInfo = buildPageInfoObjectType()
			fd.Objects = append(fd.Objects, pageInfo)
		}
		// Create *Connection and *Edge for specified types
		for _, desc := range fd.Objects {
			if _, ok := g.connectifyTypes[desc.Name]; !ok {
				continue
			}
			edge := getObjectType(fd.Objects, desc.Name+"Edge")
			if edge == nil {
				edge = buildEdgeObjectType(desc)
				fd.Objects = append(fd.Objects, edge)
			}
			connection := getObjectType(fd.Objects, desc.Name+"Connection")
			if connection == nil {
				connection = buildConnectionObjectType(desc, edge)
				fd.Objects = append(fd.Objects, connection)
			}
		}
	}
}

func (g *Generator) nodeify() {
	for _, fd := range g.Request.GraphqlFile {
		if _, ok := g.genFileNames[fd.Name]; !ok {
			continue
		}
		// Create Node interface if it does not exist
		node := getInterface(fd.Interfaces, "Node")
		if node == nil {
			node = buildNodeInterface()
			fd.Interfaces = append(fd.Interfaces, node)
		}
		// Add node root field if does not exist
		wrapExistingQueryType(fd)
		addNodeRootField(fd, node)

		// Implement Node interface for specified types
		for _, desc := range fd.Objects {
			if _, ok := g.nodeifyTypes[desc.Name]; !ok {
				continue
			}
			err := implementNode(desc, node)
			if err != nil {
				g.Error(err)
			}
		}
	}
}

func addNodeRootField(fd *graphqlc.FileDescriptorGraphql, node *graphqlc.InterfaceTypeDefinitionDescriptorProto) {
	if fd.Schema == nil {
		fd.Schema = &graphqlc.SchemaDescriptorProto{}
	}
	if fd.Schema.Query == nil {
		query := &graphqlc.ObjectTypeDefinitionDescriptorProto{Name: "RootQueryType"}
		fd.Objects = append(fd.Objects, query)
		fd.Schema.Query = query
	}
	if getFieldDefinitionDescriptorProto(fd.Schema.Query, "node") == nil {
		fd.Schema.Query.Fields = append(fd.Schema.Query.Fields, buildNodeField(node))
	}
}

func buildNodeField(node *graphqlc.InterfaceTypeDefinitionDescriptorProto) *graphqlc.FieldDefinitionDescriptorProto {
	return &graphqlc.FieldDefinitionDescriptorProto{
		Name: "node",
		Arguments: []*graphqlc.InputValueDefinitionDescriptorProto{
			&graphqlc.InputValueDefinitionDescriptorProto{
				Name: "id",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NonNullType{
						NonNullType: &graphqlc.NonNullTypeDescriptorProto{
							Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
								NamedType: &graphqlc.NamedTypeDescriptorProto{
									Name: "ID",
								},
							},
						},
					},
				},
			},
		},
		Type: &graphqlc.TypeDescriptorProto{
			Type: &graphqlc.TypeDescriptorProto_NonNullType{
				NonNullType: &graphqlc.NonNullTypeDescriptorProto{
					Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
						NamedType: &graphqlc.NamedTypeDescriptorProto{
							Name: node.Name,
						},
					},
				},
			},
		},
	}
}

func buildNodeInterface() *graphqlc.InterfaceTypeDefinitionDescriptorProto {
	return &graphqlc.InterfaceTypeDefinitionDescriptorProto{
		Name: "Node",
		Fields: []*graphqlc.FieldDefinitionDescriptorProto{
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "id",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NonNullType{
						NonNullType: &graphqlc.NonNullTypeDescriptorProto{
							Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
								NamedType: &graphqlc.NamedTypeDescriptorProto{
									Name: "ID",
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildEdgeObjectType(desc *graphqlc.ObjectTypeDefinitionDescriptorProto) *graphqlc.ObjectTypeDefinitionDescriptorProto {
	return &graphqlc.ObjectTypeDefinitionDescriptorProto{
		Name: desc.Name + "Edge",
		Fields: []*graphqlc.FieldDefinitionDescriptorProto{
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "node",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NamedType{
						NamedType: &graphqlc.NamedTypeDescriptorProto{
							Name: desc.Name,
						},
					},
				},
			},
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "cursor",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NonNullType{
						NonNullType: &graphqlc.NonNullTypeDescriptorProto{
							Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
								NamedType: &graphqlc.NamedTypeDescriptorProto{
									Name: "String",
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildConnectionObjectType(desc, edge *graphqlc.ObjectTypeDefinitionDescriptorProto) *graphqlc.ObjectTypeDefinitionDescriptorProto {
	return &graphqlc.ObjectTypeDefinitionDescriptorProto{
		Name: desc.Name + "Connection",
		Fields: []*graphqlc.FieldDefinitionDescriptorProto{
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "edge",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_ListType{
						ListType: &graphqlc.ListTypeDescriptorProto{
							Type: &graphqlc.TypeDescriptorProto{
								Type: &graphqlc.TypeDescriptorProto_NamedType{
									NamedType: &graphqlc.NamedTypeDescriptorProto{
										Name: edge.Name,
									},
								},
							},
						},
					},
				},
			},
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "PageInfo",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NonNullType{
						NonNullType: &graphqlc.NonNullTypeDescriptorProto{
							Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
								NamedType: &graphqlc.NamedTypeDescriptorProto{
									Name: "PageInfo",
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildPageInfoObjectType() *graphqlc.ObjectTypeDefinitionDescriptorProto {
	return &graphqlc.ObjectTypeDefinitionDescriptorProto{
		Name: "PageInfo",
		Fields: []*graphqlc.FieldDefinitionDescriptorProto{
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "hasPreviousPage",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NonNullType{
						NonNullType: &graphqlc.NonNullTypeDescriptorProto{
							Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
								NamedType: &graphqlc.NamedTypeDescriptorProto{
									Name: "Boolean",
								},
							},
						},
					},
				},
			},
			&graphqlc.FieldDefinitionDescriptorProto{
				Name: "hasNextPage",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NonNullType{
						NonNullType: &graphqlc.NonNullTypeDescriptorProto{
							Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
								NamedType: &graphqlc.NamedTypeDescriptorProto{
									Name: "Boolean",
								},
							},
						},
					},
				},
			},
		},
	}
}

func implementNode(desc *graphqlc.ObjectTypeDefinitionDescriptorProto, node *graphqlc.InterfaceTypeDefinitionDescriptorProto) error {
	if getInterface(desc.Implements, node.Name) != nil {
		return nil
	}

	desc.Implements = append(desc.Implements, node)
	for _, nodeFieldDesc := range node.Fields {
		fieldDesc := getFieldDefinitionDescriptorProto(desc, nodeFieldDesc.Name)
		if fieldDesc == nil {
			desc.Fields = append(desc.Fields, nodeFieldDesc)
			continue
		}
	}
	return nil
}

func getFieldDefinitionDescriptorProto(desc *graphqlc.ObjectTypeDefinitionDescriptorProto, fieldName string) *graphqlc.FieldDefinitionDescriptorProto {
	for _, fieldDesc := range desc.Fields {
		if fieldDesc.Name == fieldName {
			return fieldDesc
		}
	}
	return nil
}

func getInterface(descs []*graphqlc.InterfaceTypeDefinitionDescriptorProto, name string) *graphqlc.InterfaceTypeDefinitionDescriptorProto {
	for _, desc := range descs {
		if desc.Name == name {
			return desc
		}
	}
	return nil
}

func getObjectType(descs []*graphqlc.ObjectTypeDefinitionDescriptorProto, name string) *graphqlc.ObjectTypeDefinitionDescriptorProto {
	for _, desc := range descs {
		if desc.Name == name {
			return desc
		}
	}
	return nil
}

func wrapExistingQueryType(fd *graphqlc.FileDescriptorGraphql) {
	if fd.Schema == nil {
		return
	}
	if fd.Schema.Query == nil {
		return
	}

	for _, desc := range fd.Objects {
		if desc.Name == fd.Schema.Query.Name {
			fd.Schema.Query = desc
			return
		}
	}
}

func (g *Generator) graphqlRelayifyFileName(name string) string {
	if ext := path.Ext(name); ext == ".graphql" {
		name = name[:len(name)-len(ext)]
	}
	name += g.typeSuffix

	return name
}
