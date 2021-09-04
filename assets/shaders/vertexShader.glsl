#version 410

in vec3 position;
in vec2 uv_position;
in vec3 normal;
in vec3 tangent;

uniform mat4 cameraMatrix;
uniform mat4 projectionMatrix;
uniform mat4 transformationMatrix;

uniform vec3 lightPosition;

out vec2 pass_uv;
out vec3 pass_normal;
out vec3 pass_tangent;
out vec3 to_light;
out vec3 to_camera;

void main() {
	vec4 worldPosition = transformationMatrix * vec4(position, 1.0);

	gl_Position = projectionMatrix * cameraMatrix * worldPosition;
	pass_uv = uv_position;
	pass_normal = (transformationMatrix * vec4(normal,0.0)).xyz;
	pass_tangent = (transformationMatrix * vec4(tangent,0.0)).xyz;

	to_light = lightPosition - worldPosition.xyz;
	to_camera = (inverse(cameraMatrix) * vec4(0,0,0,1)).xyz - worldPosition.xyz;
}
