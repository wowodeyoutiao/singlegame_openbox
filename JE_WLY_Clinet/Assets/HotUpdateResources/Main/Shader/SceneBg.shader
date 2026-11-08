// Unlit shader. Simplest possible textured shader.
// - no lighting
// - no lightmap support
// - no per-material color

Shader "SelfShader/SceneBg" {
	Properties{
		_MainTex("Base (RGB)", 2D) = "white" {}
		_BgTex("BgTex",2D) = "white" {}
		_LightTex("LightTex",2D) = "white" {}
		_WaveNoiseTex("Wave Noise", 2D) = "white" {}//噪波贴图
		_WaveIndentity("Indentity", float) = 0.1//表示水波的扭曲强度
		_WaveSpeedX("WaveSpeedX", float) = 0.08//噪波贴图延X方向的移动速度
		_BgWidth("BgWidth", Range(0,1)) = 0.5
		_MiddleWidth("MiddleWidth", Range(0,1)) = 0.5
		_BgSpeed("BgSpeed", float) = 0.08
		_CenterSpeed("CenterSpeed", float) = 0.08
		_FrontSpeed("FrontSpeed", float) = 0.08
		_FrontSpeedS("FrontSpeedS", float) = 0
		_CenterSpeedS("CenterSpeedS", float) = 0
		_BgSpeedS("BgSpeedS", float) = 0
		_RimColor("Rim Color", Color) = (0.26,0.19,0.16,0)
		_MultColor("Mult Color", Color) = (1,1,1,1)

		_DisappearOffset("Disappear Offset",Range(-0.5,0.5)) = 0.5
	}

		SubShader{
		Tags{ "RenderType" = "Opaque" }
		LOD 100

		Pass{
			CGPROGRAM
			#pragma vertex vert
			#pragma fragment frag
			#pragma target 2.0
			#pragma multi_compile_fog

			#include "UnityCG.cginc"

			struct appdata_t {
				float4 vertex : POSITION;
				float2 texcoord : TEXCOORD0;
				float4 color : COLOR;
				UNITY_VERTEX_INPUT_INSTANCE_ID
			};

			struct v2f {
				float4 vertex : SV_POSITION;
				float3 texcoord : TEXCOORD0;
				float2 frontAlphacoord : TEXCOORD1;
				float2 frontColorcoord : TEXCOORD2;
				float4 color:COLOR;
				UNITY_FOG_COORDS(1)
					UNITY_VERTEX_OUTPUT_STEREO
			};

			sampler2D _MainTex;
			sampler2D _BgTex;
			sampler2D _LightTex;
			sampler2D _WaveNoiseTex;
			float4 _MainTex_ST;
			fixed _BgWidth;
			fixed _MiddleWidth;
			fixed _BgSpeed;
			fixed _CenterSpeed;
			fixed _FrontSpeed;
			fixed _BgSpeedS;
			fixed _CenterSpeedS;
			fixed _FrontSpeedS;
			fixed _WaveIndentity;
			fixed _WaveSpeedX;
			fixed _MoveSpeed;
			fixed4 _RimColor;
			fixed4 _MultColor;
			float _DisappearOffset;

			v2f vert(appdata_t v)
			{
				v2f o;
				UNITY_SETUP_INSTANCE_ID(v);
				UNITY_INITIALIZE_VERTEX_OUTPUT_STEREO(o);
				o.vertex = UnityObjectToClipPos(v.vertex);
				o.texcoord.xy = TRANSFORM_TEX(v.texcoord, _MainTex);
				o.color = v.color;

				o.texcoord.z = _DisappearOffset - v.vertex.x;
				UNITY_TRANSFER_FOG(o,o.vertex);
				return o;
			}

			fixed4 frag(v2f i) : SV_Target
			{
				fixed4 outcol;
				clip(i.texcoord.z);
				float waveX = tex2D(_WaveNoiseTex, i.texcoord.xy + float2(_Time.y * _WaveSpeedX, 0)).r;
				//float waveX2 = tex2D(_WaveNoiseTex, i.texcoord.xy + float2(0, _Time.y * 0.05)).r;
				float2 ruv = float2(i.texcoord.x + waveX * _WaveIndentity - 0.009, i.texcoord.y);
				//float2 ruv2 = float2(i.texcoord.x + waveX2 * 0.04, i.texcoord.y);

				fixed coordX = ruv.x * _MiddleWidth;
				//fixed bgcoordX = frac(ruv.x * _BgWidth - _Time.y * 0.0025 * _BgSpeed + _BgSpeedS);
				fixed centercoordX = frac(coordX - _Time.y * 0.0035 * _CenterSpeed + _CenterSpeedS);
				//fixed frontcoordX = frac(coordX - _Time.y * 0.0045 * _FrontSpeed + _FrontSpeedS);
				//fixed lightcoordX = frac(ruv2.x * _MiddleWidth + _Time.y * 0.0015 * _CenterSpeed + _CenterSpeedS);

				////背景颜色
				//outcol = tex2D(_BgTex, float2(bgcoordX, ruv.y * 0.5 + 0.5));
				////中景灰度
				//fixed4 centerAlpha = tex2D(_MainTex, float2(centercoordX, ruv.y * 0.25));
				//fixed centergray = Luminance(centerAlpha.rgb);
				////中景颜色
				//fixed4 centercol = tex2D(_MainTex, float2(centercoordX, ruv.y * 0.25 + 0.5));
				//outcol = lerp(outcol, centercol, centergray);
				////中景叠加亮度
				//fixed4 centerLight = tex2D(_LightTex, float2(lightcoordX * 8, ruv.y * 8));
				//fixed4 LightAlpha = tex2D(_BgTex, float2(centercoordX, ruv.y * 0.5));
				//fixed Lightgray = Luminance(LightAlpha.rgb);
				//centerLight = centerLight * Lightgray * 2;
				//outcol = outcol + centerLight;
				////前景灰度颜色
				//fixed4 frontAlpha = tex2D(_MainTex,float2(frontcoordX, ruv.y * 0.25 + 0.25));
				//fixed frontgray = Luminance(frontAlpha.rgb);
				////前景颜色
				//fixed4 frontcol = tex2D(_MainTex, float2(frontcoordX, ruv.y * 0.25 + 0.75));

				//outcol = tex2D(_MainTex, float2(centercoordX + 0.09, ruv.y));
				outcol = tex2D(_MainTex, ruv);

				//outcol = lerp(outcol, frontcol, frontgray);

				UNITY_APPLY_FOG(i.fogCoord, outcol);
				UNITY_OPAQUE_ALPHA(outcol.a);

				//fixed fade = sin(_Time.y * 0.5);//-1 ~ 1
				//outcol.r = outcol.r * (fade + 1)*0.5;
				//outcol -= _RimColor * (fade + 1)*0.5;


				outcol = outcol * _MultColor;
				return outcol;
			}
			ENDCG
		}
	}

}
