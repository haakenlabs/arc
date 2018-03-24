/*
Copyright (c) 2018 HaakenLabs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package scene

type Scene struct {
	LoadFunc         func() error
	OnActivateFunc   func()
	OnDeacticateFunc func()

	environment *Environment
	graph       *Graph
	cameras     []*Camera
	name        string
	loaded      bool
}

// Name returns the name of this scene.
func (s *Scene) Name() string {
	return s.name
}

// OnActivate is called when the scene transitions to the active state.
func (s *Scene) OnActivate() {
	if s.OnActivateFunc != nil {
		s.OnActivateFunc()
	}
}

// OnDeactivate is called when the scene transitions to the inactive state.
func (s *Scene) OnDeactivate() {
	if s.OnDeacticateFunc != nil {
		s.OnDeacticateFunc()
	}
}

// Load is called when the scene is being initialized.
func (s *Scene) Load() error {
	if s.loaded {
		return nil
	}

	s.graph = NewGraph(s)
	s.environment = NewEnvironment()

	if s.LoadFunc != nil {
		s.LoadFunc()
	}

	s.loaded = true

	return nil
}

// Loaded reports if the scene has been loaded.
func (s *Scene) Loaded() bool {
	return s.loaded
}

func (s *Scene) OnSceneGraphUpdate() {
	if s.graph == nil {
		return
	}

	s.cameras = s.cameras[:0]

	// Update renderer cache.
	components := s.graph.Components()
	for i := range components {
		if c, ok := components[i].(*Camera); ok {
			s.cameras = append(s.cameras, c)
		}
	}
}

func (s *Scene) Objects() []*GameObject {
	return s.graph.aCache
}

func (s *Scene) Components() []Component {
	return s.graph.cCache
}

func NewScene(name string) *Scene {
	s := &Scene{
		name: name,
	}

	return s
}
