using Cysharp.Threading.Tasks;

using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;
using System;
using System.Collections.Generic;
using System.Text;
using UnityEngine;

using static CFG_XMLConfig.XML_Config;
using static GameEvent;

public class UICom_DownBtns : UIBase<UIModelCom_DownBtns, UIViewCom_DownBtns>
{
    private UIViewCom_DownBtns uIViewCom_DownBtns;
    private SwitchOpne[] switchOpnes = new SwitchOpne[]
    {
        SwitchOpne.Play,
        SwitchOpne.Role,
        SwitchOpne.CallRole,
        SwitchOpne.No,
        SwitchOpne.Bag,
        SwitchOpne.Guild,
        SwitchOpne.ExpPool,
    };

    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }
    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
        UniGameEvent.AddEventListener<GameEvent.Refresh_LvUp>(OnRefresh_LvUp);
        UniGameEvent.AddEventListener<GameEvent.Refresh_MainBattleChanged>(OnRefresh_MainBattleChanged);
        UniGameEvent.AddEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);

        uIViewCom_DownBtns = GetView();

        for (int i = 1; i <= 7; i++)
        {
            int idx = i;
            uIViewCom_DownBtns.Btns[i - 1].btn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnDownBtnCkick(idx);
            });
        }

        OnLockCheak();
        LockInfoRefresh();
        OnBindRed();
        OnLocalRedCheck();
        OnRoleRefresh();
        OnRedCheak();

        refreshFuncNewFlag();

        NewPlayerGuideWidget.Instance.AddGuideGameObject("DownBtns_FightButton", GetView().Btns[0].trans);
    }

    private Dictionary<RedPointConst, GameObject> RedDic;

    private void OnBindRed()
    {
        RedDic = new Dictionary<RedPointConst, GameObject>()
        {
            [RedPointConst.Play] = GetView().Btns[0].objDot,
            [RedPointConst.Role] = GetView().Btns[1].objDot,
            [RedPointConst.Partner_Call_Free] = GetView().Btns[2].objDot,
            [RedPointConst.Bag] = GetView().Btns[4].objDot,
            [RedPointConst.Guild] = GetView().Btns[5].objDot,
            //[RedPointConst.ExpPool_First] = GetView().Btns[6].objDot,
        };

        foreach (var item in RedDic)
        {
            RedPointManager.Instance.RPS.Register(item.Key.ToString(), (node) =>
            {
                if (isNullGo) return;
                bool isLock = false;
                if (item.Key == RedPointConst.Play)
                    isLock = GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Play, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Play, false);
                else if (item.Key == RedPointConst.Role)
                    isLock = GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Role, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Role, false);
                else if (item.Key == RedPointConst.Partner_Call_Free)
                    isLock = GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.CallRole, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.CallRole, false);
                else if (item.Key == RedPointConst.Bag)
                    isLock = GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Bag, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Bag, false);
                else if (item.Key == RedPointConst.Guild)
                    isLock = GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Guild, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Guild, false);

                item.Value.SetActive(!isLock && node.pointCount > 0);
            });
        }

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Bag_Syn);
    }


    private void OnLocalRedCheck()
    {
        if (GuildM.Instance.IsGuildOpen())
            GuildM.Instance.CheckGuildRed();
    }

    public override void OnEnter(params object[] args)
    {
        
    }

    private void OnDownBtnCkick(int idx)
    {
        switch (idx)
        {
            case 1:
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Play))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Play))
                    return;

                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return;
                }
                UITool.Open<UIFightChoice>();
                break;

            case 2:
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Role))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Role))
                    return;
                UITool.Open<UIRole>();
                break;

            case 3:
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.CallRole))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.CallRole))
                    return;

                UITool.Open<UILoading>((ui) =>
                {
                    HttpC.GetPrivilegeDetail(1, (_r) =>
                    {
                        UITool.Close<UILoading>();
                        UITool.Open<UIPartnerCall>(null, _r);
                    }, () => { UITool.Close<UILoading>(); });
                }, new UILoading.Loading { _time = 0 });

                break;

            case 4:
                UITool.Open<UIShop>();
                break;
            case 5:
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Bag))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Bag))
                    return;

                TcpC.GetBagInfo();
                UITool.Open<UIBag>();
                break;
            case 6:
                if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Guild))
                    return;
                if (GameUtil.Instance.CheakIsLock(SwitchOpne.Guild))
                    return;

                UITool.Open<UIGuild>();
                break;
            case 7:
                //if (GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.ExpPool))
                //    return;
                //if (GameUtil.Instance.CheakIsLock(SwitchOpne.ExpPool))
                //    return;

                //UITool.Open<UIExpPool>();
                break;
            default:
                break;
        }

        for (int i = 1; i <= 7; i++)
        {
            var selGo = uIViewCom_DownBtns.Btns[i - 1].selGo;
            if (selGo != null)
                selGo.SetActive(i == idx);
        }
    }

    private void OnRoleRefresh()
    {
        float proExp = (float)RoleM.Instance.Exp / CFG_Role.GetByJobAndLevel(RoleM.Instance.Job, RoleM.Instance.Level).NeedEXP;
        uIViewCom_DownBtns.img_JyPro.fillAmount = proExp;
        if (!UnityEngine.Object.Equals(uIViewCom_DownBtns.Btns[6].effect, null))
            uIViewCom_DownBtns.Btns[6].effect.gameObject.SetActive(proExp >= 1);
        uIViewCom_DownBtns.txt_JyPro.text = GameUtil.Instance.GetPercentStr(proExp);
    }

    private void OnLockCheak()
    {
        for (int i = 0; i < switchOpnes.Length; i++)
        {
            if (switchOpnes[i] == SwitchOpne.No)
                continue;

            var isLockByLvAndStage = GameUtil.Instance.CheakIsLock(switchOpnes[i], false);
            var isLock = GameUtil.Instance.CheckOpenServerIsLock(switchOpnes[i], false);

            bool isLocked = isLockByLvAndStage || isLock;
            string name = uIViewCom_DownBtns.Btns[i].trans.name;
            //Debug.LogError($"类型：{switchOpnes[i].ToString()}，控件：{name}，是否锁定{isLocked.ToString()},curBattleID:{RoleM.Instance.MainBattleId}");
            //int lockLv = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[switchOpnes[i]];
            if (i == 0 || i == 1 || i == 4)
                uIViewCom_DownBtns.Btns[i].iconImg.material = (isLocked ? uIViewCom_DownBtns.grayMat : null);
            else if (i == 6)
            {
                uIViewCom_DownBtns.Btns[i].lockGo.SetActive(false);
                uIViewCom_DownBtns.Btns[i].sliderImg.material = (isLocked ? uIViewCom_DownBtns.grayMat : null);
                uIViewCom_DownBtns.Btns[i].expText.color = (isLocked ? Color.gray : Color.white);
                uIViewCom_DownBtns.Btns[i].expImg.material = (isLocked ? uIViewCom_DownBtns.grayMat : null);
            }
            else if (i == 2 || i == 5)
            {
                uIViewCom_DownBtns.Btns[i].iconImg.material = ((isLockByLvAndStage || isLock) ? uIViewCom_DownBtns.grayMat : null);
            }
        }
    }

    private void LockInfoRefresh()
    {
        for (int i = 0; i < switchOpnes.Length; i++)
        {
            if (i == 0 || i == 1 || i == 4 || i == 6)
            {
                StringBuilder tips = new("");
                bool islockLv = GameUtil.Instance.CheakIsLock(switchOpnes[i], false, tips);
                GetView().Btns[i].mask.gameObject.SetActive(islockLv);
                if (islockLv)
                    GetView().Btns[i].txt_lock.text = tips.ToString();//GameUtil.Instance.CheakIsLock(switchOpnes[i], false) + "级解锁";
            }
            else if (i == 2 || i == 5)
            {
                StringBuilder tips = new("");
                bool isLockServer = GameUtil.Instance.CheckOpenServerIsLock(switchOpnes[i], false);
                bool islockLv = GameUtil.Instance.CheakIsLock(switchOpnes[i], false,tips);
                GetView().Btns[i].mask.gameObject.SetActive(isLockServer || islockLv);
                if (isLockServer)
                    GetView().Btns[i].txt_lock.text = $"{GameUtil.Instance.RemainingDaysUnlockOpen(switchOpnes[i])}天后开启";
                else if (islockLv)
                    GetView().Btns[i].txt_lock.text = tips.ToString();//$"{CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[switchOpnes[i]]}级解锁";
            }
        }
    }

    /// <summary>
    /// 按钮相对于画布的位置
    /// </summary>
    /// <param name="buttonIndex">第几个按钮。 0.第一个按钮</param>
    public Vector2 GetButtonPositionRelativeCanvas(int buttonIndex)
    {
        return JUIExtensions.GetUILocalPositionRelativeCanvas(uIViewCom_DownBtns.Btns[buttonIndex].trans);
    }

    private void OnRedCheak()
    {
        //检测等级
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Into);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Partner_Call_Free);
        
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Reinforce);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Insert);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Refine);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Equip_Star);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_UpLv);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Setting);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Skill_Advance);

        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild);

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.WarFight);

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_Active);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_UpStar);

        //等级+经验
        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.ExpPool);

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.PVP_One.ToString(), RoleM.Instance.GetPVPRedNum());
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.PVP_More.ToString(), RoleM.Instance.GetGVGRedNum());

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Equip);

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.DayDo.ToString(), RedDotM.Instance.OnDayDoRedNum());

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_Award);


        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_Award);
    }

    private void OnPlayRedDot()
    {
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Fight_Click);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Fight_Award);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Boss);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Tower_Click);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_Relice);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Challenge_OnHook);
    }

    [Subscriber(ThreadMode.Main)]
    private void OnEvent_Source_Refresh(GameEvent.Refresh_InstanceInfo _e)
    {
        if (isNullGo)
            return;

        OnPlayRedDot();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnEvent_Source_Refresh(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;

        //OnLockCheak();

        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Bag_Syn);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_UpStar_1);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_Active);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Role_Title_UpStar);
        RedPointManager.Instance.RPS.Dispatch(RedPointConst.Pet_Award);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        OnRoleRefresh();

   
        //RedPointManager.Instance.RPS.Dispatch(RedPointConst.ExpPool);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_LvUp(GameEvent.Refresh_LvUp _e)
    {
        if (isNullGo)
            return;
        OnLockCheak();
        LockInfoRefresh();
        OnPlayRedDot();
        OnRedCheak();
    }

    private void OnRefresh_MainBattleChanged(GameEvent.Refresh_MainBattleChanged changed)
    {
        if (isNullGo)
            return;

        OnLockCheak();
        LockInfoRefresh();
        OnPlayRedDot();
        OnRedCheak();

        NewPlayerGuideWidget.Instance.TryShowGuide();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_CheckOpenTime(GameEvent.Refresh_ZeroTime _e)
    {
        OnLockCheak();
        LockInfoRefresh();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_LvUp>(OnRefresh_LvUp);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_MainBattleChanged>(OnRefresh_MainBattleChanged);
        RedPointManager.Instance.RPS.ClearNodeAndChildCount(RedPointConst.Guild);

        UniGameEvent.RemoveEventListener<GameEvent.UIFunctionOpenOperate>(OnRefresh_FunctionNewFlag);
    }

    public override void OnDestroy()
    {
        OnExit();
    }



    #region 功能开启角标

    private void OnRefresh_FunctionNewFlag(UIFunctionOpenOperate operate)
    {
        refreshFuncNewFlag();
    }

    private List<SwitchOpne> devEntranceGroup = new() { SwitchOpne.Title, SwitchOpne.Build, SwitchOpne.Build_EquipStar, SwitchOpne.Refined, SwitchOpne.Build_EquipInster, SwitchOpne.Build_SkillLv, SwitchOpne.Build_SkillStar, SwitchOpne.Pet };
    private List<SwitchOpne> playEntranceGroup = new() { SwitchOpne.Play, SwitchOpne.Play_Boss, SwitchOpne.Play_SlTower, SwitchOpne.Play_Relice, SwitchOpne.DefendIsland,SwitchOpne.Play_PVP, SwitchOpne.Play_GVG, SwitchOpne.WorldBoss, SwitchOpne.FightEnd, SwitchOpne.Play_QHSAndGold, SwitchOpne.Play_ExpAndSXS, SwitchOpne.Play_SY, SwitchOpne.Play_HS, SwitchOpne.ExpDungeon, SwitchOpne.GoldDungeon, SwitchOpne.WuXing, SwitchOpne.WarFight };
    private List<SwitchOpne> callPartnerEntranceGroup = new() { SwitchOpne.CallRole };
    private List<SwitchOpne> guildEntranceGroup = new() { SwitchOpne.Guild };
    private void refreshFuncNewFlag()
    {
        if (!(GameUtil.Instance.CheckOpenServerIsLock(SwitchOpne.Play, false) || GameUtil.Instance.CheakIsLock(SwitchOpne.Play, false)))
        {
            bool playFlag = FuncOpenShowWidget.Instance.FunctionIsLight(playEntranceGroup);
            GetView().Btns[0].goRedDotRoot.SetActive(!playFlag);
            GetView().Btns[0].goFunctionOpenNewDotRoot.SetActive(playFlag);
        }

        bool devFlag = FuncOpenShowWidget.Instance.FunctionIsLight(devEntranceGroup);
        GetView().Btns[1].goRedDotRoot.SetActive(!devFlag);
        GetView().Btns[1].goFunctionOpenNewDotRoot.SetActive(devFlag);

        bool callPartnerFlag = FuncOpenShowWidget.Instance.FunctionIsLight(callPartnerEntranceGroup);
        GetView().Btns[2].goRedDotRoot.SetActive(!callPartnerFlag);
        GetView().Btns[2].goFunctionOpenNewDotRoot.SetActive(callPartnerFlag);

        bool guildFlag = FuncOpenShowWidget.Instance.FunctionIsLight(guildEntranceGroup);
        GetView().Btns[5].goRedDotRoot.SetActive(!guildFlag);
        GetView().Btns[5].goFunctionOpenNewDotRoot.SetActive(guildFlag);
    }
    #endregion
}