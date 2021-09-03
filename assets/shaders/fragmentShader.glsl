#version 410

in vec2 pass_uv;
in vec3 pass_normal;
in vec3 pass_tangent;
in vec3 to_light;
in vec3 to_camera;

uniform float materialShineDamper;
uniform float materialReflectivity;

uniform sampler2D diffuseTexture;
uniform sampler2D normalMap;
uniform sampler2D specularMap;

uniform vec3 lightColour;

out vec4 frag_colour;

void main() {

	vec3 unit_to_light = normalize(to_light);
	vec3 unit_to_camera = normalize(to_camera);
	vec3 unit_normal = normalize(pass_normal);

	vec3 tangent_n = normalize(pass_tangent - (pass_tangent * dot(pass_tangent, unit_normal)));
	vec3 bitangent = cross(unit_normal, tangent_n);

	vec4 diffuse_here  = texture(diffuseTexture, pass_uv);
	vec4 normal_here   = texture(normalMap, pass_uv);
	vec4 specular_here = texture(specularMap, pass_uv);

	vec3 effective_normal = (2*normal_here.b-1) * unit_normal + (2*normal_here.r-1) * tangent_n + (2*normal_here.g-1) * bitangent;

	float light = max(dot(unit_to_light, effective_normal), 0.0);
	light = 0.15 + 0.85*light;
	vec4 diffuse_colour = vec4(lightColour * light,1.0) * diffuse_here;

	vec3 reflected_light = reflect(-unit_to_light, effective_normal);

	float specularFactor = max(0.0, dot(reflected_light, to_camera));
	specularFactor = materialReflectivity * pow(specularFactor, materialShineDamper);

	vec4 specular_light = vec4(lightColour,1.0) * specularFactor;

	if ( specular_here.b < 0.2 ) {
		discard;
	}

	frag_colour = diffuse_colour + specular_light;

	// frag_colour = vec4(normalize(pass_normal),1.0);
	// frag_colour = vec4(unit_normal, 1.0);
	// frag_colour = vec4(effective_normal,1.0);
}
