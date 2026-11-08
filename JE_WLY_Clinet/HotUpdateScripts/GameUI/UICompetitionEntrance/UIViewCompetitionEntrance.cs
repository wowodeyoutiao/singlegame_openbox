using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewCompetitionEntrance : MonoBehaviour,IUIView
{
	public LoopListView2 loopList;
	public Button btn_Close;

	public void Awake()
    {

    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UICompetitionEntrance>();
	}
}