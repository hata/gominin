package gominin


type BinarySearchList interface {
    Len() int
}

type CompareFunc func(index int) int


func BinarySearch(list BinarySearchList, compareFunc CompareFunc) int {
    start, end, pivot := 0, list.Len() - 1, 0

    for start <= end {
        pivot = (start + end) / 2

        if compareFunc(pivot) < 0 {
            start = pivot + 1
        } else if compareFunc(pivot) > 0 {
            end = pivot - 1
        } else {
            return pivot
        }
    }

    return -1
}

func LowerBound(list BinarySearchList, compareFunc CompareFunc) int {
    start, end, pivot := 0, list.Len() - 1, 0

    if compareFunc(end) < 0 {
        return end + 1
    }
    if compareFunc(start) >= 0 {
        return start
    }

    for start <= end {
        pivot = (start + end) / 2

        if compareFunc(pivot) < 0 {
            start = pivot + 1
        } else if compareFunc(pivot) > 0 {
            end = pivot - 1
        } else {
            return pivot
        }
    }

    if end == pivot - 1 && pivot >= 0 {
        return pivot
    } else if start == pivot + 1 && start < list.Len() {
        return start
    }

    return -1
}

