package relayify

import "github.com/samlitowitz/graphqlc/pkg/graphqlc"

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
			{
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
			{
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

func buildCursorType(typ string, nullable bool) *graphqlc.TypeDescriptorProto {
	switch nullable {
	case false:
		return &graphqlc.TypeDescriptorProto{
			Type: &graphqlc.TypeDescriptorProto_NonNullType{
				NonNullType: &graphqlc.NonNullTypeDescriptorProto{
					Type: &graphqlc.NonNullTypeDescriptorProto_NamedType{
						NamedType: &graphqlc.NamedTypeDescriptorProto{
							Name: typ,
						},
					},
				},
			},
		}
	default:
		return &graphqlc.TypeDescriptorProto{
			Type: &graphqlc.TypeDescriptorProto_NamedType{
				NamedType: &graphqlc.NamedTypeDescriptorProto{
					Name: typ,
				},
			},
		}
	}
}

func buildConnectionField(desc *graphqlc.FieldDefinitionDescriptorProto,
	cursorType *graphqlc.TypeDescriptorProto,
	connectionName, fieldName string) {

	desc.Name = fieldName
	desc.Type = &graphqlc.TypeDescriptorProto{
		Type: &graphqlc.TypeDescriptorProto_NamedType{
			NamedType: &graphqlc.NamedTypeDescriptorProto{
				Name: connectionName,
			},
		},
	}
	desc.Arguments = []*graphqlc.InputValueDefinitionDescriptorProto{}

	desc.Arguments = append(desc.Arguments, &graphqlc.InputValueDefinitionDescriptorProto{
		Name: "first",
		Type: &graphqlc.TypeDescriptorProto{
			Type: &graphqlc.TypeDescriptorProto_NamedType{
				NamedType: &graphqlc.NamedTypeDescriptorProto{
					Name: "Int",
				},
			},
		},
	})
	desc.Arguments = append(desc.Arguments, &graphqlc.InputValueDefinitionDescriptorProto{
		Name: "after",
		Type: cursorType,
	})
	desc.Arguments = append(desc.Arguments, &graphqlc.InputValueDefinitionDescriptorProto{
		Name: "last",
		Type: &graphqlc.TypeDescriptorProto{
			Type: &graphqlc.TypeDescriptorProto_NamedType{
				NamedType: &graphqlc.NamedTypeDescriptorProto{
					Name: "Int",
				},
			},
		},
	})
	desc.Arguments = append(desc.Arguments, &graphqlc.InputValueDefinitionDescriptorProto{
		Name: "before",
		Type: cursorType,
	})
}

func buildEdgeObjectType(desc *graphqlc.ObjectTypeDefinitionDescriptorProto, cursorType *graphqlc.TypeDescriptorProto) *graphqlc.ObjectTypeDefinitionDescriptorProto {
	return &graphqlc.ObjectTypeDefinitionDescriptorProto{
		Name: desc.Name + "Edge",
		Fields: []*graphqlc.FieldDefinitionDescriptorProto{
			{
				Name: "node",
				Type: &graphqlc.TypeDescriptorProto{
					Type: &graphqlc.TypeDescriptorProto_NamedType{
						NamedType: &graphqlc.NamedTypeDescriptorProto{
							Name: desc.Name,
						},
					},
				},
			},
			{
				Name: "cursor",
				Type: cursorType,
			},
		},
	}
}

func buildConnectionObjectType(desc, edge *graphqlc.ObjectTypeDefinitionDescriptorProto) *graphqlc.ObjectTypeDefinitionDescriptorProto {
	return &graphqlc.ObjectTypeDefinitionDescriptorProto{
		Name: desc.Name + "Connection",
		Fields: []*graphqlc.FieldDefinitionDescriptorProto{
			{
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
			{
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
			{
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
			{
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
