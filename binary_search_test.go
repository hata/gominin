package gominin

import (
    "testing"
)

type intList []int

func (list intList) Len() int {
    return len(list)
}


func TestBinarySearchSingleValue(t *testing.T) {
    target := 0
    list := intList{target}
    x := BinarySearch(list, func(index int) int { return list[index] - target })
    if x != 0 {
        t.Error("Error. BinarySearch should be able to find an item.")
    }
}

func TestBinarySearchTwoValues(t *testing.T) {
    target := 0
    list := intList{0, 1}
    x := BinarySearch(list, func(index int) int { return list[index] - target })
    if x != 0 {
        t.Error("Error. BinarySearch should be able to find an item.")
    }
}

func TestBinarySearchThreeValues(t *testing.T) {
    list := intList{0, 1, 2}
    for target := 0;target < 2;target++ {
        x := BinarySearch(list, func(index int) int { return list[index] - target })
        if x != target {
            t.Error("Error. BinarySearch should be able to find an item.")
        }
    }
}

func TestBinarySearchFiveValues(t *testing.T) {
    list := intList{0, 1, 2, 3, 4}
    for target := 0;target < 5;target++ {
        x := BinarySearch(list, func(index int) int { return list[index] - target })
        if x != target {
            t.Error("Error. BinarySearch should be able to find an item.")
        }
    }
}

func TestBinarySearchNotFound(t *testing.T) {
    list := intList{0, 1, 2, 3, 4}
    x := BinarySearch(list, func(index int) int { return list[index] - 6 })
    if x != -1 {
        t.Error("Error. BinarySearch should return not found(-1)")
    }
}

func TestLowerBound(t *testing.T) {
    target := 1
    list := intList{0, 2, 4}
    x := LowerBound(list, func(index int) int { return list[index] - target })
    if x != 1 {
        t.Error("Error. Index for 2nd item should be return.")
    }
}

func TestLowerBoundExactSame(t *testing.T) {
    target := 2
    list := intList{0, 2, 4}
    x := LowerBound(list, func(index int) int { return list[index] - target })
    if x != 1 {
        t.Error("Error. Index for 2nd item should be return.")
    }
}

func TestLowerBoundUnderMinimum(t *testing.T) {
    target := -1
    list := intList{0, 2, 4}
    x := LowerBound(list, func(index int) int { return list[index] - target })
    if x != 0 {
        t.Error("Error. Index for 1st item should be return.")
    }
}

func TestLowerBoundUnderMaximum(t *testing.T) {
    target := 5
    list := intList{0, 2, 4}
    x := LowerBound(list, func(index int) int { return list[index] - target })
    if x != 3 {
        t.Error("Error. Not found is the same as Len.")
    }
}
