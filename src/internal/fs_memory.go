package internal

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	tasksDir = "tasks"
)

type FSMemory struct {
	rootDir string

	mux        *sync.Mutex
	lastTaskId int64
}

func NewFSMemory(root string) (*FSMemory, error) {
	fsMem := &FSMemory{
		rootDir: filepath.Clean(root),
		mux:     new(sync.Mutex),
	}

	err := os.MkdirAll(filepath.Join(fsMem.rootDir, tasksDir), 0755)
	if err != nil {
		return nil, err
	}

	err = fsMem.getLastId()
	if err != nil {
		return nil, err
	}

	// assume if there was not an error before,
	// there won't be any from there
	go fsMem.loop()
	return fsMem, nil
}

func (fs *FSMemory) getLastId() error {
	lastIdFile := filepath.Join(fs.rootDir, "last.id")
	lastIdContent, err := ioutil.ReadFile(lastIdFile)
	if err != nil {
		fs.lastTaskId = 0
		return fs.updateLastId()
	}

	fs.lastTaskId, err = strconv.ParseInt(string(lastIdContent), 10, 64)
	return err
}

func (fs *FSMemory) updateLastId() error {
	id := atomic.LoadInt64(&fs.lastTaskId)
	idString := strconv.FormatInt(id, 10)
	return ioutil.WriteFile(filepath.Join(fs.rootDir, "last.id"), []byte(idString), 0755)
}

func (fs *FSMemory) loop() {
	for {
		_ = fs.updateLastId()
		time.Sleep(time.Second)
	}
}

func (fs *FSMemory) withLock(f func()) {
	fs.mux.Lock()
	f()
	fs.mux.Unlock()
}

func (fs *FSMemory) resolveTaskFilename(id int64) string {
	filename := strconv.FormatInt(id, 10) + ".task"
	relFilePath := filepath.Join(tasksDir, filename)
	filePath := filepath.Join(fs.rootDir, relFilePath)

	log.Println("Task filename:", filePath)
	return filePath
}

func (fs *FSMemory) SaveTask(task *Task) (int64, error) {
	taskBody, err := jsoniter.Marshal(task)
	if err != nil {
		return 0, err
	}

	if task.Id == -1 {
		fs.withLock(func() {
			task.Id = fs.lastTaskId + 1
			fs.lastTaskId = task.Id
		})
	}

	err = ioutil.WriteFile(fs.resolveTaskFilename(task.Id), taskBody, 0755)
	return task.Id, err
}

func (fs *FSMemory) GetTask(id int64) (*Task, error) {

	payload, err := ioutil.ReadFile(fs.resolveTaskFilename(id))
	if err != nil {
		return nil, err
	}

	return newTaskFromPayload(id, payload)
}

func (fs *FSMemory) GetAllTasks() ([]*Task, error) {
	absTaskDir := filepath.Join(fs.rootDir, tasksDir)

	log.Println("Reading:", absTaskDir)
	dirInfo, err := ioutil.ReadDir(absTaskDir)
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, 0, len(dirInfo))
	for _, fileInfo := range dirInfo {
		filename := fileInfo.Name()
		// any case where there could NOT be a dot in the filename?
		taskId, err := strconv.ParseInt(filename[:strings.IndexByte(filename, '.')], 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}

		taskFilename := filepath.Join(absTaskDir, filename)
		payload, err := ioutil.ReadFile(taskFilename)
		if err != nil {
			log.Println(err)
			continue
		}

		task, err := newTaskFromPayload(taskId, payload)
		if err != nil {
			log.Println(err)
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (fs *FSMemory) DeleteTask(id int64) error {
	err := os.Remove(fs.resolveTaskFilename(id))
	return err
}
