package helper

import "testing"

func TestResolveS3BucketNameSupportsLegacyAndCurrentEnvKeys(t *testing.T) {
	t.Setenv("PRIVATE_SE_BUCKET_NAME", "")
	t.Setenv("PRIVATE_S3_BUCKET_NAME", "genzone-private-assets")

	if got := resolveBucketName("PRIVATE_SE_BUCKET_NAME", "PRIVATE_S3_BUCKET_NAME"); got != "genzone-private-assets" {
		t.Fatalf("expected fallback bucket name to resolve, got %q", got)
	}
}
