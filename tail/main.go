package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	var name, color bool
	flag.BoolVar(&name, "n", false, "show file names")
	flag.BoolVar(&color, "c", false, "colorize")
	flag.Parse()

	args := flag.Args()

	n, err := NewFileWatcher()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	defer n.Close()

	maxlen := 0
	if name {
		for _, a := range args {
			if maxlen < len(a) {
				maxlen = len(a)
			}
		}
	}

	for i, arg := range args {
		var wr io.Writer

		wr = os.Stdout

		if color {
			wr = &DecorateWriter{
				writer: wr,
				prefix: fmt.Sprintf("\033[%dm", i%7+31),
				suffix: "\033[0m",
			}
		}

		if name {
			wr = &DecorateWriter{
				writer: wr,
				prefix: fmt.Sprintf(fmt.Sprintf("%%-%ds ", maxlen), arg),
				suffix: "",
			}
		}

		t, err := NewTailFile(arg, wr)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		defer t.Close()

		if err := n.Add(t); err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
	}

	n.Loop()
}

type DecorateWriter struct {
	writer io.Writer
	prefix string
	suffix string
}

func (w *DecorateWriter) Write(p []byte) (int, error) {

	data := make([]byte, len(p)+len(w.prefix)+len(w.suffix))
	data = append(data, []byte(w.prefix)...)
	data = append(data, p...)
	data = append(data, []byte(w.suffix)...)

	return w.writer.Write(data)
}

type Watcher interface {
	Close()
	Add(file string) error
	Events() chan fsnotify.Event
	Errors() chan error
}

type FileWatcher struct {
	watcher *fsnotify.Watcher
	files   []*TailFile
}

func NewFileWatcher() (*FileWatcher, error) {
	n := &FileWatcher{
		files: make([]*TailFile, 0),
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	n.watcher = w

	return n, nil
}

func (n *FileWatcher) Close() {
	if n.watcher != nil {
		n.watcher.Close()
		n.watcher = nil
	}
}

func (n *FileWatcher) Add(file *TailFile) error {
	fullpath, err := filepath.Abs(file.path)
	if err != nil {
		return nil
	}
	dir := filepath.Dir(fullpath)
	err = n.watcher.Add(dir)
	if err != nil {
		return nil
	}
	n.files = append(n.files, file)
	return nil
}

func (n *FileWatcher) Events() chan fsnotify.Event {
	return n.watcher.Events
}

func (n *FileWatcher) Errors() chan error {
	return n.watcher.Errors
}

func (n *FileWatcher) Loop() {

	for {
		select {
		case event, ok := <-n.Events():
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				for _, t := range n.files {
					if event.Name == t.path {
						t.OnCreate()
					}
				}
			} else if event.Op&fsnotify.Write == fsnotify.Write {
				for _, t := range n.files {
					if event.Name == t.path {
						t.OnWrite()
					}
				}
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				for _, t := range n.files {
					if event.Name == t.path {
						t.OnRemove()
					}
				}
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				for _, t := range n.files {
					if event.Name == t.path {
						t.OnRename()
					}
				}
			} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				for _, t := range n.files {
					if event.Name == t.path {
						t.OnChmod()
					}
				}
			}

		case err, ok := <-n.Errors():
			if !ok {
				return
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}

}

type TailFile struct {
	path   string
	writer io.Writer
	file   *os.File
	reader *bufio.Reader
}

func NewTailFile(path string, writer io.Writer) (*TailFile, error) {

	fullpath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		return nil, err
	}

	t := &TailFile{
		path:   fullpath,
		writer: writer,
	}

	f, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	if _, err = f.Seek(0, 2); err != nil {
		f.Close()
		return nil, err
	}
	t.file = f

	t.reader = bufio.NewReader(f)

	return t, nil
}

func (t *TailFile) Close() {
	if t.file != nil {
		t.file.Close()
		t.file = nil
	}
}

func (t *TailFile) OnCreate() error {
	return nil
}

func (t *TailFile) OnWrite() error {
	for {
		b, err := t.reader.ReadBytes(byte('\n'))
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		t.writer.Write(b)
	}

}

func (t *TailFile) OnRemove() error {
	return nil
}

func (t *TailFile) OnRename() error {
	return nil
}

func (t *TailFile) OnChmod() error {
	if t.file != nil {
		t.file.Close()
		t.file = nil
	}

	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	if _, err = f.Seek(0, 2); err != nil {
		f.Close()
		return err
	}
	t.file = f

	t.reader = bufio.NewReader(f)
	return nil
}
