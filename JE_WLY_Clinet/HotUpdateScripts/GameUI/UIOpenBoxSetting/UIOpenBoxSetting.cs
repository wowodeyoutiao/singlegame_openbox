using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using System;
using System.Collections;
using System.Collections.Generic;
using System.Reflection.Emit;

using UnityEngine;
using UnityEngine.UI;

public class UIOpenBoxSetting : UIBase<UIModelOpenBoxSetting, UIViewOpenBoxSetting>
{
    //public class Callback
    //{
    //    private Action closeFun;
    //    public Action CloseFunc => closeFun;

    //    public Callback() 
    //    {
    //    }
        
    //    public Callback(Action closeHandle)
    //    {
    //        closeFun = closeHandle;
    //    }
    //}

    //Callback callback;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_SaveSetting.onClick.AddListener(() =>
        {
            CloseAndSave();
        });
        GetView().btn_Close.onClick.AddListener(() =>
        {
            CloseAndSave();
        });
        GetView().btn_Mask.onClick.AddListener(() =>
        {
            CloseAndSave();
        });

        GetView().dropdown_Quality.options = new List<Dropdown.OptionData>()
        {
            new Dropdown.OptionData("白色品质以上"),
            new Dropdown.OptionData("绿色品质以上"),
            new Dropdown.OptionData("蓝色品质以上"),
            new Dropdown.OptionData("紫色品质以上"),
            new Dropdown.OptionData("粉色品质以上"),
            new Dropdown.OptionData("橙色品质以上"),
            new Dropdown.OptionData("红色品质以上"),
            new Dropdown.OptionData("至尊品质以上"),
        };

        for (int i = 0; i < GetView().toggles_role.Length; i++) 
        {
            int idx = i;
            GetView().toggles_role[i].onValueChanged.AddListener(isOn =>
            {
                if (isOn)
                {
                    if (GameUtil.Instance.SlotIsUnLockOrHasRole(idx))
                        slotSelect = idx;
                    else
                        SetActiveToggle(slotSelect);
                }
            });
        }

        //GetView().dropdown_Occupation.options = new List<Dropdown.OptionData>()
        //{
        //    new Dropdown.OptionData("战士"),
        //    new Dropdown.OptionData("法师"),
        //    new Dropdown.OptionData("祭祀"),
        //};

        //GetView().dropdown_Level.options = GetLevelDropdownOptions();

        //GetView().btn_SubLevel.onClick.AddListener(SubAppointLevel);
        //GetView().btn_AddLevel.onClick.AddListener(AddAppointLevel);

        //GetView().input_Level.onValueChanged.AddListener((value) =>
        //{
        //    if (string.IsNullOrEmpty(value))
        //        return;

        //    appointLevel = int.Parse(value);

        //    if (appointLevel < minLevel)
        //        appointLevel = minLevel;

        //    if (appointLevel > maxLevel)
        //        appointLevel = maxLevel;

        //    GetView().input_Level.text = appointLevel.ToString();
        //});

        //GetView().dropdown_AppointOccupation.options = GetAppointOccupationDropdownOptions();

        GetView().btn_GoFirstPay.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (ActivityM.Instance.GetAcitvityState(ActivityType.FirstPay))
            {
                HttpC.GetFirstPayInfo((req) => {
                    if (req.code == 0 && !isNullGo)
                    {
                        SaveSetting();
                        UITool.Close<UIOpenBoxSetting>();
                        UITool.Open<UIFirstPay>(null, req);
                    }
                });
            }
        });

    }

    //public List<Dropdown.OptionData> GetLevelDropdownOptions()
    //{
    //    var options = new List<Dropdown.OptionData>();
    //    int lv = 10;
    //    int maxLv = 150;
    //    options.Add(new Dropdown.OptionData($"{lv}级以上"));
    //    while (lv < maxLv) 
    //    {
    //        options.Add(new Dropdown.OptionData($"{lv += 10}级以上"));
    //    }
    //    return options;
    //}

    //public List<Dropdown.OptionData> GetAppointOccupationDropdownOptions()
    //{
    //    var options = new List<Dropdown.OptionData>()
    //    {
    //        new Dropdown.OptionData("默认全部"),
    //        new Dropdown.OptionData("战士"),
    //        new Dropdown.OptionData("法师"),
    //        new Dropdown.OptionData("祭祀")
    //    };

    //    return options;
    //}

    //private int appointLevel;
    //private int minLevel = 1;
    //private int maxLevel;

    private int slotSelect;

    public override void OnEnter(params object[] args)
    {
        //callback = args[0] as Callback;
        //GameHelper.Instance.SetDeayFrame(SetOptionState);
        SetOptionState();
    }

    //private void SubAppointLevel()
    //{
    //    appointLevel--;
    //    if (appointLevel < minLevel)
    //        appointLevel = minLevel;

    //    GetView().input_Level.text = appointLevel.ToString();
    //}

    //private void AddAppointLevel()
    //{
    //    appointLevel++;
    //    if (appointLevel > maxLevel)
    //        appointLevel = maxLevel;

    //    GetView().input_Level.text = appointLevel.ToString();
    //}

    //private int GetActiveOnToggle()
    //{
    //    for (int i = 0; i < GetView().toggles_role.Length; i++)
    //    {
    //        if (GetView().toggles_role[i].isOn) return i;
    //    }
    //    return -1;
    //}

    // 勾选一个toggle
    private void SetActiveToggle(int index)
    {
        for (int i = 0; i < GetView().toggles_role.Length; i++)
        {
            if (i == index) GetView().toggles_role[i].isOn = true;
        }
    }

    private void SetOptionState()
    {
        // 基础设置，下拉选项
        GetView().dropdown_Quality.value = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.Quality.ToString() + RoleM.Instance.Uid, 1) - 1;
        //GetView().dropdown_Level.value = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.Level.ToString() + RoleM.Instance.Uid, 10) / 10 - 1;
        //GetView().dropdown_Occupation.value = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.Occupation.ToString() + RoleM.Instance.Uid, 1) - 1;
        // 基础设置，选项开关
        GetView().toggle_Quality.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.QualityToggle.ToString() + RoleM.Instance.Uid));
        //GetView().toggle_Level.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.LevelToggle.ToString() + RoleM.Instance.Uid));
        //GetView().toggle_Occupation.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.OccupationToggle.ToString() + RoleM.Instance.Uid));
        // 高级设置，选项开关
        GetView().toggle_BagFull.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.BagFullToggle.ToString() + RoleM.Instance.Uid, 1));
        GetView().toggle_HigherThenProtagonist.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.HigherThenProtagonistEquippedToggle.ToString() + RoleM.Instance.Uid, 1));
        GetView().toggle_HigherThenPartner.isOn = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.HigherThenPartnerEquippedToggle.ToString() + RoleM.Instance.Uid, 1));

        // 指定开出某个槽位的角色的装备
        slotSelect = GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.DefaultOpenSoltCondition.ToString() + RoleM.Instance.Uid);
        SetActiveToggle(slotSelect);
        // 回收设置，选项开关
        GetView().toggle_AutoRecycle.isOn= Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.AutoOpenBoxOption.AutoRecycleToggle.ToString() + RoleM.Instance.Uid, 1));
        // 装备开箱指定
        //appointLevel = RoleM.Instance.NowOpenLevel == 0 ? RoleM.Instance.Level : RoleM.Instance.NowOpenLevel;
        //maxLevel = RoleM.Instance.Level;

        //GetView().input_Level.text = null;
        //GetView().input_Level.text = appointLevel.ToString();
        //GetView().dropdown_AppointOccupation.value = RoleM.Instance.NowOpenJob;

        GetView().txt_AppointSettingOpen.text = $"首充开启";
        //GetView().mask_AppointSetting.SetActive(!ActivityM.Instance.IsBuyFirstPayGift());
    }

    private void CloseAndSave()
    {
        SaveSetting();
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        //callback.CloseFunc?.Invoke();
        UITool.Close<UIOpenBoxSetting>();
    }

    private void SaveSetting()
    {
        // 基础设置，下拉选项
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.Quality.ToString() + RoleM.Instance.Uid, GetView().dropdown_Quality.value + 1);
        //GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.Level.ToString() + RoleM.Instance.Uid, (GetView().dropdown_Level.value + 1) * 10);
        //GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.Occupation.ToString() + RoleM.Instance.Uid, GetView().dropdown_Occupation.value + 1);
        // 基础设置，选项开关
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.QualityToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_Quality.isOn));
        //GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.LevelToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_Level.isOn));
        //GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.OccupationToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_Occupation.isOn));
        // 高级设置，选项开关
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.BagFullToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_BagFull.isOn));
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.HigherThenProtagonistEquippedToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_HigherThenProtagonist.isOn));
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.HigherThenPartnerEquippedToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_HigherThenPartner.isOn));
        // 指定开出某个槽位的角色的装备
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.DefaultOpenSoltCondition.ToString() + RoleM.Instance.Uid, slotSelect);
        // 回收设置，选项开关
        GameUtil.Instance.SetPlayerPrefs(GameConst.AutoOpenBoxOption.AutoRecycleToggle.ToString() + RoleM.Instance.Uid, Convert.ToInt32(GetView().toggle_AutoRecycle.isOn));

        // 装备开箱指定

        //RoleM.Instance.NowOpenLevel = appointLevel == RoleM.Instance.Level ? 0 : appointLevel;
        //RoleM.Instance.NowOpenJob = GetView().dropdown_AppointOccupation.value;
    }

    public override void OnExit()
    {
    }

    public override void OnDestroy()
    {

    }
}