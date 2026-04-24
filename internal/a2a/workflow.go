package a2a

type Phase struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
}

type Workflow struct {
	Repo   string  `json:"repo"`
	Ticket string  `json:"ticket"`
	Live   bool    `json:"live"`
	Phases []Phase `json:"phases"`
}

func Default(repo, ticket string, live bool) Workflow {
	return Workflow{
		Repo:   repo,
		Ticket: ticket,
		Live:   live,
		Phases: []Phase{
			{Name: "propose", Purpose: "author proposes a change"},
			{Name: "test", Purpose: "author validates the candidate change"},
			{Name: "review", Purpose: "reviewer evaluates and may block or request revision"},
			{Name: "revise", Purpose: "author revises after review"},
			{Name: "merge", Purpose: "arbiter or maintainer merges the approved change"},
			{Name: "done", Purpose: "workflow is complete"},
		},
	}
}
