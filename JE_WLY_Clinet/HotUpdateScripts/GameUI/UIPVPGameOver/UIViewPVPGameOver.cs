using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewPVPGameOver : MonoBehaviour,IUIView
{
	public GameObject obj_Win;
	public GameObject obj_Lose;

	public Transform root_Awards;
	public Button btn_yellow;
    public Text txt_yellow;

    public Button btn_blue;
    public Text txt_blue;

    public Text txt_Time;

	public Text txt_loseJF;
	public Text txt_winJF;


    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIPVPGameOver>();
	}
}