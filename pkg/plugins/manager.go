package plugins

import (
	"fmt"
	"sync"
)

// Manager manages plugin lifecycle and operations
type Manager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

// Register registers a plugin with the manager
func (m *Manager) Register(plugin Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := plugin.Name()
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already registered", name)
	}

	m.plugins[name] = plugin
	return nil
}

// Unregister removes a plugin from the manager
func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[name]; !exists {
		return fmt.Errorf("plugin %s is not registered", name)
	}

	delete(m.plugins, name)
	return nil
}

// Get retrieves a plugin by name
func (m *Manager) Get(name string) (Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List returns all registered plugins
func (m *Manager) List() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

// ListNames returns names of all registered plugins
func (m *Manager) ListNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}

	return names
}

// Count returns the number of registered plugins
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.plugins)
}

// Exists checks if a plugin is registered
func (m *Manager) Exists(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.plugins[name]
	return exists
}

// Clear removes all plugins
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.plugins = make(map[string]Plugin)
}

// RunAnalysis runs analysis on all registered plugins
func (m *Manager) RunAnalysis(resource interface{}) ([]Issue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var allIssues []Issue

	for _, plugin := range m.plugins {
		issues, err := plugin.Analyze(resource)
		if err != nil {
			// Log error but continue with other plugins
			fmt.Printf("Plugin %s analysis failed: %v\n", plugin.Name(), err)
			continue
		}

		allIssues = append(allIssues, issues...)
	}

	return allIssues, nil
}

// Install installs a plugin (alias for Register for compatibility)
func (m *Manager) Install(plugin Plugin) error {
	return m.Register(plugin)
}

// InstallByName installs a plugin by name
func (m *Manager) InstallByName(name string) error {
	// Create plugin based on name
	var plugin Plugin

	switch name {
	case "example-detector", "example":
		plugin = NewExamplePlugin()
	default:
		return fmt.Errorf("unknown plugin: %s (available: example-detector)", name)
	}

	return m.Install(plugin)
}

// Uninstall uninstalls a plugin (alias for Unregister for compatibility)
func (m *Manager) Uninstall(name string) error {
	return m.Unregister(name)
}
