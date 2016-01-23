package app

import (
	"bytes"
	"fmt"

	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"

	imageapi "github.com/openshift/origin/pkg/image/api"
)

// ErrNoMatch is the error returned by new-app when no match is found for a
// given component.
type ErrNoMatch struct {
	Value     string
	Qualifier string
	Errs      []error
}

func (e ErrNoMatch) Error() string {
	if len(e.Qualifier) != 0 {
		return fmt.Sprintf("no match for %q: %s", e.Value, e.Qualifier)
	}
	return fmt.Sprintf("no match for %q", e.Value)
}

// UsageError is the usage error message returned when no match is found.
func (e ErrNoMatch) Suggestion(commandName string) string {
	return fmt.Sprintf("%[3]s - does a Docker image with that name exist?", e.Value, commandName, e.Error())
}

// ErrPartialMatch is the error returned to new-app users when the
// best match available is only a partial match for a given component.
type ErrPartialMatch struct {
	Value string
	Match *ComponentMatch
	Errs  []error
}

func (e ErrPartialMatch) Error() string {
	return fmt.Sprintf("only a partial match was found for %q: %q", e.Value, e.Match.Name)
}

// UsageError is the usage error message returned when only a partial match is
// found.
func (e ErrPartialMatch) Suggestion(commandName string) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "* %s\n", e.Match.Description)
	fmt.Fprintf(buf, "  Use %[1]s to specify this image or template\n\n", e.Match.Argument)

	return fmt.Sprintf(`%[3]s
The argument %[1]q only partially matched the following Docker image or OpenShift image stream:

%[2]s
`, e.Value, buf.String(), cmdutil.MultipleErrors("error: ", e.Errs))
}

// ErrMultipleMatches is the error returned to new-app users when multiple
// matches are found for a given component.
type ErrMultipleMatches struct {
	Value   string
	Matches []*ComponentMatch
	Errs    []error
}

func (e ErrMultipleMatches) Error() string {
	return fmt.Sprintf("multiple images or templates matched %q: %d", e.Value, len(e.Matches))
}

// ErrNameRequired is the error returned by new-app when a name cannot be
// suggested and the user needs to provide one explicitly.
var ErrNameRequired = fmt.Errorf("you must specify a name for your app")

// CircularOutputReferenceError is the error returned by new-app when the input
// and output image stream tags are identical.
type CircularOutputReferenceError struct {
	Reference imageapi.DockerImageReference
}

func (e CircularOutputReferenceError) Error() string {
	return fmt.Sprintf("the input and output image stream tags are identical (%q)", e.Reference.DockerClientDefaults())
}
