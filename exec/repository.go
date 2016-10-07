package exec

import (
	"regexp"
)

// A repository is a set of test groups with tests.
type Repository struct {
	testScenarios []*repositoryEntry
}

type repositoryEntry struct {
	concurrency  int
	testgroup    string
	tags         []string
	testScenario *TestScenario
}

// TestFactory is a factory method which returns a test with its data.
type TestFactory func(Context) (Exec, chan Context)

func NewRepository() *Repository {
	return &Repository{
		testScenarios: make([]*repositoryEntry, 0, 0),
	}
}

func (repo *Repository) Add(szenario *TestScenario, testgroup string, concurrency int, tags ...string) {
	repo.testScenarios = append(repo.testScenarios,
		&repositoryEntry{
			testgroup:    testgroup,
			concurrency:  concurrency,
			tags:         tags,
			testScenario: szenario,
		})
}

// Run all testScenarios, which match the supplied filter criterias.
func (repo *Repository) RunTestScenarios(testgroupRegex string, nameRegex string, tagPatterns ...string) map[*repositoryEntry][]*Execution {
	scenarioExecutions := make(map[*repositoryEntry][]*Execution)
	for _, t := range repo.testScenarios {
		if matched, err := regexp.MatchString(nameRegex, t.testScenario.Name); err == nil && matched {
			if matched, err := regexp.MatchString(testgroupRegex, t.testgroup); err == nil && matched {
				if allTagsContained(t.tags, tagPatterns) {
					executions := runTestScenario(t)
					scenarioExecutions[t] = executions
				}
			}
		}
	}
	return scenarioExecutions
}

func runTestScenario(t *repositoryEntry) []*Execution {
	executions := []*Execution{}
	results := RunParallel(t.concurrency, t.testScenario.Exec, t.testScenario.ContextChannelFactory())
	for result := range results {
		executions = append(executions, result)
		println(result.String())
	}
	return executions
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
