package main

import (
    "fmt"
    "math/rand"
    "math"
    "os"
    "time"
)

type tSeqInfo struct {
    T, W, C float64
}

type tRnd struct {
    rnd []int
    index int
}

const SWAP          = 0
const REINSERTION   = 1
const OR_OPT_2      = 2
const OR_OPT_3      = 3
const TWO_OPT       = 4

var dimension int
var cost [][]float64

func ternary(s bool, t int, f int) int {if s {return t} else {return f}}

func subseq_load(s []int, seq [][]tSeqInfo) float64 {

    for i := 0; i < dimension+1; i++ {
        k := 1 - i - ternary(i == 0, 1, 0)

        seq[i][i].T = 0.0
        seq[i][i].C = 0.0
        seq[i][i].W = float64(ternary(i != 0, 1, 0))
        for j := i+1; j < dimension+1; j++ {
            var j_prev = j-1
            seq[i][j].T = cost[s[j_prev]][s[j]] + seq[i][j_prev].T
            seq[i][j].C = seq[i][j].T + seq[i][j_prev].C
            seq[i][j].W = float64(j + k)
        }
    }

    return seq[0][dimension].C
}

func remove(s []int, i int) []int {
    return append(s[:i], s[i+1:]...)
}

func swap(s []int, i int, j int) {
    var tmp = s[i]
    s[i] = s[j]
    s[j] = tmp
}

func reverse(s []int, i int, j int) {
    m := int((i+j) / 2)

    for first, last := i, j; first <= m; first,last = first+1, last-1  {
        swap(s, first, last)
    }
}

func shift(s []int, from int, to int, sz int) {
    if (from < to) {
        for i, j := from+sz-1, to+sz-1; i >= from; i,j = i-1, j-1 {
            s[j] = s[i]
        }
    } else {
        for i, j := from, to; i < from+sz; i, j = i+1, j+1 {
            s[j] = s[i]
        }
    }
}

func cpy(from []int, to []int, sz int) {
    for i := 0; i < sz; i++ {
        to[i] = from[i];
    }
}

func reinsert(s []int, i int, j int, pos int) {
    seq  := make([]int , j-i+1)
    copy(seq, s[i:j+1])
    if pos < i {
        sz := i-pos
        shift(s, pos, j+1-sz, sz)
        cpy(seq, s[pos:], j-i+1)
    } else {
        sz := pos-j-1
        shift(s, j+1, i, sz)
        cpy(seq, s[i+sz:], j-i+1)
    }
}

func sort(arr []int, r int) {
    for i := 0; i < len(arr); i++ {
        for j := 0; j < len(arr)-i-1; j++ {
            if cost[r][arr[j]] > cost[r][arr[j+1]] {
                swap(arr, j, j+1)
            }
        }
    }
}

func construct(alpha float64, s_crnt []int, rnd * tRnd) {
    s := make([]int, 0)
    s = append(s, 0)

    cL := make([]int, dimension-1)
    for i := 1; i < dimension; i++ {
        cL[i-1] = i
    }

    var r = 0
    for len(cL) > 0 {
        sort(cL, r)

        var rang =  int(float64(len(cL)) * alpha + 1.0)
        var index = rand.Intn(rang)
        r_index := rnd.index; rnd.index++
        index = rnd.rnd[r_index]
        var c = cL[index]
        r = c
        //fmt.Println(cL)
        cL = remove(cL, index)
        //fmt.Println(index, cL)
        s = append(s, c)
    }

    s = append(s, 0)

    copy(s_crnt, s)
}

func fea(s []int) bool {
    check := make([]bool, dimension)
    for i, v := range s {
        _ = i
        check[v] = true
    }

    for i, v := range check {
        _ = i
        if !v {
            return false
        }
    }

    return true
}

func search_swap(s []int, seq [][]tSeqInfo) bool {
    var cost_new, cost_concat_1, cost_concat_2, cost_concat_3, cost_concat_4 float64
    var cost_best float64 = math.MaxFloat64
    var j_prev, j_next, i_prev, i_next int
    var I int
    var J int

    for i := 1; i < dimension-1; i++ {
        i_prev = i - 1
        i_next = i + 1

        //consecutive nodes
        cost_concat_1 =                 seq[0][i_prev].T + cost[s[i_prev]][s[i_next]]
        cost_concat_2 = cost_concat_1 + seq[i][i_next].T  + cost[s[i]][s[i_next+1]]

        cost_new = seq[0][i_prev].C                                                    +           //       1st subseq
        seq[i][i_next].W               * (cost_concat_1) + cost[s[i_next]][s[i]]  +           // concat 2nd subseq
        seq[i_next+1][dimension].W   * (cost_concat_2) + seq[i_next+1][dimension].C   // concat 3rd subseq

        if cost_new < cost_best {
            cost_best = cost_new - math.SmallestNonzeroFloat64
            I = i
            J = i_next
        }

        for j := i_next+1; j < dimension; j++ {
            j_next = j + 1
            j_prev = j - 1

            cost_concat_1 =                 seq[0][i_prev].T       + cost[s[i_prev]][s[j]]
            cost_concat_2 = cost_concat_1                           + cost[s[j]][s[i_next]]
            cost_concat_3 = cost_concat_2 + seq[i_next][j_prev].T  + cost[s[j_prev]][s[i]]
            cost_concat_4 = cost_concat_3                           + cost[s[i]][s[j_next]]

            cost_new = seq[0][i_prev].C                                                 +      // 1st subseq
            cost_concat_1 +                                                             // concat 2nd subseq (single node)
            seq[i_next][j_prev].W      * cost_concat_2 + seq[i_next][j_prev].C +      // concat 3rd subseq
            cost_concat_3 +                                                             // concat 4th subseq (single node)
            seq[j_next][dimension].W * cost_concat_4 + seq[j_next][dimension].C   // concat 5th subseq

            if cost_new < cost_best {
                cost_best = cost_new - math.SmallestNonzeroFloat64
                I = i
                J = j
            }
        }
    }

    if cost_best < seq[0][dimension].C -math.SmallestNonzeroFloat64 {
        swap(s, I, J);
        subseq_load(s, seq);
        //subseq_load_b(s, seq, I);
        return true;
    }

    return false;
}

func search_two_opt(s []int, seq [][]tSeqInfo) bool {
    var cost_new, cost_concat_1, cost_concat_2 float64
    var cost_best = math.MaxFloat64 
    var rev_seq_cost float64
    var i_prev, j_next int
    var I int
    var J int

    for i := 1; i < dimension-1; i++ {
        i_prev = i - 1

        rev_seq_cost = seq[i][i+1].T
        for j := i + 2; j < dimension; j++ {
            j_next = j + 1


            rev_seq_cost += cost[s[j-1]][s[j]] * (seq[i][j].W-1.0)

            cost_concat_1 =                 seq[0][i_prev].T   + cost[s[j]][s[i_prev]]
            cost_concat_2 = cost_concat_1 + seq[i][j].T        + cost[s[j_next]][s[i]]

            cost_new = seq[0][i_prev].C                                                        +   //  1st subseq
            seq[i][j].W                * cost_concat_1 + rev_seq_cost                  +   // concat 2nd subseq (reversed seq)
            seq[j_next][dimension].W * cost_concat_2 + seq[j_next][dimension].C      // concat 3rd subseq

            if (cost_new < cost_best) {
                cost_best = cost_new - math.SmallestNonzeroFloat64
                I = i
                J = j
            }
        }
    }

    if cost_best < seq[0][dimension].C - math.SmallestNonzeroFloat64 {
        reverse(s, I, J)
        subseq_load(s, seq)
        return true
    }

    return false
}

func search_reinsertion(s []int, seq [][]tSeqInfo, opt int) bool {
    var cost_new, cost_concat_1, cost_concat_2, cost_concat_3 float64
    var cost_best = math.MaxFloat64
    var k_next, i_prev, j_next int
    var I int
    var J int
    var POS int

    for i:= 1; i < dimension-opt+1; i = i+1 {
        j := opt+i-1
        j_next = j + 1
        i_prev = i - 1

        //k -> edges 
        for k := 0; k < i_prev; k++ {
            k_next = k+1

            cost_concat_1 =                 seq[0][k].T            + cost[s[k]][s[i]]
            cost_concat_2 = cost_concat_1 + seq[i][j].T            + cost[s[j]][s[k_next]]
            cost_concat_3 = cost_concat_2 + seq[k_next][i_prev].T  + cost[s[i_prev]][s[j_next]]

            cost_new = seq[0][k].C                                                                   +   //       1st subseq
            seq[i][j].W               * cost_concat_1 + seq[i][j].C                  +   //  concat 2nd subseq (reinserted seq)
            seq[k_next][i_prev].W     * cost_concat_2 + seq[k_next][i_prev].C        +   //  concat 3rd subseq
            seq[j_next][dimension].W * cost_concat_3 + seq[j_next][dimension].C       // concat 4th subseq

            if cost_new < cost_best {
                cost_best = cost_new - math.SmallestNonzeroFloat64
                I = i
                J = j
                POS = k
            }
        }

        for k := i + opt; k < dimension; k++ {
            k_next = k + 1

            cost_concat_1 =                 seq[0][i_prev].T  + cost[s[i_prev]][s[j_next]]
            cost_concat_2 = cost_concat_1 + seq[j_next][k].T  + cost[s[k]][s[i]]
            cost_concat_3 = cost_concat_2 + seq[i][j].T       + cost[s[j]][s[k_next]]

            cost_new = seq[0][i_prev].C                                                                  +   //       1st subseq
            seq[j_next][k].W          * cost_concat_1 + seq[j_next][k].C             +   // concat 2nd subseq
            seq[i][j].W               * cost_concat_2 + seq[i][j].C                  +   // concat 3rd subseq (reinserted seq)
            seq[k_next][dimension].W * cost_concat_3 + seq[k_next][dimension].C       // concat 4th subseq

            if cost_new < cost_best {
                cost_best = cost_new - math.SmallestNonzeroFloat64;
                I = i
                J = j
                POS = k
            }
        }
    }

    if cost_best < seq[0][dimension].C - math.SmallestNonzeroFloat64 {
      //fmt.Println(cost_best, seq[0][dimension].C)
      //fmt.Println(s, I, J, POS+1)
        reinsert(s, I, J, POS+1)
      //fmt.Println(s, I, J, POS+1)
        subseq_load(s, seq)
        //subseq_load_b(s, seq, I < POS+1 ? I : POS+1);
        return true;
    }

    return false;
}

func RVND(s []int , seq [][]tSeqInfo, rnd *tRnd) {
    var neighbd_list []int
    _ = neighbd_list
    neighbd_list = []int{SWAP, TWO_OPT, REINSERTION, OR_OPT_2, OR_OPT_3}

    for len(neighbd_list) > 0 {
        r_index := rnd.index; rnd.index++


        improve := false
        index := rnd.rnd[r_index]

        switch neighbd := neighbd_list[index]; neighbd {
        case SWAP:
            improve = search_swap(s, seq)
        case REINSERTION:
            improve = search_reinsertion(s, seq, REINSERTION)
        case OR_OPT_2:
            improve = search_reinsertion(s, seq, OR_OPT_2)
        case OR_OPT_3:
            improve = search_reinsertion(s, seq, OR_OPT_3)
        case TWO_OPT:
            improve = search_two_opt(s, seq)
        }

        if !fea(s) {
            fmt.Println("qebrad")
            os.Exit(0)
        }
        //fmt.Println(index, seq[0][dimension].C)

        if improve {
            neighbd_list = []int{SWAP, TWO_OPT, REINSERTION, OR_OPT_2, OR_OPT_3}
        } else {
            neighbd_list = remove(neighbd_list, index)
        }

    }
    
}

func perturb(s_crnt []int, s_partial []int, rnd *tRnd) {
    s := make([]int, len(s_partial))

    copy(s, s_partial)

    var A_start = 1
    var A_end = 1
    var B_start = 1
    var B_end = 1

    var size_max = (dimension+1)/10
    size_max = ternary(size_max >= 2, size_max, 2)
    var size_min = 2
    //std::cout << "perturbing\n";
    //print_s(s);
    for (A_start <= B_start && B_start <= A_end) || (B_start <= A_start && A_start <= B_end) {
        /**/
        max := (dimension+1) -2 -size_max
        A_start = rand.Intn(max + 1)
        A_end = A_start + rand.Intn(size_max - size_min + 1) + size_min

        B_start = rand.Intn(max + 1)
        B_end = B_start + rand.Intn(size_max - size_min + 1) + size_min
        /**/



        //std::cout << "paa\n";

        //cout << info.rnd[info.rnd_index] << endl;
        r_index := rnd.index; rnd.index++
        A_start = rnd.rnd[r_index]
        //cout << info.rnd[info.rnd_index] << endl;
        r_index = rnd.index; rnd.index++
        A_end = A_start + rnd.rnd[r_index]
        //std::cout << "A start  " << A_start << std::endl;
        //std::cout << "A end  " << A_end << std::endl;

        //cout << info.rnd[info.rnd_index] << endl;
        r_index = rnd.index; rnd.index++
        B_start = rnd.rnd[r_index]
        //cout << info.rnd[info.rnd_index] << endl;
        r_index = rnd.index; rnd.index++
        B_end = B_start + rnd.rnd[r_index]
        //std::cout << "B start  " << B_start << std::endl;
        //std::cout << "B end  " << B_end << std::endl;
    }
    
    //cout << "A_end  " << A_end << endl << "B_end  " << B_end << endl;

    if A_start < B_start {
        reinsert(s, B_start, B_end-1, A_end)
        reinsert(s, A_start, A_end-1, B_end)
    } else {
        reinsert(s, A_start, A_end-1, B_end)
        reinsert(s, B_start, B_end-1, A_end)
    }

    //print_s(s);
    //subseq_load(solut, info);

    copy(s_crnt, s)
    //memcpy(s_crnt, s, sizeof(int)*(dimen+1));
}

//func GILS_RVND(Imax int, Iils int, R []float64) {
func GILS_RVND(rnd tRnd) {
    Imax := 10
    Iils := ternary(dimension < 100, dimension, 100)
    _ = Iils
    R := []float64{0.00, 0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17, 0.18, 0.19, 0.2, 0.21, 0.22, 0.23, 0.24, 0.25}

    s_best := make([]int, dimension+1)
    s_crnt := make([]int, dimension+1)
    s_partial := make([]int, dimension+1)

    var cost_best float64 = math.MaxFloat64
    var cost_crnt float64
    var cost_partial float64

    _, _, _ = s_best, s_crnt, s_partial
    _, _, _ = cost_best, cost_crnt, cost_partial

    seq := make([][]tSeqInfo, dimension+1)
    for i := 0; i < dimension+1; i++ {
        seq[i] = make([]tSeqInfo, dimension+1)
    }

    for i := 0; i < Imax; i++ {
        var alpha = R[rand.Intn(len(R))]
        var r_index = rnd.index; rnd.index++
        var index = rnd.rnd[r_index]
        alpha = R[index]
        fmt.Printf("[+] Search %d\n", i+1)
        fmt.Printf("\t[+] Constructing..\n");	

        construct(alpha, s_crnt, &rnd)

        cost_crnt = subseq_load(s_crnt, seq)
        //fmt.Println(s_crnt, cost_crnt)
        copy(s_partial, s_crnt)
        cost_partial = cost_crnt

        var iterILS = 0
        for iterILS < Iils {
            RVND(s_crnt, seq, &rnd)
            cost_crnt = seq[0][dimension].C

            if cost_crnt < cost_partial {
                cost_partial = cost_crnt
                copy(s_partial, s_crnt)
                iterILS = 0
            }

            perturb(s_crnt, s_partial, &rnd)
            subseq_load(s_crnt, seq)
            iterILS++
        }

        if cost_partial < cost_best {
            copy(s_best, s_partial)
            cost_best = cost_partial
        }

        fmt.Println("Current best cost: ", cost_best)
        fmt.Println(s_best)
    }

    fmt.Println("COST: ", cost_best)
    fmt.Println("SOLUCAO: ", s_best)
}

func main() {
    var rnd tRnd
    dimension, cost, rnd.rnd = loadData()
    _ = dimension
    _ = cost
    _ = rnd


    start := time.Now()
    GILS_RVND(rnd)
    elapsed := time.Since(start).Seconds()

    fmt.Println("TIME: ", elapsed)
}