package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sky/exports"
	"strconv"
	"strings"
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
	Name            string
}
type MyMonster struct {
	ID           int
	mon          *Monster
	OldCRunes    CompleteRunes
	NewCRunes    CompleteRunes
	OldAttr      MonsterAttribute
	NewAttr      MonsterAttribute
	toprunes     [7][20](*monrune)
	bestsetrune  [7][24](*monrune)
	toprunesnum  [7]int
	monsterworth exports.Worth
	monrunes     []monrune
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
	p            *Rune
	From         *Monster
	times        int
}
type monrune struct {
	r          *rune
	attrsum    float32
	sumwithset float32
}
type CompleteRunes struct {
	runes [7]*rune
}

var Setnum = []int{9, 2, 2, 4, 2, 4, 2, 2, 4, 9, 4, 4, 9, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var maxrune6 = []int{0, 2448, 63, 160, 63, 160, 63, 0, 42, 58, 80, 64, 64}
var maxrune5 = []int{0, 2050, 51, 135, 51, 135, 51, 0, 39, 47, 59, 50, 50}
var attrb = []float32{0, 0, 6.5, 0, 6.5, 0, 6.5, 0, 5, 5, 5.5, 6, 6}
var savenum int = 10
var Attrname = []string{"x", "HP", "HP%", "ATK", "ATK%", "DEF", "DEF%", "x", "SPD", "CRI", "CRIDMG", "RES", "ACC"}
var Setname = []string{"x", "祝福", "守护", "迅速", "刀刃", "激怒", "集中", "忍耐", "猛攻", "x", "绝望", "吸血", "x", "暴走", "应报", "意志", "保护", "反击", "覆灭", "斗志", "决心", "发扬", "命中", "韧性"}
var SetBaseAttrType = []int{0, 2, 6, 8, 9, 10, 12, 11, 4}
var SetBaseAttrValue = []float32{0, 15, 15, 25, 12, 40, 20, 20, 35}
var myrunes []rune

func main() {
	var game Gameinfo
	var MyMonsters []MyMonster
	dat, _ := ioutil.ReadFile("15060791.json")
	f, _ := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
	os.Stdout = f
	os.Stderr = f
	dattmp := string(dat)
	dattmp = strings.Replace(dattmp, "\r", "", -1)
	dattmp = strings.Replace(dattmp, "\n", "", -1)
	dattmp = strings.Replace(dattmp, " ", "", -1)
	dattmp = strings.Replace(dattmp, "\"runes\":{", "\"runes\": [", -1)
	for i := 0; i < 15; i++ {
		strtmp := "\"" + strconv.Itoa(i) + "\":{\"occupied_type\":"
		dattmp = strings.Replace(dattmp, strtmp, "{\"occupied_type\":", -1)
	}
	dattmp = strings.Replace(dattmp, "},\"exp_gained\":", "],\"exp_gained\":", -1)
	dat = []byte(dattmp)
	if err := json.Unmarshal(dat, &game); err != nil {
		log.Fatal(err)
	}
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
	var wor []exports.Worth
	wordat, _ := ioutil.ReadFile("../exports/worth.json")
	if err := json.Unmarshal(wordat, &wor); err != nil {
		log.Fatal(err)
	}
	for t, monster := range game.Unit_list {
		for _, w := range wor {
			str := strconv.FormatInt(monster.Unit_id, 10)
			if str == w.Unitid {
				game.Unit_list[t].Name = w.Name
			}
		}
		var montmp MyMonster
		for i, r := range monster.Runes {
			var tmp rune
			tmp.Set_id = r.Set_id
			tmp.Upgrade_curr = r.Upgrade_curr
			tmp.Slot_no = r.Slot_no
			tmp.Class = r.Class
			tmp.p = &(game.Unit_list[t].Runes[i])
			tmp.From = &game.Unit_list[t]
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
			montmp.OldCRunes.runes[tmp.Slot_no] = &(myrunes[len(myrunes)-1])
		}
		montmp.mon = &game.Unit_list[t]
		MyMonsters = append(MyMonsters, montmp)
	}

	for id := 1; id < 44; id++ {
		var monsterworth exports.Worth
		var monster MyMonster

		for _, monsterworth = range wor {
			if monsterworth.Id == int32(id) {
				break
			}
		}
		for _, monster = range MyMonsters {
			str := strconv.FormatInt(monster.mon.Unit_id, 10)
			if str == monsterworth.Unitid {
				monster.ID = id
				break
			}
		}
		monster.monsterworth = monsterworth
		for i, r := range myrunes {
			if r.times > 0 {
				continue
			}
			var tmpmonrune monrune
			sum := (float32(r.attr[1])/float32(monster.mon.Con)/15 + float32(r.attr[2])) * monsterworth.HP / attrb[2]
			sum += (float32(r.attr[3])/float32(monster.mon.Atk) + float32(r.attr[4])) * monsterworth.ATK / attrb[4]
			sum += (float32(r.attr[5])/float32(monster.mon.Def) + float32(r.attr[6])) * monsterworth.DEF / attrb[6]
			sum += float32(r.attr[8])*monsterworth.SPD/attrb[8] + float32(r.attr[9])*monsterworth.CRI/attrb[9] + float32(r.attr[10])*monsterworth.CRIDMG/attrb[10] + float32(r.attr[11])*monsterworth.RES/attrb[11] + float32(r.attr[12])*monsterworth.ACC/attrb[12]
			tmpmonrune.attrsum = sum
			sum += monsterworth.SET[r.Set_id] / float32(Setnum[r.Set_id])
			if r.Set_id == 13 {
				sum += (monsterworth.SET[4]*2 + monsterworth.SET[5] + monsterworth.SET[8]) / 3 / 4
			}
			tmpmonrune.sumwithset = sum
			tmpmonrune.r = &myrunes[i]
			monster.monrunes = append(monster.monrunes, tmpmonrune)

			var j int
			for j = 0; j < monster.toprunesnum[r.Slot_no]; j++ {
				if monster.toprunes[r.Slot_no][j] == nil || (*monster.toprunes[r.Slot_no][j]).sumwithset < sum {
					break
				}
			}
			if j < savenum {
				for k := j + 1; k < savenum; k++ {
					if k > monster.toprunesnum[r.Slot_no] {
						break
					}
					monster.toprunes[r.Slot_no][k] = monster.toprunes[r.Slot_no][k-1]
				}
				monster.toprunes[r.Slot_no][j] = &(monster.monrunes[len(monster.monrunes)-1])
				if monster.toprunesnum[r.Slot_no] < savenum {
					monster.toprunesnum[r.Slot_no]++
				}
			}
		}

		for i, r := range monster.monrunes {
			if monster.bestsetrune[r.r.Slot_no][r.r.Set_id] == nil || r.sumwithset > (*monster.bestsetrune[r.r.Slot_no][r.r.Set_id]).sumwithset {
				monster.bestsetrune[r.r.Slot_no][r.r.Set_id] = &monster.monrunes[i]
			}
		}
		for i := 1; i < 7; i++ {
			for j := 1; j < 24; j++ {
				if monster.bestsetrune[i][j] != nil && (*monster.bestsetrune[i][j]).attrsum < (*monster.toprunes[i][savenum-1]).attrsum*0.8 {
					monster.bestsetrune[i][j] = nil
					continue
				}
				for k := 0; k < savenum; k++ {
					if monster.toprunes[i][k] == monster.bestsetrune[i][j] {
						monster.bestsetrune[i][j] = nil
					}
				}
			}
		}
		fmt.Println(monster.ID, " ", monster.monsterworth.Name, " ", monster.mon.Unit_id)
		nowsum := caculate(monster.OldCRunes, &monster) //caculate(&myrunes[len(myrunes)-6], &myrunes[len(myrunes)-5], &myrunes[len(myrunes)-4], &myrunes[len(myrunes)-3], &myrunes[len(myrunes)-2], &myrunes[len(myrunes)-1], monster)
		fmt.Println("now:", nowsum.SUM, "\tDPS:", nowsum.DPS, "\tEHP:", nowsum.EHP, "\tEHPDPS:", nowsum.EHPDPS, "\tSPD:", nowsum.SPD, "\tHP:", nowsum.HP,
			"\tATK:", nowsum.ATK, "\tDEF:", nowsum.DEF, "\tCRI:", nowsum.CRI, "\tCRIDMG:", nowsum.CRIDMG, "\tACC:", nowsum.ACC, "\tRES:", nowsum.RES)
		var Cr CompleteRunes
		search(Cr, 1, &monster)
		//标记用过的符文
		for i := 1; i < 7; i++ {
			monster.NewCRunes.runes[i].times++
		}
		maxsum := caculate(monster.NewCRunes, &monster)
		fmt.Println("find:", maxsum.SUM, "  DPS:", maxsum.DPS, "  EHP:", maxsum.EHP, "  EHPDPS:", maxsum.EHPDPS, "  SPD:", maxsum.SPD, "  HP:", maxsum.HP,
			"  ATK:", maxsum.ATK, "  DEF:", maxsum.DEF, "  CRI:", maxsum.CRI, "  CRIDMG:", maxsum.CRIDMG, "  ACC:", maxsum.ACC, "  RES:", maxsum.RES)
		fmt.Println("Old:")
		for i := 1; i < 7; i++ {
			if monster.OldCRunes.runes[i] == nil {
				continue
			}
			fmt.Print(i, ": set:", Setname[(*monster.OldCRunes.runes[i].p).Set_id], "\tmain:", Attrname[(*monster.OldCRunes.runes[i].p).Pri_eff[0]], "+", (*monster.OldCRunes.runes[i].p).Pri_eff[1])
			if (*monster.OldCRunes.runes[i].p).Prefix_eff[0] > 0 {
				fmt.Print("\tPri:", Attrname[(*monster.OldCRunes.runes[i].p).Prefix_eff[0]], "+", (*monster.OldCRunes.runes[i].p).Prefix_eff[1])
			} else {
				fmt.Print("\t\t")
			}
			fmt.Print("\tAttr:")
			for j := 0; j < len((*monster.OldCRunes.runes[i].p).Sec_eff); j++ {
				fmt.Print("\t", Attrname[(*monster.OldCRunes.runes[i].p).Sec_eff[j][0]], "+", (*monster.OldCRunes.runes[i].p).Sec_eff[j][1]+(*monster.OldCRunes.runes[i].p).Sec_eff[j][3])
			}
			if monster.OldCRunes.runes[i].From == nil {
				fmt.Println("\tFrom:背包")
			} else {
				fmt.Println("\tFrom:", monster.OldCRunes.runes[i].From.Name)
			}
		}
		fmt.Println("New:")
		for i := 1; i < 7; i++ {
			if monster.NewCRunes.runes[i].p == nil {
				continue
			}
			fmt.Print(i, ": set:", Setname[(*monster.NewCRunes.runes[i].p).Set_id], "\tmain:", Attrname[(*monster.NewCRunes.runes[i].p).Pri_eff[0]], "+", (*monster.NewCRunes.runes[i].p).Pri_eff[1])
			if (*monster.NewCRunes.runes[i].p).Prefix_eff[0] > 0 {
				fmt.Print("\tPri:", Attrname[(*monster.NewCRunes.runes[i].p).Prefix_eff[0]], "+", (*monster.NewCRunes.runes[i].p).Prefix_eff[1])
			} else {
				fmt.Print("\t\t")
			}
			fmt.Print("\tAttr:")
			for j := 0; j < len((*monster.NewCRunes.runes[i].p).Sec_eff); j++ {
				fmt.Print("\t", Attrname[(*monster.NewCRunes.runes[i].p).Sec_eff[j][0]], "+", (*monster.NewCRunes.runes[i].p).Sec_eff[j][1]+(*monster.NewCRunes.runes[i].p).Sec_eff[j][3])
			}
			if monster.NewCRunes.runes[i].From == nil {
				fmt.Println("\tFrom:背包")
			} else {
				fmt.Println("\tFrom:", monster.NewCRunes.runes[i].From.Name)
			}
		}
	}
	// for i := 43; i > 0; i-- {
	// 	for _, r := range myrunes {
	// 		if r.times == i {
	// 			fmt.Print("slot:", r.Slot_no, "\rtimes:", r.times, "\rset:", Setname[r.Set_id], "\tmain:", Attrname[r.p.Pri_eff[0]], "+", r.p.Pri_eff[1])
	// 			if r.p.Prefix_eff[0] > 0 {
	// 				fmt.Print("\tPri:", Attrname[r.p.Prefix_eff[0]], "+", r.p.Prefix_eff[1])
	// 			} else {
	// 				fmt.Print("\t\t")
	// 			}
	// 			fmt.Print("\tAttr:")
	// 			for j := 0; j < len(r.p.Sec_eff); j++ {
	// 				fmt.Print("\t", Attrname[r.p.Sec_eff[j][0]], "+", r.p.Sec_eff[j][1]+r.p.Sec_eff[j][3])
	// 			}
	// 			if r.From == nil {
	// 				fmt.Println("\tFrom:背包")
	// 			} else {
	// 				fmt.Println("\tFrom:", r.From.Name)
	// 			}
	// 		}
	// 	}
	// }
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
		Cr.runes[slot] = (*mymon.toprunes[slot][i]).r
		search(Cr, slot+1, mymon)
	}
	for i := 1; i < 24; i++ {
		if mymon.bestsetrune[slot][i] == nil {
			continue
		}
		Cr.runes[slot] = (*mymon.bestsetrune[slot][i]).r
		search(Cr, slot+1, mymon)
	}
}
func caculate(Cr CompleteRunes, monster *MyMonster) MonsterAttribute {
	var sets [24]int
	var Attr MonsterAttribute
	Attr.Runes = Cr
	for i := 1; i < 7; i++ {
		if Cr.runes[i] != nil {
			sets[(*Cr.runes[i]).Set_id]++
		}
	}
	var attrsum [13]int
	for i := 1; i < 13; i++ {
		for j := 1; j < 7; j++ {
			if Cr.runes[j] != nil {
				attrsum[i] += (*Cr.runes[j]).attr[i]
			}
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
	if Attr.CRI > 85 {
		Attr.CRI = 85 + (Attr.CRI-85)/2
	}
	if Attr.ACC > 100 {
		Attr.ACC = 100
	}
	if Attr.RES > 100 {
		Attr.RES = 100
	}
	Attr.CRIDMG = monster.mon.Critical_damage + attrsum[10]
	Attr.SPD = monster.mon.Spd + attrsum[8]
	Attr.ACC = monster.mon.Accuracy + attrsum[12]
	Attr.RES = monster.mon.Resist + attrsum[11]
	switch monster.monsterworth.Type {
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
	//Special monster
	if monster.mon.Unit_id == 1325114940 {
		CRI := Attr.CRI + 15
		if CRI > 100 {
			CRI = 100
		}
		Attr.DPS = float32(Attr.ATK) * (float32(Attr.SPD)/420 + 0.5) * float32((10000 + CRI*Attr.CRIDMG)) / 10000 * float32(Attr.SPD) / 110 / 10 * 2
	}

	if sets[13] >= 4 {
		Attr.DPS = Attr.DPS * 1.27
	}

	Attr.EHP = float32(Attr.HP * (333 + Attr.DEF) / 287 / 150)
	Attr.EHPDPS = float32(math.Sqrt(float64(Attr.DPS*Attr.EHP))) * 2
	Attr.SUM = Attr.DPS*monster.monsterworth.DPS + Attr.EHP*monster.monsterworth.EHP + Attr.EHPDPS*monster.monsterworth.EHPDPS + float32(Attr.SPD)*monster.monsterworth.SPDFUNC + float32(Attr.HP)*monster.monsterworth.HPFUNC + float32(Attr.ATK)*monster.monsterworth.ATKFUNC + float32(Attr.DEF)*monster.monsterworth.DEFFUNC + float32(Attr.CRI)*monster.monsterworth.CRIFUNC
	if int32(Attr.ACC) > monster.monsterworth.ACCTAR {
		Attr.SUM += float32(monster.monsterworth.ACCTAR)*monster.monsterworth.ACC + float32(int32(Attr.ACC)-monster.monsterworth.ACCTAR)*monster.monsterworth.ACC/5
	} else {
		Attr.SUM += float32(Attr.ACC) * monster.monsterworth.ACC
	}
	if int32(Attr.RES) > monster.monsterworth.RESTAR {
		Attr.SUM += float32(monster.monsterworth.RESTAR)*monster.monsterworth.RES + float32(int32(Attr.RES)-monster.monsterworth.RESTAR)*monster.monsterworth.RES/5
	} else {
		Attr.SUM += float32(Attr.RES) * monster.monsterworth.RES
	}
	for i := 10; i < 24; i++ {
		if i == 15 && sets[i] >= 2 {
			Attr.SUM += monster.monsterworth.SET[i]
		} else {
			Attr.SUM += float32(sets[i]/Setnum[i]) * monster.monsterworth.SET[i]
		}

	}
	//fmt.Println("HP:", HP, "ATK:", ATK, "DEF:", DEF, "CRI:", CRI, "CRIDMG:", CRIDMG, "SPD:", SPD, "ACC:", ACC, "RES:", RES, "DPS:", DPS, "EHP:", EHP, "EHPDPS:", EHPDPS)
	return Attr
}
