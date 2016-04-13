package plugins

type (
	JqGrid struct {
	     StyleUI string `json:"styleUI,omitempty"`
	     Caption string `json:"caption,omitempty"`
             Url string
	     Mtype string `json:"mtype,omitempty"`
	     Datatype string `json:"datatype,omitempty"`
	     ColModel []JqGrid_colModel `json:"colModel"`
	     Page string `json:"page,omitempty"`
	     Width int `json:"width,omitempty"`
	     Height int `json:"height,omitempty"`
	     rowNum int `json:"rowNum,omitempty"`//每页显示行数
	     RowList []int  `json:"rowList,omitempty"`//显示行数列表
	     Sortname string `json:"sortname,omitempty"`
             Sortorder string `json:"sortorder,omitempty"`
	     ScrollPopUp bool `json:"scrollPopUp,omitempty"`
	     ScrollLeftOffset string `json:"scrollLeftOffset,omitempty"`
	     viewrecords bool `json:"viewrecords,omitempty"`
	     Rownumbers bool `json:"rownumbers,omitempty"` // 显示行号
	     RownumWidth int `json:"rownumWidth,omitempty"`// the width of the row numbers columns
	     Scroll int `json:"scroll,omitempty"`
	     Emptyrecords string  `json:"emptyrecords,omitempty"` //无数据时候显示文本
	     PagerDivID string `json:"pager,omitempty"`
        }

	JqGrid_colModel struct{
		Label string `json:"label,omitempty"`
		Name string `json:"name"`
		key bool `json:"key,omitempty"`
		Width int `json:"width,omitempty"`
		Sorttype string `json:"sorttype,omitempty"`
		Index string `json:"index,omitempty"`
		Formatter string `json:"formatter,omitempty"`
}

        JqGridData struct {
		Records int `json:"records"`
		Page int  `json:"page"`
		Total int `json:"total"`
		Rows interface{} `json:"rows"`
	}
)

func NewJqGrid() *JqGrid {
	return &JqGrid{
		StyleUI:"Bootstrap",
		Datatype:"json",
		Emptyrecords:"Scroll to bottom to retrieve new page",
	}
}

//
//url ：jqGrid控件通过这个参数得到需要显示的数据，具体的返回值可以使XML也可以是Json。
//datatype ：这个参数用于设定将要得到的数据类型。类型包括：json 、xml、xmlstring、local、javascript、function。
//mtype : 定义使用哪种方法发起请求，GET或者POST。
//height ：Grid的高度，可以接受数字、%值、auto，默认值为150。
//width ：Grid的宽度，如果未设置，则宽度应为所有列宽的之和；如果设置了宽度，则每列的宽度将会根据shrinkToFit选项的设置，进行设置。
//shrinkToFit ：此选项用于根据width计算每列宽度的算法。默认值为true。如果shrinkToFit为true且设置了width值，则每列宽度会根据width成比例缩放；如果shrinkToFit为false且设置了width值，则每列的宽度不会成比例缩放，而是保持原有设置，而Grid将会有水平滚动条。
//autowidth ：默认值为false。如果设为true，则Grid的宽度会根据父容器的宽度自动重算。重算仅发生在Grid初始化的阶段；如果当父容器尺寸变化了，同时也需要变化Grid的尺寸的话，则需要在自己的代码中调用setGridWidth方法来完成。
//pager ：定义页码控制条Page Bar，在上面的例子中是用一个div(<div id=”pager”></div>)来放置的。
//sortname ：指定默认的排序列，可以是列名也可以是数字。此参数会在被传递到Server端。
//viewrecords ：设置是否在Pager Bar显示所有记录的总数。
//caption ：设置Grid表格的标题，如果未设置，则标题区域不显示。
//rowNum ：用于设置Grid中一次显示的行数，默认值为20。正是这个选项将参数rows（prmNames中设置的）通过url选项设置的链接传递到Server。注意如果Server返回的数据行数超过了rowNum的设定，则Grid也只显示rowNum设定的行数。
//rowList ：一个数组，用于设置Grid可以接受的rowNum值。例如[10,20,30]。
//colNames ：字符串数组，用于指定各列的题头文本，与列的顺序是对应的。
//colModel ：最重要的数组之一，用于设定各列的参数。（稍后详述）
//prmNames ：这是一个数组，用于设置jqGrid将要向Server传递的参数名称。（稍后详述）
//jsonReader ：这又是一个数组，用来设定如何解析从Server端发回来的json数据。（稍后详述）