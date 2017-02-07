package watcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// setup creates all required files and folders for
// the tests and returns a function that is used as
// a teardown function when the tests are done.
func setup(t testing.TB) (string, func()) {
	testDir, err := ioutil.TempDir(".", "")
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(testDir, "file.txt"),
		[]byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{"file_1.txt", "file_2.txt", "file_3.txt"}

	for _, f := range files {
		filePath := filepath.Join(testDir, f)
		if err := ioutil.WriteFile(filePath, []byte{}, 0755); err != nil {
			t.Fatal(err)
		}
	}

	err = ioutil.WriteFile(filepath.Join(testDir, ".dotfile"),
		[]byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	testDirTwo := filepath.Join(testDir, "testDirTwo")
	err = os.Mkdir(testDirTwo, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(testDirTwo, "file_recursive.txt"),
		[]byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	abs, err := filepath.Abs(testDir)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatal(err)
	}
	return abs, func() {
		if os.RemoveAll(testDir); err != nil {
			t.Fatal(err)
		}
	}
}

func TestWatcherAdd(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	if err := w.Add(testDir); err != nil {
		t.Fatal(err)
	}

	if len(w.files) != 7 {
		t.Errorf("expected len(w.files) to be 7, got %d", len(w.files))
	}

	// Make sure w.names contains testDir
	if _, found := w.names[testDir]; !found {
		t.Errorf("expected w.names to contain testDir")
	}

	if _, found := w.files[testDir]; !found {
		t.Errorf("expected to find %s", testDir)
	}

	if w.files[testDir].Name() != filepath.Base(testDir) {
		t.Errorf("expected w.files[%q].Name() to be %s, got %s",
			testDir, testDir, w.files[testDir].Name())
	}

	dotFile := filepath.Join(testDir, ".dotfile")
	if _, found := w.files[dotFile]; !found {
		t.Errorf("expected to find %s", dotFile)
	}

	if w.files[dotFile].Name() != ".dotfile" {
		t.Errorf("expected w.files[%q].Name() to be .dotfile, got %s",
			dotFile, w.files[dotFile].Name())
	}

	fileRecursive := filepath.Join(testDir, "testDirTwo", "file_recursive.txt")
	if _, found := w.files[fileRecursive]; found {
		t.Errorf("expected to not find %s", fileRecursive)
	}

	fileTxt := filepath.Join(testDir, "file.txt")
	if _, found := w.files[fileTxt]; !found {
		t.Errorf("expected to find %s", fileTxt)
	}

	if w.files[fileTxt].Name() != "file.txt" {
		t.Errorf("expected w.files[%q].Name() to be file.txt, got %s",
			fileTxt, w.files[fileTxt].Name())
	}

	dirTwo := filepath.Join(testDir, "testDirTwo")
	if _, found := w.files[dirTwo]; !found {
		t.Errorf("expected to find %s directory", dirTwo)
	}

	if w.files[dirTwo].Name() != "testDirTwo" {
		t.Errorf("expected w.files[%q].Name() to be testDirTwo, got %s",
			dirTwo, w.files[dirTwo].Name())
	}
}

func TestIgnoreHiddenFilesRecursive(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()
	w.IgnoreHiddenFiles(true)

	if err := w.AddRecursive(testDir); err != nil {
		t.Fatal(err)
	}

	if len(w.files) != 7 {
		t.Errorf("expected len(w.files) to be 7, got %d", len(w.files))
	}

	// Make sure w.names contains testDir
	if _, found := w.names[testDir]; !found {
		t.Errorf("expected w.names to contain testDir")
	}

	if _, found := w.files[testDir]; !found {
		t.Errorf("expected to find %s", testDir)
	}

	if w.files[testDir].Name() != filepath.Base(testDir) {
		t.Errorf("expected w.files[%q].Name() to be %s, got %s",
			testDir, filepath.Base(testDir), w.files[testDir].Name())
	}

	fileRecursive := filepath.Join(testDir, "testDirTwo", "file_recursive.txt")
	if _, found := w.files[fileRecursive]; !found {
		t.Errorf("expected to find %s", fileRecursive)
	}

	if _, found := w.files[filepath.Join(testDir, ".dotfile")]; found {
		t.Error("expected to not find .dotfile")
	}

	fileTxt := filepath.Join(testDir, "file.txt")
	if _, found := w.files[fileTxt]; !found {
		t.Errorf("expected to find %s", fileTxt)
	}

	if w.files[fileTxt].Name() != "file.txt" {
		t.Errorf("expected w.files[%q].Name() to be file.txt, got %s",
			fileTxt, w.files[fileTxt].Name())
	}

	dirTwo := filepath.Join(testDir, "testDirTwo")
	if _, found := w.files[dirTwo]; !found {
		t.Errorf("expected to find %s directory", dirTwo)
	}

	if w.files[dirTwo].Name() != "testDirTwo" {
		t.Errorf("expected w.files[%q].Name() to be testDirTwo, got %s",
			dirTwo, w.files[dirTwo].Name())
	}
}

func TestIgnoreHiddenFiles(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()
	w.IgnoreHiddenFiles(true)

	if err := w.Add(testDir); err != nil {
		t.Fatal(err)
	}

	if len(w.files) != 6 {
		t.Errorf("expected len(w.files) to be 6, got %d", len(w.files))
	}

	// Make sure w.names contains testDir
	if _, found := w.names[testDir]; !found {
		t.Errorf("expected w.names to contain testDir")
	}

	if _, found := w.files[testDir]; !found {
		t.Errorf("expected to find %s", testDir)
	}

	if w.files[testDir].Name() != filepath.Base(testDir) {
		t.Errorf("expected w.files[%q].Name() to be %s, got %s",
			testDir, filepath.Base(testDir), w.files[testDir].Name())
	}

	if _, found := w.files[filepath.Join(testDir, ".dotfile")]; found {
		t.Error("expected to not find .dotfile")
	}

	fileRecursive := filepath.Join(testDir, "testDirTwo", "file_recursive.txt")
	if _, found := w.files[fileRecursive]; found {
		t.Errorf("expected to not find %s", fileRecursive)
	}

	fileTxt := filepath.Join(testDir, "file.txt")
	if _, found := w.files[fileTxt]; !found {
		t.Errorf("expected to find %s", fileTxt)
	}

	if w.files[fileTxt].Name() != "file.txt" {
		t.Errorf("expected w.files[%q].Name() to be file.txt, got %s",
			fileTxt, w.files[fileTxt].Name())
	}

	dirTwo := filepath.Join(testDir, "testDirTwo")
	if _, found := w.files[dirTwo]; !found {
		t.Errorf("expected to find %s directory", dirTwo)
	}

	if w.files[dirTwo].Name() != "testDirTwo" {
		t.Errorf("expected w.files[%q].Name() to be testDirTwo, got %s",
			dirTwo, w.files[dirTwo].Name())
	}
}

func TestWatcherAddRecursive(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	if err := w.AddRecursive(testDir); err != nil {
		t.Fatal(err)
	}

	// Make sure len(w.files) is 8.
	if len(w.files) != 8 {
		t.Errorf("expected 8 files, found %d", len(w.files))
	}

	// Make sure w.names contains testDir
	if _, found := w.names[testDir]; !found {
		t.Errorf("expected w.names to contain testDir")
	}

	dirTwo := filepath.Join(testDir, "testDirTwo")
	if _, found := w.files[dirTwo]; !found {
		t.Errorf("expected to find %s directory", dirTwo)
	}

	if w.files[dirTwo].Name() != "testDirTwo" {
		t.Errorf("expected w.files[%q].Name() to be testDirTwo, got %s",
			"testDirTwo", w.files[dirTwo].Name())
	}

	fileRecursive := filepath.Join(dirTwo, "file_recursive.txt")
	if _, found := w.files[fileRecursive]; !found {
		t.Errorf("expected to find %s directory", fileRecursive)
	}

	if w.files[fileRecursive].Name() != "file_recursive.txt" {
		t.Errorf("expected w.files[%q].Name() to be file_recursive.txt, got %s",
			fileRecursive, w.files[fileRecursive].Name())
	}
}

func TestWatcherAddNotFound(t *testing.T) {
	w := New()

	// Make sure there is an error when adding a
	// non-existent file/folder.
	if err := w.AddRecursive("random_filename.txt"); err == nil {
		t.Error("expected a file not found error")
	}
}

func TestWatcherRemoveRecursive(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	// Add the testDir to the watchlist.
	if err := w.AddRecursive(testDir); err != nil {
		t.Fatal(err)
	}

	// Make sure len(w.files) is 8.
	if len(w.files) != 8 {
		t.Errorf("expected 8 files, found %d", len(w.files))
	}

	// Now remove the folder from the watchlist.
	if err := w.RemoveRecursive(testDir); err != nil {
		t.Error(err)
	}

	// Now check that there is nothing being watched.
	if len(w.files) != 0 {
		t.Errorf("expected len(w.files) to be 0, got %d", len(w.files))
	}

	// Make sure len(w.names) is now 0.
	if len(w.names) != 0 {
		t.Errorf("expected len(w.names) to be empty, len(w.names): %d", len(w.names))
	}
}

func TestListFiles(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()
	w.AddRecursive(testDir)

	fileList := w.retrieveFileList()
	if fileList == nil {
		t.Error("expected file list to not be empty")
	}

	// Make sure fInfoTest contains the correct os.FileInfo names.
	fname := filepath.Join(testDir, "file.txt")
	if fileList[fname].Name() != "file.txt" {
		t.Errorf("expected fileList[%s].Name() to be file.txt, got %s",
			fname, fileList[fname].Name())
	}
}

func TestTriggerEvent(t *testing.T) {
	w := New()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case event := <-w.Event:
			if event.Name() != "triggered event" {
				t.Errorf("expected event file name to be triggered event, got %s",
					event.Name())
			}
		case <-time.After(time.Millisecond * 250):
			t.Fatal("received no event from Event channel")
		}
	}()

	go func() {
		// Start the watching process.
		if err := w.Start(time.Millisecond * 100); err != nil {
			t.Fatal(err)
		}
	}()

	w.TriggerEvent(Create, nil)

	wg.Wait()
}

func TestEventAddFile(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	// Add the testDir to the watchlist.
	if err := w.AddRecursive(testDir); err != nil {
		t.Fatal(err)
	}

	files := map[string]bool{
		"newfile_1.txt": false,
		"newfile_2.txt": false,
		"newfile_3.txt": false,
	}

	for f := range files {
		filePath := filepath.Join(testDir, f)
		if err := ioutil.WriteFile(filePath, []byte{}, 0755); err != nil {
			t.Error(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		events := 0
		for {
			select {
			case event := <-w.Event:
				if event.Op != Create {
					t.Errorf("expected event to be Create, got %s", event.Op)
				}

				files[event.Name()] = true
				events++

				if events == len(files) {
					return
				}
			case <-time.After(time.Millisecond * 250):
				for f, e := range files {
					if !e {
						t.Errorf("received no event for file %s", f)
					}
				}
				return
			}
		}
	}()

	go func() {
		// Start the watching process.
		if err := w.Start(time.Millisecond * 100); err != nil {
			t.Fatal(err)
		}
	}()

	wg.Wait()
}

// TODO: TestIgnoreFiles
func TestIgnoreFiles(t *testing.T) {}

func TestEventDeleteFile(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	// Add the testDir to the watchlist.
	if err := w.AddRecursive(testDir); err != nil {
		t.Fatal(err)
	}

	files := map[string]bool{
		"file_1.txt": false,
		"file_2.txt": false,
		"file_3.txt": false,
	}

	for f := range files {
		filePath := filepath.Join(testDir, f)
		if err := os.Remove(filePath); err != nil {
			t.Error(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		events := 0
		for {
			select {
			case event := <-w.Event:
				if event.Op != Remove {
					t.Errorf("expected event to be Remove, got %s", event.Op)
				}

				files[event.Name()] = true
				events++

				if events == len(files) {
					return
				}
			case <-time.After(time.Millisecond * 250):
				for f, e := range files {
					if !e {
						t.Errorf("received no event for file %s", f)
					}
				}
				return
			}
		}
	}()

	go func() {
		// Start the watching process.
		if err := w.Start(time.Millisecond * 100); err != nil {
			t.Fatal(err)
		}
	}()

	wg.Wait()
}

func TestEventRenameFile(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	// Add the testDir to the watchlist.
	if err := w.AddRecursive(testDir); err != nil {
		t.Fatal(err)
	}

	// Rename a file.
	if err := os.Rename(
		filepath.Join(testDir, "file.txt"),
		filepath.Join(testDir, "file1.txt"),
	); err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case event := <-w.Event:
			if event.Op != Rename {
				t.Errorf("expected event to be Rename, got %s", event.Op)
			}
		case <-time.After(time.Millisecond * 250):
			t.Fatal("received no rename event")
		}
	}()

	go func() {
		// Start the watching process.
		if err := w.Start(time.Millisecond * 100); err != nil {
			t.Fatal(err)
		}
	}()

	wg.Wait()
}

func TestEventChmodFile(t *testing.T) {
	testDir, teardown := setup(t)
	defer teardown()

	w := New()

	// Add the testDir to the watchlist.
	if err := w.Add(testDir); err != nil {
		t.Fatal(err)
	}

	files := map[string]bool{
		"file_1.txt": false,
		"file_2.txt": false,
		"file_3.txt": false,
	}

	for f := range files {
		filePath := filepath.Join(testDir, f)
		if err := os.Chmod(filePath, os.ModePerm); err != nil {
			t.Error(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		events := 0
		for {
			select {
			case event := <-w.Event:
				if event.Op != Chmod {
					t.Errorf("expected event to be Remove, got %s", event.Op)
				}

				files[event.Name()] = true
				events++

				if events == len(files) {
					return
				}
			case <-time.After(time.Millisecond * 250):
				for f, e := range files {
					if !e {
						t.Errorf("received no event for file %s", f)
					}
				}
				return
			}
		}
	}()

	go func() {
		// Start the watching process.
		if err := w.Start(time.Millisecond * 100); err != nil {
			t.Fatal(err)
		}
	}()

	wg.Wait()
}

func BenchmarkEventRenameFile(b *testing.B) {
	testDir, teardown := setup(b)
	defer teardown()

	w := New()

	// Filter for only rename events.
	w.FilterOps(Rename)

	// Add the testDir to the watchlist.
	if err := w.AddRecursive(testDir); err != nil {
		b.Fatal(err)
	}

	go func() {
		// Start the watching process.
		if err := w.Start(time.Millisecond); err != nil {
			b.Fatal(err)
		}
	}()

	var filenameFrom = filepath.Join(testDir, "file.txt")
	var filenameTo = filepath.Join(testDir, "file1.txt")

	for i := 0; i < b.N; i++ {
		// Rename a file.
		if err := os.Rename(
			filenameFrom,
			filenameTo,
		); err != nil {
			b.Error(err)
		}

		select {
		case event := <-w.Event:
			if event.Op != Rename {
				b.Errorf("expected event to be Rename, got %s", event.Op)
			}
		case <-time.After(time.Millisecond * 250):
			b.Fatal("received no rename event")
		}

		filenameFrom, filenameTo = filenameTo, filenameFrom
	}
}

func BenchmarkListFiles(b *testing.B) {
	testDir, teardown := setup(b)
	defer teardown()

	w := New()
	w.AddRecursive(testDir)

	for i := 0; i < b.N; i++ {
		fileList := w.retrieveFileList()
		if fileList == nil {
			b.Fatal("expected file list to not be empty")
		}
	}
}
