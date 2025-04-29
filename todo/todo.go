package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Options map[optFlag]bool

type optFlag string

const (
	Dtflg optFlag = "dt"
	Ucflg optFlag = "uc"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}

func (l *List) String() string {
	formatted := ""

	for k, t := range *l {
		prefix := " "
		if t.Done {
			prefix = "X "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}

	return formatted
}

func (l *List) List(opts Options) {
	fmt.Print(opts)

	formatted := ""
	switch {
	case opts[Dtflg]:
		s := addOptsToFormatted(formatted, l, Dtflg)
		fmt.Print(s)
	case opts[Ucflg]:
		s := addOptsToFormatted(formatted, l, Ucflg)
		fmt.Print(s)

	default:
		fmt.Print(l)
	}
}

func addOptsToFormatted(formatted string, l *List, opt optFlag) string {
	if opt == Ucflg {
		for k, t := range *l {
			formatted += fmt.Sprintf("%d: %s\n", k+1, t.Task)
		}
	}

	return formatted
}
