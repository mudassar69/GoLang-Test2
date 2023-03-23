package models

import (
	"errors"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `gorm:"size:100;not null" json:"title"`
	Description string `gorm:"not null"                 json:"description"`
	CreatedBy   User   `gorm:"foreignKey:UserID;"       json:"-"`
	UserID      uint   `gorm:"not null"                 json:"user_id"`
}

type AssignTask struct {
	Email   string `json:"email"`
	TaskObj Task   `json:"task"`
}

const TABLE_TASKS = "tasks"

func (t *Task) SaveTask(db *gorm.DB) (*Task, error) {

	err := db.Create(&t).Error
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Task) GetTask(db *gorm.DB) (*Task, error) {

	task := &Task{}
	err := db.Table(TABLE_TASKS).Where("title = ?", t.Title).First(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *Task) ValidateTask() error {
	if t.Title == "" {
		return errors.New("Title is required")
	}
	if t.Description == "" {
		return errors.New("Description is required")
	}
	return nil
}

func GetUserTasks(id int, db *gorm.DB) (*[]Task, error) {
	tasks := []Task{}

	err := db.Table(TABLE_TASKS).Where("user_id = ?", id).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return &tasks, nil
}

func GetTaskByID(taskID, userID int, db *gorm.DB) (*Task, error) {

	task := &Task{}
	err := db.Table(TABLE_TASKS).Where("id = ? AND user_id = ?", taskID, userID).First(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func DeleteTask(taskID, userID int, db *gorm.DB) error {
	err := db.Table(TABLE_TASKS).Where("id = ?", taskID).Delete(&Task{}).Error
	return err
}

func (t *Task) UpdateTask(taskId int, db *gorm.DB) error {
	err := db.Table(TABLE_TASKS).Where("id = ?", taskId).Updates(Task{Title: t.Title, Description: t.Description}).Error
	return err
}
