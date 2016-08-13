// copydir_test
package copydir

import "testing"

func TestCopy(t *testing.T) {
	for _, c := range []struct {
		src, dest string
		ovw       bool
		want      string
	}{
		{"not dir", "new data", false, "Source is not a directory."},
		{"sample data", "not dir", false, "Destination is not a directory."},
		{"sample data", "backup", true, ""},
		{"sample data", "backup", false, "We will not overwrite the destination."},
	} {
		err := Copy(c.src, c.dest, c.ovw)
		if err != nil {
			if err.Error() != c.want {
				t.Errorf("Copy %q to %q and overwrite %t (error %q),\nwant: %v", c.src, c.dest, c.ovw, err, c.want)
			}
		}
	}
}
