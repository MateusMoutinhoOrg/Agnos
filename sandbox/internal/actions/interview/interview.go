package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// Interview runs the interactive session over the project at path: it asks
// what is to be done, generates the questions from the declaration of the
// command that answers it, and runs that command with what it was told.
//
// Like Verify it opens a SmartIO and never calls io.Persist, and unlike every
// other action it runs no follow-up build. It writes nothing of its own: the
// SmartIO is open so the questions can offer what the project already holds —
// the commands that exist, the deps installed, the themes declared — and every
// command it dispatches runs its own action, which persists and builds for
// itself. Because List* reads disk, one SmartIO held for the whole session
// keeps answering with what the last command left behind.
func Interview(sandbox *api.Sandbox, path string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	return InterviewInternal(sandbox, io, path)
}
