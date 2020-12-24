package relayify

import (
	"fmt"
	echo "github.com/samlitowitz/graphqlc-gen-echo/pkg/graphqlc/echo"
	"github.com/samlitowitz/graphqlc/pkg/graphqlc"
)

type Generator struct {
	*echo.Generator

	config          *Config
	typeSuffix      string // Append type suffix for rename
	genFileNames    map[string]bool
	connectifyTypes map[string]bool
	nodeifyTypes    map[string]bool

	cursorType *graphqlc.TypeDescriptorProto
}

func New() *Generator {
	g := new(Generator)
	g.Generator = echo.New()
	g.LogPrefix = "graphqlc-gen-relayify"
	return g
}

func (g *Generator) CommandLineArguments(parameter string) {
	g.Generator.CommandLineArguments(parameter)

	for k, v := range g.Param {
		switch k {
		case "config":
			config, err := LoadConfig(v)
			if err != nil {
				g.Error(err)
			}
			g.config = config
		}
	}
	if g.config == nil {
		g.Fail("a configuration must be provided")
	}
	g.cursorType = buildCursorType(g.config.CursorType.Type, g.config.CursorType.Nullable)
}

func (g *Generator) BuildSchemas() {
	g.genFileNames = make(map[string]bool)
	for _, n := range g.Request.FileToGenerate {
		g.genFileNames[n] = true
	}

	g.connectifyTypes = make(map[string]bool)
	g.nodeifyTypes = make(map[string]bool)

	for _, typ := range g.config.Connectify {
		g.connectifyTypes[typ.Type] = true
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
				edge = buildEdgeObjectType(desc, g.cursorType)
				fd.Objects = append(fd.Objects, edge)
			}
			connection := getObjectType(fd.Objects, desc.Name+"Connection")
			if connection == nil {
				connection = buildConnectionObjectType(desc, edge)
				fd.Objects = append(fd.Objects, connection)
			}
		}
		// Add/overwrite connection fields
		for _, typ := range g.config.Connectify {
			for _, field := range typ.Fields {
				objDesc := getObjectType(fd.Objects, field.Type)
				if objDesc == nil {
					g.Fail(fmt.Sprintf("undefined type %q", field.Type))
				}
				fieldDesc := getFieldDefinitionDescriptorProto(objDesc, field.Field)
				if fieldDesc != nil && !field.Overwrite {
					continue
				}
				if fieldDesc == nil {
					fieldDesc = &graphqlc.FieldDefinitionDescriptorProto{}
					objDesc.Fields = append(objDesc.Fields, fieldDesc)
				}
				buildConnectionField(fieldDesc, g.cursorType, typ.Type+"Connection", field.Field)
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
