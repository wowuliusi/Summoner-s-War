package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sky/exports"
)

type Server struct {
	Occupied_type string
	Extra         string
}

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
type rune struct {
	Set_id       int
	Upgrade_curr int
	Class        int
	Slot_no      int
	attr         [13]int
	sum          float32
	p            *Rune
}

var Setnum = []int{9, 2, 2, 4, 2, 4, 2, 2, 4, 9, 4, 4, 9, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var maxrune6 = []int{0, 2448, 63, 160, 63, 160, 63, 0, 42, 58, 80, 64, 64}
var maxrune5 = []int{0, 2050, 51, 135, 51, 135, 51, 0, 39, 47, 59, 50, 50}
var attrb = []float32{0, 0, 6.5, 0, 6.5, 0, 6.5, 0, 5, 5, 5.5, 6, 6}
var monsterworth exports.Worth

func main() {
	var game Gameinfo
	var myrunes []rune
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
	var monsterid = "1998233295"
	var monsterid2 int64 = 1998233295
	var monster Monster
	var wor []exports.Worth
	wordat, _ := ioutil.ReadFile("exports/worth.json")
	if err := json.Unmarshal(wordat, &wor); err != nil {
		log.Fatal(err)
	}
	for _, monsterworth = range wor {
		if monsterworth.Unitid == monsterid {
			break
		}
	}
	for _, monster = range game.Unit_list {
		if monster.Unit_id == monsterid2 {
			for i, r := range monster.Runes {
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
			break
		}
	}
	//fmt.Println(len(myrunes))
	var toprunes [7][24](*rune)
	for i, r := range myrunes {
		sum := (float32(r.attr[1])/float32(monster.Con)/15 + float32(r.attr[2])) * monsterworth.HP / attrb[2]
		sum += (float32(r.attr[3])/float32(monster.Atk) + float32(r.attr[4])) * monsterworth.ATK / attrb[4]
		sum += (float32(r.attr[5])/float32(monster.Def) + float32(r.attr[6])) * monsterworth.DEF / attrb[6]
		sum += float32(r.attr[8])*monsterworth.SPD/attrb[8] + float32(r.attr[9])*monsterworth.CRI/attrb[9] + float32(r.attr[10])*monsterworth.CRIDMG/attrb[10] + float32(r.attr[11])*monsterworth.RES/attrb[11] + float32(r.attr[12])*monsterworth.ACC/attrb[12]
		myrunes[i].sum = sum
		if toprunes[r.Slot_no][r.Set_id] == nil {
			toprunes[r.Slot_no][r.Set_id] = &(myrunes[i])
			//fmt.Println(r.Slot_no, " ", r.Set_id, " ", sum)
		} else {
			if sum > (*(toprunes[r.Slot_no][r.Set_id])).sum {
				//fmt.Println(r.Slot_no, " ", r.Set_id, " ", sum, " ", (*(toprunes[r.Slot_no][r.Set_id])).sum)
				toprunes[r.Slot_no][r.Set_id] = &(myrunes[i])
			}
		}
	}
	nowsum := caculate(&myrunes[len(myrunes)-6], &myrunes[len(myrunes)-5], &myrunes[len(myrunes)-4], &myrunes[len(myrunes)-3], &myrunes[len(myrunes)-2], &myrunes[len(myrunes)-1], monster)
	fmt.Println("now:", nowsum)
	var maxsum float32 = 0
	var maxrune [7]int
	for r1 := 1; r1 < 24; r1++ {
		if toprunes[1][r1] != nil {
			for r2 := 1; r2 < 24; r2++ {
				if toprunes[2][r2] != nil {
					for r3 := 1; r3 < 24; r3++ {
						if toprunes[3][r3] == nil {
							continue
						}
						for r4 := 1; r4 < 24; r4++ {
							if toprunes[4][r4] == nil {
								continue
							}
							for r5 := 1; r5 < 24; r5++ {
								if toprunes[5][r5] == nil {
									continue
								}
								for r6 := 1; r6 < 24; r6++ {
									if toprunes[6][r6] == nil {
										continue
									}
									var sum float32 = caculate(toprunes[1][r1], toprunes[2][r2], toprunes[3][r3], toprunes[4][r4], toprunes[5][r5], toprunes[6][r6], monster)
									if sum > maxsum {
										maxsum = sum
										maxrune[1] = r1
										maxrune[2] = r2
										maxrune[3] = r3
										maxrune[4] = r4
										maxrune[5] = r5
										maxrune[6] = r6
									}
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("find:", maxsum)
	for i := 1; i < 7; i++ {
		fmt.Println(i, ": set:", (*toprunes[i][maxrune[i]]).Set_id, " sum:", (*toprunes[i][maxrune[i]]).sum, (*(*toprunes[i][maxrune[i]]).p).Sec_eff)
	}
}

func caculate(a1, a2, a3, a4, a5, a6 *rune, monster Monster) float32 {
	var sets [24]int
	sets[(*a1).Set_id]++
	sets[(*a2).Set_id]++
	sets[(*a3).Set_id]++
	sets[(*a4).Set_id]++
	sets[(*a5).Set_id]++
	sets[(*a6).Set_id]++

	// var sum float32
	// for i := 1; i < 24; i++ {
	// 	sum += float32(sets[i]/Setnum[i]) * monsterworth.SET[i]
	// }
	// sum += (*a1).sum + (*a2).sum + (*a3).sum + (*a4).sum + (*a5).sum + (*a6).sum

	var attrsum [13]int
	for i := 1; i < 13; i++ {
		attrsum[i] = (*a1).attr[i] + (*a2).attr[i] + (*a3).attr[i] + (*a4).attr[i] + (*a5).attr[i] + (*a6).attr[i]
	}
	attrsum[2] += (sets[1] / Setnum[1]) * 15
	attrsum[6] += (sets[2] / Setnum[2]) * 15
	attrsum[8] += (sets[3] / Setnum[3]) * 25
	attrsum[9] += (sets[4] / Setnum[4]) * 12
	attrsum[10] += (sets[5] / Setnum[5]) * 40
	attrsum[12] += (sets[6] / Setnum[6]) * 15
	attrsum[11] += (sets[7] / Setnum[7]) * 15
	attrsum[4] += (sets[8] / Setnum[8]) * 35

	HP := monster.Con*15*(100+attrsum[2])/100 + attrsum[1]
	ATK := monster.Atk*(100+attrsum[4])/100 + attrsum[3]
	DEF := monster.Def*(100+attrsum[6])/100 + attrsum[5]
	CRI := monster.Critical_rate + attrsum[9]
	if CRI > 100 {
		CRI = 100
	}
	CRIDMG := monster.Critical_damage + attrsum[10]
	SPD := monster.Spd + attrsum[8]
	var ACC float32 = float32(monster.Accuracy + attrsum[12])
	var RES float32 = float32(monster.Resist + attrsum[11])

	var DPS float32
	switch monsterworth.Type {
	case 1:
		DPS = float32(ATK) * float32((10000 + CRI*CRIDMG)) / 10000 * float32(SPD) / 110 / 10 * 2
	case 2:
		DPS = float32(ATK) * (float32(SPD)/420 + 0.5) * float32((10000 + CRI*CRIDMG)) / 10000 * float32(SPD) / 110 / 10 * 2
	case 3:
		DPS = (float32(ATK)*0.15 + float32(HP)*0.035) * float32((10000 + CRI*CRIDMG)) / 10000 * float32(SPD) / 110 / 10 * 2
	case 4:
		DPS = (float32(ATK)*0.4 + float32(DEF)*0.6) * float32((10000 + CRI*CRIDMG)) / 10000 * float32(SPD) / 110 / 10 * 2
	case 5:
		DPS = (float32(ATK)*0.2 + 2000) * float32((10000 + CRI*CRIDMG)) / 10000 * float32(SPD) / 110 / 10 * 2
	}
	if sets[13] >= 4 {
		DPS = DPS * 1.27
	}
	var EHP float32 = float32(HP * (333 + DEF) / 287 / 150)
	EHPDPS := float32(math.Sqrt(float64(DPS*EHP))) * 2
	Final := DPS*monsterworth.DPS + EHP*monsterworth.EHP + EHPDPS*monsterworth.EHPDPS + float32(SPD)*monsterworth.SPDFUNC + float32(HP)*monsterworth.HPFUNC + float32(ATK)*monsterworth.ATKFUNC + float32(DEF)*monsterworth.DEFFUNC + float32(CRI)*monsterworth.CRIFUNC
	if ACC > monsterworth.ACCTAR {
		Final += monsterworth.ACCTAR*monsterworth.ACC + (ACC-monsterworth.ACCTAR)*monsterworth.ACC/3
	} else {
		Final += ACC * monsterworth.ACC
	}
	if RES > monsterworth.RESTAR {
		Final += monsterworth.RESTAR*monsterworth.RES + (RES-monsterworth.RESTAR)*monsterworth.RES/3
	} else {
		Final += RES * monsterworth.RES
	}
	for i := 10; i < 24; i++ {
		if i == 15 && sets[i] >= 2 {
			Final += monsterworth.SET[i]
		} else {
			Final += float32(sets[i]/Setnum[i]) * monsterworth.SET[i]
		}

	}

	//fmt.Println("HP:", HP, "ATK:", ATK, "DEF:", DEF, "CRI:", CRI, "CRIDMG:", CRIDMG, "SPD:", SPD, "ACC:", ACC, "RES:", RES, "DPS:", DPS, "EHP:", EHP, "EHPDPS:", EHPDPS)
	return Final
}
