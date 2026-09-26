package signaldeps

import (
	"os"
	"os/signal"
	"syscall"

	"{{.Module}}/sandbox/deps"
)

// Bind fills deps.Deps.Signaldeps with the standard library's os/signal.
func Bind(deps *deps.Deps) {
	deps.Signaldeps.OnInterrupt = onInterrupt
}

// onInterrupt waits for the first interrupt or termination signal on a
// goroutine of its own and runs handler once. Once it has heard one it stops
// listening, so the next signal reaches the process with its default effect.
func onInterrupt(handler func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		signal.Stop(signals)
		handler()
	}()
}
