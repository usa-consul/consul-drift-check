package diff

import "sort"

// BudgetOptions controls how drift budget is evaluated.
type BudgetOptions struct {
	// MaxDriftCount is the maximum number of drifted keys allowed.
	MaxDriftCount int
	// MaxModified is the maximum number of modified keys allowed.
	MaxModified int
	// MaxMissing is the maximum number of source-only or dest-only keys allowed.
	MaxMissing int
}

// BudgetResult holds per-category spend and whether the budget was exceeded.
type BudgetResult struct {
	TotalDrift  int
	Modified    int
	Missing     int
	Exceeded    bool
	Violations  []string
}

// EvaluateBudget checks whether drift results stay within the supplied budget.
func EvaluateBudget(results []Result, opts BudgetOptions) BudgetResult {
	var modified, missing int
	for _, r := range results {
		switch r.Status {
		case "modified":
			modified++
		case "source-only", "destination-only":
			missing++
		}
	}
	total := modified + missing

	var violations []string
	if opts.MaxDriftCount > 0 && total > opts.MaxDriftCount {
		violations = append(violations, "total drift count exceeded")
	}
	if opts.MaxModified > 0 && modified > opts.MaxModified {
		violations = append(violations, "modified key count exceeded")
	}
	if opts.MaxMissing > 0 && missing > opts.MaxMissing {
		violations = append(violations, "missing key count exceeded")
	}
	sort.Strings(violations)

	return BudgetResult{
		TotalDrift: total,
		Modified:   modified,
		Missing:    missing,
		Exceeded:   len(violations) > 0,
		Violations: violations,
	}
}
