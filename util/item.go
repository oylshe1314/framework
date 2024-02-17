package util

func SplitItemNums(value1, value2 string, defaultNum uint32) (items []uint32, nums []uint32, err error) {
	if len(value2) > 0 {
		err = SplitToIntegers2(value1, ",", &items)
		if err != nil {
			return nil, nil, err
		}

		err = SplitToIntegers2(value2, ",", &nums)
		if err != nil {
			return nil, nil, err
		}

		if defaultNum > 0 {
			for i := len(nums); i < len(items); i++ {
				nums = append(nums, defaultNum)
			}
		} else {
			if len(items) > len(nums) {
				items = items[:len(nums)]
			}
		}
	} else {
		if defaultNum == 0 {
			return nil, nil, nil
		}

		err = SplitToIntegers2(value1, ",", &items)
		if err != nil {
			return nil, nil, err
		}

		nums = make([]uint32, len(items))
		for i := range items {
			nums[i] = defaultNum
		}
	}
	return
}

func SplitItemPairs(value1, value2 string, defaultNum uint32) ([][2]uint32, error) {
	var ids, nums, err = SplitItemNums(value1, value2, defaultNum)
	if err != nil {
		return nil, err
	}

	var pairs = make([][2]uint32, len(ids))
	for i := range ids {
		pairs[i][0] = ids[i]
		pairs[i][1] = nums[i]
	}
	return pairs, nil
}
