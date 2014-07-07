// keyboard
package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Finger string

type BitKey struct {
	key byte //top 3 bits - modirier index; low 5bits - key index
}

func (bk *BitKey) Key() byte {
	return bk.key & 31 //b00011111
}
func (bk *BitKey) Mod() byte {
	return bk.key >> 5
}
func (bk *BitKey) String() string {
	return fmt.Sprintf("%d %d", bk.Mod(), bk.Key())
}

func MakeBitKey(m byte, k byte) BitKey {
	return BitKey{m << 5 & (k)}
}

var AllBitKeys []byte

type Key struct {
	Mod string
	Key string
}

type FingerFinger struct {
	Fin1 string
	Fin2 string
}
type FingerKey struct {
	Fin Finger
	Key Key
}

type Configuration struct {
	Users           []string
	Groups          []string
	Fin             []string
	FinFinWeight    [][]string
	FinKeys         map[string][]string
	FinKeyWeight    [][]string
	ModFinger       map[string]string
	ModScore        map[string]int
	Actions         []string
	FixedModsGroups map[string][]string
	SameModGroups   [][]string
	SameKeyGroups   [][]string
	SeqOfActions    [][]string
	DiffKeyGroups   [][]string
}

var Fin []string
var FinFinWeight map[FingerFinger]int
var FinKeys map[Finger][]Key
var FinKeyWeight map[FingerKey]int
var ModFinger map[string]string
var ModScore map[string]int
var Actions []string
var FixedModsGroups map[string][]string
var SameModGroups [][]string
var SameKeyGroups [][]string
var DiffKeyGroups [][]string
var SeqOfActions [][]string
var SrcAllMods []string
var SrcAllKeys []string
var AllMods []string
var AllKeys []Key
var KeyFinger map[Key]Finger
var SeqOfActionsIndexed [][]int
var SameModGroupsIndexed [][]int
var SameKeyGroupsIndexed [][]int
var DiffKeyGroupsIndexed [][]int
var FixedModsGroupsIndexed map[string][]int

func init() {
	//Fin = make([]string)
	fmt.Println("Init configs")

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	Fin = config.Fin
	FinFinWeight = make(map[FingerFinger]int)
	for _, v := range config.FinFinWeight {
		FinFinWeight[FingerFinger{v[0], v[1]}], err = strconv.Atoi(v[2])

	}
	//FinKeys = config.FinKeys

	ModFinger = config.ModFinger
	ModScore = config.ModScore
	Actions = config.Actions
	FixedModsGroups = config.FixedModsGroups
	SameModGroups = config.SameModGroups
	SameKeyGroups = config.SameKeyGroups
	SeqOfActions = config.SeqOfActions
	DiffKeyGroups = config.DiffKeyGroups

	AllMods = []string{}
	for mod, _ := range ModScore {
		AllMods = append(AllMods, mod)
	}
	//generete keys with modifiers
	FinKeys = make(map[Finger][]Key)
	for f, keys := range config.FinKeys {
		modKeys := []Key{}
		for mod, _ := range ModScore {
			if f != ModFinger[mod] {
				for _, k := range keys {
					modKeys = append(modKeys, Key{mod, k})
				}
			}
		}
		FinKeys[Finger(f)] = modKeys
	}
	//generete finger=>mod+key scores
	FinKeyWeight = make(map[FingerKey]int)
	for _, v := range config.FinKeyWeight {
		f, key := v[0], v[1]
		score, _ := strconv.Atoi(v[2])
		for mod, modScore := range ModScore {
			if f != ModFinger[mod] {
				newScore := score + modScore
				if newScore < 1 {
					newScore = 1
				}
				FinKeyWeight[FingerKey{Finger(f), Key{mod, key}}] = newScore
			}
		}

	}
	//формируем список клавиш так, чтобы более доступные клавиши шли первыми
	//TODO: учитывать также и последовательности
	AllKeys = []Key{}
	var dWeightKeys map[int][]Key
	dWeightKeys = make(map[int][]Key)

	for sFinKey, iWeight := range FinKeyWeight {
		val, ok := dWeightKeys[iWeight]
		if !ok {
			dWeightKeys[iWeight] = []Key{sFinKey.Key}
		} else {
			dWeightKeys[iWeight] = append(val, sFinKey.Key)
		}
	}
	var weights []int
	for w := range dWeightKeys {
		weights = append(weights, w)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(weights)))
	for _, w := range weights {
		for _, k := range dWeightKeys[w] {
			AllKeys = append(AllKeys, k)
		}
	}

	AllBitKeys = make([]byte, len(AllKeys))

	for ki, _ := range AllKeys {
		//mod, key := kk.Mod, kk.Key
		AllBitKeys[ki] = byte(ki)

	}

	//#generete key=>finger mapping for performance
	KeyFinger = make(map[Key]Finger)
	for f, keys := range FinKeys {
		for _, k := range keys {
			KeyFinger[k] = f
		}
	}

	//#generete lSeqOfActions indexes to lActions for performance
	SeqOfActionsIndexed = make([][]int, len(SeqOfActions))
	for si, seq := range SeqOfActions {
		SeqOfActionsIndexed[si] = make([]int, len(seq))
		for ai, aa := range seq {
			found := false
			for lai, la := range Actions {
				if la == aa {
					SeqOfActionsIndexed[si][ai] = lai
					found = true
					break
				}
			}
			if !found {
				fmt.Println("--------------FAIL--------------")
			}

		}
	}
	//#generete lActionsGroups indexes to lActions  for performance
	SameModGroupsIndexed = make([][]int, len(SameModGroups))
	for si, seq := range SameModGroups {
		SameModGroupsIndexed[si] = make([]int, len(seq))
		for ai, aa := range seq {
			found := false
			for lai, la := range Actions {
				if la == aa {
					SameModGroupsIndexed[si][ai] = lai
					found = true
					break
				}
			}
			if !found {
				fmt.Println("--------------FAIL--------------")
			}
		}
	}
	//#generete lSameKeyGroups indexes to lActions  for performance
	SameKeyGroupsIndexed = make([][]int, len(SameKeyGroups))
	for si, seq := range SameKeyGroups {
		SameKeyGroupsIndexed[si] = make([]int, len(seq))
		for ai, aa := range seq {
			found := false
			for lai, la := range Actions {
				if la == aa {
					SameKeyGroupsIndexed[si][ai] = lai
					found = true
					break
				}
			}
			if !found {
				fmt.Println("--------------FAIL--------------")
			}
		}
		//fmt.Println(si, SameKeyGroupsIndexed[si])
	}

	FixedModsGroupsIndexed = make(map[string][]int)
	for mod, actions := range FixedModsGroups {
		FixedModsGroupsIndexed[mod] = make([]int, 0)
		for _, aa := range actions {
			found := false
			for lai, la := range Actions {
				if la == aa {
					FixedModsGroupsIndexed[mod] = append(FixedModsGroupsIndexed[mod], lai)
					found = true
					break
				}
			}
			if !found {
				fmt.Println("--------------FAIL--------------")
			}
		}
	}
	//#generete lSameKeyGroups indexes to lActions  for performance

	DiffKeyGroupsIndexed = make([][]int, len(DiffKeyGroups))
	for si, seq := range DiffKeyGroups {
		DiffKeyGroupsIndexed[si] = make([]int, len(seq))
		for ai, aa := range seq {
			found := false
			for lai, la := range Actions {
				if la == aa {
					DiffKeyGroupsIndexed[si][ai] = lai
					found = true
					break
				}
			}
			if !found {
				fmt.Println("--------------FAIL--------------")
			}
		}
	}

}
