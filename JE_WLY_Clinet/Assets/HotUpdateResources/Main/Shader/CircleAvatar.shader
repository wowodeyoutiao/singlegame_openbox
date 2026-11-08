Shader "UI/CircleAvatar"
{
    Properties
    {
        [PerRendererData] _MainTex ("Texture", 2D) = "white" {}
        _SoftEdge ("Soft Edge", Range(0.4,0.5)) = 0.48
        _BorderColor ("Border Color", Color) = (1,1,1,1)
        _BorderWidth ("Border Width", Float) = 2
        _AspectRatio ("Aspect Ratio", Vector) = (1,1,0,0)
    }

    SubShader
    {
        Tags
        {
            "Queue"="Transparent"
            "RenderType"="Transparent"
            "PreviewType"="Plane"
        }

        Blend SrcAlpha OneMinusSrcAlpha
        ZWrite Off
        Cull Off

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
                float4 pos : SV_POSITION;
            };

            sampler2D _MainTex;
            float4 _MainTex_ST;
            float _SoftEdge;
            float4 _BorderColor;
            float _BorderWidth;
            float4 _AspectRatio;

            v2f vert (appdata v)
            {
                v2f o;
                o.pos = UnityObjectToClipPos(v.vertex);
                o.uv = TRANSFORM_TEX(v.uv, _MainTex);
                
                // ÊÊÅä¿í¸ß±È
                o.uv = (o.uv - 0.5) * _AspectRatio.xy + 0.5;
                return o;
            }

            fixed4 frag (v2f i) : SV_Target
            {
                float2 center = float2(0.5, 0.5);
                float dist = distance(i.uv, center);
                
                // ±ßÔµÓð»¯
                float alpha = smoothstep(_SoftEdge, _SoftEdge - 0.02, dist);
                
                // ±ß¿ò¼ÆËã
                float border = smoothstep(
                    _SoftEdge - _BorderWidth/100, 
                    _SoftEdge - _BorderWidth/100 - 0.01, 
                    dist
                ) - alpha;

                fixed4 col = tex2D(_MainTex, i.uv);
                fixed4 final = col;
                final.a *= alpha;
                final.rgb = lerp(final.rgb, _BorderColor.rgb, border * _BorderColor.a);
                
                return final;
            }
            ENDCG
        }
    }
}