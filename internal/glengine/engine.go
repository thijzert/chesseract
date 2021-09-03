package engine

import (
	"context"
	"io"
	"log"

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

	eng.camera.viewDirection = mgl32.Vec3{0, 0, -1}

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
		position      mgl32.Vec3
		viewDirection mgl32.Vec3
	}
	entities []entity
	light    pointLight
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

	chessSet := []string{
		"pawn",
		"rook",
		"knight",
		"bishop",
		"queen",
		"king",
	}

	for i, name := range chessSet {
		pawn, err := eng.models.LoadModelAsset(name)
		if err != nil {
			return err
		}
		pawn.Material = eng.loadMaterialAsset(name)

		x := 1 * (float32(i) - 0.5*float32(len(chessSet)))
		eng.entities = append(eng.entities, entity{
			model:    pawn,
			position: mgl32.Vec3{x, -2, -5},
			scale:    mgl32.Vec3{1, 1, 1},
		})
	}

	eng.projectionMatrix = mgl32.Perspective(mgl32.DegToRad(70), float32(eng.runConfig.WindowWidth)/float32(eng.runConfig.WindowHeight), 0.1, 1000)

	for !eng.window.ShouldClose() && eng.ctx.Err() == nil {
		eng.prepareDisplay()

		for _, e := range eng.entities {
			eng.drawEntity(e, program)
		}

		eng.GUI.DrawGUI()

		glfw.PollEvents()
		eng.updateDisplay(eng.window)

		// TODO: decouple from the render loop
		eng.updatePhysics()
	}

	eng.window.Destroy()

	return eng.ctx.Err()
}

func (eng *Engine) updatePhysics() {
	eng.entities[1].position[0] -= 0.002
	eng.entities[1].rotation[2] += 0.1
	eng.entities[0].position[0] += 0.005
	eng.entities[0].position[2] -= 0.01
	eng.entities[0].rotation[1] += 0.1
	eng.entities[2].rotation[1] -= 0.05

	if eng.window.GetKey(glfw.KeyW) == glfw.Press {
		eng.camera.position[2] -= 0.2
	}
	if eng.window.GetKey(glfw.KeyA) == glfw.Press {
		eng.camera.position[0] -= 0.2
	}
	if eng.window.GetKey(glfw.KeyS) == glfw.Press {
		eng.camera.position[2] += 0.2
	}
	if eng.window.GetKey(glfw.KeyD) == glfw.Press {
		eng.camera.position[0] += 0.2
	}
}

func (eng *Engine) getCameraMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(eng.camera.position, eng.camera.position.Add(eng.camera.viewDirection), mgl32.Vec3{0, 1, 0})
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

func (eng *Engine) drawEntity(e entity, program glProgram) {
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

	program.UniformMatrix4(U_PROJECTION, eng.projectionMatrix)
	program.UniformMatrix4(U_CAMERA, eng.getCameraMatrix())
	program.UniformMatrix4(U_TRANSFORM, e.getTransformation())

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
