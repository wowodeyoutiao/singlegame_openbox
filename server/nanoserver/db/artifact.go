package db

import (
	"fmt"
	"nanoserver/db/model"
)

func GetArtifactByUid(uid int64) (*model.Artifact, error) {
	artifact := &model.Artifact{UID: uid}

	has, err := database.Get(artifact)
	if err != nil {
		return nil, fmt.Errorf("get artifact by uid failed: %w", err)
	}

	if !has {
		return nil, nil
	}

	return artifact, nil
}

func InsertArtifact(artifact *model.Artifact) error {
	_, err := database.Insert(artifact)
	return err
}

func UpdateArtifact(artifact *model.Artifact) error {
	_, err := database.Where("uid = ?", artifact.UID).Update(artifact)
	return err
}
