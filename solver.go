package main

import (
	config "./configs"
	svg "./svgdraw"
	"fmt"
	"runtime"
	"sort"
	"time"
)

/*
TODO:
- print action sequences in term of keys

*/

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func UniqInt(col []int) []int {
	m := map[int]struct{}{}
	for _, v := range col {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
		}
	}
	list := make([]int, len(m))

	i := 0
	for v := range m {
		list[i] = v
		i++
	}
	return list
}

func intMin(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}

}
func intMax(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}

}

func keyInArray(a config.Key, arr []config.Key) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

func intInArray(a int, arr []int) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

func byteInArray(a byte, arr []byte) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

type dataJobs struct {
	i       int
	val     []byte
	explain bool
}

type dataResults struct {
	i   int
	val int
}

func Evaluate(id int, jobs <-chan dataJobs, results chan<- dataResults) {
	for j := range jobs {
		lActKeysState := j.val
		explain := j.explain
		sum := 0
		aksLen := len(lActKeysState)

		for ai, ki := range lActKeysState {
			key := config.AllKeys[ki]
			f := config.KeyFinger[key]
			iFingerScore := config.FinKeyWeight[config.FingerKey{f, key}]
			sum = sum + iFingerScore
			if explain {
				fmt.Println(config.Actions[ai], key, iFingerScore)
			}
		}
		if explain {
			fmt.Println("ater FinKeyWeight sum = ", sum)
		}
		for _, seq := range config.SeqOfActionsIndexed {
			prevFinger := config.Finger("nil")
			for _, ai := range seq {
				if ai < aksLen {
					key := config.AllKeys[lActKeysState[ai]]
					f := config.KeyFinger[key]

					iRepScore := 0
					if prevFinger != config.Finger("nil") {
						finfin := config.FingerFinger{string(prevFinger), string(f)}
						iRepScore = config.FinFinWeight[finfin]
						// below trying to reflect that reaching
						// to a key within one finger also have different
						// efforts. so D+F should be greater then Q+B
						if f != prevFinger { //got nothing for same finger
							iRepScore = iRepScore + config.FinKeyWeight[config.FingerKey{prevFinger, key}]
							iRepScore = iRepScore + config.FinKeyWeight[config.FingerKey{f, key}]
						}
						if explain {
							fmt.Println("iRepScore", prevFinger, f, iRepScore)
						}
					}
					prevFinger = f

					sum = sum + iRepScore
					//if aksLen == 2 && (key.Key == "F" || key.Key == "D") {
					//	fmt.Println(sum, iFingerScore, iRepScore)
					//}
				} else {
					prevFinger = config.Finger("nil")
				}
			}

			//fmt.Println(seq, prevFinger)

		}
		if explain {
			fmt.Println("ater SeqOfActions sum = ", sum)
		}
		results <- dataResults{j.i, sum}
	}
}

//func Evaluate(lActKeysState []config.Key) int {
//	sum := 0
//	aksLen := len(lActKeysState)

//	for _, group := range config.ActionsGroupsIndexed {
//		prevMod := "nil"
//		for _, ai := range group {
//			if ai < aksLen {
//				key := lActKeysState[ai]
//				m := key.Mod
//				if prevMod == "nil" {
//					prevMod = m
//				}
//				if prevMod != m {
//					return 0
//				}
//				prevMod = m
//			}
//		}
//	}

//	for _, group := range config.SameKeyGroupsIndexed {
//		prevKey := "nil"
//		for _, ai := range group {
//			if ai < aksLen {
//				key := lActKeysState[ai]
//				k := key.Key
//				if prevKey == "nil" {
//					prevKey = k
//				}
//				if prevKey != k {
//					return 0
//				}
//				prevKey = k
//			}
//		}
//	}

//	for _, key := range lActKeysState {
//		f := config.KeyFinger[key]
//		iFingerScore := config.FinKeyWeight[config.FingerKey{f, key}]
//		sum = sum + iFingerScore
//	}

//	for _, seq := range config.SeqOfActionsIndexed {
//		prevFinger := config.Finger("nil")
//		for _, ai := range seq {
//			if ai < aksLen {
//				key := lActKeysState[ai]
//				f := config.KeyFinger[key]
//				iFingerScore := config.FinKeyWeight[config.FingerKey{f, key}]
//				iRepScore := 0
//				if prevFinger != config.Finger("nil") {
//					finfin := config.FingerFinger{string(prevFinger), string(f)}
//					iRepScore = config.FinFinWeight[finfin]
//				}
//				prevFinger = f
//				sum = sum + iFingerScore + iRepScore
//			}
//		}
//		fmt.Println(seq, prevFinger)
//	}

//	return sum
//}

func GenerateMoves(sUsedKeySets [][]byte, actionIndex int) [][]byte {
	//defer timeTrack(time.Now(), "GenerateMoves")
	skipKeyCount := 0
	lAllStates := make([][]byte, 0, len(sUsedKeySets))
	for _, usedKeys := range sUsedKeySets {
		nextKeys := make([]byte, len(usedKeys)+1, len(usedKeys)+1)
		copy(nextKeys, usedKeys)
		for _, ki := range config.AllBitKeys {
			if !byteInArray(ki, usedKeys) {
				nextKeys[len(usedKeys)] = ki
				aksLen := len(nextKeys)
				state := nextKeys
				mod := config.AllKeys[ki].Mod

				skipKey := false
				if !skipKey {
					act, ok := config.FixedModsGroupsIndexed[mod]
					if ok && intInArray(actionIndex, act) {
						skipKey = true
						skipKeyCount++
					}
				}
				if !skipKey {
					for gi := 0; gi < len(config.SameKeyGroupsIndexed); gi++ {
						//for _, group := range config.SameKeyGroupsIndexed {
						group := config.SameKeyGroupsIndexed[gi]
						prevKey := "nil"
						for ai := 0; ai < len(group); ai++ {
							//for _, ai := range group {
							action := group[ai]
							if action < aksLen {
								key := config.AllKeys[state[action]]
								k := key.Key
								if prevKey == "nil" {
									prevKey = k
								}
								if prevKey != k {
									skipKey = true
									skipKeyCount++
									break
								}
								prevKey = k
							}
						}
						if skipKey {
							break
						}
					}
				}
				if !skipKey {
					for gi := 0; gi < len(config.DiffKeyGroupsIndexed); gi++ {
						//for _, group := range config.DiffKeyGroupsIndexed {
						group := config.DiffKeyGroupsIndexed[gi]
						diffKeys := make(map[string]struct{})
						for ai := 0; ai < len(group); ai++ {
							//for _, ai := range group {
							action := group[ai]
							if action < aksLen {
								key := config.AllKeys[state[action]]
								k := key.Key
								if _, ok := diffKeys[k]; !ok {
									diffKeys[k] = struct{}{}
								} else {
									skipKey = true
									skipKeyCount++
									break
								}
							}
						}
						if skipKey {
							break
						}
					}
				}
				if !skipKey {
					for gi := 0; gi < len(config.SameModGroupsIndexed); gi++ {
						group := config.SameModGroupsIndexed[gi]
						//for _, group := range config.SameModGroupsIndexed {
						prevMod := "nil"
						for ai := 0; ai < len(group); ai++ {
							//for _, ai := range group {
							action := group[ai]
							if action < aksLen {
								key := config.AllKeys[state[action]]
								m := key.Mod
								if prevMod == "nil" {
									prevMod = m
								}
								if prevMod != m {
									skipKey = true
									skipKeyCount++
									break
								}
								prevMod = m
							}
						}
						if skipKey {
							break
						}
					}
				}

				if !skipKey {
					lAllStates = append(lAllStates, nextKeys)
					nextKeys = make([]byte, len(usedKeys)+1, len(usedKeys)+1)
					copy(nextKeys, usedKeys)
				}
			}
		}
	}
	fmt.Println("len(lAllStates):", len(lAllStates), "skipKeyCount: ", skipKeyCount)
	return lAllStates
}

func CutoffMoves(lAllStates [][]byte, v []int,
	topScoreCount int, cutOffNum int, toprint bool) [][]byte {
	//defer timeTrack(time.Now(), "CutoffMoves")

	states := make([][]byte, 0, len(lAllStates))

	scoreVals := UniqInt(v)

	sort.Sort(sort.IntSlice(scoreVals))

	topScoreVals := scoreVals[intMax(0, len(scoreVals)-topScoreCount):]
	for _, tS := range topScoreVals {
		cutOff := 0
		for si := 0; si < len(lAllStates); si++ {
			state := lAllStates[si]
			//for si, state := range lAllStates {
			score := v[si]
			if score == tS {
				states = append(states, state)
				if toprint {
					fmt.Println("Score: ", score)
					for ai, ki := range state {
						actName := config.Actions[ai]
						if actName == "ControlGroupRecall0" ||
							actName == "ControlGroupAppend0" ||
							actName == "ControlGroupAssign0" ||
							actName == "CameraSave0" ||
							actName == "CameraView0" ||
							actName == "BurrowDown" {
							fmt.Println()
						}
						fmt.Print(actName, config.AllKeys[ki], ", ")
					}
					fmt.Println("\n---------------------------")
				}
				cutOff++
			}
			if cutOff >= cutOffNum {
				break
			}
		}
	}

	return states
}

//func main() {
//	config.Init()

//	prevStates := make([][]config.Key, 1)
//	prevStates[0] = make([]config.Key, 0, len(config.Actions)+1)

//	lALen := len(config.Actions)
//	for ai, aa := range config.Actions {
//		sUsedKeySets := prevStates
//		fmt.Println(ai, aa, lALen-1, len(sUsedKeySets))

//		lAllStates := GenerateMoves(sUsedKeySets, ai)
//		if len(lAllStates) == 0 {
//			fmt.Println("Out of keys!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
//			break
//		}
//		v := make([]int, len(lAllStates))
//		for si, state := range lAllStates {
//			v[si] = Evaluate(state)
//		}
//		prevStates = CutoffMoves(lAllStates, v, 10, 5000)
//	}

//}

func main() {
	var m runtime.MemStats
	runtime.GOMAXPROCS(runtime.NumCPU())
	//return
	results := make(chan dataResults, 1000)
	jobs := make(chan dataJobs, 1000)

	for w := 1; w <= runtime.NumCPU(); /*runtime.NumCPU()*2*/ w++ {
		go Evaluate(w, jobs, results)
	}

	prevStates := make([][]byte, 1)
	prevStates[0] = make([]byte, 0, len(config.Actions)+1)
	lALen := len(config.Actions)
	for ai, aa := range config.Actions {
		sUsedKeySets := prevStates
		fmt.Println(ai, aa, lALen-1, len(sUsedKeySets))

		lAllStates := GenerateMoves(sUsedKeySets, ai)
		if len(lAllStates) == 0 {
			panic("Out of keys!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			break
		}
		v := make([]int, len(lAllStates))
		sendCount := 0
		for si, state := range lAllStates {
			jobs <- dataJobs{si, state, false}
			sendCount++

			if sendCount > 500 {
				for i := 0; i < sendCount; i++ {
					d := <-results
					v[d.i] = d.val
				}
				sendCount = 0
			}
		}
		if sendCount > 0 {
			for i := 0; i < sendCount; i++ {
				d := <-results
				v[d.i] = d.val
			}
		}
		//fmt.Println("            len(results)", len(results), "            len(jobs)   ", len(jobs))
		//for si, state := range lAllStates {
		//	fmt.Println(v[si], " === ", state)
		//}
		if ai == lALen-1 {
			prevStates = CutoffMoves(lAllStates, v, 2, 3, true)
			for _, state := range prevStates {
				jobs <- dataJobs{0, state, true}
				_ = <-results
				fmt.Println("\n---------------------------")
			}
		} else {
			//prevStates = CutoffMoves(lAllStates, v, lALen-ai/4, 1000+ai*2, false)
			prevStates = CutoffMoves(lAllStates, v, lALen-ai/4, 1000+ai*2, false)
			//prevStates = CutoffMoves(lAllStates, v, 2, 3, true)
		}
		//for _, state := range prevStates {
		//	for _, ki := range state {
		//		fmt.Print(config.AllKeys[ki])
		//	}
		//	fmt.Println(state)
		//}
		//if ai == 2 {
		//	break
		//}
		runtime.ReadMemStats(&m)
		fmt.Printf("memstat,%d,%d,%d,%d,%d,%d,%d\n", m.HeapSys, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, m.TotalAlloc, m.Mallocs, m.Frees)
	}
	close(jobs)

	for _, m := range config.AllMods {
		keyAction := make(map[string]string)
		for ai, ki := range prevStates[len(prevStates)-1] {
			actName := config.Actions[ai]
			keymod := config.AllKeys[ki]
			key := keymod.Key
			mod := keymod.Mod
			if mod == m {
				keyAction[key] = actName
			}
		}
		svg.DrawKeys(keyAction, "keyboard_"+m+".svg")
	}

}
