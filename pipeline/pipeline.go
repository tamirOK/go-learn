package pipeline

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// Wrap stage, so it can be cancelled.
func wrapper(stage Stage, in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for v := range stage(in) {
			select {
			case <-done:
				return
			default:
				out <- v
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	elements := make([]interface{}, 0)

	// build pipeline
	for _, stage := range stages {
		in = wrapper(stage, in, done)
	}

	// consume values from pipeline and listen for cancellation
	for v := range in {
		select {
		case <-done:
			return done
		default:
			elements = append(elements, v)
		}
	}
	// convert slice to buffered channel
	results := make(Bi, len(elements))

	for _, elem := range elements {
		results <- elem
	}

	close(results)
	return results
}
