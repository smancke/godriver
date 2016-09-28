package exec

import (
	"regexp"
)

// A repository is a set of test groups with tests.
type Repository struct {
	tests []*test
}

// TestFactory is a factory method which returns a test with its data.
type TestFactory func(Context) (Exec, chan Context)

type test struct {
	name      string
	testgroup string
	tags      []string
	factory   TestFactory
	cntx      Context
}

func NewRepository() *Repository {
	return &Repository{
		tests: make([]*test, 0, 0),
	}
}

func (repo *Repository) AddTest(testgroup string, name string, factory TestFactory, cntx Context, tags ...string) {
	repo.tests = append(repo.tests, &test{
		name:      name,
		testgroup: testgroup,
		factory:   factory,
		tags:      tags,
		cntx:      cntx,
	})
}

// Run all tests, which match the supplied filter criterias.
func (repo *Repository) RunTests(testgroupRegex string, nameRegex string, tagPatterns ...string) {
	for _, t := range repo.tests {
		if matched, err := regexp.MatchString(nameRegex, t.name); err == nil && matched {
			if matched, err := regexp.MatchString(testgroupRegex, t.testgroup); err == nil && matched {
				if allTagsContained(t.tags, tagPatterns) {
					runTest(t)
				}
			}
		}
	}
}

func runTest(t *test) {
	test, data := t.factory(t.cntx)
	results := RunParallel(t.cntx.ExecutionConcurrency(), test, data)

	for result := range results {
		println(result.String())
	}
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
