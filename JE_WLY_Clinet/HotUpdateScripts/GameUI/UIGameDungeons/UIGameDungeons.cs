using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using UnityEngine;

public class UIGameDungeons : UIBase<UIModelGameDungeons, UIViewGameDungeons>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);

        SetJoystickInit();
        GetView().btnExti.onClick.AddListener(exitDungeons);

        for (int i = 0; i < 13; i++)
        {
            int SkillIdx = i;

            var touch = JUIExtensions.touch(GetView().SkillBtns[i].rect);
            touch.AddClickListener((p) =>
            {
                if (UITool.Check<UIItemShow>())
                {
                    return;
                }
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnSkillClick(SkillIdx);
            });

            touch.AddPointerDownListener((p) =>
            {
                int cfgId = SkillM.Instance.GetSkillEquipBySolt(SkillIdx);
                if (cfgId != 0)
                    _TouchTime = GameHelper.Instance.AddInvoke(new TimeItem(GameConst.TouchTime, () =>
                    {
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = (int)CFG_Skill.GetById(cfgId).SkillIcon });
                    }));
            });

            touch.AddPointerUpListener((p) =>
            {
                ClearTouch();
            });

            touch.AddDragBeginListener(OnDragBegin);
            touch.AddDragEndListener(OnDragEnd);
        }

        GetView().btn_SkillPanel.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnSkillPanelOpen(!GetView().obj_SkillPanelOn.activeSelf);
        });

        GetView().btn_Attack.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGameFightSelect>();
        });

        GetView().btn_AutoPlay.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnAutoPlayClick(!GameRoomM.Instance.AutoPlay);
        });


        GetView().touch_Skill.AddDragBeginListener(OnDragBegin);
        GetView().touch_Skill.AddDragEndListener(OnDragEnd);

        GetView().btn_FightAward.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIFightAward>(null);
        });

        OnSkillShow();
        OnAutoLockSet();
        OnAutoPlayClick(GameRoomM.Instance.AutoPlay);

        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More)
            UITool.LoadUII<UICom_GVGBattle>(GetView().root_Base, (ui) =>
            {
                _uICom_GVGBattle = ui as UICom_GVGBattle;
            });
        else if (GameRoomM.Instance.IsDayDo())
            UITool.LoadUII<UIGameDayDoPro>(GetView().root_Base, (ui) =>
            {
                _uIGameDayDoPro = ui as UIGameDayDoPro;
            });

        // else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Guide)
        //     RoleM.Instance.SetGuide(new UIGuide.Guide { ID = 3, GuidePos = Vector3.zero });

        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Relic)
            UITool.LoadUII<UIGameRelicPro>(GetView().root_Base, (ui) =>
            {
                _uIUIGameRelicPro = ui as UIGameRelicPro;

                _uIUIGameRelicPro.OnRelicInfoShow();
            });
        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Guild)
            UITool.LoadUII<UIGameGuildPro>(GetView().root_Base, (ui) =>
            {
                _uIGameGuildPro = ui as UIGameGuildPro;
            });


        GameRoomM.Instance.ShortAttactToPlayer = GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_One;

        GetView().obj_Touch.SetActive(GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.PVP_One);
        //UITool.LoadUII<UICom_DownBtns>(GetView().root_Base, (ui) =>
        //{
        //    DownUI = ui as UICom_DownBtns;
        //});

        GetView().btn_FightAward.gameObject.SetActive(GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_Peace || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace);

        GetView().btn_Chat.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIChat>(null, 2);
        });

        GetView().btn_Chat.gameObject.SetActive(GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Guild || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace);


        GetView().btn_Pos.gameObject.SetActive(GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace);
        GetView().btn_Pos.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnPosSend();
        });
    }

    private void OnPosSend()
    {
        if (!RoleM.Instance.IsJoinedGuild())
        {
            UITool.ShowTips("请加入行会后再进行发言!");
            return;
        }

        if (!GameUtil.Instance.ComCoinToCheck(CFG_XMLConfig.xmlConfig.BaseValue.chatConsume.itemid, CFG_XMLConfig.xmlConfig.BaseValue.chatConsume.num, null))
        {
            UITool.ShowTips("金币不足!");
            return;
        }

        string chatStr = $"{GameRoomM.Instance.CFGRoom.Name}[{MainRole._MainRole.pointNow.X},{MainRole._MainRole.pointNow.Y}]<color=#00000000>{GameRoomM.Instance.CFGRoom.MapId}|{MainRole._MainRole.pointNow.X}|{MainRole._MainRole.pointNow.Y}</color>";
        HttpC.SendChatMessage(ChatType.GUILD, chatStr);
    }

    public override void OnLayerChange()
    {
        if (TopUI == null)
        {
            UITool.LoadUII<UICom_TopRes>(GetView().root_Base, (ui) =>
            {
                TopUI = ui as UICom_TopRes;
                TopUI.OnSetLayer(_canvasOrder);
            });
        }
    }

    private UICom_DownBtns DownUI = null;
    private UICom_TopRes TopUI = null;
    private UICom_GVGBattle _uICom_GVGBattle;
    private UIGameDayDoPro _uIGameDayDoPro;
    private UIGameRelicPro _uIUIGameRelicPro = null;
    private UIGameGuildPro _uIGameGuildPro = null;

    //自动战斗VIP状态
    private void OnAutoLockSet()
    {
        //GetView().obj_AutoPlayLock.SetActive(false);//RoleM.Instance.Vip < CFG_XMLConfig.xmlConfig.BaseValue.autoPlayVipLv);
        if (!GameRoomM.Instance.AutoPlayLock)
        {
            TcpC.C2S_AutoPlay();
        }
    }

    public void OnSkillPanelOpen(bool isOpne)
    {
        GetView().obj_SkillPanelOn.SetActive(isOpne);
        GetView().obj_SkillPaneloOff.SetActive(!isOpne);
        GetView().touch_Skill.gameObject.SetActive(isOpne);
    }

    private void OnAutoPlayClick(bool isPaly)
    {
        GameRoomM.Instance.AutoPlay = isPaly;
        GetView().obj_AutoPlayOn.SetActive(GameRoomM.Instance.AutoPlay);
        GetView().obj_AutoPlayAni.SetActive(GameRoomM.Instance.AutoPlay);
        GetView().obj_AutoPlayIcon.SetActive(!GameRoomM.Instance.AutoPlay);
        //GetView().txt_AutoPlay.text = isPaly ? "停止\n自动" : "自动\n战斗";
    }


    //技能面板滑动
    Vector2 TouchPos;
    private void OnDragBegin(UnityEngine.EventSystems.PointerEventData pEvent)
    {
        TouchPos = pEvent.position;
    }

    private void OnDragEnd(UnityEngine.EventSystems.PointerEventData pEvent)
    {
        if (Mathf.Abs(TouchPos.y - pEvent.position.y) > 3)
        {
            if (TouchPos.y - pEvent.position.y < 0)
                BaseRotate -= 90;
            else
                BaseRotate += 90;

            BaseRotate %= 360;
            SkillRotal(BaseRotate);
        }
    }

    private int BaseRotate = 0;
    private void SkillRotal(int _rot)
    {
        GetView().root_Skill.transform.DOKill();
        GetView().root_Skill.transform.DORotate(Vector3.forward * _rot, 0.3f).SetEase(Ease.Linear).OnComplete(() =>
        {
            for (int i = 1; i <= 12; i++)
            {
                GetView().SkillBtns[i].rect.localEulerAngles = Vector3.forward * (-_rot);
            }
        });
        int idx = Mathf.Abs((_rot > 0 ? _rot - 360 : _rot) / 90) + 1;
        GetView().txt_TouchNum.text = idx.ToString();
    }

    private void OnTimeShow()
    {
        ClearTime();
        bool isPlayedCDWarning = false;
        ShowTime = GameHelper.Instance.AddInvoke(new TimeItem(GameRoomM.Instance.GameFightTime, 1, () =>
        {
            if (isNullGo)
                return;

            OnSkillCDRefresh();

            var cdTime = GameRoomM.Instance.GameFightTime - GameUtil.Instance.GetCurrentTimeSeconds();

            if (!isPlayedCDWarning && cdTime <= 5)
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_BattleCDWarning);
            }

            GetView().txt_time.text = $"倒计时:<color=#f2cd6a>{GameUtil.Instance.GetTimeSToHMS(GameRoomM.Instance.GameFightTime - GameUtil.Instance.GetCurrentTimeSeconds())}</color>";

            if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace)
            {
                long _time = GameRoomM.Instance.MosterBossRefreshTime - GameUtil.Instance.GetCurrentTimeSeconds();
                GetView().txt_MonsterTip.text = _time > 0 ? $"距离BOSS刷新:<color=#ff6868>{GameUtil.Instance.GetTimeSToHMS(_time)}</color>" : "BOSS已刷新";
            }
            if (_uIGameDayDoPro != null)
                _uIGameDayDoPro.OnTimeShow();

            if (_uIUIGameRelicPro != null)
                _uIUIGameRelicPro.OnRelicInfoShow();

            if (GameRoomM.Instance.GameFightTime - GameUtil.Instance.GetCurrentTimeSeconds() <= 0)
            {
                ClearTime();
            }
        }, true));

    }

    private TimeItem ShowTime;

    private void ClearTime()
    {
        if (ShowTime != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(ShowTime);
            ShowTime = null;
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_GetNewSkill(GameEvent.Sync_GetNewSkill _e)
    {
        if (isNullGo)
            return;
        SkillM.Instance.OnInitGameSkill(false);
        OnSkillShow();
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_SkillUpgrade(GameEvent.Sync_SkillUpgrade _e)
    {
        if (isNullGo)
            return;
        SkillM.Instance.OnInitGameSkill(false);
        OnSkillShow();
    }

    private void OnSkillShow()
    {
        //特技暂时隐藏
        for (int i = 0; i < 13; i++)
        {
            int cfgId = SkillM.Instance.GetSkillEquipBySolt(i);

            var skillObj = GetView().SkillBtns[i];
            skillObj.img_icon.enabled = cfgId != 0 && !GameRoomM.Instance.GetCanNoAttackGame();
            if (cfgId != 0 && !GameRoomM.Instance.GetCanNoAttackGame())
            {
                var cfgInfo = CFG_Skill.GetById(cfgId);
                GameUtil.Instance.SetItemIcon((int)cfgInfo.SkillIcon, skillObj.img_icon);
                skillObj.obj_on.SetActive(SkillM.Instance.GetGameSkillOpen(cfgId, i));
                skillObj.img_cd.gameObject.SetActive(false);
                skillObj.obj_lock.SetActive(false);

                var skillInfo = SkillM.Instance.GetSkillInfoByCFGID(cfgId);
                if (skillInfo == default)
                    skillInfo = PartnerM.Instance.GetSkillInfoBySlotSkillID(cfgId);

                skillObj.txt_Reinforce.text = cfgInfo.SkillLevel == 0 ? "" : cfgInfo.SkillLevel.ToString();
                skillObj.txt_Star.text = skillInfo.Advanced == 0 ? "" : skillInfo.Advanced.ToString();
            }
            else
            {
                int openLevel = CFG_XMLConfig.xmlConfig.SkillSlotDic[i];
                bool isLock = i == 0 || (openLevel > 0 && RoleM.Instance.Level < openLevel);
                skillObj.obj_on.SetActive(false);
                skillObj.img_cd.gameObject.SetActive(true);
                skillObj.obj_Dead.SetActive(false);
                skillObj.obj_lock.SetActive(isLock);
                skillObj.txt_Time.text = isLock ? (i == 0 ? "特技暂未开放" : $"{openLevel}级\n解锁") : "";
            }
        }
    }

    private TimeItem _TouchTime;
    private void ClearTouch()
    {
        if (_TouchTime != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(_TouchTime);
            _TouchTime = null;
        }
    }

    //技能按钮点击
    private void OnSkillClick(int skillIdx)
    {
        int cfgId = SkillM.Instance.GetSkillEquipBySolt(skillIdx);
        if (cfgId != 0 && !GameRoomM.Instance.GetCanNoAttackGame())
        {
            SkillM.Instance.SetGameSkillOpen(cfgId, !SkillM.Instance.GetGameSkillOpen(cfgId, skillIdx), skillIdx);
            GetView().SkillBtns[skillIdx].obj_on.SetActive(SkillM.Instance.GetGameSkillOpen(cfgId, skillIdx));
        }
    }

    public void SetJoystickInit()
    {
        //设置摇杆位移和绑定跟随
        GetView()._ETCJoystick.onMove.AddListener(OnJoyMove);
        GetView()._ETCJoystick.onMoveStart.AddListener(OnJoyMoveStart);
        GetView()._ETCJoystick.onMoveEnd.AddListener(OnJoyMoveEnd);
    }

    private void OnJoyMoveStart()
    {
        if (MainRole._MainRole.targetID != 0)
        {
            if (RoleM.Instance.CheakBuffId((int)CFG_BuffEffect.BufferID.XuanYun, (int)CFG_BuffEffect.BufferID.BingDong, (int)CFG_BuffEffect.BufferID.CaoFeng))
                return;
            MainRole._MainRole.targetID = 0;
        }

        if (MainRole._MainRole.targetCollectID != 0)
            MainRole._MainRole.targetCollectID = 0;

        if (MainRole._MainRole._SendMovePos != default)
            MainRole._MainRole._SendMovePos = default;
    }

    private void OnJoyMove(Vector2 moveVec)
    {
        if (RoleM.Instance.CheakBuffId((int)CFG_BuffEffect.BufferID.XuanYun, (int)CFG_BuffEffect.BufferID.BingDong, (int)CFG_BuffEffect.BufferID.CaoFeng))
            return;
        Direction _dir = GetDirection(moveVec);
        MainRole._MainRole._joyDir = _dir;
    }

    private void OnJoyMoveEnd()
    {
        MainRole._MainRole._joyDir = Direction.NONE;
    }

    private Direction GetDirection(Vector2 position)
    {
        //Debug.Log("GetDirection position:" + position);
        if (position.x > 0.1)
        {
            if (position.y > 0.1) return Direction.UP_RIGHT;
            if (position.y < -0.1) return Direction.DOWN_RIGHT;
            return Direction.RIGHT;
        }
        else if (position.x < -0.1)
        {
            if (position.y > 0.1) return Direction.UP_LEFT;
            if (position.y < -0.1) return Direction.DOWN_LEFT;
            return Direction.LEFT;
        }
        return position.y > 0 ? Direction.UP : Direction.DOWN;
    }

    //退出
    private void exitDungeons()
    {
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace)
        {
            TcpC.C2S_Exit();
            return;
        }
        string showStr = "提前退出副本会直接判定挑战失败，是否确认退出？";
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_Peace)
            showStr = "离开副本之后当前获得的战利品将会直接发放到包裹中，副本剩余时间会保留，是否确定离开?";
        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Guild || GameRoomM.Instance.IsDayDo())
            showStr = "是否确定离开?";

        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = showStr,
            YesFun = () =>
            {
                if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_Peace)
                    SaveCDTime();
                TcpC.C2S_Exit();
                ShuShuSend(GameRoomM.Instance.CFGRoom.Roomtype);
            },
            BtnYes = "直接退出",
            BtnNo = "继续游戏",
        });
    }

    void SaveCDTime()
    {
        string cdTimeKey = GameConst.PlayerPrefsKey.HangUpCDTime.ToString() + "1_" + GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear;
        if (GameRoomM.Instance.CFGRoom.ID > 5)
            cdTimeKey = GameConst.PlayerPrefsKey.HangUpCDTime.ToString() + "2_" + GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear;

        GameUtil.Instance.SetPlayerPrefs(cdTimeKey, GameUtil.Instance.GetCurrentTimeSeconds());
    }

    public override void OnEnter(params object[] args)
    {
        OnUniEventAddOrMove(true);
        //OnRoleRefresh();
        OnTimeShow();

        BaseRotate = 0;
        SkillRotal(BaseRotate);
        OnSkillPanelOpen(true);
        GetView().txt_MonsterTip.text = "";
        OnMonsterShow();
        OnChatUI();

        NewPlayerGuideWidget.Instance.AddGuideGameObject("UIGameDungeons_AutoPlay", GetView().btn_AutoPlay.transform);
        //NewPlayerGuideWidget.Instance.AddGuideGameObject("DownBtns_FightButton", GetView().btn_AutoPlay.transform);
    }

    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(OnSync_UserBase);

        OnUniEventAddOrMove(false);
        ClearTouch();

        if (DownUI != null)
            DownUI.OnDestroy();

        if (TopUI != null)
            TopUI.OnDestroy();

        if (_uICom_GVGBattle != null)
            _uICom_GVGBattle.OnDestroy();

        if (_uIGameDayDoPro != null)
            _uIGameDayDoPro.OnDestroy();

        if (_uIGameGuildPro != null)
            _uIGameGuildPro.OnDestroy();
    }
    private void OnUniEventAddOrMove(bool isAdd)
    {
        if (isAdd)
        {
            UniGameEvent.AddEventListener<GameEvent.Fight_Die>(OnFight_Die);
            UniGameEvent.AddEventListener<GameEvent.Fight_MonsterAdd>(OnFight_MonsterAdd);
        }
        else
        {
            UniGameEvent.RemoveEventListener<GameEvent.Fight_Die>(OnFight_Die);
            UniGameEvent.RemoveEventListener<GameEvent.Fight_MonsterAdd>(OnFight_MonsterAdd);
        }
    }

    private void OnFight_Die(GameEvent.Fight_Die _e)
    {
        if (isNullGo)
            return;
        OnMonsterShow();
    }

    private void OnFight_MonsterAdd(GameEvent.Fight_MonsterAdd _e)
    {
        if (isNullGo)
            return;
        OnMonsterShow();
    }

    private void OnMonsterShow()
    {
        if (GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.Tower)
            return;
        GetView().txt_MonsterTip.text = "剩余怪物数量:" + HotUpdateScripts.Game.MonsterManager.Instance.MonesterCount;
    }

    private void OnChatUI()
    {
        if (!GetView().btn_Chat.gameObject.activeSelf)
            return;

        //GetView().root_ChatNew.DOKill();
        //GetView().root_ChatNew.DOScale(0, 0.2f).SetDelay(2f);

        var chatInfos = ChatM.Instance.GetFightChats();

        for (int i = 0; i < chatInfos.Count; i++)
        {
            var root = GetView().root_Chat.GetChild(i);
            var chatInfo = chatInfos[i];
            root.gameObject.SetActive(!string.IsNullOrEmpty(chatInfo.Context));
            if (string.IsNullOrEmpty(chatInfo.Context))
                return;

            string content = "<color=#ffffff00>占位符</color>" + ChatM.Instance.GetChatByType(chatInfo.ChatType, chatInfo.Context, chatInfo.Name);
            JUIExtensions.txt(root).text = content;
            bool isEmjoy = ChatM.Instance.CheckIsEmoji(chatInfo.Context);
            root.Find("emjoy").gameObject.SetActive(isEmjoy);
            if (isEmjoy)
            {
                content += "\n";
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ChatAtlas, chatInfo.Context, JUIExtensions.img(root.Find("emjoy")));
            }
            (string atlas, string icon) = ChatM.Instance.GetChatTypeImg(chatInfo.ChatType);
            ResouseHelper.Instance.LoadSpriteByAtlas(atlas, icon, JUIExtensions.img(root.Find("img_Tip")));
        }

    }

    [Subscriber(ThreadMode.Main)]
    private void OnRefreshChatUI(GameEvent.Sync_ChatPush _e)
    {
        if (isNullGo)
            return;
        OnChatUI();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;
        OnAutoLockSet();
    }

    private void OnSkillCDRefresh()
    {
        for (int i = 0; i < 13; i++)
        {
            int cfgId = SkillM.Instance.GetSkillEquipBySolt(i);
            var skillObj = GetView().SkillBtns[i];
            if (cfgId != 0)
            {
                int _time = (int)(SkillM.Instance.GetGameSkillDealy(cfgId) - GameUtil.Instance.GetCurrentTimeSeconds());
                skillObj.txt_Time.text = "" + _time;
                skillObj.img_cd.gameObject.SetActive(_time > 0);

                long uid = SkillM.Instance.GetUiDBySkillId(cfgId);
                skillObj.obj_Dead.SetActive(uid != 0 && GameRoomM.Instance.GetPlayerHP(uid) == 0 && GameRoomM.Instance._GameState == GameRoomM.GameState.Play);

            }
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_UserBase(GameEvent.Refresh_ConsecutiveKills _e)
    {
        if (isNullGo)
            return;
        string iconStr = _e.count == 10 ? "word_Battle_lz_1" : _e.count == 50 ? "word_Battle_lz_2" : "word_Battle_lz_3";

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, iconStr, GetView().img_lz);
        GetView().img_lz.gameObject.SetActive(_e.count > 0);
        GetView().img_lz.transform.DOKill();
        GetView().img_lz.transform.DOScale(0, 0.2f).SetDelay(2f).OnComplete(() =>
        {
            GetView().img_lz.gameObject.SetActive(false);
            GetView().img_lz.transform.localScale = Vector3.one;
        });
    }

    private void ShuShuSend(CFG_Room.RoomType roomType)
    {
        switch (roomType)
        {
            case CFG_Room.RoomType.PVE_Peace:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = "挂机升级战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.PVE_NoPeace:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "群雄战场",
                    position = "群雄战场战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.PVP_One:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "即时挑战模式",
                    position = "即时挑战",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.PVP_More:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "回合挑战模式",
                    position = "回合挑战",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.Tower:
                break;
            case CFG_Room.RoomType.WorldBoss:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "世界boss",
                    position = "世界boss战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.EXP:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = "经验副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.Gold:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = "金币副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.QHS:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = "强化石副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.SY:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = "书页副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.SXS:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = "升星石副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.HS:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "副本",
                    position = "魂石副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.Guide:
                break;
            case CFG_Room.RoomType.Fight:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = "战役战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.Boss:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = "首领战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.SLTower:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "试炼塔",
                    position = "试炼塔战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.Relic:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "挑战",
                    position = "遗迹战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.Guild:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "社交",
                    component = "行会",
                    position = "行会副本战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.YZDD:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "玩法",
                    component = "竞技",
                    position = "一战到底战斗界面",
                    event_click = $"退出战斗",
                });
                break;
            case CFG_Room.RoomType.WuXing:
                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "五行",
                    component = "五行战斗界面",
                    position = "",
                    event_click = $"退出战斗",
                });
                break;
            default:
                break;
        }
    }

}

