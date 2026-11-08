// Upgrade NOTE: replaced 'mul(UNITY_MATRIX_MVP,*)' with 'UnityObjectToClipPos(*)'

// - Unlit


Shader "SelfShader/FishWithWave" {
Properties {
	_MainTex ("Base layer (RGB)", 2D) = "white" {}
	_MainColor("Base Layer Color", Color) = (1,1,1,1)
	_DetailTex ("layer (RGB)", 2D) = "white" {}
	_DetailColor ("layer Color", Color) = (1,1,1,1)
	_ModelUpVector("model Up vector",vector) = (0,1,0)

	_IceTex("Ice", 2D) = "white" {}
	_IceStrength ("ice strength", Range(0, 1)) = 0 //0.75

	_ScrollX ("layer Scroll speed X", Float) = 1
	_ScrollY ("layer Scroll speed Y", Float) = 0.0

	_DetailScale("Detail Scale", Float) = 0.5

	
	_MMultiplier ("Layer Multiplier", Float) = 1.0

	_RimColor("Rim Color",Color) = (0,0,0,0)


	_ShadowPlane("PlanePos", vector) = (-0.1, -0.5, 85, 85)
	_ShadowProjDir("LightDir", vector) = (-0.3, -1.13, 1.12, 0.0)
	//_ShadowProjDir("LightDir", vector) = (3.3, -2.41, 1.64, 0.0)
	_ShadowFalloff("ShadowFalloff", float) = 1
	_ZTestAddValue("ZTest Add Value",float) = 0
	 [Enum(UnityEngine.Rendering.CompareFunction)] _ZTest("ZTest", Float) = 4
}


	SubShader{
		Tags{ "IgnoreProjector" = "True" "RenderType" = "Transparent" "Queue" = "Transparent" }
		ZTest[_ZTest]
		//Blend SrcAlpha OneMinusSrcAlpha
		//Cull Off Lighting Off ZWrite Off Fog { Color (0,0,0,0) }

		LOD 100

		Pass{
	//	ZWrite Off
		Blend SrcAlpha OneMinusSrcAlpha
		CGPROGRAM

		#pragma multi_compile LIGHTMAP_OFF LIGHTMAP_ON
		#pragma exclude_renderers molehill
		#include "UnityCG.cginc"

		#pragma vertex vert
		#pragma fragment frag
		#pragma fragmentoption ARB_precision_hint_fastest	

		
		sampler2D _MainTex;
		sampler2D _DetailTex;

		float4 _MainTex_ST;
		float4 _DetailTex_ST;

		half _ScrollX;
		half _ScrollY;
		half _MMultiplier;
		half _DetailScale;

		half3 _ModelUpVector;

		fixed4 _MainColor;
		fixed4 _DetailColor;
		fixed4 _RimColor;

		sampler2D _IceTex;
		float4 _IceTex_ST;
		float _IceStrength;

		struct v2f {
			half4 pos : SV_POSITION;
			half4 uv : TEXCOORD0;
			fixed4 color : TEXCOORD1;
			fixed percent : TEXCOORD2;
			half2 iceuv : TEXCOORD3;
		};
		float _ZTestAddValue;
		v2f vert(appdata_full v)
		{
			v2f o;

			 float4 clipPos = UnityObjectToClipPos(v.vertex);
			//计算深度时额外增加一个值：         
			#if UNITY_REVERSED_Z
				clipPos.z += _ZTestAddValue;
				clipPos.z = min(clipPos.z, clipPos.w * UNITY_NEAR_CLIP_VALUE);//限定在范围内
			#else
				clipPos.z -= _ZTestAddValue;
				clipPos.z = max(clipPos.z, clipPos.w * UNITY_NEAR_CLIP_VALUE);//限定在范围内
			#endif
			o.pos = clipPos;
			// o.pos = UnityObjectToClipPos(v.vertex);
			o.uv.xy = TRANSFORM_TEX(v.texcoord.xy,_MainTex);
			half4 worldPos = mul(unity_ObjectToWorld, v.vertex);

			fixed2 s2 = fixed2(worldPos.z / 8.66 + 1, 1);
			o.uv.zw = (worldPos.xz * _DetailScale / s2 + half2(_ScrollX, _ScrollY) * _Time);//frac(worldPos.xz );

			o.iceuv = TRANSFORM_TEX(v.texcoord.xy, _IceTex);

			o.color = _MMultiplier.xxxx;
			//o.percent = clamp(dot(_ModelUpVector,UnityObjectToWorldNormal(v.normal) - 0.1f), 0, 1);// 
			o.percent = abs( dot(_ModelUpVector, UnityObjectToWorldNormal(v.normal)) ) - 0.1;
			return o;
		}

		fixed4 frag (v2f i) : COLOR
		{
			fixed4 o;
			fixed4 tex = tex2D (_MainTex, i.uv.xy);
			fixed4 tex2 = tex2D (_DetailTex, i.uv.zw);
			fixed4 ice = tex2D (_IceTex, i.iceuv);
			

			o = (tex * _MainColor) + (tex2 * _DetailColor * i.percent * _MMultiplier) +_RimColor;
			

			//fixed fade = 1 - (sin(_Time.y * 1) + 1) * 0.5 * 0.3;
			fixed fade = 0.8 - (sin(_Time.y * 1.5)* 0.2);
			o.rgb = lerp(o.rgb, ice.rgb, _IceStrength * fade) * 1.0 ;
			o.a = 1;
			return o;
		}
		ENDCG 
	}

	Pass
	{
		ZWrite On
		ColorMask 0
		//Cull Off

		CGPROGRAM

#pragma vertex vert
#pragma fragment frag

		float4 _ShadowPlane;
		float4 _ShadowProjDir;
		float _ShadowFalloff;
		float _ZTestAddValue;
		struct appdata
		{
			float4 vertex : POSITION;
		};

		struct v2f
		{
			float4 vertex : SV_POSITION;
		};

		v2f vert(appdata v)
		{
			//v2f o;
			//float3 worldpos = mul(unity_ObjectToWorld, v.vertex).xyz;
			//fixed scal = (8.66 + _ShadowPlane.z) / (8.66 + worldpos.z);
			//fixed s2 = worldpos.z / 3.82 + 1;
			//worldpos = fixed3((worldpos.x + _ShadowPlane.x * s2) * scal, (worldpos.y + _ShadowPlane.y * s2) * scal , _ShadowPlane.z);
			//o.vertex = mul(UNITY_MATRIX_VP, float4(worldpos, 1.0));
			////o.scale = smoothstep(1, 0, 1 / scal) + 0.2;
			//return o;

			v2f o;
			float3 worldpos = mul(unity_ObjectToWorld, v.vertex).xyz;
			// #if UNITY_REVERSED_Z
			// 	worldpos.z += _ZTestAddValue;
			// 	worldpos.z = min(worldpos.z, 1 * UNITY_NEAR_CLIP_VALUE);//限定在范围内
			// #else
			// 	worldpos.z -= _ZTestAddVal;
			// 	worldpos.z = max(worldpos.z, 1 * UNITY_NEAR_CLIP_VALUE);//限定在范围内
			// #endif
			float s2 = worldpos.z / 3.77 + 1;
			float3 lightdir = normalize(_ShadowProjDir);
			worldpos = worldpos + _ShadowFalloff * s2  * lightdir.xyz;//
			o.vertex = mul(UNITY_MATRIX_VP, float4(worldpos, 1.0));
			return o;
		}

		float4 frag(v2f i) : SV_Target
		{
			return fixed4(0,0,0,0);
		}

		ENDCG
	}

	Pass
	{
		Blend SrcAlpha  OneMinusSrcAlpha
		ZWrite Off
		Cull Back
		ColorMask RGB


		CGPROGRAM

#pragma vertex vert
#pragma fragment frag

		float4 _ShadowPlane;
		float4 _ShadowProjDir;
		float _ShadowFalloff;
		float _ZTestAddValue;
		struct appdata
		{
			float4 vertex : POSITION;
		};

		struct v2f
		{
			float4 vertex : SV_POSITION;
			float3 scale : TEXCOORD0;
		};

		v2f vert(appdata v)
		{
			//v2f o;
			//float3 worldpos = mul(unity_ObjectToWorld, v.vertex).xyz;
			//fixed scal = (8.66 + _ShadowPlane.z) / (8.66 + worldpos.z);
			//fixed s2 = worldpos.z / 3.82 + 1;
			//worldpos = fixed3((worldpos.x + _ShadowPlane.x * s2) * scal, (worldpos.y + _ShadowPlane.y * s2) * scal , _ShadowPlane.z);
			//o.vertex = mul(UNITY_MATRIX_VP, float4(worldpos, 1.0));
			////o.scale = smoothstep(1, 0, 1 / scal) + 0.2;
			//return o;

			v2f o;
			float3 worldpos = mul(unity_ObjectToWorld, v.vertex).xyz;
			// #if UNITY_REVERSED_Z
			// 	worldpos.z += _ZTestAddValue;
			// 	worldpos.z = min(worldpos.z, 1 * UNITY_NEAR_CLIP_VALUE);//限定在范围内
			// #else
			// 	worldpos.z -= _ZTestAddVal;
			// 	worldpos.z = max(worldpos.z,1 * UNITY_NEAR_CLIP_VALUE);//限定在范围内
			// #endif
			float s2 = worldpos.z / 3.77 + 1;
			float3 lightdir = normalize(_ShadowProjDir);
			worldpos = worldpos + _ShadowFalloff * s2  * lightdir.xyz;//
			o.vertex = mul(UNITY_MATRIX_VP, float4(worldpos, 1.0));
			return o;
		}

		float4 frag(v2f i) : SV_Target
		{
			float4 color;
			color.xyz = float3(0.0, 0.0, 0.0) ;
			color.w = 0.3 * step(0.001, _ShadowFalloff);// i.scale;
			return color;
		}

		ENDCG
	}

	}
	Fallback "Transparent/VertexLit"
}
