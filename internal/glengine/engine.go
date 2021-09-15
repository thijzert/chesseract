package engine

import (
	"context"
	"io"
	"log"
	"math"
	"sync"
	"time"

	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var PackageVersion string

type RunConfig struct {
	Logger       *log.Logger
	WindowWidth  int
	WindowHeight int
}

func DefaultConfig() RunConfig {
	return RunConfig{
		Logger:       log.Default(),
		WindowWidth:  1280,
		WindowHeight: 720,
	}.mindTheGaps()
}

func (r RunConfig) mindTheGaps() RunConfig {
	if r.Logger == nil {
		r.Logger = log.New(io.Discard, "discarded", 0)
	}

	if r.WindowWidth < 5 {
		r.WindowWidth = 1280
	}
	if r.WindowHeight < 5 {
		r.WindowHeight = 720
	}

	return r
}

func (r RunConfig) NewEngine(ctx context.Context) *Engine {
	r = r.mindTheGaps()
	eng := Engine{
		ctx:       ctx,
		runConfig: r,
		models:    &modelManager{},
		GUI:       &GUIManager{},
	}

	eng.camera.pitch = 0.5
	eng.camera.yaw = -0.15
	eng.camera.radius = 10

	return &eng
}

type Engine struct {
	ctx              context.Context
	runConfig        RunConfig
	window           *glfw.Window
	models           *modelManager
	GUI              *GUIManager
	projectionMatrix mgl32.Mat4
	camera           struct {
		pitch  float64
		yaw    float64
		radius float64
	}
	mu              sync.Mutex
	currentEntities []rawEntity
	nextEntities    []Entity
	entitiesDirty   bool
	Entities        []Entity
	loadedModels    map[string]rawModel
	light           pointLight
}

func (eng *Engine) log(format string, items ...interface{}) {
	eng.runConfig.Logger.Printf(format, items...)
}

func (eng *Engine) Run() error {
	eng.log("Hello, world")

	runtime.LockOSThread()

	var err error
	eng.window, err = eng.initGlfw()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	err = eng.initOpenGL()
	if err != nil {
		return err
	}

	eng.light.position = mgl32.Vec3{9, 3, 0}
	eng.light.colour = mgl32.Vec3{1, 1, 1}

	program, err := eng.loadSimpleProgram("vertexShader", "fragmentShader")
	if err != nil {
		return err
	}

	eng.GUI.Program, err = eng.loadSimpleProgram("guiVertex", "guiFragment")
	if err != nil {
		return err
	}
	eng.GUI.Quad = eng.models.DefaultQuad()

	eng.projectionMatrix = mgl32.Perspective(mgl32.DegToRad(70), float32(eng.runConfig.WindowWidth)/float32(eng.runConfig.WindowHeight), 0.1, 1000)

	frameCap := int64(30) // TODO: make configurable
	frameTime := time.Duration(int64(1*time.Second) / frameCap)
	frameTimeAvg := float64(frameTime)
	frameTime -= frameTime / 100
	frameRobin := 0

	tLast := time.Now()
	for !eng.window.ShouldClose() && eng.ctx.Err() == nil {
		glfw.PollEvents()

		// TODO: decouple from the render loop
		eng.updatePhysics()

		eng.prepareDisplay()

		err := eng.updateEntities()
		if err != nil {
			// TODO: figure out if we can handle this gracefully
			//       (e.g. just not load that particular model)
			return err
		}

		eng.mu.Lock()
		for _, e := range eng.currentEntities {
			eng.drawEntity(e, program)
		}
		eng.mu.Unlock()

		eng.GUI.DrawGUI()

		dFrame := time.Since(tLast)
		if dFrame < frameTime {
			time.Sleep(frameTime - dFrame)
		}
		dFrame = time.Since(tLast)
		frameTimeAvg = 0.75*frameTimeAvg + 0.25*float64(dFrame)

		eng.updateDisplay(eng.window)

		frameRobin = (frameRobin + 1) & 0x1f
		if frameRobin == 0 {
			eng.log("Frame rate: %.1f", float64(time.Second)/frameTimeAvg)
		}
		tLast = time.Now()
	}

	eng.window.Destroy()

	return eng.ctx.Err()
}

func (eng *Engine) ClearEntities() {
	eng.Entities = eng.Entities[:0]
}

func (eng *Engine) SwapEntities() {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	eng.nextEntities = eng.nextEntities[:0]
	eng.nextEntities = append(eng.nextEntities, eng.Entities...)

	eng.entitiesDirty = true
}

func (eng *Engine) updateEntities() error {
	if !eng.entitiesDirty {
		return nil
	}

	eng.mu.Lock()
	defer eng.mu.Unlock()

	// Load all models into the cache
	if eng.loadedModels == nil {
		eng.loadedModels = make(map[string]rawModel)
	}
	for _, ent := range eng.nextEntities {
		if _, ok := eng.loadedModels[ent.ModelName]; !ok {
			model, err := eng.models.LoadModelAsset(ent.ModelName)
			if err != nil {
				return err
			}
			model.Material = eng.loadMaterialAsset(ent.ModelName)

			eng.loadedModels[ent.ModelName] = model
		}
	}

	// Clear the currentEntities list and fill it without reallocating every time
	eng.currentEntities = eng.currentEntities[:0]
	for _, ent := range eng.nextEntities {
		model := eng.loadedModels[ent.ModelName]
		entity := rawEntity{
			model:     model,
			position:  ent.Position,
			rotation:  ent.Rotation,
			scale:     ent.Scale,
			tileIndex: ent.TileIndex,
		}
		eng.currentEntities = append(eng.currentEntities, entity)
	}

	eng.entitiesDirty = false

	return nil
}

func (eng *Engine) updatePhysics() {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	if eng.window.GetKey(glfw.KeyW) == glfw.Press {
		eng.camera.pitch += 0.05
	}
	if eng.window.GetKey(glfw.KeyS) == glfw.Press {
		eng.camera.pitch -= 0.05
	}
	if eng.window.GetKey(glfw.KeyA) == glfw.Press {
		eng.camera.yaw -= 0.05
	}
	if eng.window.GetKey(glfw.KeyD) == glfw.Press {
		eng.camera.yaw += 0.05
	}
	if eng.window.GetKey(glfw.KeyZ) == glfw.Press {
		eng.camera.radius -= 0.35
	}
	if eng.window.GetKey(glfw.KeyX) == glfw.Press {
		eng.camera.radius += 0.35
	}

	eng.camera.pitch = clamp(eng.camera.pitch, -0.8, 1.3)
	eng.camera.yaw = math.Mod(eng.camera.yaw, math.Pi)
	eng.camera.radius = clamp(eng.camera.radius, 4, 100)
}

func clamp(val, min, max float64) float64 {
	if val < min {
		val = min
	}
	if val > max {
		val = max
	}
	return val
}

func (eng *Engine) getCameraMatrix() mgl32.Mat4 {
	y := math.Sin(eng.camera.pitch) * eng.camera.radius
	rxz := math.Cos(eng.camera.pitch) * eng.camera.radius

	x := math.Sin(eng.camera.yaw) * rxz
	z := math.Cos(eng.camera.yaw) * rxz

	position := mgl32.Vec3{float32(x), float32(y), float32(z)}
	lookingAt := mgl32.Vec3{0, 1, 0}
	up := mgl32.Vec3{0, 1, 0}
	return mgl32.LookAtV(position, lookingAt, up)
}

func (eng *Engine) Shutdown() error {
	eng.log("Engine shutdown started")
	eng.models.Cleanup()
	return nil
}

// initGlfw initialises glfw and returns a window we can use
func (eng *Engine) initGlfw() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	//glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(eng.runConfig.WindowWidth, eng.runConfig.WindowHeight, "Chesseract", nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()

	return window, nil
}

// initOpenGL initialises OpenGL, and creates an initialised program
func (eng *Engine) initOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	eng.log("OpenGL version %s", version)

	return nil
}

func (eng *Engine) prepareDisplay() error {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.15, 0.05, 0.1, 1)

	return nil
}

func (eng *Engine) updateDisplay(window *glfw.Window) error {
	window.SwapBuffers()

	return nil
}

func (eng *Engine) drawEntity(e rawEntity, program glProgram) {
	if e.model.Material.Transparency {
		gl.Disable(gl.CULL_FACE)
	} else {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	}

	program.Start()

	if e.model.Material.Texture != 0 {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, uint32(e.model.Material.Texture))
	}
	if e.model.Material.NormalMap != 0 {
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, uint32(e.model.Material.NormalMap))
	}
	if e.model.Material.SpecularMap != 0 {
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, uint32(e.model.Material.SpecularMap))
	}

	program.UniformFloat(U_MATERIAL_SHINE_DAMPER, e.model.Material.ShineDamper)
	program.UniformFloat(U_MATERIAL_REFLECTIVITY, e.model.Material.Reflectivity)
	program.UniformVec2(U_TILE_SIZE, float32(e.model.Material.TileSize[0]), float32(e.model.Material.TileSize[1]))

	program.UniformMatrix4(U_PROJECTION, eng.projectionMatrix)
	program.UniformMatrix4(U_CAMERA, eng.getCameraMatrix())
	program.UniformMatrix4(U_TRANSFORM, e.getTransformation())
	program.UniformFloat(U_TILE_INDEX, float32(e.tileIndex))

	program.UniformVec3(U_TMP_LIGHT_POS, eng.light.position)
	program.UniformVec3(U_TMP_LIGHT_COLOUR, eng.light.colour)

	gl.BindVertexArray(e.model.VAO)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	gl.DrawElements(gl.TRIANGLES, e.model.NVertices, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)
	gl.BindVertexArray(0)

	program.Stop()
}
