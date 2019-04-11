package concurrent_map

import "sync/atomic"

type BucketStatus uint8

const (
	BUCKET_STATUS_NORMAL      BucketStatus = 0
	BUCKET_STATUS_UNDERWEIGHT BucketStatus = 1
	BUCKET_STATUS_OVERWEIGHT  BucketStatus = 2
)

type PairRedistributor interface {
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	Redistribe(status BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

type myPairRedistributor struct {
	loadFactor            float64
	upperThreshold        uint64
	overweightBucketCount uint64
	emptyBucketCount      uint64
}

func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) PairRedistributor {

}

var bucketCountTemplate = `Bucket count:
	pairTotal: %d
	bucketNumber: %d
	average: %f
	upperThreshold: %d
	emptyBucketCount: %d
`

func (pr *myPairRedistributor) UpdateThreshold(pairTotal uint64, bucketNumber int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))

	if average < 100 {
		average = 100
	}

	// defer func() {
	// 	fmt.Printf(bucketCountTemplate,
	// 		pairTotal,
	// 		bucketNumber,
	// 		average,
	// 		atomic.LoadUint64(&pr.upperThreshold),
	// 		atomic.LoadUint64(&pr.emptyBucketCount))
	// }()

	atomic.StoreUint64(&pr.upperThreshold, uint64(average*pr.loadFactor))
}

var bucketStatusTemplate = `Check bucket status: 
    pairTotal: %d
    bucketSize: %d
    upperThreshold: %d
    overweightBucketCount: %d
    emptyBucketCount: %d
    bucketStatus: %d	

`

func (pr *myPairRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus) {
	// defer func() {
	// 	fmt.Printf(bucketStatusTemplate,
	// 		pairTotal,
	// 		bucketSize,
	// 		atomic.LoadUint64(&pr.upperThreshold),
	// 		atomic.LoadUint64(&pr.overweightBucketCount),
	// 		atomic.LoadUint64(&pr.emptyBucketCount),
	// 		bucketStatus)
	// }()

	if bucketSize > DEFAULT_BUCKET_MAX_SIZE ||
		bucketSize >= atomic.LoadUint64(&pr.upperThreshold) {
		atomic.AddUint64(&pr.overweightBucketCount, 1)
		bucketStatus = BUCKET_STATUS_OVERWEIGHT
		return
	}
	if bucketSize == 0 {
		atomic.AddUint64(&pr.emptyBucketCount, 1)
	}
	return
}

var redistributionTemplate = `Redistributing: 
    bucketStatus: %d
    currentNumber: %d
    newNumber: %d

`

func (pr *myPairRedistributor) Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool) {
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber

	// defer func() {
	// 	fmt.Printf(redistributionTemplate,
	// 		bucketStatus,
	// 		currentNumber,
	// 		newNumber)
	// }()

	switch bucketStatus {
	case BUCKET_STATUS_OVERWEIGHT:
		if atomic.LoadUint64(&pr.overweightBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber << 1
	case BUCKET_STATUS_UNDERWEIGHT:
		if currentNumber < 100 || atomic.LoadUint64(&pr.emptyBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber >> 1
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}
	if newNumber == currentNumber {
		atomic.StoreUint64(&pr.overweightBucketCount, 0)
		atomic.StoreUint64(&pr.emptyBucketCount, 0)
		return nil, false
	}

	var pairs []Pair
	for _, b := range buckets {
		for e := b.GetFirstPair(); e != nil; e = e.Next() {
			pairs = append(pairs, e)
		}
	}
	if newNumber > currentNumber {
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}
		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		buckets = make([]Bucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}
	var count int
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		b.Put(p, nil)
		count++
	}
	atomic.StoreUint64(&pr.overweightBucketCount, 0)
	atomic.StoreUint64(&pr.emptyBucketCount, 0)
	return buckets, true
}
