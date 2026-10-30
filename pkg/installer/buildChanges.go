package installer

// BuildChangesFields is the fields parameter of GET /app/rest/builds/id:<id> that BuildChanges is decoded from: the build's
// revisions and its changes, each with its VCS root instance.
const BuildChangesFields = "revisions(revision(vcs-root-instance(id))),changes($locator(count:10000),change(id,version,vcsRootInstance(id)))"

type vcsRootInstance struct {
	Id string `json:"id"`
}

// Change is a VCS change attributed to a build: Id is TeamCity's internal change id (used as a boundary for the
// sinceChange locator) and Version is the VCS revision (commit SHA).
type Change struct {
	Id              int64           `json:"id"`
	Version         string          `json:"version"`
	VcsRootInstance vcsRootInstance `json:"vcsRootInstance"`
}

type BuildChanges struct {
	Revisions struct {
		Revision []struct {
			VcsRootInstance vcsRootInstance `json:"vcs-root-instance"`
		} `json:"revision"`
	} `json:"revisions"`
	Changes struct {
		Change []Change `json:"change"`
	} `json:"changes"`
}

// CheckedOutChanges returns the changes of the build, newest first, from the VCS roots the build is built from.
//
// TeamCity also attributes to a build the commits of the project's versioned settings repository (e.g.
// intellij-teamcity-config) when settings changes are shown in builds. That repository is not a revision of the build,
// and its commits come first in the changes, so they would become the last commit of a range that is bisected in the
// product repository, where they don't exist. A build without revisions keeps all its changes.
func (b *BuildChanges) CheckedOutChanges() []Change {
	roots := make(map[string]struct{}, len(b.Revisions.Revision))
	for _, revision := range b.Revisions.Revision {
		roots[revision.VcsRootInstance.Id] = struct{}{}
	}
	changes := make([]Change, 0, len(b.Changes.Change))
	for _, change := range b.Changes.Change {
		if _, ok := roots[change.VcsRootInstance.Id]; ok || len(roots) == 0 {
			changes = append(changes, change)
		}
	}
	return changes
}

// CheckedOutCommits returns the commits of CheckedOutChanges.
func (b *BuildChanges) CheckedOutCommits() []string {
	changes := b.CheckedOutChanges()
	commits := make([]string, len(changes))
	for i, change := range changes {
		commits[i] = change.Version
	}
	return commits
}
