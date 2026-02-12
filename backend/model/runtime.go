package model

type TempArtifact struct {
	Path string `json:"path"`
}

type TempArtifactIndex struct {
	Artifacts []TempArtifact `json:"artifacts"`
}
