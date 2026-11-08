// Upgrade NOTE: replaced 'mul(UNITY_MATRIX_MVP,*)' with 'UnityObjectToClipPos(*)'


Shader "Custom/UICircular"
{
     Properties
    {
        [PerRendererData] _MainTex ("Texture", 2D) = "white" {}
        _Radius ("Radius", Range(0,1)) = 0.5
        _Softness ("Softness", Range(0,0.5)) = 0.01
    }
    
    SubShader
    {
        Tags { "Queue"="Transparent" "RenderType"="Transparent" }
        
        Blend SrcAlpha OneMinusSrcAlpha
        
        Pass
        {
            CGPROGRAM
            #pragma vertex vert
            #pragma fragment frag
            #include "UnityCG.cginc"
            
            struct appdata
            {
                float4 vertex : POSITION;
                float2 uv : TEXCOORD0;
            };
            
            struct v2f
            {
                float2 uv : TEXCOORD0;
                float4 vertex : SV_POSITION;
            };
            
            sampler2D _MainTex;
            float _Radius;
            float _Softness;
            
            v2f vert (appdata v)
            {
                v2f o;
                o.vertex = UnityObjectToClipPos(v.vertex);
                o.uv = v.uv;
                return o;
            }
            
            fixed4 frag (v2f i) : SV_Target
            {
                float2 center = float2(0.5, 0.5);
                float distance = length(i.uv - center);
                float alpha = smoothstep(_Radius + _Softness, _Radius - _Softness, distance);
                fixed4 col = tex2D(_MainTex, i.uv);
                col.a *= alpha;
                return col;
            }
            ENDCG
        }
    }
}
