using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPetAward : MonoBehaviour,IUIView
{
	public Transform obj_eff7;
    public Transform obj_eff8;

	public Transform root_Pet;
	public Button btn_back;
	public Text txt_Name;

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPetAward>();
	}
}