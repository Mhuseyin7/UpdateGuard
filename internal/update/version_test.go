package update

import "testing"

func TestClassifySemver(t *testing.T) {
	cases := []struct {
		from, to string
		want     VersionChange
	}{{"1.2.3", "1.2.4", Patch}, {"v1.2.3", "1.3.0", Minor}, {"1.2.3", "2.0.0", Major}, {"postgres:16", "postgres:17", Unknown}}
	for _, tt := range cases {
		if got := ClassifySemver(tt.from, tt.to); got != tt.want {
			t.Errorf("%s → %s = %s", tt.from, tt.to, got)
		}
	}
}
