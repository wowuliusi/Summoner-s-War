package exports

// Worth exported from worth.xlsx
type Worth struct {
	Id        int32     `json:"Id"`
	MonsterHP int32     `json:"MonsterHP"`
	Weight    float32   `json:"Weight"`
	HP        float32   `json:"HP"`
	ATK       float32   `json:"ATK"`
	DEF       float32   `json:"DEF"`
	SPD       float32   `json:"SPD"`
	CRI       float32   `json:"CRI"`
	CRIDMG    float32   `json:"CRIDMG"`
	RES       float32   `json:"RES"`
	ACC       float32   `json:"ACC"`
	SET       []float32 `json:"SET"`
}
