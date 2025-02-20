package ecs

type storage struct {
	archetypes []archetype
	tables     []table
}

func newStorage(capacity uint32, reg *registry) storage {
	tables := make([]table, 0, 128)
	tables = append(tables, newTable(capacity, reg))
	archetypes := make([]archetype, 0, 128)
	archetypes = append(archetypes,
		archetype{
			tables: []*table{&tables[0]},
		},
	)
	return storage{
		archetypes: archetypes,
		tables:     tables,
	}
}
