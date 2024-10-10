package config

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/filters"
	bkconfig "github.com/moby/buildkit/cmd/buildkitd/config"
)

// BuilderGCRule represents a GC rule for buildkit cache
type BuilderGCRule struct {
	All           bool            `json:",omitempty"`
	Filter        BuilderGCFilter `json:",omitempty"`
	ReservedSpace string          `json:",omitempty"`
	MaxUsedSpace  string          `json:",omitempty"`
	MinFreeSpace  string          `json:",omitempty"`
}

func (x *BuilderGCRule) UnmarshalJSON(data []byte) error {
	var v struct {
		BuilderGCRule

		// Deprecated option is now equivalent to ReservedSpace.
		KeepStorage string `json:",omitempty"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.KeepStorage != "" && v.ReservedSpace == "" {
		v.ReservedSpace = v.KeepStorage
	}
	*x = v.BuilderGCRule
	return nil
}

// BuilderGCFilter contains garbage-collection filter rules for a BuildKit builder
type BuilderGCFilter filters.Args

// MarshalJSON returns a JSON byte representation of the BuilderGCFilter
func (x *BuilderGCFilter) MarshalJSON() ([]byte, error) {
	f := filters.Args(*x)
	keys := f.Keys()
	sort.Strings(keys)
	arr := make([]string, 0, len(keys))
	for _, k := range keys {
		values := f.Get(k)
		for _, v := range values {
			arr = append(arr, k+"="+v)
		}
	}
	return json.Marshal(arr)
}

// UnmarshalJSON fills the BuilderGCFilter values structure from JSON input
func (x *BuilderGCFilter) UnmarshalJSON(data []byte) error {
	var arr []string
	f := filters.NewArgs()
	if err := json.Unmarshal(data, &arr); err != nil {
		// backwards compat for deprecated buggy form
		err := json.Unmarshal(data, &f)
		*x = BuilderGCFilter(f)
		return err
	}
	for _, s := range arr {
		name, value, _ := strings.Cut(s, "=")
		name = strings.ToLower(strings.TrimSpace(name))
		value = strings.TrimSpace(value)
		f.Add(name, value)
	}
	*x = BuilderGCFilter(f)
	return nil
}

// BuilderGCConfig contains GC config for a buildkit builder
type BuilderGCConfig struct {
	Enabled            bool            `json:",omitempty"`
	Policy             []BuilderGCRule `json:",omitempty"`
	DefaultKeepStorage string          `json:",omitempty"`
}

// BuilderHistoryConfig contains history config for a buildkit builder
type BuilderHistoryConfig struct {
	MaxAge     bkconfig.Duration `json:",omitempty"`
	MaxEntries int64             `json:",omitempty"`
}

// BuilderEntitlements contains settings to enable/disable entitlements
type BuilderEntitlements struct {
	NetworkHost      *bool `json:"network-host,omitempty"`
	SecurityInsecure *bool `json:"security-insecure,omitempty"`
}

// BuilderConfig contains config for the builder
type BuilderConfig struct {
	GC           BuilderGCConfig       `json:",omitempty"`
	Entitlements BuilderEntitlements   `json:",omitempty"`
	History      *BuilderHistoryConfig `json:",omitempty"`
}
