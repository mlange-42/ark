{{- define "query" -}}
package ecs

type cursor struct {
	table    int
	index    uintptr
	maxIndex int64
}

{{range makeRange 1 8}}
{{- $lower := lowerLetters . -}}
{{- $upper := upperLetters . -}}
{{- $generics := join "" " any, " " any" $upper -}}
{{- $genericsShort := join "" ", " "" $upper -}}
{{- $return := join "*" ", *" "" $upper -}}
{{- $mask := join "id" ", id" "" $upper -}}

// Query{{.}} is a filter for two components.
type Query{{.}}[{{$generics}}] struct {
	world      *World
	mask       Mask
	cursor     cursor
	{{- range $upper}}
	component{{.}} *componentStorage
	{{- end}}
	{{- range $upper}}
	column{{.}}    *column
	{{- end}}
}

// NewQuery{{.}} creates a new [Query{{.}}].
func NewQuery{{.}}[{{$generics}}](world *World) Query{{.}}[{{$genericsShort}}] {
	{{- range $upper}}
	id{{.}} := ComponentID[{{.}}](world)
	{{- end}}

	return Query{{.}}[{{$genericsShort}}]{
		world:      world,
		mask:       All({{$mask}}),
		{{- range $upper}}
		component{{.}}: &world.storage.components[id{{.}}.id],
		{{- end}}
		cursor: cursor{
			table:    -1,
			index:    0,
			maxIndex: -1,
		},
	}
}

func (q *Query{{.}}[{{$genericsShort}}]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

func (q *Query{{.}}[{{$genericsShort}}]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		table := &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[table.archetype]
		if !archetype.mask.Contains(&q.mask) || table.entities.Len() == 0 {
			continue
		}
		{{- range $upper}}
		q.column{{.}} = q.component{{.}}.columns[q.cursor.table]
		{{- end}}

		q.cursor.index = 0
		q.cursor.maxIndex = int64(table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	return false
}

func (q *Query{{.}}[{{$genericsShort}}]) Get() ({{$return}}) {
	return {{range $i, $v := $upper}}{{if $i}},
	    {{end}}(*{{$v}})(q.column{{$v}}.Get(q.cursor.index)){{end}}
}
{{end -}}
{{end}}