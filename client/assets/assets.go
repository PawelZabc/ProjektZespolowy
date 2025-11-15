package assets

import (
	"fmt"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:generate go run ../utils/assetgen/main.go

const assetsBasePath = "client/assets"

var Manager *AssetManager

type AssetManager struct {
	Models   map[string]rl.Model
	Textures map[string]rl.Texture2D
	Sounds   map[string]rl.Sound
	Music    map[string]rl.Music
	Fonts    map[string]rl.Font
	Shaders  map[string]rl.Shader
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		Models:   make(map[string]rl.Model),
		Textures: make(map[string]rl.Texture2D),
		Sounds:   make(map[string]rl.Sound),
		Music:    make(map[string]rl.Music),
		Fonts:    make(map[string]rl.Font),
		Shaders:  make(map[string]rl.Shader),
	}
}

func (am *AssetManager) LoadModel(filename string) (rl.Model, error) {
	// cache
	if model, exists := am.Models[filename]; exists {
		return model, nil
	}

	path := filepath.Join(assetsBasePath, "models", filename)
	model := rl.LoadModel(path)

	if model.MeshCount == 0 {
		return rl.Model{}, fmt.Errorf("failed to load model: %s", filename)
	}

	am.Models[filename] = model
	return model, nil
}

func (am *AssetManager) LoadTexture(filename string) (rl.Texture2D, error) {

	if texture, exists := am.Textures[filename]; exists {
		return texture, nil
	}

	path := filepath.Join(assetsBasePath, "graphics", filename)
	texture := rl.LoadTexture(path)

	if texture.ID == 0 {
		return rl.Texture2D{}, fmt.Errorf("failed to load texture: %s", filename)
	}

	am.Textures[filename] = texture
	return texture, nil
}

// LoadSound loads a sound from disk
func (am *AssetManager) LoadSound(filename string) (rl.Sound, error) {
	// Check cache
	if sound, exists := am.Sounds[filename]; exists {
		return sound, nil
	}

	path := filepath.Join(assetsBasePath, "audio", filename)
	sound := rl.LoadSound(path)

	if sound.Stream.Buffer == nil {
		return rl.Sound{}, fmt.Errorf("failed to load sound: %s", filename)
	}

	am.Sounds[filename] = sound
	return sound, nil
}

func (am *AssetManager) LoadMusic(filename string) (rl.Music, error) {
	if music, exists := am.Music[filename]; exists {
		return music, nil
	}

	path := filepath.Join(assetsBasePath, "audio", filename)
	music := rl.LoadMusicStream(path)

	if music.Stream.Buffer == nil {
		return rl.Music{}, fmt.Errorf("failed to load music: %s", filename)
	}

	am.Music[filename] = music
	return music, nil
}

// LoadFont loads a font from disk
func (am *AssetManager) LoadFont(filename string, fontSize int32) (rl.Font, error) {
	cacheKey := fmt.Sprintf("%s_%d", filename, fontSize)

	// Check cache
	if font, exists := am.Fonts[cacheKey]; exists {
		return font, nil
	}

	path := filepath.Join(assetsBasePath, "fonts", filename)
	font := rl.LoadFontEx(path, fontSize, nil, 0)

	if font.Texture.ID == 0 {
		return rl.Font{}, fmt.Errorf("failed to load font: %s", filename)
	}

	am.Fonts[cacheKey] = font
	return font, nil
}

// LoadShader loads a shader from disk
func (am *AssetManager) LoadShader(vsFilename, fsFilename string) (rl.Shader, error) {
	cacheKey := vsFilename + "|" + fsFilename

	// Check cache
	if shader, exists := am.Shaders[cacheKey]; exists {
		return shader, nil
	}

	vsPath := filepath.Join(assetsBasePath, "shaders", vsFilename)
	fsPath := filepath.Join(assetsBasePath, "shaders", fsFilename)

	shader := rl.LoadShader(vsPath, fsPath)

	if shader.ID == 0 {
		return rl.Shader{}, fmt.Errorf("failed to load shader: %s, %s", vsFilename, fsFilename)
	}

	am.Shaders[cacheKey] = shader
	return shader, nil
}

// GetModel returns a cached model or error if not loaded
func (am *AssetManager) GetModel(filename string) (rl.Model, bool) {
	model, exists := am.Models[filename]
	return model, exists
}

// GetTexture returns a cached texture or error if not loaded
func (am *AssetManager) GetTexture(filename string) (rl.Texture2D, bool) {
	texture, exists := am.Textures[filename]
	return texture, exists
}

// GetSound returns a cached sound or error if not loaded
func (am *AssetManager) GetSound(filename string) (rl.Sound, bool) {
	sound, exists := am.Sounds[filename]
	return sound, exists
}

// Unload frees all loaded assets
func (am *AssetManager) Unload() {
	for _, model := range am.Models {
		rl.UnloadModel(model)
	}
	for _, texture := range am.Textures {
		rl.UnloadTexture(texture)
	}
	for _, sound := range am.Sounds {
		rl.UnloadSound(sound)
	}
	for _, music := range am.Music {
		rl.UnloadMusicStream(music)
	}
	for _, font := range am.Fonts {
		rl.UnloadFont(font)
	}
	for _, shader := range am.Shaders {
		rl.UnloadShader(shader)
	}

	// Clear maps
	am.Models = make(map[string]rl.Model)
	am.Textures = make(map[string]rl.Texture2D)
	am.Sounds = make(map[string]rl.Sound)
	am.Music = make(map[string]rl.Music)
	am.Fonts = make(map[string]rl.Font)
	am.Shaders = make(map[string]rl.Shader)
}

// UnloadModel unloads a specific model from cache
func (am *AssetManager) UnloadModel(filename string) {
	if model, exists := am.Models[filename]; exists {
		rl.UnloadModel(model)
		delete(am.Models, filename)
	}
}

// UnloadTexture unloads a specific texture from cache
func (am *AssetManager) UnloadTexture(filename string) {
	if texture, exists := am.Textures[filename]; exists {
		rl.UnloadTexture(texture)
		delete(am.Textures, filename)
	}
}

// Helper function to get file extension
func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return strings.ToLower(ext)
}

// PreloadAssets loads multiple assets at once with error reporting
func (am *AssetManager) PreloadAssets(
	models []string,
	textures []string,
	sounds []string,
) []error {
	var errors []error

	for _, filename := range models {
		if _, err := am.LoadModel(filename); err != nil {
			errors = append(errors, fmt.Errorf("model %s: %w", filename, err))
		}
	}

	for _, filename := range textures {
		if _, err := am.LoadTexture(filename); err != nil {
			errors = append(errors, fmt.Errorf("texture %s: %w", filename, err))
		}
	}

	for _, filename := range sounds {
		if _, err := am.LoadSound(filename); err != nil {
			errors = append(errors, fmt.Errorf("sound %s: %w", filename, err))
		}
	}

	return errors
}

// =============================================================================
// Asset Groups for organized loading
// =============================================================================

// AssetGroup represents a named group of assets
type AssetGroup struct {
	Name     string
	Models   []string
	Textures []string
	Sounds   []string
	Music    []string
}

// LoadGroup loads all assets in a group
func (am *AssetManager) LoadGroup(group AssetGroup) []error {
	var errors []error

	for _, filename := range group.Models {
		if _, err := am.LoadModel(filename); err != nil {
			errors = append(errors, err)
		}
	}

	for _, filename := range group.Textures {
		if _, err := am.LoadTexture(filename); err != nil {
			errors = append(errors, err)
		}
	}

	for _, filename := range group.Sounds {
		if _, err := am.LoadSound(filename); err != nil {
			errors = append(errors, err)
		}
	}

	for _, filename := range group.Music {
		if _, err := am.LoadMusic(filename); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// UnloadGroup unloads all assets in a group
func (am *AssetManager) UnloadGroup(group AssetGroup) {
	for _, filename := range group.Models {
		am.UnloadModel(filename)
	}

	for _, filename := range group.Textures {
		am.UnloadTexture(filename)
	}

	for _, filename := range group.Sounds {
		if sound, exists := am.Sounds[filename]; exists {
			rl.UnloadSound(sound)
			delete(am.Sounds, filename)
		}
	}

	for _, filename := range group.Music {
		if music, exists := am.Music[filename]; exists {
			rl.UnloadMusicStream(music)
			delete(am.Music, filename)
		}
	}
}
