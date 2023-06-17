package cluster

import (
	"github.com/geowa4/ocm-workon/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	Elevations        []Elevation `json:",omitempty"`
}

type Elevation struct {
	gorm.Model
	RecordedClusterID string `json:"-"`
	RecordedCluster   RecordedCluster
	Source            string
	Reason            string
}

func RecordElevation(baseDir string, clusterId string, source string, reason string) error {
	db, err := makeDb(baseDir)
	if err != nil {
		return err
	}
	elevation := &Elevation{
		RecordedCluster: RecordedCluster{ID: clusterId},
		Source:          source,
		Reason:          reason,
	}
	db.Save(elevation)
	return nil
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
	db.Save(cluster)
	return nil
}

func makeDb(baseDir string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(baseDir+utils.PathSep+"workon.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&Elevation{}, &RecordedCluster{}); err != nil {
		return nil, err
	}
	return db, nil
}

func findRecordingsSince(baseDir, timeAgo string, queryFunc func(*gorm.DB, time.Time)) (err error) {
	db, err := makeDb(baseDir)
	if err != nil {
		return
	}
	timeAgoAsDuration, err := time.ParseDuration(timeAgo)
	if err != nil {
		return
	}
	sinceWhen := time.Now().Add(-1 * timeAgoAsDuration)
	queryFunc(db, sinceWhen)
	return
}

func FindRecordedClustersSince(baseDir, timeAgo string) (clusters []RecordedCluster, err error) {
	err = findRecordingsSince(baseDir, timeAgo, func(db *gorm.DB, sinceWhen time.Time) {
		db.Where("updated_at > ?", sinceWhen).
			Order("updated_at").
			Find(&clusters)
	})
	return
}

func FindElevationsSince(baseDir, timeAgo string) (elevations []Elevation, err error) {
	err = findRecordingsSince(baseDir, timeAgo, func(db *gorm.DB, sinceWhen time.Time) {
		db.Joins("RecordedCluster").Where("elevations.updated_at > ?", sinceWhen).
			Order("elevations.updated_at").
			Find(&elevations)
	})
	return
}
