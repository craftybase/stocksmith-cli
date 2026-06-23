package selfupdate

import "testing"

func TestAssetNames(t *testing.T) {
	c := Config{BinaryName: "craftybase", GOOS: "darwin", GOARCH: "arm64"}
	archive, checksums := c.assetNames("v0.3.0")
	if archive != "craftybase_0.3.0_darwin_arm64.tar.gz" {
		t.Errorf("archive = %q", archive)
	}
	if checksums != "craftybase_0.3.0_checksums.txt" {
		t.Errorf("checksums = %q", checksums)
	}
}

func TestIsDevVersion(t *testing.T) {
	for _, v := range []string{"", "dev", "garbage", "none"} {
		if !isDevVersion(v) {
			t.Errorf("isDevVersion(%q) = false, want true", v)
		}
	}
	for _, v := range []string{"0.2.0", "v0.2.0", "1.0.0-rc1"} {
		if isDevVersion(v) {
			t.Errorf("isDevVersion(%q) = true, want false", v)
		}
	}
}

func TestUpdateAvailable(t *testing.T) {
	if !updateAvailable("0.2.0", "v0.3.0") {
		t.Error("0.2.0 -> v0.3.0 should be available")
	}
	if updateAvailable("0.3.0", "v0.3.0") {
		t.Error("equal versions: not available")
	}
	if updateAvailable("0.4.0", "v0.3.0") {
		t.Error("current ahead of latest: not available")
	}
}
