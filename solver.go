package main

import (
	config "./configs"
	misc "./misc"
	svg "./svgdraw"
	"fmt"
	"runtime"
	"sort"
)

/*
TODO:
- print action sequences in term of keys
- check configs integrity
- custom memory managment for states
- rename to state in all functions
*/

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
		for ai := 0; ai < aksLen; ai++ {
			ki := lActKeysState[ai]
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

		for si := 0; si < len(config.SeqOfActionsIndexed); si++ {
			seq := config.SeqOfActionsIndexed[si]
			prevFinger := config.Finger("nil")
			for ai := 0; ai < len(seq); ai++ {
				action := seq[ai]
				if action < aksLen {
					key := config.AllKeys[lActKeysState[action]]
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

func GenerateMoves(sUsedKeySets [][]byte, actionIndex int) [][]byte {
	//defer timeTrack(time.Now(), "GenerateMoves")
	skipKeyCount1 := 0
	skipKeyCount2 := 0
	skipKeyCount3 := 0
	skipKeyCount4 := 0
	lAllStates := make([][]byte, 0, len(sUsedKeySets))
	for ui := 0; ui < len(sUsedKeySets); ui++ {
		usedKeys := sUsedKeySets[ui]

		nextKeys := make([]byte, len(usedKeys)+1, len(usedKeys)+1)
		copy(nextKeys, usedKeys)
		for ki := 0; ki < len(config.AllBitKeys); ki++ {
			keyIndex := config.AllBitKeys[ki]

			if !misc.ByteInArray(keyIndex, usedKeys) {
				nextKeys[len(usedKeys)] = keyIndex
				aksLen := len(nextKeys)
				state := nextKeys
				mod := config.AllKeys[keyIndex].Mod

				skipKey := false

				if !skipKey {
					m, ok := config.FixedModsGroupsIndexed[actionIndex]
					if ok && mod != m {
						skipKey = true
						skipKeyCount1++
						//if actionIndex == 45 && skipKeyCount1 < 50 {
						//	fmt.Println("skipKeyCount1")
						//	PrintState(state)
						//}
					}
				}
				if !skipKey {
					for gi := 0; gi < len(config.SameKeyGroupsIndexed); gi++ {
						group := config.SameKeyGroupsIndexed[gi]
						prevKey := "nil"
						for ai := 0; ai < len(group); ai++ {
							action := group[ai]
							if action < aksLen {
								key := config.AllKeys[state[action]]
								k := key.Key
								if prevKey == "nil" {
									prevKey = k
								}
								if prevKey != k {
									skipKey = true
									skipKeyCount2++
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
						group := config.DiffKeyGroupsIndexed[gi]
						diffKeys := make(map[string]struct{})
						for ai := 0; ai < len(group); ai++ {
							action := group[ai]
							if action < aksLen {
								key := config.AllKeys[state[action]]
								k := key.Key
								if _, ok := diffKeys[k]; !ok {
									diffKeys[k] = struct{}{}
								} else {
									skipKey = true
									skipKeyCount3++
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
						prevMod := "nil"
						for ai := 0; ai < len(group); ai++ {
							action := group[ai]
							if action < aksLen {
								key := config.AllKeys[state[action]]
								m := key.Mod
								if prevMod == "nil" {
									prevMod = m
								}
								if prevMod != m {
									skipKey = true
									skipKeyCount4++
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
	fmt.Println("len(lAllStates):", len(lAllStates), "skipKeyCount: ", skipKeyCount1, skipKeyCount2, skipKeyCount3, skipKeyCount4)
	return lAllStates
}

func CutoffMoves(lAllStates [][]byte, v []int,
	topScoreCount int, cutOffNum int, toprint bool) [][]byte {
	//defer timeTrack(time.Now(), "CutoffMoves")

	states := make([][]byte, 0, len(lAllStates))

	scoreVals := misc.UniqInt(v)

	sort.Sort(sort.IntSlice(scoreVals))

	topScoreVals := scoreVals[misc.IntMax(0, len(scoreVals)-topScoreCount):]
	for ti := 0; ti < len(topScoreVals); ti++ {
		tS := topScoreVals[ti]
		cutOff := 0
		for si := 0; si < len(lAllStates); si++ {
			state := lAllStates[si]
			score := v[si]
			if score == tS {
				states = append(states, state)
				if toprint {
					fmt.Println("Score: ", score)
					for ai := 0; ai < len(state); ai++ {
						ki := state[ai]
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

func PrintState(state []byte) {
	for ai := 0; ai < len(state); ai++ {
		ki := state[ai]
		actName := config.Actions[ai]
		keymod := config.AllKeys[ki]
		key := keymod.Key
		mod := keymod.Mod
		fmt.Print(actName, "(", mod, "+", key, ") ")
	}
	fmt.Println("")
}

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
	for ai := 0; ai < lALen; ai++ {
		aa := config.Actions[ai]
		sUsedKeySets := prevStates
		fmt.Println(ai, aa, lALen-1, len(sUsedKeySets))

		lAllStates := GenerateMoves(sUsedKeySets, ai)
		if len(lAllStates) == 0 {
			panic("Out of keys!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			break
		}
		v := make([]int, len(lAllStates))
		sendCount := 0
		for si := 0; si < len(lAllStates); si++ {
			state := lAllStates[si]
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
		if ai == lALen-1 {
			prevStates = CutoffMoves(lAllStates, v, 2, 3, true)
			//for _, state := range prevStates {
			//	jobs <- dataJobs{0, state, true}
			//	_ = <-results
			//	fmt.Println("\n---------------------------")
			//}
		} else {
			//prevStates = CutoffMoves(lAllStates, v, lALen-ai/5, 10000+ai*20, false)
			prevStates = CutoffMoves(lAllStates, v, lALen-ai/3, 50+ai*2, false)
			//prevStates = CutoffMoves(lAllStates, v, 2, 3, true)
		}
		fmt.Println("len(prevStates) === ", len(prevStates))
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

	bestState := prevStates[len(prevStates)-1]
	for si := 0; si < len(config.SeqOfActionsIndexed); si++ {
		seq := config.SeqOfActionsIndexed[si]
		seqActions := make([]string, 0, len(seq))
		seqKeys := make([]config.Key, 0, len(seq))
		for ai := 0; ai < len(seq); ai++ {
			action := seq[ai]
			actName := config.Actions[action]
			key := config.AllKeys[bestState[action]]
			//f := config.KeyFinger[key]
			seqActions = append(seqActions, actName)
			seqKeys = append(seqKeys, key)
		}
		fmt.Println(seqActions)
		fmt.Println(seqKeys)

	}

	for _, m := range config.AllMods {
		keyAction := make(map[string]string)
		for ai, ki := range bestState {
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
