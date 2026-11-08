using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewGuildDonation : MonoBehaviour,IUIView
{
	public Button Btn_Close;
	public Button Btn_Go_GuildDonation;
	public Button Btn_Go_SkillPray;
	public Button Btn_Go_ChallengeLeader;

	public void Awake()
    {
		
    }

	public void Init(GameObject go)
	{
        
    }

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIGuildDonation>();
	}
}