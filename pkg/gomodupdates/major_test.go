package gomodupdates

import (
	"context"
	"reflect"
	"testing"
)

type fakeVersionLister map[string][]string

func (l fakeVersionLister) Versions(_ context.Context, modulePath string) ([]string, error) {
	return l[modulePath], nil
}

func TestFindMajorUpdate(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		version  string
		versions fakeVersionLister
		want     MajorUpdate
	}{
		{
			name:    "slash major from v1",
			path:    "github.com/hashicorp/golang-lru",
			version: "v1.0.2",
			versions: fakeVersionLister{
				"github.com/hashicorp/golang-lru/v2": {"v2.0.0-rc.1", "v2.0.7"},
			},
			want: MajorUpdate{Path: "github.com/hashicorp/golang-lru/v2", Version: "v2.0.7"},
		},
		{
			name:    "already suffixed module checks next major",
			path:    "github.com/containerd/containerd/v2",
			version: "v2.2.2",
			versions: fakeVersionLister{
				"github.com/containerd/containerd/v3": {"v3.0.0"},
			},
			want: MajorUpdate{Path: "github.com/containerd/containerd/v3", Version: "v3.0.0"},
		},
		{
			name:    "slash major with dotted module root",
			path:    "gotest.tools/v3",
			version: "v3.5.2",
			versions: fakeVersionLister{
				"gotest.tools/v4": {"v4.0.0"},
			},
			want: MajorUpdate{Path: "gotest.tools/v4", Version: "v4.0.0"},
		},
		{
			name:    "slash major with dotted vanity path",
			path:    "go.yaml.in/yaml/v4",
			version: "v4.0.0",
			versions: fakeVersionLister{
				"go.yaml.in/yaml/v5": {"v5.0.0"},
			},
			want: MajorUpdate{Path: "go.yaml.in/yaml/v5", Version: "v5.0.0"},
		},
		{
			name:    "stops after missing next major",
			path:    "github.com/Microsoft/hcsshim",
			version: "v1.0.0",
			versions: fakeVersionLister{
				"github.com/Microsoft/hcsshim/v3": {"v3.0.0"},
			},
			want: MajorUpdate{},
		},
		{
			name:    "gopkg major",
			path:    "gopkg.in/yaml.v3",
			version: "v3.0.1",
			versions: fakeVersionLister{
				"gopkg.in/yaml.v4": {"v4.0.0"},
			},
			want: MajorUpdate{Path: "gopkg.in/yaml.v4", Version: "v4.0.0"},
		},
		{
			name:    "ignores prereleases",
			path:    "github.com/example/mod",
			version: "v1.0.0",
			versions: fakeVersionLister{
				"github.com/example/mod/v2": {"v2.0.0-alpha.1"},
			},
			want: MajorUpdate{},
		},
		{
			name:    "plus incompatible starts after current major",
			path:    "github.com/docker/cli",
			version: "v29.4.0+incompatible",
			versions: fakeVersionLister{
				"github.com/docker/cli/v2":  {"v2.0.0"},
				"github.com/docker/cli/v30": {"v30.0.0"},
			},
			want: MajorUpdate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindMajorUpdate(context.Background(), tt.path, tt.version, tt.versions, 10)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestLatestStableUsesSemanticVersionOrder(t *testing.T) {
	got := latestStable([]string{
		"not-a-version",
		"v2.0.9",
		"v2.0.10",
		"v2.1.0-rc.1",
	})
	if got != "v2.0.10" {
		t.Fatalf("expected v2.0.10, got %s", got)
	}
}
