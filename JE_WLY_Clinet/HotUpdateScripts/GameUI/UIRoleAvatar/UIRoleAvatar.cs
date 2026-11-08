using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using SuperScrollView;

using System;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIRoleAvatar : UIBase<UIModelRoleAvatar, UIViewRoleAvatar>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    UIViewRoleAvatar _uiView;

    public override void OnInit()
    {
        _uiView = GetView();

        _uiView.btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRoleAvatar>();
        });

        for (int i = 0; i < _uiView.btn_Selects.Length; i++)
        {
            int idx = i;
            _uiView.btn_Selects[i].gameObject.SetActive(i < btnsName.Length);
            if (i < btnsName.Length)
            {
                Text text1 = _uiView.btn_Selects[i].transform.Find("Label").GetComponent<Text>();
                Text text2 = _uiView.btn_Selects[i].transform.Find("Select/Label").GetComponent<Text>();
                text1.text = btnsName[i];
                text2.text = btnsName[i];
                _uiView.btn_Selects[i].onClick.AddListener(()=>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    if(buttonSelectIndex == idx) return;
                    OnSelectBtnClick(idx);
                });
            }
        }

        _uiView.avatarList.InitGridView(0, ItemRefresh);

        _uiView.btn_Use.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UseAvatarOrFrame();
        });

        _uiView.btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 50);
        });
    }


    private LoopGridViewItem ItemRefresh(LoopGridView list, int index, int row, int column)
    {
        int _index = index;
        var item = list.NewListViewItem("Item");
        Transform t = item.transform;
        bool isAvatar = buttonSelectIndex == 0;
        var cfg = isAvatar ? avatars[index] : frames[index];
        bool isUsing = isAvatar ? RoleM.Instance.CurrentAvatarId == cfg.ID : RoleM.Instance.CurrentFrameId == cfg.ID;
        bool isHave = isAvatar ? RoleM.Instance.HasAvatar(cfg.ID) : RoleM.Instance.HasFrame(cfg.ID);
        t.Find("using").gameObject.SetActive(isUsing);
        t.Find("select").gameObject.SetActive(_index == itemSelectIndex);
        t.Find("quality").gameObject.SetActive(isAvatar);
        t.Find("frame").gameObject.SetActive(!isAvatar);

        PlayerAvatarInfo info;

        if (isAvatar)
        {
            RoleM.Instance.TryGetAvatar(cfg.ID, out info);
            
            Image quality = JUIExtensions.img(t, "quality");
            ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, GameConst.HeadIconQualityPrefix + cfg.Quality, quality);

            Image head = JUIExtensions.img(t, "quality/icon");
            string headName = cfg.GetIcon(RoleM.Instance.Job);
            ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, headName, head); 
        }
        else
        {
            RoleM.Instance.TryGetFrame(cfg.ID, out info);

            Image frame = JUIExtensions.img(t, "frame");
            string frameName = cfg.GetIcon(RoleM.Instance.Job);
            ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, frameName, frame);
        }

        string timeRemaining;
        if (cfg.Period == 0)
            timeRemaining = "有效期：永久";
        else if (isHave)
            timeRemaining = $"有效期：{GameUtil.Instance.RemainingDays(info.expire_time * 1000) - 1}天";
        else
            timeRemaining = $"有效期：{cfg.Period}天";

        JUIExtensions.txt(t, "txt_Time").text = timeRemaining;

        Button btn = JUIExtensions.btn(t);
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if(_index != itemSelectIndex)
            {
                OnItemSelect(_index);
            }
        });

        return item;
    }

    private void OnItemSelect(int index)
    {
        itemSelectIndex = index;
        _uiView.avatarList.RefreshAllShownItem();

        bool isAvatar = buttonSelectIndex == 0;
        var cfg = isAvatar ? avatars[itemSelectIndex] : frames[itemSelectIndex];
        _uiView.txt_Name.text = GameUtil.Instance.GetStrByQuality(cfg.Quality, cfg.Name);
        _uiView.outline_Name.effectColor = GameUtil.Instance.GetQualityOutlineColor(cfg.Quality);
        _uiView.img_Frame.gameObject.SetActive(!isAvatar);
        _uiView.img_Head.gameObject.SetActive(isAvatar);
        _uiView.img_Quality.gameObject.SetActive(isAvatar);

        string iconName = cfg.GetIcon(RoleM.Instance.Job);
        if (isAvatar)
        {
            ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, GameConst.HeadIconQualityPrefix + cfg.Quality, _uiView.img_Quality);
            ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, iconName, _uiView.img_Head);
        }
        else
        {
            ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, iconName, _uiView.img_Frame);
        }

        _uiView.txt_GetFrom.text = "获取途径:  " + cfg.Get;
        _uiView.txt_Attribute.text = cfg.GetAttributeString();
        SetUseButtonState();
    }

    private void SetUseButtonState()
    {
        bool isAvatar = buttonSelectIndex == 0;
        var cfg = isAvatar ? avatars[itemSelectIndex] : frames[itemSelectIndex];
        bool isUsing = isAvatar ? RoleM.Instance.CurrentAvatarId == cfg.ID : RoleM.Instance.CurrentFrameId == cfg.ID;
        bool isHave = isAvatar ? RoleM.Instance.HasAvatar(cfg.ID) : RoleM.Instance.HasFrame(cfg.ID);
        bool showUseButton = isHave && !isUsing;
        _uiView.btn_Use.gameObject.SetActive(showUseButton);
        _uiView.txt_UseOrGetState.gameObject.SetActive(!showUseButton);
        if (!showUseButton)
        {
            string useOrGetState = isUsing ? "<color=#40B023>使用中</color>" : (isHave ? "" : "<color=#FF5855>未获得</color>");
            _uiView.txt_UseOrGetState.text = useOrGetState;
        }
    }

    private void OnSelectBtnClick(int index)
    {
        buttonSelectIndex = index;

        if (index == 0 || index == 1)
        {
            itemSelectIndex = 0;
            RefreshList();
            OnItemSelect(0);
        }

        for (int i = 0; i < _uiView.btn_Selects.Length; i++)
        {
            Transform t = _uiView.btn_Selects[i].transform;
            t.Find("Select").gameObject.SetActive(i == index);
        }
    }

    private void RefreshList()
    {
        if (buttonSelectIndex == 0)
        {
            _uiView.avatarList.SetListItemCount(avatars.Count);
            _uiView.avatarList.RefreshAllShownItem();
        }
        else if (buttonSelectIndex == 1)
        {
            _uiView.avatarList.SetListItemCount(frames.Count);
            _uiView.avatarList.RefreshAllShownItem();
        }
    }

    private string[] btnsName = { "头像", "头像框", /*"光环", "时装"*/ };

    private int buttonSelectIndex;
    private int itemSelectIndex = -1;

    private List<CFGItem_RoleIcon> avatars = new List<CFGItem_RoleIcon>();
    private List<CFGItem_RoleIcon> frames = new List<CFGItem_RoleIcon>();

    public override void OnEnter(params object[] args)
    {
        UniGameEvent.AddEventListener<GameEvent.UserAvatarPush>(OnRefresh_UserAvatarPush);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnRefresh_ZeroTimeTime);
        UniGameEvent.AddEventListener<GameEvent.Refresh_UseAvatar>(OnRefrersh_Avatar);

        avatars = CFG_RoleIcon.GetListByType(AvatarType.AVATAR);
        frames = CFG_RoleIcon.GetListByType(AvatarType.FRAME);

        SortList(avatars, true);
        SortList(frames, false);

        OnSelectBtnClick(0);
    }

    public void SortList(List<CFGItem_RoleIcon> cfgs, bool isAvatar) 
    {
        int defaultId = isAvatar ? CFG_XMLConfig.xmlConfig.BaseValue.Avatar : CFG_XMLConfig.xmlConfig.BaseValue.AvatarFrame;
        cfgs.Sort((x, y) =>
        {
            int isDefauleX = Convert.ToInt32(x.ID == defaultId);
            int isDefaultY = Convert.ToInt32(y.ID == defaultId);

            //int hasX = Convert.ToInt32(isAvatar ? RoleM.Instance.HasAvatar(x.ID) : RoleM.Instance.HasFrame(x.ID));
            //int hasY = Convert.ToInt32(isAvatar ? RoleM.Instance.HasAvatar(y.ID) : RoleM.Instance.HasFrame(y.ID));

            if (isDefauleX != isDefaultY)
                return isDefauleX.CompareTo(isDefauleX);
            //else if (hasX != hasY)
            //    return hasY.CompareTo(hasX);
            else if (x.Quality != y.Quality)
                return y.Quality.CompareTo(x.Quality);
            else if (x.ID != y.ID)
                return x.ID.CompareTo(y.ID);
            else return 0;
        });
    }

    private void UseAvatarOrFrame()
    {
        bool isAvatar = buttonSelectIndex == 0;
        var cfg = isAvatar ? avatars[itemSelectIndex] : frames[itemSelectIndex];
        TcpC.C2S_UseAvatarRequest(isAvatar ? AvatarType.AVATAR : AvatarType.FRAME, cfg.ID);
    }

    private void OnRefresh_UserAvatarPush(GameEvent.UserAvatarPush _e)
    {
        if (isNullGo) return;
        if (buttonSelectIndex == 0 || buttonSelectIndex == 1)
            RefreshList();
    }

    /// <summary>
    /// 头像头像框使用返回
    /// </summary>
    private void OnRefrersh_Avatar(GameEvent.Refresh_UseAvatar _e)
    {
        if (isNullGo) return;
        if (buttonSelectIndex == 0 || buttonSelectIndex == 1)
        {
            RefreshList();
            // 刷新使用按钮状态
            SetUseButtonState();
        }
    }

    private void OnRefresh_ZeroTimeTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        if (buttonSelectIndex == 0 || buttonSelectIndex == 1)
            RefreshList();
    }

    public override void OnExit()
    {
        UniGameEvent.RemoveEventListener<GameEvent.UserAvatarPush>(OnRefresh_UserAvatarPush);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnRefresh_ZeroTimeTime);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_UseAvatar>(OnRefrersh_Avatar);
    }

    public override void OnDestroy()
    {

    }
}