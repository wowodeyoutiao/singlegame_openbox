using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using System;
using System.Collections.Generic;
using UnityEngine.UI;

public class UIGameFightOver : UIBase<UIModelGameFightOver, UIViewGameFightOver>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_mask.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBack();
        });

        GetView().btn_blue.onClick.AddListener(OnBlueBtnClick);

        GetView().btn_yellow.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnBack();
        });

        for (int i = 0; i < GetView().root_Awards1.childCount - 1; i++)
        {
            int idx = i;
            GetView().root_Awards1.GetChild(i + 1).GetComponent<Button>()?.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = gameResult.Items[idx].CfgId });
            });
        }

        for (int i = 0; i < GetView().root_Awards2.childCount - 1; i++)
        {
            int idx = i;
            GetView().root_Awards2.GetChild(i + 1).GetComponent<Button>()?.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = otherAwards[idx].CfgId });
            });
        }
    }

    private void OnBlueBtnClick()
    {
        if (GuidanceM.Instance.isGuiding)
            return;

        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        var RoomType = CFG_Room.GetById(gameResult.RoomId).Roomtype;

        if (GetView().txt_blue.text == "继续挑战" || GetView().txt_blue.text == "下一关")
        {
            if (gameResult.RoomId != 0)
            {
                if (RoomType == CFG_Room.RoomType.PVP_More)
                {
                    TurnBasedFightCtrl.Instance.GoNextFightPVPMore();
                    return;
                }
                else if (RoomType == CFG_Room.RoomType.WuXing)
                {
                    TurnBasedFightCtrl.Instance.GoNextFightWuXing();
                    return;
                }
            }

            int NextBattleId = gameResult.BattleId;
            int Copytype = CFG_Battle.GetById(NextBattleId).Copytype;

            if (Copytype == (int)CFG_Room.RoomType.Fight)
            {
                NextBattleId = CFG_Battle.GetNextFightById(NextBattleId).UID;
            }
            else
            {
                List<CFGItem_Battle> CFG_Battles = Copytype == (int)CFG_Room.RoomType.Boss ? CFG_Battle.BossList : CFG_Battle.GetSLTowerList();
                if (Copytype == (int)CFG_Room.RoomType.Boss)
                    CFG_Battles = CFG_Battle.BossList;
                else if (Copytype == (int)CFG_Room.RoomType.YZDD)
                    CFG_Battles = CFG_Battle.YZDDList;
                else
                    CFG_Battles = CFG_Battle.GetSLTowerList();

                if (Copytype == (int)CFG_Room.RoomType.YZDD)
                {
                    for (int i = 0; i < CFG_Battles.Count; i++)
                    {
                        if (CFG_Battles[i].UID > NextBattleId)
                        {
                            NextBattleId = CFG_Battles[i].UID;
                            break;
                        }
                    }
                }
                else
                {
                    for (int i = (CFG_Battles.Count - 1); i >= 0; i--)
                    {
                        if (CFG_Battles[i].UID > NextBattleId)
                        {
                            NextBattleId = CFG_Battles[i].UID;
                            break;
                        }
                    }
                }
            }
            if (NextBattleId == gameResult.BattleId)
            {
                UITool.ShowTips("您已挑战到最高关卡了！");
                return;
            }

            GameHelper.Instance.StarGame(NextBattleId);
        }
        else if (GetView().txt_blue.text == "再领一倍")
        {
            //if (!GameUtil.Instance.ComCoinToCheck(uesdID, uesdNum, null, false))
            //{
            //    UITool.ShowTips("货币不足");
            //    return;
            //}
            //string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();

            //if (GameUtil.Instance.GetPlayerPrefs("DayDoAward") == timeStr)
            //{
            //    TcpC.C2S_DoubleReward(true);
            //}
            //else
            //{
            //    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            //    {
            //        ShowStr = $"是否确定消耗{GameUtil.Instance.FormatNumWithUnit(uesdNum)}{GameUtil.Instance.GetItemName(uesdID)}进行双倍领取？",
            //        isTogg = true,
            //        TogNoFun = () =>
            //        {
            //            GameUtil.Instance.SetPlayerPrefs("DayDoAward", "");
            //        },
            //        TogYesFun = () =>
            //        {
            //            GameUtil.Instance.SetPlayerPrefs("DayDoAward", timeStr);
            //        },
            //        YesFun = delegate { TcpC.C2S_DoubleReward(true); },
            //    });
            //}
        }
        else if (GetView().txt_blue.text == "赋能")
        {
            if (gameResult.RoomId != 0)
            {
                if (RoomType == CFG_Room.RoomType.WuXing)
                {
                    TurnBasedFightCtrl.Instance.GoWuXingFuNeng();
                    return;
                }
            }
        }
        else
        {
            //TODo 我要变强
            OnBack(() => { UITool.Open<UIStrengthen>(); });
        }
    }

    private void OnBack(Action Jump = null)
    {
        UITool.Close<UIGuide>();
        UITool.Close<UIGameFightOver>();
        if (gameResult.OverType != GameOverType.GameOverTypeCopySweep && gameResult.RoomId != 0)
        {
            if (CFG_Room.GetById(gameResult.RoomId).Roomtype == CFG_Room.RoomType.PVP_More)
            {
                ClearJa();
                TurnBasedFightCtrl.Instance.OnExitGame(() =>
                {
                    UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.PvP);
                });
                return;
            }
            else if (CFG_Room.GetById(gameResult.RoomId).Roomtype == CFG_Room.RoomType.WuXing)
            {
                ClearJa();
                TurnBasedFightCtrl.Instance.OnExitGame(() =>
                {
                    bool isMoveLayer = gameResult.OverType == GameOverType.GameOverTypeMDead;
                    UITool.Open<UIWuXingTrials>(null, new UIWuXingTrials.Params
                    {
                        isPlayMoveLayerAnimation = isMoveLayer,
                    });
                });
                return;
            }
        }

        if (GameRoomM.Instance._GameState == GameRoomM.GameState.Over)
        {
            ClearJa();
            FightController.Instance.OnExitGame(Jump);
        }
    }

    private TimeItem Ja;
    private void ClearJa()
    {
        if (GameHelper.Instance != null && Ja != null)
        {
            GameHelper.Instance.RemoveInvoke(Ja);
            Ja = null;
        }
    }

    GameResultBattlePush gameResult;
    List<ItemInfo> otherAwards;

    //private int uesdID;
    //private int uesdNum;

    public override void OnEnter(params object[] args)
    {
        if (GuidanceM.Instance.isGuiding)
            UniGameEvent.Send(new GameEvent.Refresh_Guide { completeType = GuidanceCompleteType.Battle, value = 1 });

        

        gameResult = (GameResultBattlePush)args[0];
        if (args.Length > 1)
            otherAwards = (List<ItemInfo>)args[1];

        bool isDayDo = gameResult.RoomId != 0 && GameRoomM.Instance.IsDayDo((int)CFG_Room.GetById(gameResult.RoomId).Roomtype);

        bool IsWin = gameResult.OverType == GameOverType.GameOverTypeMDead || gameResult.OverType == GameOverType.GameOverTypeCopySweep || isDayDo;
        bool isSweep = gameResult.OverType == GameOverType.GameOverTypeCopySweep;

        GetView().obj_lose.gameObject.SetActive(!IsWin);
        GetView().obj_win.gameObject.SetActive(IsWin);

        GetView().txt_CloseTip.text = isSweep ? "点击空白处关闭界面" : "点击空白处可退出挑战";

        bool isFightEnd = false;    // 是不是一战到底？
        bool isWuXing = false;      // 是不是五行副本？需要显示首通奖励。
        bool isShowUseTime = false;

        string award1Tip = ""; string award2Tip = "";

        if (IsWin)
        {
            string Str = gameResult.OverType == GameOverType.GameOverTypeCopySweep ? "恭喜扫荡获得" : "恭喜挑战获得";

            if (gameResult.OverType != GameOverType.GameOverTypeCopySweep)
            {
                var RoomType = CFG_Room.GetById(gameResult.RoomId).Roomtype;
                if (gameResult.BattleId != 0)
                {
                    var CFGBattle = CFG_Battle.GetById(gameResult.BattleId);

                    if (RoomType == CFG_Room.RoomType.Fight)
                        Str = "已通关:战役" + CFGBattle.Chapter + "-" + CFGBattle.Layers + "\n" + Str;
                    else if (RoomType == CFG_Room.RoomType.Tower)
                        Str = "已通关:试炼塔" + CFGBattle.Layers + "层" + "\n" + Str;
                    else if (RoomType == CFG_Room.RoomType.Boss)
                        Str = "已通关:首领" + CFGBattle.Layers + "层" + "\n" + Str;
                    else if (RoomType == CFG_Room.RoomType.Relic)
                        Str = "已通关:遗迹" + CFGBattle.Chapter + "-" + CFGBattle.Layers + "\n" + Str;
                    else if (RoomType == CFG_Room.RoomType.YZDD)
                    {
                        isFightEnd = true;
                        if (gameResult.BattleId <= FightPlayM.Instance.GetInfo().to_end.BattleId)
                            Str = "已通关:一战到底 第" + CFGBattle.Layers + "关";
                        else
                            Str = "已通关:一战到底 第" + CFGBattle.Layers + "关" + "\n" + Str;
                    }
                }
                else
                {
                    if (RoomType == CFG_Room.RoomType.WuXing)
                    {
                        Str = "";
                        isWuXing = true;
                        award1Tip = "当前挑战已获得物品";
                        award2Tip = "首通奖励";
                    }
                }
            }

            isShowUseTime = IsWin && gameResult.OverType != GameOverType.GameOverTypeCopySweep && gameResult.BattleId != 0 && isFightEnd;

            GetView().root_useTime.gameObject.SetActive(isShowUseTime);
            if (isShowUseTime)
            {
                GetView().txt_useTime.text = GameUtil.Instance.GetTimeToHMS(gameResult.CompletionTime);
                GetView().obj_newRecord.gameObject.SetActive(FightPlayM.Instance.IsFightEndNewRecord(gameResult.BattleId, gameResult.CompletionTime));
            }

            GetView().txt_WinDesc.text = Str;
        }
        else
        {
            string Str;
            var RoomType = CFG_Room.GetById(gameResult.RoomId).Roomtype;

            if (RoomType == CFG_Room.RoomType.WuXing)
            {
                isWuXing = true;
                Str = "很遗憾本次挑战失败，请提升元素属性之后再来挑战！";
            }
            else
            {
                Str = "很遗憾本次挑战失败\n请继续努力!";
            }
            GetView().txt_LoseDesc.text = Str;
        }

        GetView().txt_Awards1Tip.text = award1Tip;
        GetView().txt_Awards2Tip.text = award2Tip;

        bool isTurnBased = gameResult.RoomId != 0 && CFG_Room.GetById(gameResult.RoomId).Roomtype == CFG_Room.RoomType.PVP_More;

        // 一战到底通关是否收到过这个关卡的奖励
        bool isFightEndReceivedAwards = isFightEnd && gameResult.BattleId <= FightPlayM.Instance.GetInfo().to_end.BattleId;

        // 五行试炼需要显示首通奖励和分割线
        GetView().obj_line.SetActive(isWuXing && IsWin);
        GetView().obj_space.gameObject.SetActive(!isTurnBased && !isShowUseTime);

        if (isSweep)
        {
            GetView().btn_blue.gameObject.SetActive(false);
            GetView().txt_point.gameObject.SetActive(false);
        }
        else if (isDayDo)
        {
            //GetView().txt_blue.text = "再领一倍";
            GetView().btn_blue.gameObject.SetActive(false);
            GetView().txt_num.gameObject.SetActive(false);
            GetView().txt_point.gameObject.SetActive(false);

            //uesdID = CFG_XMLConfig.xmlConfig.DayRoomDic[CFG_Room.GetById(gameResult.RoomId).Roomtype].doubleaccessID;
            //uesdNum = CFG_XMLConfig.xmlConfig.DayRoomDic[CFG_Room.GetById(gameResult.RoomId).Roomtype].doubleaccessNum;

            //GetView().txt_num.text = GameUtil.Instance.GetCheakNumStr(uesdID, uesdNum,
            //    "x" + GameUtil.Instance.FormatNumWithUnit(uesdNum));
        }
        else if (isTurnBased)
        {
            GetView().btn_blue.gameObject.SetActive(false);
            GetView().txt_point.gameObject.SetActive(true);
            string score = TurnBasedM.Instance.turnBasedBattleData.Result.Score < 0 ?
               "积分" + TurnBasedM.Instance.turnBasedBattleData.Result.Score.ToString() : "积分+" + TurnBasedM.Instance.turnBasedBattleData.Result.Score;
            GetView().txt_point.text = score;
            //GetView().txt_num.gameObject.SetActive(false);
        }
        else if (isWuXing)
        {
            if (IsWin)
                GetView().txt_blue.text = "下一关";
            else
                GetView().txt_blue.text = "赋能";

            GetView().btn_blue.gameObject.SetActive(true);
            GetView().txt_point.gameObject.SetActive(false);
            GetView().txt_num.gameObject.SetActive(false);
        }
        else
        {
            GetView().txt_point.gameObject.SetActive(false);
            GetView().txt_num.gameObject.SetActive(false);
            string nextStr = isFightEnd ? "下一关" : "继续挑战";
            GetView().txt_blue.text = gameResult.OverType == GameOverType.GameOverTypeMDead ? nextStr : "我要变强";
            bool isShowNextButton = !isFightEnd || CFG_Battle.YZDDList[CFG_Battle.YZDDList.Count - 1].UID > gameResult.BattleId;
            if (CFG_Battle.GetById(gameResult.BattleId).Copytype == (int)CFG_Room.RoomType.Relic)
                GetView().btn_blue.gameObject.SetActive(gameResult.OverType != GameOverType.GameOverTypeCopySweep && isShowNextButton && !IsWin);
            else
                GetView().btn_blue.gameObject.SetActive(gameResult.OverType != GameOverType.GameOverTypeCopySweep && isShowNextButton);
        }

        // 一战到底已通关此关卡，并且已收到过此通关奖励。
        if (isFightEndReceivedAwards)
        {
            CFGItem_Battle.AwardItem[] items = CFG_Battle.GetById(gameResult.BattleId).Settlement;
            for (int i = 0; i < GetView().root_Awards1.childCount - 1; i++)
            {
                var go = GetView().root_Awards1.GetChild(i + 1);
                go.gameObject.SetActive(items != null && i < items.Length);
                if (items != null && i < items.Length)
                {
                    GameUtil.Instance.SetItemIcon(items[i].id, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform.Find("icon/nature")));
                    JUIExtensions.txt(go, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(items[i].num);
                    go.Find("received").gameObject.SetActive(true);
                }
            }
        }
        else // 正常的通关奖励
        {
            for (int i = 0; i < GetView().root_Awards1.childCount - 1; i++)
            {
                var go = GetView().root_Awards1.GetChild(i + 1);
                go.gameObject.SetActive(gameResult.Items != null && i < gameResult.Items.Count);
                if (gameResult.Items != null && i < gameResult.Items.Count)
                {
                    GameUtil.Instance.SetItemIcon(gameResult.Items[i].CfgId, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform.Find("icon/nature")));
                    JUIExtensions.txt(go, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(gameResult.Items[i].Num);
                    go.Find("received").gameObject.SetActive(false);
                }
            }
        }

        if (isWuXing)
        {
            for (int i = 0; i < GetView().root_Awards2.childCount - 1; i++)
            {
                var go = GetView().root_Awards2.GetChild(i + 1);
                go.gameObject.SetActive(otherAwards != null && i < otherAwards.Count);
                if (otherAwards != null && i < otherAwards.Count)
                {
                    GameUtil.Instance.SetItemIcon(otherAwards[i].CfgId, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform.Find("icon/nature")));
                    JUIExtensions.txt(go, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(otherAwards[i].Num);
                    go.Find("received").gameObject.SetActive(false);
                }
            }
        }

        GetView().txt_yellow.text = gameResult.OverType == GameOverType.GameOverTypeCopySweep ? "确定" : "退出";


        GetView().root_LuckAward.gameObject.SetActive(gameResult.FateCoinReward != null && gameResult.FateCoinReward.Num > 0);
        if (gameResult.FlipCardReward != null)
        {
            GetView().txt_Time.text = "";

            
            UITool.Open<UIGameOverAwards>(null, gameResult.BattleId, gameResult.FlipCardReward.CfgId, gameResult.FlipCardReward.Num);
        }
        //if (!GuidanceM.Instance.isGuiding)
        else
        {
            string StrTime = gameResult.OverType == GameOverType.GameOverTypeCopySweep ? "秒后将自动关闭" : "秒后将自动退出";
            int Count = 10;
            Ja = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () =>
            {
                if (isNullGo) return;

                GetView().txt_Time.text = Count + StrTime;
                Count--;
                if (Count <= 0)
                    OnBack();
            }));
        }

        if (gameResult.BattleId > 0 && CFG_Battle.GetById(gameResult.BattleId).Copytype == (int)CFG_Room.RoomType.Boss && GameRoomM.Instance._GameState != GameRoomM.GameState.No)
        {
            TcpC.C2S_InstanceInfo();
        }

        NewPlayerGuideWidget.Instance.AddGuideGameObject("GameFightOver_ExitGame", GetView().btn_yellow.transform);
    }

    public override void OnExit()
    {
        ClearJa();
        otherAwards = null;
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}