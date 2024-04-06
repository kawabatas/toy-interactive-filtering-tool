package toyinteractivefilteringtool

import (
	"bufio"
	"io"
	"sync"
)

// Ctx contains all the important data. while you can easily access
// data in this struct from anwyehre, only do so via channels
type Ctx struct {
	result       string
	loopCh       chan struct{}
	queryCh      chan string
	drawCh       chan []Match
	pagingCh     chan PagingRequest
	mutex        sync.Mutex
	query        []rune
	selectedLine int
	lines        []Match
	current      []Match

	wait *sync.WaitGroup
}

type Match struct {
	line    string
	matches [][]int
}

func NewCtx() *Ctx {
	return &Ctx{
		"",
		make(chan struct{}),      // loopCh
		make(chan string),        // queryCh
		make(chan []Match),       // drawCh
		make(chan PagingRequest), // pagingCh
		sync.Mutex{},
		[]rune{},
		1,
		[]Match{},
		nil,
		&sync.WaitGroup{},
	}
}

func (c *Ctx) Result() string {
	return c.result
}

func (c *Ctx) AddWaitGroup() {
	c.wait.Add(1)
}

func (c *Ctx) ReleaseWaitGroup() {
	c.wait.Done()
}

func (c *Ctx) WaitDone() {
	c.wait.Wait()
}

func (c *Ctx) LoopCh() chan struct{} {
	return c.loopCh
}

func (c *Ctx) QueryCh() chan string {
	return c.queryCh
}

func (c *Ctx) DrawCh() chan []Match {
	return c.drawCh
}

func (c *Ctx) PagingCh() chan PagingRequest {
	return c.pagingCh
}

func (c *Ctx) Terminate() {
	close(c.loopCh)
}

func (c *Ctx) ExecQuery(v string) {
	c.queryCh <- v
}

func (c *Ctx) DrawMatches(m []Match) {
	c.drawCh <- m
}

func (c *Ctx) Refresh() {
	c.DrawMatches(nil)
}

func (c *Ctx) Buffer() []Match {
	// Copy lines so it's safe to read it
	lcopy := make([]Match, len(c.lines))
	copy(lcopy, c.lines)
	return lcopy
}

func (c *Ctx) ReadBuffer(input io.Reader) {
	rdr := bufio.NewReader(input)
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			break
		}
		c.lines = append(c.lines, Match{line, nil})
	}
}

func (c *Ctx) NewView() *View {
	return &View{c}
}

func (c *Ctx) NewFilter() *Filter {
	return &Filter{c}
}

func (c *Ctx) NewInput() *Input {
	return &Input{c}
}
