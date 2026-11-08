using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;
using System;
using UnityEngine;
using UnityEngine.UI;

public class UIHallBox : UIBase<UIModelHallBox, UIViewHallBox>
{
    private UIViewHallBox uIViewHallBox;

    private int curOpenNum = 1;
    private int maxOpenNum = 1;
    private int boxLevel = 0;
    private CFGItem_BoxLevel cfgBoxLevelItem;

    private bool isAuto = false;
    private Animation autoAnimation;

    private Action<bool, int> callBack;
    private Action openCallBack;
    private Action closeUICallBack;

    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }
 
    public override void OnInit()
    {
        uIViewHallBox = GetView();

        autoAnimation = uIViewHallBox.atuoBtn.GetComponentInChildren<Animation>();

        uIViewHallBox.lvBtn.onClick.AddListener(OnClickLvBtn);
        uIViewHallBox.leftBtn.onClick.AddListener(OnClickLeftBtn);
        uIViewHallBox.rightBtn.onClick.AddListener(OnClickRightBtn);
        uIViewHallBox.settingBtn.onClick.AddListener(OnClicSettingBtn);
        uIViewHallBox.atuoBtn.onClick.AddListener(OnClicAutoBtn);
        uIViewHallBox.okBtn.onClick.AddListener(OnClicOkBtn);

        openEnable = true;
        var boxBtn = uIViewHallBox.boxTrans.GetComponent<Button>();
        boxBtn.onClick.AddListener(OnClickBoxButton);

        uIViewHallBox.lvRedGo.SetActive(false);

        RedPointManager.Instance.RPS.Register(RedPointConst.BoxLevel1.ToString(), (node) =>
        {
            if (uIViewHallBox != null && uIViewHallBox.lvRedGo != null)
                uIViewHallBox.lvRedGo.SetActive(node.pointCount > 0);
        });
    }

    // lv
    private void OnClickLvBtn()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Open<UIBoxUpgrade>();
    }

    // setting
    private void OnClicSettingBtn()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        UITool.Open<UIOpenBoxSetting>();
    }
    private void OnClicOkBtn()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        SetActive(false);
        closeUICallBack?.Invoke();
    }

    #region box
    private bool openEnable;

    private void OnClickBoxButton()
    {
        if (!openEnable || RoleM.Instance.TreasuerBox == 0)
            return;

        openEnable = false;
        openCallBack?.Invoke();
        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "宝箱按钮",
            position = "宝箱界面",
            event_click = $"单次开宝箱",
        });
    }

    public void OpenFail()
    {
        openEnable = true;
    }

    public void ShowBoxSpine()
    {
        SpineUtil.Instance.PlayOnce(uIViewHallBox.boxTrans, "dead", () =>
        {
            openEnable = true;
            SpineUtil.Instance.PlayLoop(uIViewHallBox.boxTrans, "stand");
        });
    }
    #endregion

    #region left|right
    private void OnClickLeftBtn()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        OnBoxOpneNumChange(false);
    }

    private void OnClickRightBtn()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        OnBoxOpneNumChange(true);
    }

    private void OnBoxOpneNumChange(bool isAdd)
    {
        if (isAdd)
            curOpenNum++;
        else
            curOpenNum--;

        curOpenNum = Mathf.Clamp(curOpenNum, 1, maxOpenNum);
        uIViewHallBox.openNumText.text = curOpenNum.ToString();

        callBack?.Invoke(isAuto, curOpenNum);
    }
    #endregion

    #region auto 
    private void OnClicAutoBtn()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        if (!isAuto && RoleM.Instance.TreasuerBox == 0)
            return;

        isAuto = !isAuto && RoleM.Instance.TreasuerBox > 0;
        autoAnimation.enabled = isAuto;

        callBack?.Invoke(isAuto, curOpenNum);
        string openState = isAuto ? "开" : "关";
        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "宝箱按钮",
            position = "宝箱界面",
            event_click = $"自动开宝箱[{openState}]",
        });
    }
    #endregion

    private void ShowBoxLevel()
    {
        uIViewHallBox.lvText.text = "LV." + RoleM.Instance.BoxNowlevel;
    }

    private void GetCurrentLevelBoxConfig()
    {
        cfgBoxLevelItem = CFG_BoxLevel.GetById(RoleM.Instance.BoxNowlevel);
        maxOpenNum = cfgBoxLevelItem.Maxtime;
    }

    private void SetDefaultOpenNum()
    {
        curOpenNum = maxOpenNum;
        uIViewHallBox.openNumText.text = curOpenNum.ToString();
    }

    private void SetOpenBoxNumOnBoxLevelUp()
    {
        if (RoleM.Instance.BoxNowlevel > boxLevel)
        {
            boxLevel = RoleM.Instance.BoxNowlevel;
            SetDefaultOpenNum();
        }
    }

    public void DisableAutoOpenBox()
    {
        isAuto = false;
        if (autoAnimation != null)
            autoAnimation.enabled = isAuto;
    }

    public void RefreshNum()
    {
        if (RoleM.Instance.TreasuerBox == 0)
            OnClicAutoBtn();

        uIViewHallBox.numText.text = RoleM.Instance.TreasuerBox.ToString();
    }

    public void SetData(bool isAuto, Action<bool, int> callBack, Action boxOpenCallBack, Action closeUICallBack)
    {
        RefreshNum();
        this.callBack = callBack;
        this.openCallBack = boxOpenCallBack;
        this.closeUICallBack = closeUICallBack;
        openEnable = true;

        GetCurrentLevelBoxConfig();
        ShowBoxLevel();
        SetDefaultOpenNum();

        CheckRed();

        autoAnimation.enabled = isAuto;
        if(!isAuto)
            SpineUtil.Instance.PlayLoop(uIViewHallBox.boxTrans, "stand");

        this.isAuto = isAuto;
    }

    private void CheckRed()
    {
        RoleM.Instance.CheckBoxLevelRed();
    }

    public void SetActive(bool active)
    {
        uIViewHallBox.gameObject.SetActive(active);
    }

    public void RefreshUserBoxInfo()
    {
        GetCurrentLevelBoxConfig();
        ShowBoxLevel();
        SetOpenBoxNumOnBoxLevelUp();
        // 宝箱红点
        CheckRed();
    }

    public void RefreshUserBasePush()
    {
        GetCurrentLevelBoxConfig();
        CheckRed();
    }

    public bool ActiveSelf()
    {
        return uIViewHallBox.gameObject.activeSelf;
    }

    public override void OnDestroy()
    {
        uIViewHallBox = null;
        RedPointManager.Instance.RPS.Register(RedPointConst.BoxLevel1.ToString(), null);
    }
}