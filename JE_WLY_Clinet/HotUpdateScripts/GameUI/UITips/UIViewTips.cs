using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

/// <summary>
/// Auto Generate Class!!!
/// </summary>
public class UIViewTips :MonoBehaviour, IUIView
{
	public GameObject obj_item;
    public Transform root_Hide;

    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
    }

    private void OnDestroy()
    {
        UITool.RemoveUI<UITips>();
    }
}