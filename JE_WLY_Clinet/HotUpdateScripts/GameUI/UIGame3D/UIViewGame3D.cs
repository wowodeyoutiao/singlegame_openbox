using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System;

public class UIViewGame3D : MonoBehaviour,IUIView
{

	public Transform root_Monter;
	public GameObject obj_Monster;
	public RectTransform rect_Monster;

	public Transform root_PlayerHub;
	public GameObject obj_PlayerHub;
	public GameObject root_Collect;

	public RectTransform rect_HPChnage;
	public RectTransform root_FlyHP;

	public RectTransform root_BossMonterHub;

	public GameObject tip_enemy;

	public Slider slider_Collect;
    public Text txt_SiliderCollect;

	public GameObject root_Pos;
	public Text txt_Pos;

    public Dictionary<string, GameObject> EffObjDic = new Dictionary<string, GameObject>();

	public Action UpateCallBack;

	//Ð¡µØÍ¼
	public Transform root_Map;
    public Transform root_MapBG;
    public Transform root_MyPos;
    public Transform root_BossPos;
    public Transform root_TeamPos;
    public Transform root_EnemyPlayerPos;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{
	}
	public void OnDestroy()
	{
		UpateCallBack = null;
        UITool.RemoveUI<UIGame3D>();
	}

	public void LateUpdate()
	{
		UpateCallBack?.Invoke();

    }
}