package filelist

import (
	"fmt"
	"hash"
	"time"
)

type Pair struct {
	calculated bool
	FileOne    File
	FileTwo    File
}

type HashCalculationError struct {
	err  error
	pair *Pair
}

func NewPair(f1 File, f2 File) Pair {
	return Pair{
		FileOne: f1,
		FileTwo: f2,
	}
}

func (p *Pair) Calculate(h func() hash.Hash) error {
	err := p.FileOne.calculate(h())
	if err != nil {
		return err
	}

	err = p.FileTwo.calculate(h())
	if err != nil {
		return err
	}

	p.calculated = true

	return nil
}

func (p *Pair) Compare() bool {
	if !p.calculated {
		panic("file hashes not calculated")
	}

	return p.FileOne.FileHash == p.FileTwo.FileHash
}

func (p *Pair) CalculateAndCompare(h func() hash.Hash) (bool, error) {
	err := p.Calculate(h)
	if err != nil {
		return false, err
	}

	return p.Compare(), nil
}

func CalculateAllPairs(ps *[]Pair, h func() hash.Hash, threads int) {
	if threads <= 0 {
		threads = 1
	}

	ipt := make(chan *Pair, threads)
	errs := make(chan HashCalculationError, threads)

	//Start error checker
	go func() {
		for {
			hce, ok := <-errs

			if !ok {
				return
			}

			fmt.Println("Error calculating hash for: ", hce.pair.FileOne.GetSubPath())
			fmt.Println("\t", hce.err)
		}
	}()

	//Start worker threads
	for i := 0; i < threads; i++ {
		go calculateWorker(ipt, errs, h)
	}

	//Send pairs to workers
	for i := range *ps {
		ipt <- &(*ps)[i]
	}

	//Wait for pairs to finish processing before closing
	if len(ipt) != 0 {
		waitForIptClose(ipt, 3*time.Second)
	}

	//Once done, close ipt channel to indicate to workers that we're done
	close(ipt)

	//Wait for errors to be processed before closing error channel
	if len(errs) != 0 {
		waitForErrsClose(errs, 1*time.Second)
	}

	close(errs)
}

func calculateWorker(in chan *Pair, errs chan HashCalculationError, h func() hash.Hash) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic occurred inside worker - worker dead.")
			fmt.Println("\t", r)
			return
		}
	}()

	for {
		p, ok := <-in

		//Check if the channel is closed
		if !ok {
			return
		}

		err := p.Calculate(h)
		if err != nil {
			errs <- HashCalculationError{
				err:  err,
				pair: p,
			}
		}
	}
}

func waitForErrsClose(c chan HashCalculationError, timeout time.Duration) {
	t := time.NewTimer(timeout)

	for {
		if len(c) == 0 {
			return
		}

		select {
		case <-t.C:
			fmt.Println("Timeout waiting for error channel to finish processing")
			return
		default:
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func waitForIptClose(c chan *Pair, timeout time.Duration) {
	t := time.NewTimer(timeout)

	for {
		if len(c) == 0 {
			return
		}

		select {
		case <-t.C:
			fmt.Println("Timeout waiting for input channel to finish processing")
			return
		default:
			time.Sleep(25 * time.Millisecond)
		}
	}
}
