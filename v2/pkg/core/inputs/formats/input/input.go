package input

import (
	"github.com/pkg/errors"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/pkg/core/inputs/formats"
	"github.com/projectdiscovery/nuclei/v2/pkg/core/inputs/formats/json"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/contextargs"
)

// InputProvider is an interface implemented by format nuclei input provider
type InputProvider struct {
	format    formats.Format
	inputFile string
	count     int64
}

// NewInputProvider creates a new input provider
//
// TODO: Currently the provider does not cache results. If we see this
// parsing everytime being slow we can maybe switch, but this uses less memory
// so it has been implemented this way.
func NewInputProvider(inputFile, inputMode string) (*InputProvider, error) {
	var format formats.Format
	switch inputMode {
	case "jsonl":
		format = json.New()
	default:
		return nil, errors.Errorf("invalid input mode %s", inputMode)
	}

	// Do a first pass over the input to identify any errors
	// and get the count of the input file as well
	count := int64(0)
	parseErr := format.Parse(inputFile, func(request *formats.RawRequest) bool {
		count++
		return false
	})
	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "could not parse input file")
	}
	return &InputProvider{format: format, inputFile: inputFile, count: count}, nil
}

// Count returns the number of items for input provider
func (i *InputProvider) Count() int64 {
	return i.count
}

// Scan iterates the input and each found item is passed to the
// callback consumer.
func (i *InputProvider) Scan(callback func(value *contextargs.MetaInput) bool) {
	err := i.format.Parse(i.inputFile, func(request *formats.RawRequest) bool {
		return callback(&contextargs.MetaInput{
			RawRequest: request,
		})
	})
	if err != nil {
		gologger.Warning().Msgf("Could not parse input file: %s\n", err)
	}
}

// Set adds item to input provider
func (i *InputProvider) Set(value string) {}
