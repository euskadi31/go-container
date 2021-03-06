// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package service

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyService struct {
	Name string
}

type Config struct {
}

func TestContainer(t *testing.T) {
	c := New()

	assert.False(t, c.Has("test.bad.service.name"))

	assert.Equal(t, []string{}, c.GetKeys())

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	c.Extend("my.service", func(s *MyService, c Container) *MyService {
		s.Name = "My Service"

		return s
	})

	assert.True(t, c.Has("my.service"))

	assert.Equal(t, []string{"my.service"}, c.GetKeys())

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	myService1 := c.Get("my.service").(*MyService)

	myService2 := c.Get("my.service").(*MyService)

	assert.Equal(t, myService1, myService2)

	assert.Equal(t, "My Service", myService1.Name)

	assert.Panics(t, func() {
		c.Set("my.service", func(c Container) interface{} {
			return &MyService{}
		})
	})

	assert.Panics(t, func() {
		c.Extend("my.service", func(s *MyService, c Container) *MyService {
			s.Name = "My Service 2"

			return s
		})
	})

	assert.Panics(t, func() {
		c.Extend("not.exists.service", func(s *MyService, c Container) *MyService {
			s.Name = "My Service 3"

			return s
		})
	})

	assert.Panics(t, func() {
		c.Get("test.bad.service.name")
	})

	var myService3 *MyService

	c.Fill("my.service", &myService3)

	assert.Equal(t, myService2, myService3)

	assert.Panics(t, func() {
		var bad string

		c.Fill("my.service", &bad)
	})

	assert.False(t, c.Has("my.static.value"))

	c.SetValue("my.static.value", "bar")

	assert.Equal(t, "bar", c.Get("my.static.value"))

	assert.True(t, c.Has("my.static.value"))

	assert.Panics(t, func() {
		c.SetValue("my.static.value", "bar")
	})
}

func TestContainerWithExtendBug(t *testing.T) {
	c := New()

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	c.Set("config", func(c Container) interface{} {
		cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		cmd.String("config", "", "")
		// nolint:gosec
		_ = cmd.Parse(os.Args[1:])

		return cmd
	})

	c.Extend("my.service", func(s *MyService, c Container) *MyService {
		_ = c.Get("config")

		s.Name = "My Service"

		return s
	})

	assert.True(t, c.Has("my.service"))

	_ = c.Get("my.service").(*MyService)

}

func BenchmarkContainerGet(b *testing.B) {
	c := New()

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	for n := 0; n < b.N; n++ {
		c.Get("my.service")
	}
}

func BenchmarkContainerFill(b *testing.B) {
	c := New()

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	for n := 0; n < b.N; n++ {
		var myService3 *MyService

		c.Fill("my.service", &myService3)
	}
}

func BenchmarkContainerGetPreInit(b *testing.B) {
	c := New()

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	c.Get("my.service")

	for n := 0; n < b.N; n++ {
		c.Get("my.service")
	}
}

func BenchmarkContainerGetWithExtend(b *testing.B) {
	c := New()

	c.Set("my.service", func(c Container) interface{} {
		return &MyService{}
	})

	c.Set("config", func(c Container) interface{} {
		return &Config{}
	})

	c.Extend("my.service", func(s *MyService, c Container) *MyService {
		_ = c.Get("config")

		s.Name = "My Service"

		return s
	})

	for n := 0; n < b.N; n++ {
		c.Get("my.service")
	}
}
