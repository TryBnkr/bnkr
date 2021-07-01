package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/services"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
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
	err := dal.FindQueueByTypeAndObject(queue, typ, q.ID)

	// If this operation already queued then do nothing
	if !(errors.Is(err, sql.ErrNoRows)) {
		return
	}

	que, _ := CreateQueue(typ, q.ID)

	if q.Process == "restore" {
		job := &types.NewJobDTO{}
		if err := dal.FindJobById(job, q.ID); err != nil {
			// TODO send mail here & handle error
			DeleteQueue(que.ID)
			return
		}

		backup := &types.NewBackupDTO{}
		if err := dal.FindBackupById(backup, job.Backup); err != nil {
			// TODO send mail here & handle error
			DeleteQueue(que.ID)
			return
		}

		if err := services.Repo.RestoreBackup(backup, job); err != nil {
			// TODO send mail here & handle error
			DeleteQueue(que.ID)
			return
		}
		// TODO handle error
		DeleteQueue(que.ID)
	} else {
		backup := &types.NewBackupDTO{}

		if err := dal.FindBackupById(backup, q.ID); err != nil {
			DeleteQueue(que.ID)
			return
		}

		_, err := services.Repo.CreateNewJob(backup, true)

		if err != nil {
			DeleteQueue(que.ID)
			return
		}

		DeleteQueue(que.ID)
	}
}

func CreateQueue(t string, o int) (*dal.Queue, error) {
	q := &dal.Queue{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Type:      t,
		Object:    o,
	}

	if _, err := dal.CreateQueue(q); err != nil {
		return nil, err
	}

	return q, nil
}

func DeleteQueue(id int) error {
	if _, err := dal.DeleteQueue(id); err != nil {
		return errors.New("unable to delete queue")
	}

	return nil
}
