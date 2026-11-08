using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPVPOneReady : MonoBehaviour,IUIView
{
	public Transform root_Player;
	public Button btn_Cancel;
	public GameObject obj_Over;
	public Text txt_Time;

	public List<Transform> root_Playes;

	public void Awake()
    {

    }
	public void Init(GameObject go)
	{
		root_Playes = new List<Transform>();

        for (int i = 0; i < 6; i++)
        {
			root_Playes.Add(root_Player.GetChild(i));
		}
	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPVPOneReady>();
	}
}