package actions

import (
	"testing"
	"os"
	"net/url"

	"github.com/stretchr/testify/assert"
)

const pathToTemplateRepoFile = "../resources/repository_list.json"

func TestGetTemplates(t *testing.T) {
	tests := map[string]struct {
		in         string
		wantedType []Template
		wantedErr  error
	}{
		"success case: input good path": {
			in:         pathToTemplateRepoFile,
			wantedType: []Template{},
			wantedErr:  nil,
		},
		"fail case: input /path/to/nowhere": {
			in:         "/path/to/nowhere",
			wantedType: nil,
			wantedErr:  new(os.PathError),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetTemplates(test.in)
			assert.IsType(t, test.wantedType, got)
			assert.IsType(t, test.wantedErr, err)
		})
	}
}

func TestReadJSONFile(t *testing.T) {
	tests := map[string]struct {
		in   string
		want []TemplateRepo
		err  error
	}{
		"success case: input good path": {
			in: pathToTemplateRepoFile,
			want: []TemplateRepo{
				{
					Description: "Standard Codewind templates.",
					URL:         "https://raw.githubusercontent.com/kabanero-io/codewind-templates/master/devfiles/index.json",
				},
			},
			err: nil,
		},
		"fail case: input /path/to/nowhere": {
			in:   "/path/to/nowhere",
			want: nil,
			err:  new(os.PathError),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ReadJSONFile(test.in)
			assert.Equal(t, test.want, got)
			assert.IsType(t, test.err, err)
		})
	}
}

func TestGetTemplateDescriptions(t *testing.T) {
	tests := map[string]struct {
		in         TemplateRepo
		wantedType []TemplateDescription
		wantedErr  error
	}{
		"success case: input good repo details": {
			in: TemplateRepo{
				Description: "Standard Codewind templates.",
				URL:         "https://raw.githubusercontent.com/kabanero-io/codewind-templates/master/devfiles/index.json",
			},
			wantedType: []TemplateDescription{},
			wantedErr:  nil,
		},
		"fail case: input nil": {
			in:         TemplateRepo{},
			wantedType: nil,
			wantedErr:  new(url.Error),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetTemplateDescriptions(test.in)
			assert.IsType(t, test.wantedType, got)
			assert.IsType(t, test.wantedErr, err)
		})
	}
}

func TestFormatTemplate(t *testing.T) {
	in := TemplateDescription{
		DisplayName: "Codewind template",
		Description: "A starter template",
		Language:    "go",
		Location:    "/url/to/somewhere",
		ProjectType: "docker",
	}
	want := Template{
		Label:       "Codewind template",
		Description: "A starter template",
		Language:    "go",
		URL:         "/url/to/somewhere",
		ProjectType: "docker",
	}
	got := FormatTemplate(in)
	assert.Equal(t, want, got)
}
