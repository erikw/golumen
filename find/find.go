package find

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
)

const defaultWorkerCount = 8

type Finder struct {
	logger *slog.Logger
	follow bool
}

type walkTask struct {
	path   string
	isRoot bool
}

type taskQueue struct {
	mu      sync.Mutex
	cond    *sync.Cond
	tasks   []walkTask
	pending int
	closed  bool
}

type findCollector struct {
	logger       *slog.Logger
	searchRegexp *regexp.Regexp
	follow       bool
	visitedDirs  map[string]struct{}
	matches      []string
	rootErr      error
	mu           sync.Mutex
}

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func New(logger *slog.Logger, follow bool) *Finder {
	return &Finder{logger: logger, follow: follow}
}

func (f *Finder) Find(path string, pattern string) (matches []string, err error) {
	f.logger.Debug(fmt.Sprintf("💡 Shedding light to %s for %s:", path, pattern))

	var r *regexp.Regexp
	var rErr error
	if r, rErr = regexp.Compile(pattern); rErr != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex pattern \"%s\": %v\n", pattern, rErr.Error())
		os.Exit(1)
	}

	fc := findCollector{
		logger:       f.logger,
		searchRegexp: r,
		follow:       f.follow,
		visitedDirs:  make(map[string]struct{}),
	}

	err = fc.walk(path)
	if err != nil {
		return nil, err
	}

	return fc.sortedMatches(), nil
}

func (fc *findCollector) walk(path string) error {
	queue := newTaskQueue()
	queue.push(walkTask{path: path, isRoot: true})

	var workers sync.WaitGroup

	for i := 0; i < defaultWorkerCount; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for {
				task, ok := queue.pop()
				if !ok {
					return
				}

				fc.processTask(task, queue)
			}
		}()
	}

	workers.Wait()

	return fc.rootError()
}

func (fc *findCollector) processTask(task walkTask, queue *taskQueue) {
	defer queue.done()

	if err := fc.walkPath(task, queue); err != nil {
		fc.setRootError(err)
	}
}

func (fc *findCollector) walkPath(task walkTask, queue *taskQueue) error {
	info, err := os.Lstat(task.path)
	if err != nil {
		return fc.handlePathError(task, "Could not stat path. Skipping.", err)
	}

	baseName := filepath.Base(task.path)
	skip := !task.isRoot && fc.blockPath(baseName)
	fc.logger.Debug("Walking path", "path", task.path, "skip", skip)
	if skip {
		return nil
	}

	if fc.patternMatch(baseName) {
		fc.addMatch(task.path)
	}

	if info.IsDir() {
		return fc.walkDirectory(task.path, task.isRoot, queue)
	}

	if info.Mode()&fs.ModeSymlink == 0 || !fc.follow {
		return nil
	}

	targetInfo, err := os.Stat(task.path)
	if err != nil {
		return fc.handlePathError(task, "Could not resolve symlink target. Skipping.", err)
	}

	if !targetInfo.IsDir() {
		return nil
	}

	return fc.walkDirectory(task.path, task.isRoot, queue)
}

func (fc *findCollector) walkDirectory(path string, isRoot bool, queue *taskQueue) error {
	task := walkTask{path: path, isRoot: isRoot}

	resolvedPath, err := resolvedPath(path)
	if err != nil {
		return fc.handlePathError(task, "Could not resolve directory path. Skipping.", err)
	}

	if fc.markVisited(resolvedPath) {
		fc.logger.Debug("Skipping already visited directory", "path", path, "resolvedPath", resolvedPath)
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fc.handlePathError(task, "Could not enter directory. Skipping.", err)
	}

	for _, entry := range entries {
		queue.push(walkTask{
			path: filepath.Join(path, entry.Name()),
		})
	}

	return nil
}

func newTaskQueue() *taskQueue {
	queue := &taskQueue{}
	queue.cond = sync.NewCond(&queue.mu)
	return queue
}

func (q *taskQueue) push(task walkTask) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tasks = append(q.tasks, task)
	q.pending++
	q.cond.Signal()
}

func (q *taskQueue) pop() (walkTask, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.tasks) == 0 && !q.closed {
		q.cond.Wait()
	}

	if len(q.tasks) == 0 {
		return walkTask{}, false
	}

	last := len(q.tasks) - 1
	task := q.tasks[last]
	q.tasks = q.tasks[:last]
	return task, true
}

func (q *taskQueue) done() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.pending--
	if q.pending == 0 {
		q.closed = true
		q.cond.Broadcast()
	}
}

func (fc *findCollector) handlePathError(task walkTask, msg string, err error) error {
	if task.isRoot {
		return err
	}

	fc.logger.Debug(msg, "path", task.path, "error", err.Error())
	return nil
}

func (fc *findCollector) markVisited(resolvedPath string) bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, seen := fc.visitedDirs[resolvedPath]; seen {
		return true
	}

	fc.visitedDirs[resolvedPath] = struct{}{}
	return false
}

func (fc *findCollector) addMatch(path string) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.matches = append(fc.matches, path)
}

func (fc *findCollector) setRootError(err error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if fc.rootErr == nil {
		fc.rootErr = err
	}
}

func (fc *findCollector) rootError() error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	return fc.rootErr
}

func (fc *findCollector) sortedMatches() []string {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	matches := append([]string(nil), fc.matches...)
	sort.Strings(matches)
	return matches
}

func resolvedPath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	resolvedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", err
	}

	return filepath.Clean(resolvedPath), nil
}

func (fc *findCollector) blockPath(baseName string) bool {
	_, ok := defaultBlockedPaths[baseName]
	fc.logger.Debug("Consider to block file", "baseName", baseName, "block", ok)
	return ok
}

func (fc *findCollector) patternMatch(baseName string) bool {
	m := fc.searchRegexp.MatchString(baseName)
	fc.logger.Debug("Trying to match file name to pattern", "baseName", baseName, "pattern", fc.searchRegexp.String(), "match", m)
	return m
}
