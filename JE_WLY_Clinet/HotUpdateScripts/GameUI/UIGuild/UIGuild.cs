using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.MyNet;
using MetaFramework.UI;
using pb;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static CFG_RankingAward;
using static CFG_XMLConfig.XML_Config;

public class UIGuild : UIBase<UIModelGuild, UIViewGuild>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    #region tag enum

    /// <summary>
    /// 工会子面板类型
    /// </summary>
    public enum GuildChildPanelType
    {
        GuildInfo,      // 工会信息
        GuildMember,    // 工会成员
        GuildSkill,     // 工会技能
        GuildStore,      // 工会商店
        //GuildMap,       // 工会地图
        GuildList,      // 工会列表
    }

    /// <summary>
    /// 底部页签按钮
    /// </summary>
    public enum FunctionType
    {
        Guild,  // 行会
        Friend, // 好友
        Fxpand, // 拓展
    }

    #endregion

    // 标题名
    public string[] titleName = { "行会信息", "成员列表", "行会技能", "行会商店", /*"行会地图",*/ "行会列表" };
    // 行会功能列表（按钮名字）
    public string[] guildTagsName = { "信息", "成员", "技能", "商店", /*"地图",*/ "列表" };
    // 好友功能列表（按钮名字）
    public string[] friendTagsName = { };
    // 拓展功能列表（按钮名字）
    public string[] expandTagsName = { };

    // 行会子面板
    private UIGuildInfo         uiGuildInfo;
    private UIGuildMemberList   uiGuildMemberList;
    private UIGuildSkill        uiGuildSkill;
    private UIGuildStore        uiGuildStore;
    private UIGuildMap          uiGuildMap;
    private UIGuildList         uiGuildList;


    public override void OnInit()
    {
        for (int i = 0; i < GetView().btns_Subitem.Length; i++)
        {
            int index = i;
            GetView().btns_Subitem[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnSubitemButtonClick(index);
            });
        }

        for (int i = 0; i < GetView().btns_Bottom.Length; i++)
        {
            int index = i;
            GetView().btns_Bottom[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnButtomButtonClick(index);
            });
        }

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuild>();
        });

        FuncOpenShowWidget.Instance.HideFuncNewFlag(new List<SwitchOpne>() { SwitchOpne.Guild });
    }

    int lastClick_ButtomButtonIndex = -1;
    int last_ClickSubitemButtonIndex = -1;
    // 暂时不用
    int defultSubitemButtonSelect;

    bool isJoinGuild;

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(Event_UserBase);

        isJoinGuild = RoleM.Instance.IsJoinedGuild();

        int bottomBtnIndex = 0;
        defultSubitemButtonSelect = 0;

        if(args.Length > 0)
            bottomBtnIndex = Convert.ToInt32(args[0]);

        if(args.Length > 1)
            defultSubitemButtonSelect = Convert.ToInt32(args[1]);
        OnBindRed();
        OnButtomButtonClick(bottomBtnIndex);
    }

    private Dictionary<RedPointConst, GameObject> RedDic;


    private void OnBindRed()
    {
        RedDic = new Dictionary<RedPointConst, GameObject>
        {
            [RedPointConst.Guild_Tag]       = GetView().btns_Bottom[0].transform.Find("dot").gameObject,
            //[RedPointConst.GuildFirend_Tag] = GetView().btns_Bottom[1].transform.Find("dot").gameObject,
            //[RedPointConst.GuildExtend_Tag] = GetView().btns_Bottom[2].transform.Find("dot").gameObject,

            [RedPointConst.Guild_Info]      = GetView().btns_Subitem[0].transform.Find("dot").gameObject,
            [RedPointConst.Guild_Member]    = GetView().btns_Subitem[1].transform.Find("dot").gameObject,
            [RedPointConst.Guild_Skill]     = GetView().btns_Subitem[2].transform.Find("dot").gameObject,
            [RedPointConst.Guild_Store]     = GetView().btns_Subitem[3].transform.Find("dot").gameObject,
            //[RedPointConst.Guild_Map]       = GetView().btns_Subitem[4].transform.Find("dot").gameObject,
        };

        foreach (var item in RedDic)
        {
            RedPointManager.Instance.RPS.Register(item.Key.ToString(), (node) =>
            {
                if (item.Value)
                    item.Value.SetActive(node.pointCount > 0);
            });
        }
    }

    private void ClearBindRed()
    {
        if (RedDic != null)
        {
            foreach (var item in RedDic)
            {
                RedPointManager.Instance.RPS.Register(item.Key.ToString(), null);
            }
        }
    }

    /// <summary>
    /// 底部按钮点击事件
    /// </summary>
    /// <param name="index">非 -1 的值表示要跳转到指定的子界面</param>
    private void OnButtomButtonClick(int index, int subIndex = -1)
    {
        if (index == 1 || index == 2) 
        {
            UITool.ShowTips("功能暂未开放");
            return;
        }

        if (lastClick_ButtomButtonIndex == index)
        {
            return;
        }

        lastClick_ButtomButtonIndex = index;

        for (int i = 0; i < GetView().btns_Bottom.Length; i++)
        {
            GetView().btns_Bottom[i].transform.Find("Select").gameObject.SetActive(i == index);
        }

        FunctionType functionType = (FunctionType)index;

        GetView().root_GuildChildPanel.gameObject.SetActive(functionType == FunctionType.Guild);
        GetView().root_FriendChildPanel.gameObject.SetActive(functionType == FunctionType.Friend);
        GetView().root_ExpandChildPanel.gameObject.SetActive(functionType == FunctionType.Fxpand);

        SetSubitemButtonsName(functionType);

        switch (functionType)
        {
            case FunctionType.Guild:
                OnlyShowGuildListButton(!isJoinGuild);
                int subitemButtonIndex = isJoinGuild ? 0 : 4;
                subitemButtonIndex = subIndex == -1 ? subitemButtonIndex : subIndex;
                OnSubitemButtonClick(subitemButtonIndex);

                RedPointManager.Instance.SetRedClose(RedPointConst.Guild_InGuild.ToString(), () =>
                {
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_InGuild.ToString(), RedPointManager.Instance.GetRedNum(RedPointConst.Guild_InGuild.ToString()));
                });

                break;
            case FunctionType.Friend:
                break;
            case FunctionType.Fxpand:
                break;
            default:
                break;
        }
    }

    /// <summary>
    /// 设置左侧按钮名称
    /// </summary>
    /// <param name="type"></param>
    private void SetSubitemButtonsName(FunctionType type)
    {
        string[] names = null;
        switch (type)
        {
            case FunctionType.Guild:
                names = guildTagsName;
                break;
            case FunctionType.Friend:
                names = friendTagsName;
                break;
            case FunctionType.Fxpand:
                names = expandTagsName;
                break;
            default:
                break;
        }

        for (int i = 0; i < GetView().btns_Subitem.Length; i++)
        {
            GetView().btns_Subitem[i].gameObject.SetActive(i < names.Length);

            if (i < names.Length)
            {
                GetView().btns_Subitem[i].transform.Find("Label").GetComponent<Text>().text = names[i];
                GetView().btns_Subitem[i].transform.Find("Select/Label").GetComponent<Text>().text = names[i];
            }
        }
    }

    /// <summary>
    /// 子页签按钮点击事件
    /// </summary>
    /// <param name="index"></param>
    private void OnSubitemButtonClick(int index)
    {

        if (last_ClickSubitemButtonIndex == index)
            return;

        last_ClickSubitemButtonIndex = index;

        for (int i = 0; i < GetView().btns_Subitem.Length; i++)
        {
            GetView().btns_Subitem[i].transform.Find("Select").gameObject.SetActive(i == index);
        }

        FunctionType functionType = (FunctionType)lastClick_ButtomButtonIndex;

        switch (functionType)
        {
            case FunctionType.Guild:
                OpenGuildChildPanel((GuildChildPanelType)index);
                break;
            case FunctionType.Friend:
                break;
            case FunctionType.Fxpand:
                break;
            default:
                break;
        }
    }

    /// <summary>
    /// 打开行会子面板
    /// </summary>
    /// <param name="panelType"></param>
    private void OpenGuildChildPanel(GuildChildPanelType panelType)
    {
        switch (panelType)
        {
            case GuildChildPanelType.GuildInfo:
                GetView().txt_title.text = "";
                if (uiGuildInfo == null)
                {
                    UITool.LoadUII<UIGuildInfo>(GetView().root_GuildChildPanel, (ui) =>
                    {
                        uiGuildInfo = ui as UIGuildInfo;
                        uiGuildInfo.SetTitleListener((guildName) =>
                        {
                            GetView().txt_title.text = guildName;
                        });
                        uiGuildInfo.OnShow(true);
                    });
                }
                else
                {
                    uiGuildInfo.OnShow(true);
                }
                break;
            case GuildChildPanelType.GuildMember:
                GetView().txt_title.text = "成员列表";
                if (uiGuildMemberList == null)
                {
                    UITool.LoadUII<UIGuildMemberList>(GetView().root_GuildChildPanel, (ui) =>
                    {
                        uiGuildMemberList = ui as UIGuildMemberList;
                        uiGuildMemberList.OnShow(true);
                    });
                }
                else
                {
                    uiGuildMemberList.OnShow(true);
                }
                break;
            case GuildChildPanelType.GuildSkill:
                GetView().txt_title.text = "行会技能";
                if (uiGuildSkill == null)
                {
                    UITool.LoadUII<UIGuildSkill>(GetView().root_GuildChildPanel, (ui) =>
                    {
                        uiGuildSkill = ui as UIGuildSkill;
                        uiGuildSkill.OnShow(true);
                    });
                }
                else
                {
                    uiGuildSkill.OnShow(true);
                }
                break;
            case GuildChildPanelType.GuildStore:
                GetView().txt_title.text = "行会商店";
                if (uiGuildStore == null)
                {
                    UITool.LoadUII<UIGuildStore>(GetView().root_GuildChildPanel, (ui) =>
                    {
                        uiGuildStore = ui as UIGuildStore;
                        uiGuildStore.OnShow(true);
                    });
                }
                else
                {
                    uiGuildStore.OnShow(true);
                }

                RedPointManager.Instance.SetRedClose(RedPointConst.Guild_Store.ToString(), () =>
                {
                    int num = RedPointManager.Instance.GetRedNum(RedPointConst.Guild_Store.ToString());
                    if (num == 1)
                    {
                        num = GuildM.Instance.IsGuildStoreCanBuy() ? 1 : 0;
                    }
                    RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_Store.ToString(), num);
                });

                break;
            //case GuildChildPanelType.GuildMap:
            //    GetView().txt_title.text = "行会地图";
            //    if (uiGuildMap == null)
            //    {
            //        UITool.LoadUII<UIGuildMap>(GetView().root_GuildChildPanel, (ui) =>
            //        {
            //            uiGuildMap = ui as UIGuildMap;
            //            uiGuildMap.OnShow(true);
            //        });
            //    }
            //    else
            //    {
            //        uiGuildMap.OnShow(true);
            //    }
            //    break;
            case GuildChildPanelType.GuildList:
                GetView().txt_title.text = "行会列表";
                if (uiGuildList == null)
                {
                    UITool.LoadUII<UIGuildList>(GetView().root_GuildChildPanel, (ui) =>
                    {
                        uiGuildList = ui as UIGuildList;
                        uiGuildList.OnShow(true);
                    });
                }
                else
                {
                    uiGuildList.OnShow(true);
                }
                break;
            default:
                break;
        }

        if (uiGuildInfo != null && panelType != GuildChildPanelType.GuildInfo)
            uiGuildInfo.OnShow(false);

        if (uiGuildMemberList != null && panelType != GuildChildPanelType.GuildMember)
            uiGuildMemberList.OnShow(false);

        if (uiGuildSkill != null && panelType != GuildChildPanelType.GuildSkill)
            uiGuildSkill.OnShow(false);

        if (uiGuildStore != null && panelType != GuildChildPanelType.GuildStore)
            uiGuildStore.OnShow(false);

        //if (uiGuildMap != null && panelType != GuildChildPanelType.GuildMap)
        //    uiGuildMap.OnShow(false);

        if (uiGuildList != null && panelType != GuildChildPanelType.GuildList)
            uiGuildList.OnShow(false);
    }

    /// <summary>
    /// 是否只显示行会列表按钮
    /// </summary>
    /// <param name="value"></param>
    private void OnlyShowGuildListButton(bool value)
    {
        for (int i = 0; i < GetView().btns_Subitem.Length; i++)
        {
            if (i < guildTagsName.Length - 1)
                GetView().btns_Subitem[i].gameObject.SetActive(!value);
            else if (i == guildTagsName.Length - 1)
                GetView().btns_Subitem[i].gameObject.SetActive(true);
            else
                GetView().btns_Subitem[i].gameObject.SetActive(false);
        }
    }

    //[Subscriber(ThreadMode.Main)]
    // 加入或退出工会后创新界面按钮。
    private void Event_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        if (isJoinGuild != RoleM.Instance.IsJoinedGuild())
        {
            isJoinGuild = RoleM.Instance.IsJoinedGuild();

            if (lastClick_ButtomButtonIndex == 1 || lastClick_ButtomButtonIndex == 2)
                return;

            OnlyShowGuildListButton(!isJoinGuild);
            int subitemButtonIndex = isJoinGuild ? 0 : 4;
            OnSubitemButtonClick(subitemButtonIndex);
        }
    }


    [Subscriber(ThreadMode.Main)]
    private void SelectSubPage(SelectionSubPage selection)
    {
        if (isNullGo)
            return;

        if (lastClick_ButtomButtonIndex == selection.bottomButtonIndex)
            OnSubitemButtonClick(selection.subButtonIndex);
        else
            OnButtomButtonClick(selection.bottomButtonIndex, selection.subButtonIndex);
    }

    [Subscriber(ThreadMode.Main)]
    private void SetTitleName(Title title)
    {
        if (isNullGo)
            return;

        if (title.panelType == GuildChildPanelType.GuildInfo)
            GetView().txt_title.text = title.titleName;
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(Event_UserBase);

        if (uiGuildInfo != null)
            uiGuildInfo.OnShow(false);

        if (uiGuildMemberList != null)
            uiGuildMemberList.OnShow(false);

        if (uiGuildSkill != null)
            uiGuildSkill.OnShow(false);

        if (uiGuildStore != null)
            uiGuildStore.OnShow(false);

        if (uiGuildMap != null)
            uiGuildMap.OnShow(false);

        if (uiGuildList != null)
            uiGuildList.OnShow(false);

        lastClick_ButtomButtonIndex = -1;
        last_ClickSubitemButtonIndex = -1;
        ClearBindRed();
    }

    public override void OnDestroy()
    {
        
    }

    public struct SelectionSubPage
    {
        public int bottomButtonIndex;
        public int subButtonIndex;
    }

    public struct Title
    {
        public string titleName;
        public GuildChildPanelType panelType;
    }
}