package exports

// Worth exported from worth.xlsx
type Worth struct {
	Id      int32     `json:"Id"`
	Name    string    `json:"Name"`
	Unitid  string    `json:"Unitid"`
	Weight  float32   `json:"Weight"`
	Type    int32     `json:"Type"`
	EHP     float32   `json:"EHP"`
	DPS     float32   `json:"DPS"`
	EHPDPS  float32   `json:"EHPDPS"`
	SPDFUNC float32   `json:"SPDFUNC"`
	ACCTAR  int32     `json:"ACCTAR"`
	RESTAR  int32     `json:"RESTAR"`
	HPFUNC  float32   `json:"HPFUNC"`
	ATKFUNC float32   `json:"ATKFUNC"`
	DEFFUNC float32   `json:"DEFFUNC"`
	CRIFUNC float32   `json:"CRIFUNC"`
	HP      float32   `json:"HP"`
	ATK     float32   `json:"ATK"`
	DEF     float32   `json:"DEF"`
	SPD     float32   `json:"SPD"`
	CRI     float32   `json:"CRI"`
	CRIDMG  float32   `json:"CRIDMG"`
	RES     float32   `json:"RES"`
	ACC     float32   `json:"ACC"`
	SET     []float32 `json:"SET"`
}
