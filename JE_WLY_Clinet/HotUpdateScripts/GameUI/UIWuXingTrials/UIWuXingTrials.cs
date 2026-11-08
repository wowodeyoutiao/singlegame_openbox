using DG.Tweening;

using HotUpdateScripts.Config;
using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using SuperScrollView;

using System.Collections.Generic;
using System.Linq;

using UnityEngine;

public class UIWuXingTrials : UIBase<UIModelWuXingTrials, UIViewWuXingTrials>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().loopLayers.InitListView(0, OnItemRefresh);

        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 45);
        });

        GetView().btn_Rank.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIRank>(null, CFG_RankingAward.RankClass.WuXing);
        });

        GetView().btn_Gift.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGiftBag>(null, ActivityType.Element);
        });

        GetView().btn_FuNeng.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIFuNeng>();
        });

        GetView().btn_AddChallengeTimes.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameUtil.Instance.CommonChallengeCountExchange(CFG_Item.ItemID.WXTZJZ);
        });

        GetView().btn_Challenge.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (!isCanChallenge) return;
            OnChallenge();
        });

        GetView().btn_Sweep.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnSweep();
        });

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIWuXingTrials>();
        });
    }

    private Color disabledTowerItemColor => new Color(0.5f, 0.5f, 0.5f);
    private Color disabledTowerTextColor => new Color(0.6f, 0.6f, 0.6f);

    private readonly float[] cloudsDefaultPosX = { -256, 228, -40, -342, 28 };
    private readonly float[] cloudsMoveToPosX = { -722, 580, -733, -578, 730 };

    const int itemHeight = 258;
    const int scrollOffset = 435;

    // 总层数
    int totalConfigLayers;
    // 可用的层数
    int availableLayers;
    // 最后一次通关的层数
    int passedLayer;

    bool isPlayMoveLayer;
    // 是否继续下一关卡？ true, 打开布阵界面
    bool isContinueNextLevel;
    // 是否打开赋能界面
    bool isOpenFuNengUI;
    // 是否需要扫荡
    bool isNeedSweep;
    // 是否完成了所有的挑战
    bool isCompletedChallenge => passedLayer >= availableLayers;

    Sequence layerMoveSeq, cloudsMoveSeq;
    List<Tweener> sweepTweens = new List<Tweener>();

    public struct Params
    {
        public bool isPlayMoveLayerAnimation;
        public bool isContinueChallenge;
        public bool isOpenFuNeng;
    }

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        RegisterEvent();

        if (args.Length > 0)
        {
            Params _params = (Params)args[0];
            isPlayMoveLayer = _params.isPlayMoveLayerAnimation;
            isContinueNextLevel = _params.isContinueChallenge;
            isOpenFuNengUI = _params.isOpenFuNeng;
        }

        availableLayers = CFG_XMLConfig.xmlConfig.Element.copynum;
        totalConfigLayers = CFG_ElementTrials.GetLayerCount();
        ResetCloudsPosition();
        GetView().levelDescCanvas.alpha = 0;
        GetView().btn_Challenge.transform.localScale = Vector3.one;
        GetView().btn_Challenge.gameObject.SetActive(true);
        GetView().btn_Sweep.gameObject.SetActive(false);
        GetView().scrollTower.vertical = false;
        GetView().loopLayers.SetListItemCount(0);
        GetView().loopLayers.RefreshAllShownItem();

        ShowResetPlayTime();
        TcpC.C2S_ElementInfoRequest();
        OnRankDataRequest();
    }

    #region 塔层

    private LoopListViewItem2 OnItemRefresh(LoopListView2 list, int index)
    {
        var item = list.NewListViewItem("Item");
        Transform t = item.transform;
        t.Find("img_level").gameObject.SetActive(index > 0);
        t.Find("icons").gameObject.SetActive(index > 0);

        JUIExtensions.img(t, "img_Layer").color = index > passedLayer ? Color.white : disabledTowerItemColor;
        JUIExtensions.img(t, "img_level").color = index > passedLayer ? Color.white : disabledTowerItemColor;
        JUIExtensions.txt(t, "img_level/txt").color = index > passedLayer ? Color.white : disabledTowerTextColor;
        JUIExtensions.txt(t, "img_level/txt").text = $"第{index}层";

        var img_passed = JUIExtensions.img(t, "icons/icon_passed");
        var img_lock = JUIExtensions.img(t, "icons/icon_lock");
        var skeleton = JUIExtensions.skeleton(t.Find("icons/spine_unlock"));
        var btn_sweep = JUIExtensions.btn(t, "icons/btn_sweep");
        var canvasg_sweep = JUIExtensions.canvasGroup(t, "icons/btn_sweep");

        btn_sweep.onClick.RemoveAllListeners();
        btn_sweep.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (TurnBasedM.Instance.elementTimes < 1)
            {
                UITool.ShowTips("挑战次数不足");
                return;
            }

            UITool.Open<UIGameSweepAndBuy>(null, UIGameSweepAndBuy.OpType.WuXingSweep);
        });

        img_passed.color = Color.white;
        skeleton.color = Color.white;

        bool showSweep = index == passedLayer;
        bool showPassed = index < passedLayer;
        bool showSpine = index == passedLayer + 1;
        bool showNoPassed = index > passedLayer + 1;

        if (isNeedSweep)
        {

            if (showNoPassed)
            {
                skeleton.gameObject.SetActive(false);
                btn_sweep.gameObject.SetActive(false);

            }
            else if (showSpine)
            {
                skeleton.gameObject.SetActive(!isCompletedChallenge);
                btn_sweep.gameObject.SetActive(isCompletedChallenge);

                if (!isCompletedChallenge)
                {
                    if(TurnBasedM.Instance.needSweep == 1)
                        SpineUtil.Instance.PlayLoop(skeleton, "loop");
                    else
                    {
                        SpineUtil.Instance.PlayOnce(skeleton, "start", () =>
                        {
                            SpineUtil.Instance.PlayLoop(skeleton, "loop");
                        });
                    }
                }
            }
            else if (showSweep)
            {
                skeleton.color = Color.white;
                skeleton.gameObject.SetActive(false);
                canvasg_sweep.alpha = 0;
                btn_sweep.gameObject.SetActive(true);

                //SpineUtil.Instance.PlayOnce(skeleton, "start", null);
                Tweener tweener1 = skeleton.DOColor(new Color(1, 1, 1, 0), 0.2f).SetDelay(0.45f);
                Tweener tweener2 = canvasg_sweep.DOFade(1, 0.2f).SetDelay(0.65f);

                AddLayerSweepTween(tweener1);
                AddLayerSweepTween(tweener2);
            }
            else if (showPassed)
            {
                skeleton.gameObject.SetActive(false);
                btn_sweep.gameObject.SetActive(false);
            }

            img_passed.gameObject.SetActive(showPassed);
            img_lock.gameObject.SetActive(showNoPassed);
        }
        else
        {
            skeleton.gameObject.SetActive(showSpine && !isCompletedChallenge);

            if (showSpine && !isCompletedChallenge)
            {
                if (isPlayMoveLayer)
                {
                    SpineUtil.Instance.PlayOnce(skeleton, "start", () =>
                    {
                        SpineUtil.Instance.PlayLoop(skeleton, "loop");
                    });
                }
                else
                {
                    SpineUtil.Instance.PlayLoop(skeleton, "loop");
                }
            }

            btn_sweep.gameObject.SetActive(showSweep);
            img_passed.gameObject.SetActive(showPassed);
            img_lock.gameObject.SetActive(showNoPassed || (showSpine && isCompletedChallenge));
        }

        return item;
    }

    private void RegisterEvent()
    {
        UniGameEvent.AddEventListener<GameEvent.WuXingTrialsInfo>(Sync_ElementInfo);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_ZeroRefresh);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnEvent_Sweep);
    }

    private void UnregisterEvent()
    {
        UniGameEvent.RemoveEventListener<GameEvent.WuXingTrialsInfo>(Sync_ElementInfo);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_ZeroRefresh);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnEvent_Sweep);
    }

    private void Sync_ElementInfo(GameEvent.WuXingTrialsInfo _e)
    {
        if (isNullGo) return;
        //Debug.Log("<color=cyan>五行试炼 刷新data</color>");
        isNeedSweep = TurnBasedM.Instance.needSweep == 1;

        if (isNeedSweep)
            passedLayer = 0;
        else
            passedLayer = TurnBasedM.Instance.currentPassedLayer;

        GetView().awardRoot.parent.transform.localScale = Vector3.one;
        GetView().awardRoot.parent.gameObject.SetActive(true);

        ShowAward();
        ShowRemainChallengeTimes();
        ShowButton();
        ShowResetPlayTime();

        GetView().loopLayers.SetListItemCount(totalConfigLayers + 1);
        GetView().loopLayers.RefreshAllShownItem();

        LocateLayer();

        if (isOpenFuNengUI)
        {
            UITool.Open<UIFuNeng>();
            isOpenFuNengUI = false;
        }
    }

    private void OnEvent_Sweep(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        if (_e.code == 0 && _e.cid == CmdId.C2S_ElementSweep)
        {
            //Debug.Log("扫荡成功！");
            TurnBasedM.Instance.needSweep = 2;
            OnRankDataRequest();
            PlaySweep();
        }
    }

    private void OnSync_ZeroRefresh(GameEvent.Refresh_ZeroTime _e)
    {
        ShowResetPlayTime();
        TcpC.C2S_ElementInfoRequest();
        OnRankDataRequest();
        //HttpC.OnGetShopInfo(YesFun: (v) =>
        //{
        //    RefreshAddButton();
        //});
    }

    /// <summary>
    /// 显示当前层的怪物信息
    /// </summary>
    /// <param name="fadeInFadeOut">播放淡入淡出动画</param>
    private void ShowLayerInfo(bool fadeInFadeOut = false)
    {
        int redLayer = passedLayer;
        if (passedLayer < availableLayers)
            redLayer++;

        Transform layerItem = GetView().loopLayers.GetShownItemByItemIndex(redLayer).transform;
        Vector2 itemPos = JUIExtensions.GetUILocalPositionRelativeCanvas(layerItem);
        GetView().txt_LevelDesc.text = $"<color=#ffd41c><size=23>通关小建议：</size></color>\n\n{CFG_ElementTrials.GetByLayer(redLayer).Suggestion}";
        GetView().levelDescCanvas.transform.localPosition = new Vector2(-210, itemPos.y + itemHeight * 0.5f);

        if (fadeInFadeOut)
            GetView().levelDescCanvas.DOFade(1, 0.5f);
        else
            GetView().levelDescCanvas.alpha = 1;
    }

    // 定位
    private void LocateLayer()
    {
        if(TurnBasedM.Instance.needSweep == 1)
        {
            // 从第一层开始扫荡
            isPlayMoveLayer = false;
            GetView().loopLayers.MovePanelToItemIndex(1, scrollOffset);
            ShowLayerInfo(false);
            return;
        }

        if (isPlayMoveLayer)
        {
            GetView().loopLayers.MovePanelToItemIndex(passedLayer, scrollOffset);
            if(passedLayer < availableLayers)
            {
                isCanChallenge = false;
                GetView().btn_Challenge.transform.localScale = Vector3.zero;
                // 等解锁动画播放完后开始移动
                ClearMoveLayerItem();

                MoveClouds();
                moveLayerItem = GameHelper.Instance.AddInvoke(new TimeItem(0.5f, () =>
                {
                    PlayMoveToLayer(passedLayer + 1);
                    isPlayMoveLayer = false;
                }));

                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_TowerUnlock);
            }
            else // 已经通关所有层,永远停止在最后一层
            {
                isCanChallenge = true;
                GetView().btn_Challenge.transform.localScale = Vector3.one;
                ShowLayerInfo(false);
            }
        }
        else
        {
            isCanChallenge = true;
            GetView().btn_Challenge.transform.localScale = Vector3.one;

            if (passedLayer < availableLayers)
                GetView().loopLayers.MovePanelToItemIndex(passedLayer + 1, scrollOffset);
            else
                GetView().loopLayers.MovePanelToItemIndex(passedLayer, scrollOffset);

            ShowLayerInfo(false);
        }
    }

    TimeItem moveLayerItem;

    private void ClearMoveLayerItem()
    {
        if (GameHelper.Instance != null && moveLayerItem != null)
        {
            GameHelper.Instance.RemoveInvoke(moveLayerItem);
            moveLayerItem = null;
        }
    }

    public float NeedMoveOffset(int itemIndex, float offset = scrollOffset)
    {
        ItemPosStruct itemPos = GetView().loopLayers.GetFirstShownItemIndexAndOffset();
        float itemOffset = (itemIndex - itemPos.mItemIndex) * itemHeight;
        float moveOffset = itemOffset + itemPos.mItemOffset - offset;
        moveOffset = Mathf.Max(0, moveOffset);
        return -moveOffset;
    }

    public Tweener MoveToLayer(int layer)
    {
        float needOffset = NeedMoveOffset(layer);
        return GameUtil.Instance.MoveScrollViewContentByOffset(GetView().loopLayers, needOffset, 1.5f);
    }

    bool isCanChallenge;

    private void MoveClouds()
    {
        StopCloudsTween();
        cloudsMoveSeq = DOTween.Sequence();
        MoveClouds(false);
        MoveClouds(true);
    }

    public void MoveClouds(bool back)
    {
        for (int i = 0; i < 5; i++)
        {
            float posx, insertPos;
            Ease ease;
            if (back)
            {
                ease = Ease.OutQuad;
                insertPos = 2f;
                posx = cloudsDefaultPosX[i];
            }
            else
            {
                ease = Ease.InQuad;
                insertPos = 0.35f;
                posx = cloudsMoveToPosX[i];
            }

            cloudsMoveSeq.Insert(insertPos, GetView().clouds[i].DOLocalMoveX(posx, 1f).SetEase(ease));
        }
    }

    private void PlayMoveToLayer(int layer, float showPlayBtnDelay = 0.1f)
    {
        StopLayerTween();
        layerMoveSeq = DOTween.Sequence();
        GetView().btn_Challenge.transform.localScale = Vector3.zero;

        layerMoveSeq.Append(MoveToLayer(layer));
        if (isNeedSweep)
        {
            // 扫荡动画，解锁状态改变
            layerMoveSeq.InsertCallback(1.5f, () =>
            {
                passedLayer = TurnBasedM.Instance.currentPassedLayer;
                
                ShowAward();
                StopLayerSweepTween();
                GetView().awardRoot.parent.transform.localScale = Vector3.zero;
                GetView().awardRoot.parent.gameObject.SetActive(true);
                AddLayerSweepTween(GetView().awardRoot.parent.DOScale(1, 0.25f));
                GetView().loopLayers.RefreshAllShownItem();
            });
        }
        layerMoveSeq.AppendCallback(() =>
        {
            ShowLayerInfo(true);
        });
        layerMoveSeq.AppendInterval(showPlayBtnDelay);
        layerMoveSeq.Append(GetView().btn_Challenge.transform.DOScale(1, 0.3f).SetEase(Ease.OutQuad));
        layerMoveSeq.AppendInterval(0.3f);
        layerMoveSeq.AppendCallback(() =>
        {
            isCanChallenge = true;
            // 继续下一关
            if (!isNeedSweep && isContinueNextLevel && passedLayer < availableLayers)
            {
                if (isCompletedChallenge)
                {
                    UITool.ShowTips("已达到当前最高层");
                    return;
                }

                if (TurnBasedM.Instance.elementTimes < 1)
                {
                    UITool.ShowTips("挑战次数不足");
                    return;
                }

                int layerId = TurnBasedM.Instance.currentPassedLayer + 1;
                UITool.Open<UIFightFormation>(null, CFG_Room.RoomType.WuXing);
            }

            if (isNeedSweep) isNeedSweep = false;
        });
    }

    private void ResetCloudsPosition()
    {
        GetView().clouds[0].anchoredPosition3D = new Vector3(-256, 49, 0);
        GetView().clouds[1].anchoredPosition3D = new Vector3(228, -107, 0);
        GetView().clouds[2].anchoredPosition3D = new Vector3(-40, -272, 0);
        GetView().clouds[3].anchoredPosition3D = new Vector3(-342, 117, 0);
        GetView().clouds[4].anchoredPosition3D = new Vector3(28, -9.5f, 0);
    }

    private void AddLayerSweepTween(Tweener tweener)
    {
        sweepTweens.Add(tweener);
    }

    private void StopLayerSweepTween()
    {
        foreach (var tweener in sweepTweens)
        {
            tweener.Kill();
        }
        sweepTweens.Clear();
    }

    private void StopLayerTween()
    {
        if (layerMoveSeq != null)
        {
            layerMoveSeq.Kill();
            layerMoveSeq = null;
        }
    }

    private void StopCloudsTween()
    {
        if (cloudsMoveSeq != null)
        {
            cloudsMoveSeq.Kill();
            cloudsMoveSeq = null;
        }
    }

    #endregion

    #region 排行

    private HTTP_RankDetail RankListInfo;

    private void OnRankDataRequest()
    {
        HttpC.GetRankList(CFG_RankingAward.RankClass.WuXing, (_r) =>
        {
            if (isNullGo)
                return;

            RankListInfo = _r.data;
            OnRankShow();
        });
    }

    private void OnRankShow()
    {
        for (int i = 0; i < 3; i++)
        {
            bool have = i < RankListInfo.rank_all.Length;
            GetView().rankItems[i].txt_Name.text = have ? RankListInfo.rank_all[i].name : "虚位以待";
            GetView().rankItems[i].txt_Layer.text = have ? $"{RankListInfo.rank_all[i].value / 100000000000}层" : "0层";
        }
    }

    #endregion

    #region 奖励预览
    private void ShowAward()
    {
        int redLayer = passedLayer;
        if (passedLayer < availableLayers)
            redLayer++;

        CFGItem_ElementTrials.Item[] awards = CFG_ElementTrials.GetByLayer(redLayer).Settlement;
        for (int i = 0; i < GetView().awardRoot.childCount; i++)
        {
            int index = i;
            Transform item = GetView().awardRoot.GetChild(i);
            bool isShow = i < awards.Length;
            item.gameObject.SetActive(i < awards.Length);

            var btn = JUIExtensions.btn(item);
            btn.onClick.RemoveAllListeners();
            if (isShow) 
            {
                GameUtil.Instance.SetItemIcon(awards[i].id, JUIExtensions.img(item.Find("icon")), JUIExtensions.img(item));
                JUIExtensions.txt(item, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(awards[i].num);
                btn.onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = awards[index].id });
                });
            }
        }
    }

    /// <summary>
    /// 扫荡奖励
    /// </summary>
    /// <param name="layer">层数</param>
    public List<ItemInfo> SweepAwards(int layer)
    {
        Dictionary<int, ItemInfo> awards = new Dictionary<int, ItemInfo>();

        for (int i = 1; i <= layer; i++)
        {
            var items = CFG_ElementTrials.GetByLayer(i).Settlement;
            foreach (var item in items)
            {
                if (awards.ContainsKey(item.id))
                {
                    awards[item.id].Num += item.num;
                }
                else
                {
                    ItemInfo info = new ItemInfo
                    {
                        CfgId = item.id,
                        Num = item.num
                    };
                    awards.Add(item.id, info);
                }
            }
        }

        return awards.Values.ToList();
    }

    #endregion

    #region 挑战信息

    /// <summary>
    /// 剩余挑战次数。
    /// </summary>
    private void ShowRemainChallengeTimes()
    {
        string colorStr = TurnBasedM.Instance.elementTimes > 0 ? "6FF37D" : "DA4E38";
        GetView().txt_ChallengeTimes.text = $"今日剩余挑战次数  <color=#{colorStr}>{TurnBasedM.Instance.elementTimes}</color>";
    }

    private void ShowButton()
    {
        GetView().img_Challenge.SetGray(isCompletedChallenge);
        GetView().btn_Challenge.gameObject.SetActive(TurnBasedM.Instance.needSweep != 1);
        GetView().btn_Sweep.gameObject.SetActive(TurnBasedM.Instance.needSweep == 1);
        //RefreshAddButton();
    }

    //private void RefreshAddButton()
    //{
    //    bool canBuy = TurnBasedM.Instance.GetCanBuyChallengeCount_WuXing() > 0;
    //    GetView().btn_AddChallengeTimes.interactable = canBuy;
    //}

    private void AddChallenge()
    {
        int roomType = (int)CFG_Room.GetById(TurnBasedM.wuXingRoomId).Roomtype;
        UITool.Open<UIPVPCountBuy>(null, roomType);
    }

    TimeItem resetPlayTimeItem;

    private void ClearResetTime()
    {
        if (GameHelper.Instance != null && resetPlayTimeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(resetPlayTimeItem);
            resetPlayTimeItem = null;
        }
    }

    private void ShowResetPlayTime()
    {
        ClearResetTime();

        int seconds = GameUtil.Instance.GetNextTimeSecond(CFG_XMLConfig.xmlConfig.Element.starttime, CFG_XMLConfig.xmlConfig.Element.day) + 1;
        resetPlayTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(seconds, 1, () => 
        {
            if (isNullGo) return;

            seconds--;
            GetView().txt_ResetTime.text = $"距离玩法重置剩余：<color=#FFFD44>{GameUtil.Instance.GetTimeToDHM(seconds)}</color>";
            if (seconds == 0)
            {
                ClearResetTime();
                GetView().txt_ResetTime.text = "";
                TcpC.C2S_ElementInfoRequest();
                OnRankDataRequest();
            }
        },
        true));
    }

    private void OnChallenge()
    {
        if (isCompletedChallenge)
        {
            UITool.ShowTips("已达到当前最高层");
            return;
        }

        if (TurnBasedM.Instance.elementTimes < 1)
        {
            UITool.ShowTips("挑战次数不足");
            return;
        }

        int layerId = TurnBasedM.Instance.currentPassedLayer + 1;
        int levelLimit = CFG_ElementTrials.GetById(layerId).Levellimit;

        if(RoleM.Instance.Level < levelLimit)
        {
            UITool.ShowTips($"角色{levelLimit}级可挑战");
            return;
        }

        UITool.Open<UIFightFormation>(null, CFG_Room.RoomType.WuXing);
    }

    private void OnSweep()
    {
        GetView().levelDescCanvas.DOFade(0, 0.3f);
        TcpC.C2S_ElementSweepRequest();


        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "五行试炼",
            position = "",
            event_click = "扫荡",
        });
    }

    private void PlaySweep()
    {
        var awards = SweepAwards(TurnBasedM.Instance.currentPassedLayer);
        GetView().btn_Sweep.gameObject.SetActive(false);
        GetView().btn_Challenge.gameObject.SetActive(true);
        GetView().awardRoot.parent.gameObject.SetActive(false);

        MoveClouds();
        PlayMoveToLayer(TurnBasedM.Instance.currentPassedLayer + 1, 0.5f);
        layerMoveSeq.InsertCallback(1.6f, () =>
        {
            var gameResult = new GameResultBattlePush
            {
                RoomId = TurnBasedM.Instance.roomId,
                BattleId = 0,
                OverType = GameOverType.GameOverTypeCopySweep,
                Items = awards
            };
            UITool.Open<UIGameFightOver>(null, gameResult);
        });
    }

    #endregion

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UnregisterEvent();

        ClearMoveLayerItem();
        StopLayerTween();
        StopCloudsTween();
        StopLayerSweepTween();
        ClearResetTime();
        isCanChallenge = false;

        isPlayMoveLayer = false;
        isContinueNextLevel = false;
        isOpenFuNengUI = false;
    }

    public override void OnDestroy()
    {
        ClearMoveLayerItem();
        ClearResetTime();
    }
}