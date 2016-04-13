package plugins

type (
	EditableGrid struct {
		MDatas editableDatas
	}

        editableDatas struct {
		Columns []EditableColumn  `json:"metadata,omitempty"`
		Datas []EditableData  `json:"data,omitempty"`
	}

	EditableColumn struct {
		FieldName string  `json:"name"`
		Label string  `json:"label,omitempty"`
		Datatype string `json:"datatype,omitempty"`
		Editable bool `json:"editable,omitempty"`
		Values ComboxEditor `json:"values,omitempty"`
	}

	ComboxEditor map[string]map[string]string

	EditableData struct {
		ID int `json:"id"`
		Value interface{} `json:"values"`
	}
)

func NewEditableGrid()*EditableGrid{
   return &EditableGrid{
	   MDatas:editableDatas{},
   }
}

func (e *editableDatas)AddColumn(cols... EditableColumn) {
	e.Columns=append(e.Columns,cols...)

}

func (e *editableDatas)AddData(ds interface{}){
	for i,d:=range ds.([]interface{}){
		e.Datas=append(e.Datas,EditableData{
			ID:i,
			Value: d,
		})
	}
}
