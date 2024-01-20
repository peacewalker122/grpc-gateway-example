package main

import (
	"context"
	"errors"
	"sync"
)

type Database interface {
	Create(ctx context.Context, val string) error
	Find(ctx context.Context, index int) (string, error)
	FindAll(ctx context.Context) ([]string, error)
	Update(ctx context.Context, kindex int, val string) error
	Delete(ctx context.Context, index int) error
}

type database struct {
	data []string

	rw sync.RWMutex
}

func NewDatabase() *database {
	return &database{
		data: make([]string, 0),
		rw:   sync.RWMutex{},
	}
}

func (d *database) FindAll(ctx context.Context) ([]string, error) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	res := make([]string, len(d.data))

	for i, v := range d.data {
		res[i] = v
	}

	return res, nil
}

func (d *database) Create(ctx context.Context, val string) error {
	d.rw.Lock()
	defer d.rw.Unlock()

	d.data = append(d.data, val)

	return nil
}

func (d *database) Find(ctx context.Context, key int) (string, error) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	if key > len(d.data)-1 {
		return "", errors.New("index out of range")
	}

	return d.data[key], errors.New("not found")
}

func (d *database) Update(ctx context.Context, index int, val string) error {
	d.rw.Lock()
	defer d.rw.Unlock()

	if index > len(d.data)-1 {
		return errors.New("index out of range")
	}

	d.data[index] = val

	return nil
}

func (d *database) Delete(ctx context.Context, index int) error {
	d.rw.Lock()
	defer d.rw.Unlock()

	if index > len(d.data)-1 {
		return errors.New("index out of range")
	}

	d.data = append(d.data[:index], d.data[index+1:]...)

	return nil
}

var _ Database = (*database)(nil)
