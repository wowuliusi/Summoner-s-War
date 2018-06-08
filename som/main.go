package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sky/exports"
)

type Rune struct {
	Occupied_type int
	Extra         int
	Sell_value    int
	Pri_eff       []int
	Prefix_eff    []int
	Slot_no       int
	Rank          int
	Occupied_id   int
	Sec_eff       [][]int
	Wizard_id     int
	Upgrade_curr  int
	Rune_id       int
	Base_value    int
	Class         int
	Set_id        int
	Upgrade_limit int
}

type Gameinfo struct {
	Runes     []Rune
	Unit_list []Monster
}

type Monster struct {
	Resist          int
	Spd             int
	Critical_rate   int
	Runes           []Rune
	Unit_id         int64
	Accuracy        int
	Atk             int
	Critical_damage int
	Def             int
	Con             int
}
type MyMonster struct {
	ID        int64
	mon       Monster
	OldCRunes CompleteRunes
	NewCRunes CompleteRunes
	OldAttr   MonsterAttribute
	NewAttr   MonsterAttribute
}
type MonsterAttribute struct {
	HP     int
	ATK    int
	DEF    int
	CRI    int
	CRIDMG int
	SPD    int
	ACC    int
	RES    int
	DPS    float32
	EHP    float32
	EHPDPS float32
	SUM    float32
	Runes  CompleteRunes
}
type rune struct {
	Set_id       int
	Upgrade_curr int
	Class        int
	Slot_no      int
	attr         [13]int
	sum          float32
	p            *Rune
}
type CompleteRunes struct {
	Runes [7]*rune
}

var Setnum = []int{9, 2, 2, 4, 2, 4, 2, 2, 4, 9, 4, 4, 9, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var maxrune6 = []int{0, 2448, 63, 160, 63, 160, 63, 0, 42, 58, 80, 64, 64}
var maxrune5 = []int{0, 2050, 51, 135, 51, 135, 51, 0, 39, 47, 59, 50, 50}
var attrb = []float32{0, 0, 6.5, 0, 6.5, 0, 6.5, 0, 5, 5, 5.5, 6, 6}
var savenum int = 10
var Attrname = []string{"x", "HP", "HP%", "ATK", "ATK%", "DEF", "DEF%", "x", "SPD", "CRI", "CRIDMG", "RES", "ACC"}

var toprunes [7][10](*rune)
var bestsetrune [7][24](*rune)
var toprunesnum [7]int
var myrunes []rune
var monsterworth exports.Worth

func main() {
	var game Gameinfo
	var MyMonsters []MyMonster
	dat, _ := ioutil.ReadFile("15060791.json")
	//xlFile, _ := xlsx.OpenFile("11.xlsx")
	// for _, sheet := range xlFile.Sheets {
	// 	fmt.Printf("Sheet Name: %s\n", sheet.Name)
	// 	for _, row := range sheet.Rows {
	// 		for _, cell := range row.Cells {
	// 			text := cell.String()
	// 			fmt.Printf("%s\n", text)

	// 		}
	// 	}
	// }
	if err := json.Unmarshal(dat, &game); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(len(game.Runes))
	for i, r := range game.Runes {
		var tmp rune
		tmp.Set_id = r.Set_id
		tmp.Upgrade_curr = r.Upgrade_curr
		tmp.Slot_no = r.Slot_no
		tmp.Class = r.Class
		tmp.p = &(game.Runes[i])
		if tmp.Class == 6 {
			tmp.attr[r.Pri_eff[0]] += maxrune6[r.Pri_eff[0]]
		} else {
			tmp.attr[r.Pri_eff[0]] += maxrune5[r.Pri_eff[0]]
		}
		tmp.attr[r.Prefix_eff[0]] += r.Prefix_eff[1]
		for _, j := range r.Sec_eff {
			tmp.attr[j[0]] += j[1] + j[3]
		}
		myrunes = append(myrunes, tmp)
	}
	var monsterid = "2350714359"
	var monsterid2 int64 = 2350714359
	var monster Monster
	var wor []exports.Worth
	wordat, _ := ioutil.ReadFile("../exports/worth.json")
	if err := json.Unmarshal(wordat, &wor); err != nil {
		log.Fatal(err)
	}
	for _, monsterworth = range wor {
		if monsterworth.Unitid == monsterid {
			break
		}
	}
	var k int
	for k, monster = range game.Unit_list {
		if monster.Unit_id == monsterid2 {
			var montmp MyMonster
			for i, r := range monster.Runes {
				var tmp rune
				tmp.Set_id = r.Set_id
				tmp.Upgrade_curr = r.Upgrade_curr
				tmp.Slot_no = r.Slot_no
				tmp.Class = r.Class
				tmp.p = &(game.Unit_list[k].Runes[i])
				if tmp.Class == 6 {
					tmp.attr[r.Pri_eff[0]] += maxrune6[r.Pri_eff[0]]
				} else {
					tmp.attr[r.Pri_eff[0]] += maxrune5[r.Pri_eff[0]]
				}
				tmp.attr[r.Prefix_eff[0]] += r.Prefix_eff[1]
				for _, j := range r.Sec_eff {
					tmp.attr[j[0]] += j[1] + j[3]
				}
				myrunes = append(myrunes, tmp)
				montmp.OldCRunes.Runes[tmp.Slot_no] = &(myrunes[len(myrunes)-1])
			}
			montmp.mon = monster
			MyMonsters = append(MyMonsters, montmp)
			break
		}
	}
	//fmt.Println(len(myrunes))

	for i, r := range myrunes {
		sum := (float32(r.attr[1])/float32(monster.Con)/15 + float32(r.attr[2])) * monsterworth.HP / attrb[2]
		sum += (float32(r.attr[3])/float32(monster.Atk) + float32(r.attr[4])) * monsterworth.ATK / attrb[4]
		sum += (float32(r.attr[5])/float32(monster.Def) + float32(r.attr[6])) * monsterworth.DEF / attrb[6]
		sum += float32(r.attr[8])*monsterworth.SPD/attrb[8] + float32(r.attr[9])*monsterworth.CRI/attrb[9] + float32(r.attr[10])*monsterworth.CRIDMG/attrb[10] + float32(r.attr[11])*monsterworth.RES/attrb[11] + float32(r.attr[12])*monsterworth.ACC/attrb[12]
		myrunes[i].sum = sum
		var j int
		for j = 0; j < toprunesnum[r.Slot_no]; j++ {
			if toprunes[r.Slot_no][j] == nil || (*toprunes[r.Slot_no][j]).sum < sum {
				break
			}
		}
		if j < savenum {
			for k := j + 1; k < savenum; k++ {
				if k > toprunesnum[r.Slot_no] {
					break
				}
				toprunes[r.Slot_no][k] = toprunes[r.Slot_no][k-1]
			}
			toprunes[r.Slot_no][j] = &(myrunes[i])
			if toprunesnum[r.Slot_no] < savenum {
				toprunesnum[r.Slot_no]++
			}
		}
	}

	for i, r := range myrunes {
		if bestsetrune[r.Slot_no][r.Set_id] == nil || r.sum > (*bestsetrune[r.Slot_no][r.Set_id]).sum {
			bestsetrune[r.Slot_no][r.Set_id] = &myrunes[i]
		}
	}
	for i := 1; i < 7; i++ {
		for j := 1; j < 24; j++ {
			if bestsetrune[i][j] != nil && (*bestsetrune[i][j]).sum < (*toprunes[i][savenum-1]).sum*0.9 {
				bestsetrune[i][j] = nil
			}
		}
	}
	nowsum := caculate(MyMonsters[0].OldCRunes, &MyMonsters[0]) //caculate(&myrunes[len(myrunes)-6], &myrunes[len(myrunes)-5], &myrunes[len(myrunes)-4], &myrunes[len(myrunes)-3], &myrunes[len(myrunes)-2], &myrunes[len(myrunes)-1], monster)
	fmt.Println("now:", nowsum.SUM, "  DPS:", nowsum.DPS, "  EHP:", nowsum.EHP, "  EHPDPS:", nowsum.EHPDPS, "  SPD:", nowsum.SPD)
	for i := 1; i < 7; i++ {
		fmt.Print(i, ": set:", (*MyMonsters[0].OldCRunes.Runes[i].p).Set_id, "  main:", Attrname[(*MyMonsters[0].OldCRunes.Runes[i].p).Pri_eff[0]], "+", (*MyMonsters[0].OldCRunes.Runes[i].p).Pri_eff[1])
		if (*MyMonsters[0].OldCRunes.Runes[i].p).Prefix_eff[0] > 0 {
			fmt.Print("  Pri:", Attrname[(*MyMonsters[0].OldCRunes.Runes[i].p).Prefix_eff[0]], "+", (*MyMonsters[0].OldCRunes.Runes[i].p).Prefix_eff[1])
		}
		fmt.Println("  sum:", MyMonsters[0].OldCRunes.Runes[i].sum,
			(*MyMonsters[0].OldCRunes.Runes[i].p).Sec_eff)
	}
	var Cr CompleteRunes
	search(Cr, 1, &MyMonsters[0])
	maxsum := caculate(MyMonsters[0].NewCRunes, &MyMonsters[0])
	fmt.Println("find:", maxsum.SUM, "  DPS:", maxsum.DPS, "  EHP:", maxsum.EHP, "  EHPDPS:", maxsum.EHPDPS, "  SPD:", maxsum.SPD)
	for i := 1; i < 7; i++ {
		fmt.Print(i, ": set:", (*MyMonsters[0].NewCRunes.Runes[i].p).Set_id, "  main:", Attrname[(*MyMonsters[0].NewCRunes.Runes[i].p).Pri_eff[0]], "+", (*MyMonsters[0].NewCRunes.Runes[i].p).Pri_eff[1])
		if (*MyMonsters[0].NewCRunes.Runes[i].p).Prefix_eff[0] > 0 {
			fmt.Print("  Pri:", Attrname[(*MyMonsters[0].NewCRunes.Runes[i].p).Prefix_eff[0]], "+", (*MyMonsters[0].NewCRunes.Runes[i].p).Prefix_eff[1])
		}
		fmt.Println("  sum:", MyMonsters[0].NewCRunes.Runes[i].sum,
			(*MyMonsters[0].NewCRunes.Runes[i].p).Sec_eff)
	}
}
func search(Cr CompleteRunes, slot int, mymon *MyMonster) {
	if slot == 7 {
		Attr := caculate(Cr, mymon)
		if Attr.SUM > mymon.NewAttr.SUM {
			mymon.NewAttr = Attr
			mymon.NewCRunes = Cr
		}
		return
	}
	for i := 0; i < savenum; i++ {
		Cr.Runes[slot] = toprunes[slot][i]
		search(Cr, slot+1, mymon)
	}
	for i := 1; i < 24; i++ {
		if bestsetrune[slot][i] == nil {
			continue
		}
		Cr.Runes[slot] = bestsetrune[slot][i]
		search(Cr, slot+1, mymon)
	}
}
func caculate(Cr CompleteRunes, monster *MyMonster) MonsterAttribute {
	var sets [24]int
	var Attr MonsterAttribute
	Attr.Runes = Cr
	for i := 1; i < 7; i++ {
		sets[(*Cr.Runes[i]).Set_id]++
	}
	var attrsum [13]int
	for i := 1; i < 13; i++ {
		for j := 1; j < 7; j++ {
			attrsum[i] += (*Cr.Runes[j]).attr[i]
		}
	}
	attrsum[2] += (sets[1] / Setnum[1]) * 15
	attrsum[6] += (sets[2] / Setnum[2]) * 15
	attrsum[8] += (sets[3] / Setnum[3]) * 25
	attrsum[9] += (sets[4] / Setnum[4]) * 12
	attrsum[10] += (sets[5] / Setnum[5]) * 40
	attrsum[12] += (sets[6] / Setnum[6]) * 15
	attrsum[11] += (sets[7] / Setnum[7]) * 15
	attrsum[4] += (sets[8] / Setnum[8]) * 35

	Attr.HP = monster.mon.Con*15*(100+attrsum[2])/100 + attrsum[1]
	Attr.ATK = monster.mon.Atk*(100+attrsum[4])/100 + attrsum[3]
	Attr.DEF = monster.mon.Def*(100+attrsum[6])/100 + attrsum[5]
	Attr.CRI = monster.mon.Critical_rate + attrsum[9]
	if Attr.CRI > 100 {
		Attr.CRI = 100
	}
	Attr.CRIDMG = monster.mon.Critical_damage + attrsum[10]
	Attr.SPD = monster.mon.Spd + attrsum[8]
	Attr.ACC = monster.mon.Accuracy + attrsum[12]
	Attr.RES = monster.mon.Resist + attrsum[11]
	switch monsterworth.Type {
	case 1:
		Attr.DPS = float32(Attr.ATK) * float32((10000 + Attr.CRI*Attr.CRIDMG)) / 10000 * float32(Attr.SPD) / 110 / 10 * 2
	case 2:
		Attr.DPS = float32(Attr.ATK) * (float32(Attr.SPD)/420 + 0.5) * float32((10000 + Attr.CRI*Attr.CRIDMG)) / 10000 * float32(Attr.SPD) / 110 / 10 * 2
	case 3:
		Attr.DPS = (float32(Attr.ATK)*0.15 + float32(Attr.HP)*0.035) * float32((10000 + Attr.CRI*Attr.CRIDMG)) / 10000 * float32(Attr.SPD) / 110 / 10 * 2
	case 4:
		Attr.DPS = (float32(Attr.ATK)*0.4 + float32(Attr.DEF)*0.6) * float32((10000 + Attr.CRI*Attr.CRIDMG)) / 10000 * float32(Attr.SPD) / 110 / 10 * 2
	case 5:
		Attr.DPS = (float32(Attr.ATK)*0.2 + 2000) * float32((10000 + Attr.CRI*Attr.CRIDMG)) / 10000 * float32(Attr.SPD) / 110 / 10 * 2
	case 6:
		Attr.DPS = float32(Attr.ATK) * float32((10000 + Attr.CRI*Attr.CRIDMG)) / 10000 / 10 * 2
	}
	if sets[13] >= 4 {
		Attr.DPS = Attr.DPS * 1.27
	}
	Attr.EHP = float32(Attr.HP * (333 + Attr.DEF) / 287 / 150)
	Attr.EHPDPS = float32(math.Sqrt(float64(Attr.DPS*Attr.EHP))) * 2
	Attr.SUM = Attr.DPS*monsterworth.DPS + Attr.EHP*monsterworth.EHP + Attr.EHPDPS*monsterworth.EHPDPS + float32(Attr.SPD)*monsterworth.SPDFUNC + float32(Attr.HP)*monsterworth.HPFUNC + float32(Attr.ATK)*monsterworth.ATKFUNC + float32(Attr.DEF)*monsterworth.DEFFUNC + float32(Attr.CRI)*monsterworth.CRIFUNC
	if int32(Attr.ACC) > monsterworth.ACCTAR {
		Attr.SUM += float32(monsterworth.ACCTAR)*monsterworth.ACC + float32(int32(Attr.ACC)-monsterworth.ACCTAR)*monsterworth.ACC/3
	} else {
		Attr.SUM += float32(Attr.ACC) * monsterworth.ACC
	}
	if int32(Attr.RES) > monsterworth.RESTAR {
		Attr.SUM += float32(monsterworth.RESTAR)*monsterworth.RES + float32(int32(Attr.RES)-monsterworth.RESTAR)*monsterworth.RES/3
	} else {
		Attr.SUM += float32(Attr.RES) * monsterworth.RES
	}
	for i := 10; i < 24; i++ {
		if i == 15 && sets[i] >= 2 {
			Attr.SUM += monsterworth.SET[i]
		} else {
			Attr.SUM += float32(sets[i]/Setnum[i]) * monsterworth.SET[i]
		}

	}
	//fmt.Println("HP:", HP, "ATK:", ATK, "DEF:", DEF, "CRI:", CRI, "CRIDMG:", CRIDMG, "SPD:", SPD, "ACC:", ACC, "RES:", RES, "DPS:", DPS, "EHP:", EHP, "EHPDPS:", EHPDPS)
	return Attr
}
