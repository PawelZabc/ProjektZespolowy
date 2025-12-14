#version 330

// Input vertex attributes (from vertex shader)
in vec3 fragPosition;
in vec2 fragTexCoord;
//in vec4 fragColor;
in vec3 fragNormal;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// NOTE: Add here your custom variables

#define     MAX_LIGHTS              4
#define     LIGHT_DIRECTIONAL       0
#define     LIGHT_POINT             1

struct MaterialProperty {
    vec3 color;
    int useSampler;
    sampler2D sampler;
};

struct Light {
    int enabled;
    int type;
    vec3 position;
    vec3 target;
    vec4 color;
};

// Input lighting values
uniform Light lights[MAX_LIGHTS];
uniform vec4 ambient;
uniform vec3 viewPos;

void main()
{
    // Texel color fetching from texture sampler
    vec4 texelColor = texture(texture0, fragTexCoord);
    vec3 lightDot = vec3(0.0);
    vec3 normal = normalize(fragNormal);
    vec3 viewD = normalize(viewPos - fragPosition);
    vec3 specular = vec3(0.0);

    // NOTE: Implement here your fragment shader code

    for (int i = 0; i < MAX_LIGHTS; i++)
    {
        if (lights[i].enabled == 1)
        {
            vec3 light = vec3(0.0);
            float attenuation = 1.0;

            if (lights[i].type == LIGHT_DIRECTIONAL)
            {
                light = -normalize(lights[i].target - lights[i].position);
            }

            if (lights[i].type == LIGHT_POINT)
            {
                // Compute light direction and distance
                light = lights[i].position - fragPosition;
                float distance = length(light);
                light = normalize(light);

                // Apply inverse-square attenuation
                attenuation = 1.0 / (distance * distance);

                // Optional: clamp to max distance
                float maxDistance = 10.0;
                attenuation = 1.0 / (distance * distance);  // inverse-square
                attenuation *= smoothstep(maxDistance, 0.0, distance);
            }

            float intensity = 5.0;
            float NdotL = max(dot(normal, light), 0.0);
            lightDot += lights[i].color.rgb * NdotL * attenuation * intensity;

            float specCo = 0.0;
            if (NdotL > 0.0) 
                specCo = pow(max(0.0, dot(viewD, reflect(-(light), normal))), 16.0); // 16 refers to shine

            specular += specCo * attenuation * intensity;
        }
    }

    finalColor = (texelColor * ((colDiffuse + vec4(specular, 1.0)) * vec4(lightDot, 1.0)));
    finalColor += texelColor * (ambient / 10.0) * colDiffuse;

    // Gamma correction
    finalColor = pow(finalColor, vec4(1.0/2.2));
}
