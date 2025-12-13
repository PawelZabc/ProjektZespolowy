package entities

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightType int32

const (
	LightTypeDirectional LightType = iota
	LightTypePoint
)

type Light struct {
	shader    rl.Shader
	lightType LightType
	Position  rl.Vector3
	target    rl.Vector3
	color     rl.Color
	Enabled   int32
	// shader locations
	enabledLoc int32
	typeLoc    int32
	posLoc     int32
	targetLoc  int32
	colorLoc   int32
}

const maxLightsCount = 4

var lightCount = 0

func NewLight(
	lightType LightType,
	position, target rl.Vector3,
	color rl.Color,
	shader rl.Shader) Light {
	light := Light{
		shader: shader,
	}
	if lightCount < maxLightsCount {
		light.Enabled = 1
		light.lightType = lightType
		light.Position = position
		light.target = target
		light.color = color
		light.enabledLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].enabled", lightCount))
		light.typeLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].type", lightCount))
		light.posLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].position", lightCount))
		light.targetLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].target", lightCount))
		light.colorLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].color", lightCount))
		light.UpdateValues()
		lightCount++
	}
	return light
}

func (lt *Light) UpdateValues() {
	// Send to shader light enabled state and type
	rl.SetShaderValue(lt.shader, lt.enabledLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.Enabled)), 4), rl.ShaderUniformInt)
	rl.SetShaderValue(lt.shader, lt.typeLoc, unsafe.Slice((*float32)(unsafe.Pointer(&lt.lightType)), 4), rl.ShaderUniformInt)

	// Send to shader light position values
	rl.SetShaderValue(lt.shader, lt.posLoc, []float32{lt.Position.X, lt.Position.Y, lt.Position.Z}, rl.ShaderUniformVec3)

	// Send to shader light target target values
	rl.SetShaderValue(lt.shader, lt.targetLoc, []float32{lt.target.X, lt.target.Y, lt.target.Z}, rl.ShaderUniformVec3)

	// Send to shader light color values
	rl.SetShaderValue(lt.shader, lt.colorLoc,
		[]float32{float32(lt.color.R) / 255, float32(lt.color.G) / 255, float32(lt.color.B) / 255, float32(lt.color.A) / 255},
		rl.ShaderUniformVec4)
}
