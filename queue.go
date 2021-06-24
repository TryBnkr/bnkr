package main

import (
	"errors"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/services"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"gorm.io/gorm"
)

func listenForQueues() {
	// Run pool of 5 queue workers so we can run 5 jobs at the same time
	for i := 0; i < 5; i++ {
		go func() {
			for {
				job := <-app.Queue
				queueJob(job)
			}
		}()
	}
}

func queueJob(q types.NewQueueDTO) {
	typ := "job"
	if q.Process != "restore" {
		typ = "backup"
	}

	queue := &dal.Queue{}
	result := dal.FindQueueTypeAndObject(&queue, typ, q.ID)

	// If this operation already queued then do nothing
	if !(errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		return
	}

	que, _ := CreateQueue(typ, &q.ID)

	if q.Process == "restore" {
		job := &types.NewJobDTO{}
		if err := dal.FindJobsById(&job, q.ID).Error; err != nil {
			// TODO send mail here & handle error
			DeleteQueue(&que.ID)
			return
		}

		backup := &types.NewBackupDTO{}
		if err := dal.FindBackupsById(&backup, job.Backup).Error; err != nil {
			// TODO send mail here & handle error
			DeleteQueue(&que.ID)
			return
		}

		if err := services.Repo.RestoreBackup(backup, job); err != nil {
			// TODO send mail here & handle error
			DeleteQueue(&que.ID)
			return
		}
		// TODO handle error
		DeleteQueue(&que.ID)
	} else {
		backup := &types.NewBackupDTO{}

		if err := dal.FindBackupsById(&backup, q.ID).Error; err != nil {
			DeleteQueue(&que.ID)
			return
		}

		_, err := services.Repo.CreateNewJob(backup, true)

		if err != nil {
			DeleteQueue(&que.ID)
			return
		}

		DeleteQueue(&que.ID)
	}
}

func CreateQueue(t string, o *uint) (*dal.Queue, error) {
	q := &dal.Queue{
		Type:   t,
		Object: o,
	}

	if err := dal.CreateQueue(q).Error; err != nil {
		return nil, err
	}

	return q, nil
}

func DeleteQueue(id *uint) error {
	if res := dal.DeleteQueue(id); res.RowsAffected == 0 {
		return errors.New("unable to delete queue")
	}

	return nil
}
