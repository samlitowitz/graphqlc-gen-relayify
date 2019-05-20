package generator

import (
	"strings"

	"github.com/samlitowitz/graphqlc/pkg/graphqlc"

	"github.com/samlitowitz/graphqlc-gen-relayify/pkg/graphqlc-gen-relayify/generator"
	echo "github.com/samlitowitz/graphqlc/pkg/echo/generator"
)

type Generator struct {
	*echo.Generator

	Param map[string]string // Command-line parameters

	config       *generator.Config
	genFileNames map[string]bool
	nodeifyTypes map[string]bool
}

func New() *Generator {
	g := new(Generator)
	g.Generator = echo.New()
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
		}
	}
	if g.config == nil {
		g.Fail("a configuration must be provided")
	}
}

func (g *Generator) BuildSchemas() {
	g.genFileNames = make(map[string]bool)
	for _, n := range g.Request.FileToGenerate {
		g.genFileNames[n] = true
	}

	g.nodeifyTypes = make(map[string]bool)
	for _, typ := range g.config.Nodeify {
		g.nodeifyTypes[typ] = true
	}

	if len(g.nodeifyTypes) > 0 {
		g.nodeify()
	}

	// Create PageInfo type if does not exist
	// Create *Connection and *Edge for specified types
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
