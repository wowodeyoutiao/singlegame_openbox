using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;

public class UIViewOpenBoxSetting : MonoBehaviour, IUIView
{
    #region 基础设置
    public Toggle toggle_Quality;
    //public Toggle toggle_Level;
    //public Toggle toggle_Occupation;

    public Dropdown dropdown_Quality;
    //public Dropdown dropdown_Level;
    //public Dropdown dropdown_Occupation;
    #endregion

    #region 高级设置
    public Toggle toggle_BagFull;
    public Toggle toggle_HigherThenProtagonist;
    public Toggle toggle_HigherThenPartner;
    #endregion

    #region 装备指定
    //public InputField input_Level;
    //public Button btn_SubLevel;
    //public Button btn_AddLevel;
    //public Dropdown dropdown_AppointOccupation;

    public ToggleGroup toggleGroup;
    [HideInInspector] public Toggle[] toggles_role;
    public GameObject mask_AppointSetting;
    public Text txt_AppointSettingOpen;
    public Button btn_GoFirstPay;
    public GameObject btn_GuideClick;
    #endregion

    #region 回收设置
    public Toggle toggle_AutoRecycle;
    #endregion

    public Button btn_SaveSetting;
    public Button btn_Close;
    public Button btn_Mask;

    public void Awake()
    {
        int count = 6;
        toggles_role = new Toggle[count];
        for (int i = 0; i < count; i++)
        {
            toggles_role[i] = JUIExtensions.toggle(toggleGroup.transform.GetChild(i + 1));
        }
    }


    public void Init(GameObject go)
    {

    }


    public void OnDestroy()
    {
        UITool.RemoveUI<UIOpenBoxSetting>();
    }
}