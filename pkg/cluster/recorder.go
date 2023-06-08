package cluster

import (
	"github.com/geowa4/ocm-workon/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type RecordedCluster struct {
	gorm.Model
	Environment       string
	Name              string
	ID                string
	ExternalID        string
	InfraID           string
	HiveShard         string
	ManagementCluster string
	ServiceCluster    string
}

func NewRecordedCluster(environment string, ncd *NormalizedClusterData) *RecordedCluster {
	return &RecordedCluster{
		Environment:       environment,
		Name:              ncd.Name,
		ID:                ncd.InternalID,
		ExternalID:        ncd.ExternalID,
		InfraID:           ncd.InfraID,
		HiveShard:         ncd.HiveShard,
		ManagementCluster: ncd.ManagementCluster,
		ServiceCluster:    ncd.ServiceCluster,
	}
}

func (cluster *RecordedCluster) RecordAccess(baseDir string) error {
	db, err := makeDb(baseDir)
	if err != nil {
		return err
	}
	save(db, cluster)
	return nil
}

func makeDb(baseDir string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(baseDir+utils.PathSep+"workon.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&RecordedCluster{}); err != nil {
		return nil, err
	}
	return db, nil
}

func save(db *gorm.DB, ncd *RecordedCluster) {
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(ncd)
}

func FindClustersUpdatedSinceTwoWeeksAgo(baseDir string) (clusters []RecordedCluster, err error) {
	db, err := makeDb(baseDir)
	if err != nil {
		return
	}
	twoWeeksAgo := time.Now().Add(-24 * time.Hour * 14)
	db.Where("updated_at > ?", twoWeeksAgo).Find(&clusters)
	return
}
