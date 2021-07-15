package assets

// AssetsEmbedded indicates whether or not static assets are embedded in this binary
func AssetsEmbedded() bool {
	return assetsEmbedded
}

// GetAssets gets an embedded static asset
func GetAsset(name string) ([]byte, error) {
	return getAsset(name)
}
