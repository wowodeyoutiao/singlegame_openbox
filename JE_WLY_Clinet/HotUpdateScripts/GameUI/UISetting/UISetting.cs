using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.MyNet;

using MetaFramework.UI;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UISetting : UIBase<UIModelSetting, UIViewSetting>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        //JEvent.defaultEvent.Register(this);

        // 个人信息
        GetView().Btn_ChangeName.onClick.AddListener(OnChangeName);
        GetView().Btn_ChangeServer.onClick.AddListener(OnChangeServer);
        GetView().Btn_CopyID.onClick.AddListener(OnCopyId);

        // 音效控制
        GetView().toggle_Music.onValueChanged.AddListener(OnMusicToggleClick);
        GetView().toggle_Sound.onValueChanged.AddListener(OnSoundToggleClick);

        // 显示设置
        // 显示模式
        GetView().toggle_SimpleMode.onValueChanged.AddListener(OnSimpleModeToggleClick);
        GetView().toggle_FluencyMode.onValueChanged.AddListener(OnFluencyModeToggleClick);
        GetView().toggle_CompleteMode.onValueChanged.AddListener(OnCompleteModeToggleClick);
        GetView().btn_Help_DisplayMode.onClick.AddListener(OnHelpDisplayModeClick);
        // 画质
        GetView().toggle_LowQuality.onValueChanged.AddListener(OnLowQualityToggleClick);
        GetView().toggle_HighQuality.onValueChanged.AddListener(OnHighQualityToggleClick);
        GetView().btn_Help_ImageQuality.onClick.AddListener(OnHelpImageQualityModeClick);
        // 帧率
        GetView().toggle_30FPS.onValueChanged.AddListener(On30FPSToggleClick);
        GetView().toggle_60FPS.onValueChanged.AddListener(On60FPSToggleClick);
        GetView().btn_Help_FrameRateSetting.onClick.AddListener(OnHelpFPSClick);
        // 特效显示
        GetView().toggle_ShieldingEveryone.onValueChanged.AddListener(OnShieldingEveryoneToggleClick);
        GetView().toggle_ShieldingOther.onValueChanged.AddListener(OnShieldingOtherToggleClick);
        GetView().btn_Help_EffectSetting.onClick.AddListener(OnHelpEffectSettingClick);

        GetView().Btn_Close.onClick.AddListener(OnCloseClick);
    }

    public override void OnEnter(params object[] args)
    {
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);

        ShowPlayerName(RoleM.Instance.Nickname);
        GetView().txt_Server.text = GetServerName();
        GetView().txt_ID.text = RoleM.Instance.Uid.ToString();

        var cfgHead = CFG_RoleIcon.GetById(RoleM.Instance.CurrentAvatarId);
        var cfgFrame = CFG_RoleIcon.GetById(RoleM.Instance.CurrentFrameId);
        ResouseHelper.Instance.LoadSpriteByAtlas(cfgHead.Path, cfgHead.GetIcon(RoleM.Instance.Job), GetView().img_Head);
        ResouseHelper.Instance.LoadSpriteByAtlas(cfgFrame.Path, cfgFrame.GetIcon(RoleM.Instance.Job), GetView().img_HeadFrame);
        ResouseHelper.Instance.LoadSpriteByAtlas(cfgHead.Path, GameConst.HeadIconQualityPrefix + cfgHead.Quality, GetView().img_HeadQuality);

        LoadSetting();
    }

    // 显示玩家昵称
    private void ShowPlayerName(string playerName)
    {
        GetView().txt_Name.text = playerName;
    }

    // 获取当前服务器名称
    private string GetServerName()
    {
        var SelectSever = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx];
        return SelectSever.name;
    }

    private void LoadSetting()
    {
        GetView().toggle_Music.isOn = AudioManager.Instance.isOpenBgMusic;
        GetView().toggle_Sound.isOn = AudioManager.Instance.isOpenSoundEffect;

        int displayMode = GameUtil.Instance.GetPlayerPrefsInt(GameConst.SystemSettingOption.DisplayMode.ToString(), 1);
        int imageQuality = GameUtil.Instance.GetPlayerPrefsInt(GameConst.SystemSettingOption.ImageQuality.ToString(), 0);
        int fpsSetting = GameUtil.Instance.GetPlayerPrefsInt(GameConst.SystemSettingOption.TargetFrameRate.ToString(), 30);
        int effectSetting = GameUtil.Instance.GetPlayerPrefsInt(GameConst.SystemSettingOption.EffectSetting.ToString(), -1);

        switch (displayMode)
        {
            case 0:
                GetView().toggle_SimpleMode.isOn = true;
                break;
            case 1:
                GetView().toggle_FluencyMode.isOn = true;
                break;
            case 2:
                GetView().toggle_CompleteMode.isOn = true;
                break;
            default:
                break;
        }

        switch (imageQuality)
        {
            case 0:
                GetView().toggle_LowQuality.isOn = true;
                break;
            case 1:
                GetView().toggle_HighQuality.isOn = true;
                break;
            default:
                break;
        }

        switch (fpsSetting)
        {
            case 30:
                GetView().toggle_30FPS.isOn = true;
                break;
            case 60:
                GetView().toggle_60FPS.isOn = true;
                break;
            default:
                break;
        }

        switch (effectSetting)
        {
            case -1:
                GetView().toggleGroup_EffectSetting.SetAllTogglesOff();
                break;
            case 0:
                GetView().toggle_ShieldingEveryone.isOn = true;
                break;
            case 1:
                GetView().toggle_ShieldingOther.isOn = true;
                break;
            default:
                break;
        }
    }

    #region 音效控制
    // 背景音乐
    private void OnMusicToggleClick(bool isOn)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        AudioManager.Instance.OnMusic(isOn);
    }

    // 游戏音效
    private void OnSoundToggleClick(bool isOn)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        AudioManager.Instance.OnSound(isOn);
    }
    #endregion

    #region 显示设置
    // 简约模式
    private void OnSimpleModeToggleClick(bool isOn)
    {
        if (isOn) 
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.DisplayMode.ToString(), 0);
        }
    }

    // 流畅模式
    private void OnFluencyModeToggleClick(bool isOn)
    {
        if (isOn)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.DisplayMode.ToString(), 1);
        }
    }

    // 完整模式
    private void OnCompleteModeToggleClick(bool isOn)
    {
        if (isOn)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.DisplayMode.ToString(), 2);
        }
    }

    private void OnHelpDisplayModeClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIComRule>(null, 21);
    }
    #endregion

    #region 画质设置
    private void OnLowQualityToggleClick(bool isOn)
    {
        if (isOn)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.ImageQuality.ToString(), 0);
        }
    }
    
    private void OnHighQualityToggleClick(bool isOn)
    {
        if (isOn)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.ImageQuality.ToString(), 1);
        }
    }

    private void OnHelpImageQualityModeClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIComRule>(null, 22);
    }
    #endregion

    #region 帧率设置
    private void On30FPSToggleClick(bool isOn)
    {
        if (isOn)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            SDKCtr.Instance.SetPreferredFramesPerSecond(30);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.TargetFrameRate.ToString(), 30);
        }
    }

    private void On60FPSToggleClick(bool isOn)
    {
        if (isOn)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            SDKCtr.Instance.SetPreferredFramesPerSecond(60);
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.TargetFrameRate.ToString(), 60);
        }
    }

    private void OnHelpFPSClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIComRule>(null, 23);
    }
    #endregion

    #region 特效设置
    private void OnShieldingEveryoneToggleClick(bool isOn)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (isOn)
        {
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.EffectSetting.ToString(), 0);
        }
        else if(!GetView().toggleGroup_EffectSetting.AnyTogglesOn())
        {
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.EffectSetting.ToString(), -1);
        }
    }
    
    private void OnShieldingOtherToggleClick(bool isOn)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (isOn)
        {
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.EffectSetting.ToString(), 1);
        }
        else if (!GetView().toggleGroup_EffectSetting.AnyTogglesOn())
        {
            GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.EffectSetting.ToString(), -1);
        }
    }

    private void OnHelpEffectSettingClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIComRule>(null, 24);
    }

    #endregion
    // 改名
    private void OnChangeName()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIChangeName>();
    }

    // 换服
    private void OnChangeServer()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        GameHelper.Instance.BackToLogin(true);
        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "设置界面",
            position = "设置按钮",
            event_click = $"选服",
        });
        //TDAnalyticsUtil.Logout();
    }


    // 将用户ID复制到剪切板
    private void OnCopyId()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        SDKCtr.Instance.CopyString(RoleM.Instance.Uid.ToString(), () =>
        {
            UITool.ShowTips("账号ID已复制到剪切板");
        });
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        ShowPlayerName(RoleM.Instance.Nickname);
    }

    private void OnCloseClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Close<UISetting>();
    }

    public override void OnExit()
    {
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
    }

    public override void OnDestroy()
    {
        //JEvent.defaultEvent.Unregister(this);
    }
}