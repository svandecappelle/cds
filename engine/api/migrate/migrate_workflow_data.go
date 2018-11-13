package migrate

import (
	"context"
	"github.com/go-gorp/gorp"

	"github.com/ovh/cds/engine/api/cache"
	"github.com/ovh/cds/engine/api/project"
	"github.com/ovh/cds/engine/api/workflow"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/log"
)

// MigrateToWorkflowData migrates all workflow from WorkflowNode to Node
func MigrateToWorkflowData(DBFunc func() *gorp.DbMap, store cache.Store) {
	log.Info("Start migrate MigrateToWorkflowData")
	defer func() {
		log.Info("End migrate MigrateToWorkflowData")
	}()

	for {
		db := DBFunc()
		var IDs []int64
		query := "SELECT id FROM workflow WHERE workflow_data IS NOT NULL AND to_delete = false AND root_node_id is not null AND id= 745"
		if _, err := db.Select(&IDs, query); err != nil {
			log.Error("MigrateToWorkflowData> Unable to select workflows id: %v", err)
			return
		}
		if len(IDs) == 0 {
			return
		}

		jobs := make(chan int64, 100)
		results := make(chan int64, 100)

		// 5 workers
		for w := 1; w <= 5; w++ {
			go migrationWorker(db, store, jobs, results)
		}

		for _, ID := range IDs {
			jobs <- ID
		}
		close(jobs)
		for a := 0; a < len(IDs); a++ {
			<-results
		}
	}
}

func migrationWorker(db *gorp.DbMap, store cache.Store, jobs <-chan int64, results chan<- int64) {
	for ID := range jobs {
		if err := migrateWorkflowData(db, store, ID); err != nil {
			log.Error("MigrateToWorkflowData> Unable to migrate workflow data %d: %v", ID, err)
		}
		results <- ID
	}
}

func migrateWorkflowData(db *gorp.DbMap, store cache.Store, ID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return sdk.WrapError(err, "MigrateToWorkflowData> unable to start transaction")
	}
	defer tx.Rollback() // nolint

	query := "SELECT id FROM workflow WHERE id=$1 FOR UPDATE NOWAIT"
	if _, err := tx.Exec(query, ID); err != nil {
		return nil
	}

	p, err := project.LoadProjectByWorkflowID(tx, store, nil, ID,
		project.LoadOptions.WithPlatforms,
		project.LoadOptions.WithPipelines,
		project.LoadOptions.WithEnvironments,
		project.LoadOptions.WithApplicationWithDeploymentStrategies)
	if err != nil {
		return sdk.WrapError(err, "migrateWorkflowData> Unable to load project from workflow %d", ID)
	}

	w, err := workflow.LoadByID(tx, store, p, ID, nil, workflow.LoadOptions{})
	if err != nil {
		return sdk.WrapError(err, "migrateWorkflowData> Unable to load workflow %d", ID)
	}

	oldW := *w

	for i := range w.WorkflowData.Joins {
		j := &w.WorkflowData.Joins[i]
		for k := range j.JoinContext {
			parentContext := &j.JoinContext[k]
			n := w.WorkflowData.NodeByID(parentContext.NodeID)
			parentContext.ParentName = n.Name
		}
	}

	if err := workflow.Update(context.Background(), tx, store, w, &oldW, p, nil); err != nil {
		return sdk.WrapError(err, "migrateWorkflowData> Unable to update join for %d", ID)
	}

	return tx.Commit()
}
