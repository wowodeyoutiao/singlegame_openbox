using System.Collections.Generic;
using UnityEngine;
using MetaFramework.UI;
using JEngine.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.Event;
using System.Text;

public class UIDayDo : UIBase<UIModelDayDo, UIViewDayDo>
{
    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }

    private CFG_Room.RoomType[] AllRoomInfo = new CFG_Room.RoomType[] { 
        CFG_Room.RoomType.QHS,
        CFG_Room.RoomType.SXS,
        CFG_Room.RoomType.HS,
        CFG_Room.RoomType.SY,
        CFG_Room.RoomType.Gold,
        CFG_Room.RoomType.EXP,
    };

    public static Dictionary<CFG_Room.RoomType, CFG_XMLConfig.XML_Config.SwitchOpne> RoomToLockDic = new Dictionary<CFG_Room.RoomType, CFG_XMLConfig.XML_Config.SwitchOpne>
    {
        [CFG_Room.RoomType.QHS] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_QHSAndGold,
        [CFG_Room.RoomType.SXS] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_ExpAndSXS,
        [CFG_Room.RoomType.HS] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_HS,
        [CFG_Room.RoomType.SY] = CFG_XMLConfig.XML_Config.SwitchOpne.Play_SY,
        [CFG_Room.RoomType.Gold] = CFG_XMLConfig.XML_Config.SwitchOpne.GoldDungeon,
        [CFG_Room.RoomType.EXP] = CFG_XMLConfig.XML_Config.SwitchOpne.ExpDungeon,
    };

    public override void OnInit()
    {
        GetView().btn_AutoPlay.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if(!PrivilegeM.Instance.CheckQuicklyDaily())
            {
                if (GameUtil.Instance.CheakIsLock(RoomToLockDic[OnSelectIdx]))
                    return;

                if (GetView().obj_vip.activeSelf)
                {
                    UITool.ShowTips("开启月卡，可永久享用扫荡功能！");
                    return;
                }
            }

            if (RoleM.Instance.GetCanTimeByRoomType((int)OnSelectIdx) == 0)
            {
                if (CFG_VIP.IsMaxVipLevel(RoleM.Instance.Vip))
                {
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = "今日可挑战次数不足，是否提升VIP等级增加每日可挑战次数？",
                        YesFun = () => {
                            UITool.Close<UIFightChoice>();
                            UITool.Open<UIVIP>();
                        },
                    });
                }
                else
                {
                    UITool.ShowTips("今日挑战次数已用完");
                }
                return;
            }

            TcpC.C2S_CopySweep(SelectCFG.ID);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "副本",
                position = SelectCFG.Name,
                event_click = "扫荡",
            });
        });

        GetView().btn_Paly.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (GameUtil.Instance.CheakIsLock(RoomToLockDic[OnSelectIdx]))
                return;

            if (RoleM.Instance.GetCanTimeByRoomType((int)OnSelectIdx) == 0)
            {
                if (CFG_VIP.IsMaxVipLevel(RoleM.Instance.Vip))
                {
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = "今日可挑战次数不足，是否提升VIP等级增加每日可挑战次数？",
                        YesFun = () => {
                            UITool.Close<UIFightChoice>();
                            UITool.Open<UIVIP>();
                        },
                    });
                }
                else 
                {
                    UITool.ShowTips("今日挑战次数已用完");
                }
                return;
            }

            TcpC.GetGameJoin(SelectCFG.ID,0);

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "副本",
                position = SelectCFG.Name,
                event_click = "挑战",
            });
        });

        for (int i = 0; i < GetView().root_Btns.childCount; i++)
        {
            int idx = i;
            JUIExtensions.btn(GetView().root_Btns.GetChild(i)).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnRoomSelect(AllRoomInfo[idx]);
            });

            var cfgInfo = CFG_Room.GetByRoomType(AllRoomInfo[idx], -1);
            JUIExtensions.txt(GetView().root_Btns.GetChild(i).Find("txt_des")).text = cfgInfo.Suggest;
        }

        GetView().btn_Add.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            int itemId = 0;
            switch (OnSelectIdx)
            {
                case CFG_Room.RoomType.EXP:
                    itemId = 25080;
                    break;
                case CFG_Room.RoomType.Gold:
                    itemId = 25079;
                    break;
                case CFG_Room.RoomType.QHS:
                    itemId = 25075;
                    break;
                case CFG_Room.RoomType.SY:
                    itemId = 25078;
                    break;
                case CFG_Room.RoomType.SXS:
                    itemId = 25076;
                    break;
                case CFG_Room.RoomType.HS:
                    itemId = 25077;
                    break;
                default:
                    break;
            }

            if (itemId != 0)
                GameUtil.Instance.CommonChallengeCountExchange(itemId);
        });
    }

    private CFG_Room.RoomType OnSelectIdx ;
    private CFGItem_Room SelectCFG;

    private void OnRedRefrsh()
    {
        for (int i = 0; i < GetView().root_Btns.childCount; i++)
        {
            var idx = AllRoomInfo[i];

            bool isLockDayLimt = GameUtil.Instance.CheckOpenServerIsLock(RoomToLockDic[idx], false);
            StringBuilder strTips = new();
            bool isLockLv = GameUtil.Instance.CheakIsLock(RoomToLockDic[idx], false, strTips);

            Transform root = GetView().root_Btns.GetChild(i);
            root.Find("dot").gameObject.SetActive(!isLockDayLimt && !isLockLv && RoleM.Instance.GetCanTimeByRoomType((int)idx) > 0);

            root.Find("lock").gameObject.SetActive(isLockDayLimt || isLockLv);
            JUIExtensions.txt(root, "lock/txt").text = isLockDayLimt ? $"{GameUtil.Instance.RemainingDaysUnlockOpen(RoomToLockDic[idx])}天后开启" : strTips.ToString();//CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[RoomToLockDic[idx]] + "级解锁";
        }
    }

    private void OnRoomSelect(CFG_Room.RoomType roomType,bool isRefresh = false)
    {
        if (GameUtil.Instance.CheckOpenServerIsLock(RoomToLockDic[roomType]) || GameUtil.Instance.CheakIsLock(RoomToLockDic[roomType]))
            return;

        OnSelectIdx = roomType;
        SelectCFG = CFG_Room.GetByRoomType(roomType, -1);

        if (!isRefresh)
        {
            for (int i = 0; i < GetView().root_Btns.childCount; i++)
            {
                GetView().root_Btns.GetChild(i).Find("select").gameObject.SetActive(AllRoomInfo[i] == roomType);
            }

            GetView().txt_Title.text = SelectCFG.Name;
            GetView().txt_RoomDes.text = SelectCFG.CopyDes;

            for (int i = 0; i < GetView().root_Award.childCount; i++)
            {
                var go = GetView().root_Award.GetChild(i);
                go.gameObject.SetActive(i < SelectCFG.Droppreview.Length);
                if (i < SelectCFG.Droppreview.Length)
                {
                    int cfgId = SelectCFG.Droppreview[i];
                    GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(go.Find("icon")), JUIExtensions.img(go));
                    JUIExtensions.btn(go).onClick.RemoveAllListeners();
                    JUIExtensions.btn(go).onClick.AddListener(() => {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                    });
                }
            }
            GameUtil.Instance.SetMonsterIcon((int)SelectCFG.Show, GetView().img_Monster);

            GetView().obj_vip.SetActive(!PrivilegeM.Instance.CheckQuicklyDaily());
            GetView().txt_vip.text = $"月卡开启";

            //GetView().obj_lock1.SetActive(GameUtil.Instance.CheakIsLock(RoomToLockDic[roomType],false));
            //GetView().txt_lock1.text = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[RoomToLockDic[roomType]] + "级开启";
            //GetView().obj_lock2.SetActive(GameUtil.Instance.CheakIsLock(RoomToLockDic[roomType], false));
            //GetView().txt_lock2.text = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[RoomToLockDic[roomType]] + "级开启";

        }

        GetView().txt_pro.text = $"今日可挑战次数:{RoleM.Instance.GetCanTimeByRoomType((int)roomType)}/{CFG_XMLConfig.xmlConfig.DayRoomDic[roomType].challengeTimeNum + + CFG_VIP.GetAddRoomCount(RoleM.Instance.Vip,(int) roomType)}";
    }

    public void OnShow(bool isShow)
    {
        if (isShow)
        {
            OnEnter();
        }
        else
        {
            OnExit();
        }
        uiGo.SetActive(isShow);
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);

        OnSelectIdx = AllRoomInfo[0];
        OnRoomSelect(OnSelectIdx);
        OnRedRefrsh();
    }

    [Subscriber(ThreadMode.Main)]
    private void OnDaily_MapInfo(GameEvent.Daily_MapInfo _e)
    {
        if (isNullGo)
            return;
        OnRoomSelect(OnSelectIdx,true);
        OnRedRefrsh();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_CheckOpenTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo) return;
        TcpC.C2S_DailyMapInfo();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);

    }
}