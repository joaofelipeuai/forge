package forge

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// HotReload represents the hot reload functionality
type HotReload struct {
	enabled     bool
	watchDirs   []string
	extensions  []string
	lastModTime map[string]time.Time
	mu          sync.RWMutex
	onChange    func()
	debounce    time.Duration
	lastChange  time.Time
}

// NewHotReload creates a new hot reload instance
func NewHotReload() *HotReload {
	return &HotReload{
		enabled:     false,
		watchDirs:   []string{"."},
		extensions:  []string{".go", ".html", ".css", ".js", ".json"},
		lastModTime: make(map[string]time.Time),
		debounce:    500 * time.Millisecond,
	}
}

// Enable enables hot reload
func (hr *HotReload) Enable() {
	hr.enabled = true
}

// Disable disables hot reload
func (hr *HotReload) Disable() {
	hr.enabled = false
}

// AddWatchDir adds a directory to watch
func (hr *HotReload) AddWatchDir(dir string) {
	hr.watchDirs = append(hr.watchDirs, dir)
}

// AddExtension adds a file extension to watch
func (hr *HotReload) AddExtension(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	hr.extensions = append(hr.extensions, ext)
}

// SetOnChange sets the callback function for file changes
func (hr *HotReload) SetOnChange(fn func()) {
	hr.onChange = fn
}

// SetDebounce sets the debounce duration
func (hr *HotReload) SetDebounce(duration time.Duration) {
	hr.debounce = duration
}

// Start starts the hot reload watcher
func (hr *HotReload) Start() {
	if !hr.enabled {
		return
	}

	log.Println("🔥 Hot reload enabled - watching for file changes...")

	// Initial scan
	hr.scanFiles()

	// Start watching
	go hr.watch()
}

// watch continuously watches for file changes
func (hr *HotReload) watch() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if hr.enabled {
			hr.checkForChanges()
		}
	}
}

// scanFiles scans all files and records their modification times
func (hr *HotReload) scanFiles() {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for _, dir := range hr.watchDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if hr.shouldWatch(path) {
				hr.lastModTime[path] = info.ModTime()
			}

			return nil
		})
	}
}

// checkForChanges checks if any watched files have changed
func (hr *HotReload) checkForChanges() {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	changed := false
	now := time.Now()

	for _, dir := range hr.watchDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if hr.shouldWatch(path) {
				lastMod, exists := hr.lastModTime[path]
				if !exists || info.ModTime().After(lastMod) {
					hr.lastModTime[path] = info.ModTime()
					if exists { // Don't trigger on first scan
						changed = true
						log.Printf("🔄 File changed: %s", path)
					}
				}
			}

			return nil
		})
	}

	// Debounce changes
	if changed && now.Sub(hr.lastChange) > hr.debounce {
		hr.lastChange = now
		if hr.onChange != nil {
			go hr.onChange()
		}
	}
}

// shouldWatch checks if a file should be watched based on its extension
func (hr *HotReload) shouldWatch(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, watchExt := range hr.extensions {
		if ext == strings.ToLower(watchExt) {
			return true
		}
	}
	return false
}

// Hot reload middleware for development
func HotReloadMiddleware(hr *HotReload) MiddlewareFunc {
	return func(c *Context) error {
		if hr.enabled {
			// Add hot reload script to HTML responses
			c.Response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Response.Header().Set("Pragma", "no-cache")
			c.Response.Header().Set("Expires", "0")
		}
		return c.Next()
	}
}

// Development server with hot reload
func (f *Forge) ListenWithHotReload(addr string, watchDirs ...string) error {
	hr := NewHotReload()
	hr.Enable()
	
	// Add watch directories
	if len(watchDirs) > 0 {
		hr.watchDirs = watchDirs
	}
	
	// Set up hot reload callback
	hr.SetOnChange(func() {
		log.Println("🔄 Files changed - restart recommended")
		// In a real implementation, this could trigger a server restart
		// For now, we just log the change
	})
	
	// Add hot reload middleware
	f.Use(HotReloadMiddleware(hr))
	
	// Start hot reload watcher
	hr.Start()
	
	fmt.Printf("🔨 Forge v%s server running on %s (Hot Reload: ON)\n", Version, addr)
	return f.Listen(addr)
}

// Hot reload for templates
func (te *TemplateEngine) EnableHotReload() {
	hr := NewHotReload()
	hr.Enable()
	hr.AddWatchDir(te.baseDir)
	hr.AddExtension(te.extension)
	
	hr.SetOnChange(func() {
		log.Println("🔄 Templates changed - reloading...")
		if err := te.LoadTemplates(); err != nil {
			log.Printf("❌ Error reloading templates: %v", err)
		} else {
			log.Println("✅ Templates reloaded successfully")
		}
	})
	
	hr.Start()
}

// Static file hot reload
func StaticHotReload(staticDir string) MiddlewareFunc {
	hr := NewHotReload()
	hr.Enable()
	hr.AddWatchDir(staticDir)
	hr.AddExtension(".css")
	hr.AddExtension(".js")
	hr.AddExtension(".html")
	
	lastChange := time.Now()
	
	hr.SetOnChange(func() {
		lastChange = time.Now()
		log.Printf("🔄 Static files changed: %s", staticDir)
	})
	
	hr.Start()
	
	return func(c *Context) error {
		// Add timestamp to static file URLs for cache busting
		c.Set("static_version", lastChange.Unix())
		return c.Next()
	}
}