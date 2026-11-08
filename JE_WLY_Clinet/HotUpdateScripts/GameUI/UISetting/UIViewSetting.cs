using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIViewSetting : MonoBehaviour, IUIView
{
    // **    个人信息     **//
    public Image img_HeadQuality;
    public Image img_HeadFrame;
    public Image img_Head;
    public Text txt_Name;
    public Text txt_Server;
    public Text txt_ID;
    public Button Btn_ChangeName;
    public Button Btn_ChangeServer;
    public Button Btn_CopyID;

    // **    声音控制     **//
    public Toggle toggle_Music;    // 背景音乐
    public Toggle toggle_Sound;    // 游戏音效
    
    // **    显示设置     **//
    public Toggle toggle_SimpleMode;    // 简约模式
    public Toggle toggle_FluencyMode;   // 流畅模式
    public Toggle toggle_CompleteMode;  // 完整模式
    public Button btn_Help_DisplayMode;

    public Toggle toggle_LowQuality;  // 低画质
    public Toggle toggle_HighQuality;  // 高画质
    public Button btn_Help_ImageQuality;

    public Toggle toggle_30FPS;
    public Toggle toggle_60FPS;
    public Button btn_Help_FrameRateSetting;

    public ToggleGroup toggleGroup_EffectSetting;
    public Toggle toggle_ShieldingEveryone;  // 屏蔽所有人(包括自己)
    public Toggle toggle_ShieldingOther;    // 屏蔽其他人
    public Button btn_Help_EffectSetting;

    public Button Btn_Close;

	public void Awake()
    {
        
    }
	public void Init(GameObject go)
	{

	}

    private void OnDestroy()
    {
		UITool.RemoveUI<UISetting>();
    }
}