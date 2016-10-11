package exec

import (
	"regexp"
)

// A repository is a set of test groups with tests.
type Repository struct {
	testScenarios []*repositoryEntry
	runResults    []*repositoryRunResult
}

type repositoryEntry struct {
	concurrency  int
	testGroup    string
	tags         []string
	testScenario *TestScenario
}

type repositoryRunResult struct {
	scenario   *repositoryEntry
	executions []*Execution
}

// TestFactory is a factory method which returns a test with its data.
type TestFactory func(Context) (Exec, chan Context)

func NewRepository() *Repository {
	return &Repository{
		testScenarios: make([]*repositoryEntry, 0, 0),
	}
}

func (repo *Repository) Add(scenario *TestScenario, testGroup string, concurrency int, tags ...string) {
	repo.testScenarios = append(repo.testScenarios,
		&repositoryEntry{
			testGroup:    testGroup,
			concurrency:  concurrency,
			tags:         tags,
			testScenario: scenario,
		})
}

// Run all testScenarios, which match the supplied filter criteria.
func (repo *Repository) RunTestScenarios(testGroupRegex string, nameRegex string, tagPatterns ...string) {
	runResults := make([]*repositoryRunResult, 0, 0)
	for _, t := range repo.testScenarios {
		if matched, err := regexp.MatchString(nameRegex, t.testScenario.Name); err == nil && matched {
			if matched, err := regexp.MatchString(testGroupRegex, t.testGroup); err == nil && matched {
				if allTagsContained(t.tags, tagPatterns) {
					executions := t.runTestScenario()
					runResults = append(runResults, &repositoryRunResult{t, executions})
				}
			}
		}
	}
	repo.runResults = runResults
}

func (repo *Repository) GetErrorExecutions() []*Execution {
	if repo.runResults == nil {
		return nil
	}

	errorExecs := []*Execution{}
	for _, result := range repo.runResults {
		errorExecs = append(result.getErrorExecutions())
	}
	return errorExecs
}

func (t *repositoryEntry) runTestScenario() []*Execution {
	executions := []*Execution{}
	results := RunParallel(t.concurrency, t.testScenario.Exec, t.testScenario.ContextChannelFactory())
	for result := range results {
		executions = append(executions, result)
		println(result.String())
	}
	return executions
}

func (r *repositoryRunResult) getErrorExecutions() []*Execution {
	errorExecs := []*Execution{}
	for _, exec := range r.executions {
		if exec.err != nil {
			errorExecs = append(errorExecs, exec)
		}
	}
	return errorExecs
}

func allTagsContained(tags []string, tagRegex []string) bool {
	for _, tRegex := range tagRegex {
		matchesOneTag := false
		for _, t := range tags {
			if matched, err := regexp.MatchString(tRegex, t); err == nil && matched {
				matchesOneTag = true
			}
		}
		if !matchesOneTag {
			return false
		}
	}
	return true
}
