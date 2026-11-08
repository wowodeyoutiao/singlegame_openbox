using Cysharp.Threading.Tasks;
using DG.Tweening;
using HallBattle;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static GameEvent;

public enum HallBattleState
{
    None = 0,
    Rest, // 休息
    Battle, // 战斗
    SettlementFinish, // boos结算
    InSettlement, // 正在结算
    Rest2Battle, // 去战斗
    Battle2Rest, // 去休息
}

public enum HallBattleSpecialSkillType
{
    None = 0,
    Waiting, // 等待
    ShowTip, // 显示提示
    Show, // 显示特技
}

public partial class UIHall
{
    public float timeInterval = 0.03f;

    private HallBattleState hallBattleState = HallBattleState.None;
    public float timeScale = 1f;

    // 状态机
    private HallBattleStateMachine<UIHall> hallBattleFsm;

    // 关卡数据
    private CFGItem_MainBattle[] mainBattleCfgs = new CFGItem_MainBattle[5];
    public CFGItem_MainBattle curMainBattleCfg;
    private bool passAllLevel;

    // 扫荡
    public bool isMopUp;

    private void BattleClear()
    {
        ClearBattleTimer();
    }

    private void OnInitHallBattle()
    {
        hallBattleState = HallBattleState.None;
        passAllLevel = false;

        RefreshConfigData();
        RefreshCurConfigData();
        InitLevel();
        InitSpecialSkill();
        LoadButton();
        AddBattleTimer();
        InitBossHp();
        InitSettlement();

        hallBattleFsm = new HallBattleStateMachine<UIHall>(this);

        UniGameEvent.AddEventListener<MainBattleResponse>(Sever_MainBattleResponse);
        UniGameEvent.AddEventListener<GameEvent.HallBattleBossBeHit>(HallBattleBossBeHitEnt);
        UniGameEvent.AddEventListener<GameEvent.HallBattleDropItem>(HallBattleDropItemEnt);
        UniGameEvent.AddEventListener<GameEvent.HallBattleRoleAttack>(HallBattleRoleAttackEnt);
        UniGameEvent.AddEventListener<GameEvent.HallBallShowSpecialSkillName>(HallBallShowSpecialSkillNameEnt);
        UniGameEvent.AddEventListener<GameEvent.HallBallPauseOrPlay>(PauseSpecialSkill);
    }

    // 刷新关卡的配置数据（5关）
    private void RefreshConfigData()
    {
        if (mainBattleCfgs == null || mainBattleCfgs[0] == null || mainBattleCfgs[4].ID <= RoleM.Instance.MainBattleId)
            mainBattleCfgs = CFG_MainBattle.GetCfgsById(RoleM.Instance.MainBattleId);

        passAllLevel = CFG_MainBattle.IsPassAll(RoleM.Instance.MainBattleId);
    }

    // 刷新关卡的配置数据（1关）
    private void RefreshCurConfigData(int cfgId = -1, int index = -1)
    {
        if (index > 0)
            curMainBattleCfg = mainBattleCfgs[index];
        else
        {
            if (cfgId < 0)
            {
                if (passAllLevel)
                    cfgId = RoleM.Instance.MainBattleId - 1;
                else
                    cfgId = RoleM.Instance.MainBattleId + 1;
            }

            curMainBattleCfg = CFG_MainBattle.GetById(cfgId);
        }
    }

    // 体力判断
    private bool CheckEnergy()
    {
        if (curMainBattleCfg == null || RoleM.Instance.Energy < curMainBattleCfg.cost)
        {
            Debug.Log($"攻打本关卡，需要体力： {curMainBattleCfg.cost}");
            UITool.ShowTips($"体力不足");

            if (hallBattleState != HallBattleState.Rest)
                ChangeState(HallBattleState.Rest);
            return false;
        }
        return true;
    }


    // 改变状态
    public bool ChangeState(HallBattleState state)
    {
        if (hallBattleState == state)
            return false;

        if (state == HallBattleState.Battle && !CheckEnergy())
            return false;

        hallBattleState = state;
        if (state == HallBattleState.Battle)
            hallBattleFsm.ChangeState<HallBattleAttackState>();
        else
        {
            SetPetTalkActive(true);
            hallBattleFsm.ChangeState<HallBattleRestState>();
        }

        // 休息图标
        uIViewHall.levelRestItem.Find("FireImg").gameObject.SetActive(state == HallBattleState.Rest);
        uIViewHall.levelRestItem.Find("ZPImg").gameObject.SetActive(state == HallBattleState.Battle);

        if (state == HallBattleState.Battle) //在大厅要检查是否有新功能开启
        {
            //检查是否有新功能开启需要表现
            FuncOpenShowWidget.Instance.TryAddShowData();
            FuncOpenShowWidget.Instance.TryShow();
            uIViewHall.boxShowUIMaskGo.SetActive(false);
        }else
        {
            uIViewHall.boxShowUIMaskGo.transform.SetAsLastSibling();
            uIViewHall.boxShowUIMaskGo.SetActive(RoleM.Instance.OpenBoxItems != null && RoleM.Instance.OpenBoxItems.Count > 0);
        }

        CheckBox();
        RefreshCurLevelUI();

        return true;
    }

    public bool IsState(HallBattleState state)
    {
        return hallBattleState == state;
    }

    #region 转场特效
    public void ShowZhuanLoading(Action callBack = null)
    {
        uIViewHall.ZhuanLoadingTrans.SetAsLastSibling();
        uIViewHall.ZhuanLoadingTrans.gameObject.SetActive(true);
        SpineUtil.Instance.PlayOnce(uIViewHall.ZhuanLoadingTrans, "animation", () =>
        {
            callBack?.Invoke();
            if (!isNullGo && uIViewHall.ZhuanLoadingTrans != null)
                uIViewHall.ZhuanLoadingTrans.gameObject.SetActive(false);
        }, timeScale);
    }
    #endregion

    #region button
    private int comboTimes;
    private float comboTimeInterval = 0;
    private Sequence comboSequence;
    private long oldComboAddTime = 0;

    private void LoadButton()
    {
        ButtonAddListener(uIViewHall.backHomeBtn, BattleBackToRest); // 战场->休息室
        ButtonAddListener(uIViewHall.boxBtn, () =>
        {
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.OpenBox))
                return;

            UITool.Open<UIBoxUpgrade>();
        }); // 宝箱升级
        ButtonAddListener(uIViewHall.attackBtn, OnClickAttackButton); // 出战
    }

    // 点击攻击
    public void OnClickAttackButton()
    {
        if (hallBattleState == HallBattleState.Rest)
        {
            isMopUp = false;
            // 通关情况下，刷新前一关
            if (passAllLevel)
            {
                curMainBattleCfg = mainBattleCfgs[mainBattleCfgs.Length - 2];
                if (!CheckEnergy()) return;
                ChangeState(HallBattleState.Battle);
            }
            else
            {
                // 攻击下一关
                for (int i = 0; i < uIViewHall.levelItems.Length; i++)
                {
                    var cfg = mainBattleCfgs[i];
                    if (cfg.ID == RoleM.Instance.MainBattleId + 1)
                    {
                        if (cfg.flag == (int)MainBattleFlayType.Boss)
                            curMainBattleCfg = CFG_MainBattle.GetById(RoleM.Instance.MainBattleId);
                        else
                            curMainBattleCfg = cfg;

                        if (!CheckEnergy()) return;
                        ChangeState(HallBattleState.Battle);
                        break;
                    }
                }
            }
        }
    }
    #endregion

    #region 特技技能
    private Image[] norSpecialImgs = new Image[5];
    private Image[] selSpecialImgs = new Image[5];
    private Text specialTimeText;
    private Text specialTipText;
    private Transform specialWaitEffect;
    private GameObject perfectGo;
    private Animator specialWaitAnimator;

    private HallBattleSpecialSkillType specialSkillType;

    // 配置
    private int superskillCD; // 秒
    private int missskillCD; // 秒
    private int superskilltime; // 秒
    private int superskillatktimes; // 当角色攻击次数为6次后，开启特技
    private int specialSkillCD; // 这个为superskillCD 或者 missskillCD

    private string[] specialSkillNorIconNames = new string[] { "btn_Main_fkan_2", "btn_Main_sjan_2", "btn_Main_yqan_2", "btn_Main_dcan_2" };
    private string[] specialSkillSelIconNames = new string[] { "btn_Main_fkan_1", "btn_Main_sjan_1", "btn_Main_yqan_1", "btn_Main_dcan_1" };

    private float specialSkillTimeInterval;

    private List<int> specialSkillCombos = new List<int>();
    private CFGItem_SkillMainBattle cFGItem_SkillMainBattle;
    private int specialSkillIndex;
    private bool isNorSpecialSkill = true;

    private bool pauseSpecialSkill = false;
    // 打点
    private bool isFirstSpecialSkill;

    // 自动战斗
    private bool autoHallBattle;


    private void InitSpecialSkill()
    {
        for(int i = 0; i < uIViewHall.skillBtnsTrans.childCount; i++)
        {
            var index = i;
            var btn = uIViewHall.skillBtnsTrans.GetChild(i).GetComponent<Button>();
            ButtonAddListener(btn, () => { OnClickSpecialSkillButton(index); });

            NewPlayerGuideWidget.Instance.AddGuideGameObject($"Hall_SpecialSkill_{i+1}", btn.transform);

            NewPlayerGuideWidget.Instance.AddGuideGameObject($"Hall_SpecialSkill_Pos{i + 1}", btn.transform.GetChild(0).transform);
        }

        for(int i = 0; i < 5; i++)
        {
            var trans = uIViewHall.skillShowTrans.Find("Icons/Image_" + (i + 1));
            norSpecialImgs[i] = trans.GetComponent<Image>();
            selSpecialImgs[i] = trans.GetChild(0).GetComponent<Image>();
        }

        specialTimeText = uIViewHall.skillShowTrans.Find("TimeText").GetComponent<Text>();
        specialTipText = uIViewHall.skillShowTrans.Find("TipText").GetComponent<Text>();
        specialWaitEffect = uIViewHall.skillShowTrans.Find("SkillWaitEffect");
        perfectGo = uIViewHall.skillShowTrans.Find("SkillPerfect").gameObject;
        specialWaitAnimator = specialWaitEffect.GetComponent<Animator>();

        var autoGo1 = uIViewHall.autoSpecialSkillGo.transform.Find("Image_1").gameObject;
        var autoGo2 = uIViewHall.autoSpecialSkillGo.transform.Find("Auto").gameObject;
        ButtonAddListener(uIViewHall.autoSpecialSkillBtn, () =>
        {
            autoHallBattle = !autoHallBattle;
            if (autoHallBattle)
            {
                if (specialSkillType == HallBattleSpecialSkillType.Show && cFGItem_SkillMainBattle != null)
                {
                    HideSpecialSkill();
                    specialTipText.text = "自动释放特技!";
                    // 胜利，播放技能
                    hallRole.DoSpecialSkill(cFGItem_SkillMainBattle.ID);
                }
            }
            else
            {
                if (specialSkillType == HallBattleSpecialSkillType.Show)
                    ShowSpecialSkill();
            }

            autoGo1.SetActive(!autoHallBattle);
            autoGo2.SetActive(autoHallBattle);
        });

        specialSkillType = HallBattleSpecialSkillType.None;
        SetSpecialSkillCD();

        autoHallBattle = false;
        ShowAutoSpecialSkill(false);
    }

    private void PauseSpecialSkill(HallBallPauseOrPlay play)
    {
        pauseSpecialSkill = play.pause;
    }

    private void SetSpecialSkillCD()
    {
        if(RoleM.Instance.MainBattleId <= 21)
        {
            superskillCD = CFG_XMLConfig.xmlConfig.MainSkill.superskillCD - 2;
            missskillCD = CFG_XMLConfig.xmlConfig.MainSkill.missskillCD - 2;
            superskilltime = CFG_XMLConfig.xmlConfig.MainSkill.superskilltime;
            superskillatktimes = CFG_XMLConfig.xmlConfig.MainSkill.superskillatktimes - 1;
        }else
        {
            superskillCD = CFG_XMLConfig.xmlConfig.MainSkill.superskillCD;
            missskillCD = CFG_XMLConfig.xmlConfig.MainSkill.missskillCD;
            superskilltime = CFG_XMLConfig.xmlConfig.MainSkill.superskilltime;
            superskillatktimes = CFG_XMLConfig.xmlConfig.MainSkill.superskillatktimes;
        }
    }

    private void ShowSpecialSkillUI()
    {
        SetSpecialSkillCD();
        specialSkillIndex = 0;
        HideSpecialSkill();
        specialSkillCombos.Clear();
        if(isNorSpecialSkill)
            cFGItem_SkillMainBattle = CFG_SkillMainBattle.GetOneNorSpecialSkill();
        else
            cFGItem_SkillMainBattle = CFG_SkillMainBattle.GetOneBigSpecialSkill();

        if (cFGItem_SkillMainBattle == null || string.IsNullOrEmpty(cFGItem_SkillMainBattle.Combo))
        {
            specialSkillType = HallBattleSpecialSkillType.ShowTip;
            isNorSpecialSkill = !isNorSpecialSkill;
            return;
        }

        var ComboList = cFGItem_SkillMainBattle.Combo.Split('+');
        UniGameEvent.Send(new GameEvent.Refresh_NewPlayerGuide() { args = ComboList });

        specialSkillType = HallBattleSpecialSkillType.Show;
        isNorSpecialSkill = !isNorSpecialSkill;

        if (autoHallBattle)
        {
            specialTipText.text = "自动释放特技!";
            // 胜利，播放技能
            hallRole.DoSpecialSkill(cFGItem_SkillMainBattle.ID);
            return;
        }

        for (int i = 0; i < 5; i++)
        {
            norSpecialImgs[i].gameObject.SetActive(i < ComboList.Length);
            selSpecialImgs[i].gameObject.SetActive(false);
            if (i < ComboList.Length)
            {
                var temp = int.Parse(ComboList[i].ToString());
                specialSkillCombos.Add(temp);

                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.HallAtlas, specialSkillNorIconNames[temp - 1], norSpecialImgs[i]);
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.HallAtlas, specialSkillSelIconNames[temp - 1], selSpecialImgs[i]);
            }
        }

        specialTimeText.text = Mathf.CeilToInt(superskilltime) + "s";
    }

    public void HideSpecialSkill()
    {
        for (int i = 0; i < norSpecialImgs.Length; i++)
            norSpecialImgs[i].gameObject.SetActive(false);

        specialTimeText.text = string.Empty;
        specialTipText.text = string.Empty;
        specialSkillType = HallBattleSpecialSkillType.None;
        specialWaitEffect.gameObject.SetActive(false);
        perfectGo.SetActive(false);
    }

    private void HallBattleRoleAttackEnt(HallBattleRoleAttack args)
    {
        if(args.attackTimes >= superskillatktimes && specialSkillType == HallBattleSpecialSkillType.None)
        {
            isFirstSpecialSkill = true;
            ShowSpecialSkill();
        }
    }

    private void OnClickSpecialSkillButton(int index)
    {
        if (autoHallBattle)
            return;

        if (hallBattleState != HallBattleState.Battle)
            return;

        if (specialSkillCombos == null || specialSkillCombos.Count == 0)
            return;

        if (specialSkillType != HallBattleSpecialSkillType.Show)
            return;

        if (specialSkillIndex == 0)
            TrySendShuShu();

        if (specialSkillCombos[0] != index + 1)
        {
            ShowSpecialSkillTip(false);

            SendToShuShu(false, true, false);

            TyrSendFirstErrorToShuShu(index + 1);
        }
        else
        {
            selSpecialImgs[specialSkillIndex].gameObject.SetActive(true);

            specialSkillCombos.RemoveAt(0);
            if (specialSkillCombos.Count == 0)
            {
                SendToShuShu(true, false, false);

                pauseSpecialSkill = false;

                ShowSpecialSkillTip(true);
                // 按完按钮，重新开始游戏
                UniGameEvent.Send(new GameEvent.HallBallPauseOrPlay() { pause = false });
                // 胜利，播放技能
                hallRole.DoSpecialSkill(cFGItem_SkillMainBattle.ID);
            }
            specialSkillIndex++;
        }
    }

    private void ShowSpecialSkill()
    {
        specialSkillTimeInterval = 0;
        specialWaitEffect.gameObject.SetActive(false);
        perfectGo.SetActive(false);
        // 显示特技
        ShowSpecialSkillUI();

        trySendToShuShuWhenShowSpecialSKill();
    }

    private void ShowSpecialSkillTip(bool isSuc)
    {
        HideSpecialSkill();

        specialSkillTimeInterval = 0;
        specialSkillCD = isSuc ? superskillCD : missskillCD;

        perfectGo.SetActive(isSuc);
        specialWaitEffect.gameObject.SetActive(isSuc);

        if(isSuc)
        {
            var index = UnityEngine.Random.Range(1, 7);
            specialWaitAnimator.Play("HallSpecialSkill" + index, 0, 0);
            specialTimeText.text = $"{Mathf.CeilToInt(specialSkillCD)}s"; // 后出现特技组合，请准备!
        }else
        {
            specialTipText.text = $"{Mathf.CeilToInt(specialSkillCD)}s后出现特技组合，请准备!";
        }

        specialSkillType = HallBattleSpecialSkillType.ShowTip;
    }

    private void SpecialSkillTick(float dTime)
    {
        if (pauseSpecialSkill)
            return;

        if (specialSkillType == HallBattleSpecialSkillType.None)
            return;

        if (hallRole == null || !hallRole.IsState(HallRoleState.Attack))
            return;

        // 显示
        if(specialSkillType == HallBattleSpecialSkillType.Show)
        {
            var timeCount = superskilltime;
            if (autoHallBattle)
                timeCount = 1;

            if (specialSkillTimeInterval >= timeCount)
            {
                ShowSpecialSkillTip(autoHallBattle);
                SendToShuShu(false, false, true);
            }
            else
            {
                specialSkillTimeInterval += dTime;

                if(autoHallBattle)
                    specialTimeText.text = "";
                else
                {
                    var temp = superskilltime - specialSkillTimeInterval;
                    specialTimeText.text = Mathf.CeilToInt(temp) + "s";
                }  
            }
        }

        // 显示tip
        if (specialSkillType == HallBattleSpecialSkillType.ShowTip)
        {
            if (specialSkillTimeInterval >= specialSkillCD)
            {
                specialTimeText.text = "";
                ShowSpecialSkill();
            }
            else
            {
                specialSkillTimeInterval += dTime;

                var temp = specialSkillCD - specialSkillTimeInterval;
                if (specialSkillCD == superskillCD)
                    specialTimeText.text = $"{Mathf.CeilToInt(temp)}s"; // 后出现特技组合，请准备!
                else
                    specialTipText.text = $"{Mathf.CeilToInt(temp)}s后出现特技组合，请准备!";
            }
        }
    }

    private void HallBallShowSpecialSkillNameEnt(HallBallShowSpecialSkillName args)
    {
        HallBallShowSpecialSkillNameEnt(args.id).Forget();
    }

    private async UniTask HallBallShowSpecialSkillNameEnt(int id)
    {
        var skillMainBattleCfg = CFG_SkillMainBattle.GetById(id);
        if (skillMainBattleCfg == null || string.IsNullOrEmpty(skillMainBattleCfg.Image))
            return;

        var poolKey = "HallBallShowSpecialSkillName";
        if (!PoolHelper.Instance.CheckPool(poolKey))
        {
            var go = GameObject.Instantiate(uIViewHall.specialSkillNameGo);
            PoolHelper.Instance.AddPool(poolKey, go);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
        var nameGo = PoolHelper.Instance.GetPoolObj(poolKey);
        nameGo.name = poolKey;
        nameGo.transform.SetParent(uIViewHall.transform);
        nameGo.transform.localPosition = new Vector3(-170, 250, 0);
        nameGo.transform.localScale = Vector3.one;
        nameGo.SetActive(false);

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.HallBattleAtlas, skillMainBattleCfg.Image, JUIExtensions.img(nameGo.transform));

        var sp = nameGo.GetComponent<UIHallSpecialSKillName>();
        if (sp == null)
            sp = nameGo.AddComponent<UIHallSpecialSKillName>();

        sp.Show();
    }

    private void SendToShuShu(bool isSuc, bool clickFail, bool timeOut)
    {
        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties.Add("is_success", $"{isSuc}");
        properties.Add("main_id", curMainBattleCfg.ID);
        properties.Add("is_first_combo", isFirstSpecialSkill);
        properties.Add("click_fail_combo", clickFail);

        var key = "is_first_hand_combo" + RoleM.Instance.Uid;
        var res = GameUtil.Instance.GetPlayerPrefs(key);

        // 新手过完后，第一次手动点击，非时间过了
        if (string.IsNullOrEmpty(res) && !pauseSpecialSkill && !timeOut)
        {
            properties.Add("is_first_hand_combo", true);
            GameUtil.Instance.SetPlayerPrefs(key, "is_first_hand_combo");
            //Debug.LogError("is_first_hand_combo");
        }
        //if(isSuc)
        //    Debug.LogError("combo_Success");
        //else if(clickFail)
        //    Debug.LogError("combo_fail");

        TDAnalyticsUtil.Track("combo", properties, true);

        if (isFirstSpecialSkill)
            isFirstSpecialSkill = false;
    }

    private void TrySendShuShu()
    {
        var key = "is_first_hand_combo" + RoleM.Instance.Uid;
        var res = GameUtil.Instance.GetPlayerPrefs(key);
        if (string.IsNullOrEmpty(res))
            return;
        //Debug.LogError("combo_clicknumberone");
        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties.Add("is_success", true);
        properties.Add("click_combo", 1);

        TDAnalyticsUtil.Track("combo_clicknumberone", properties, true);
    }

    private void trySendToShuShuWhenShowSpecialSKill()
    {
        //Debug.LogError("combo_appear");
        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties.Add("count", 1);

        TDAnalyticsUtil.Track("combo_appear", properties, true);
    }

    private void TyrSendFirstErrorToShuShu(int index)
    {
        var key = "is_first_hand_error_combo" + RoleM.Instance.Uid;
        var res = GameUtil.Instance.GetPlayerPrefs(key);

        // 非首次出错不上报
        if (!string.IsNullOrEmpty(res))
            return;

        GameUtil.Instance.SetPlayerPrefs(key, "is_first_error_combo");
        //Debug.LogError($"首次上报错误位置！{index}");

        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties.Add("failno", index);
        TDAnalyticsUtil.Track("combo_firstfail", properties, true);
    }

    private void CheckAutoSpecialSkill()
    {
        //if (RoleM.Instance.MainBattleId > 0)
        //{
        //    ShowAutoSpecialSkill(true);
        //}
    }

    private void ShowAutoSpecialSkill(bool active)
    {
        uIViewHall.autoSpecialSkillGo.SetActive(active);
    }
    #endregion

    #region 连击
    private void Combo(float dTime)
    {
        //if (isBattleAutoAttack || comboTimes == 0)
        //    return;

        //comboTimeInterval += dTime;
        //if(comboTimeInterval >= 1)
        //    ResetCombo();
    }

    // 连击重置
    public void ResetCombo()
    {
        comboTimes = 0;
        comboTimeInterval = 0;

        hideComboUI();
    }

    // 连击UI
    public void ShowComboUI()
    {
        //if(playTimeScale == 2)
        //{
        //    hideComboUI();
        //    return;
        //}

        var curTime = GameUtil.Instance.GetCurrentTimeSeconds();
        if(oldComboAddTime == 0 || curTime - oldComboAddTime > 0.4f)
        {
            oldComboAddTime = curTime;
            return;
        }

        comboTimes++;
        comboTimeInterval = 0;

        if (comboTimes == 1)
            return;

        //if(comboTimes > 2)
        //    AudioManager.Instance.PlaySound(AudioKey.Sound.Combo);

        uIViewHall.comboEffectGo.SetActive(comboTimes >= 5);

        uIViewHall.comboText.text = comboTimes.ToString();

        ClearComboSequence();

        uIViewHall.comboTrans.localScale = Vector3.one;
        if (!uIViewHall.comboGo.activeSelf)
        {
            uIViewHall.comboGo.SetActive(true);

            uIViewHall.comboTrans.localScale = Vector3.one * 2f;
            comboSequence = DOTween.Sequence();
            comboSequence.Append(uIViewHall.comboTrans.DOScale(0.8f, 0.1f));
            comboSequence.Append(uIViewHall.comboTrans.DOScale(1f, 0.01f));
            comboSequence.Play();
        }
        else
        {
            var scale = 1.5f;
            if(comboTimes >= 4)
                scale = 2.6f;

            comboSequence = DOTween.Sequence();
            comboSequence.Append(uIViewHall.comboText.transform.DOScale(scale, 0.05f));
            comboSequence.Append(uIViewHall.comboText.transform.DOScale(1f, 0.01f));
            comboSequence.Play();
        } 
    }

    public void hideComboUI()
    {
        //uIViewHall.comboText.text = string.Empty;
        //uIViewHall.comboGo.SetActive(false);
        //uIViewHall.comboEffectGo.SetActive(false);

        //ClearComboSequence();
    }

    private void ClearComboSequence()
    {
        if (comboSequence != null)
        {
            comboSequence.Kill();
            comboSequence = null;
        }
    }
    #endregion

    #region timer
    private TimeItem battleTimeItem;
    private void AddBattleTimer()
    {
        ClearBattleTimer();
        battleTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(990000, timeInterval, () =>
        {
            hallBattleFsm?.Update(timeInterval);
            Combo(timeInterval);
            SpecialSkillTick(timeInterval);

            // 挑衅
            ShowChallenge(timeInterval);
        }));
    }

    private void ClearBattleTimer()
    {
        if (battleTimeItem != null && GameHelper.Instance != null && !isNullGo)
        {
            GameHelper.Instance.RemoveInvoke(battleTimeItem);
            battleTimeItem = null;
        }
    }
    #endregion

    #region Role
    private string[] roleNames = new string[5];
    public HallRole hallRole;

    // 角色创建
    private void SetRoleSpine()
    {
        for (int i = 0; i < uIViewHall.roleTranss.Length; i++)
        {
            string poolKey = null, poolPath = null;
            int idx = i;

            if (i == 0)
            {
                var go = uIViewHall.roleTranss[0].gameObject;
                hallRole = go.GetComponent<HallRole>();
                if (hallRole == null)
                    hallRole = go.AddComponent<HallRole>();

                hallRole.InitWithInfo(true, RoleM.Instance.Job, tScale: timeScale);
                hallRole.AddHpChangeTrans(uIViewHall.hpChangeNodeTrans);

                poolKey = "MainRole" + RoleM.Instance.Job + RoleM.Instance.Sex;
                poolPath = GameConst.GetModelByJobAndSex((JobType)RoleM.Instance.Job, (SexType)RoleM.Instance.Sex);
            }
            else
            {
                var soltInfo = PartnerM.Instance.GetOnePartnerSlotInfo(idx);
                if (soltInfo.PartnerUid != 0)
                {
                    var partnerInfo = PartnerM.Instance.GetOnePartnerInfoByUid(soltInfo.PartnerUid);
                    poolKey = "Partner_" + partnerInfo.CfgId;
                    poolPath = GameConst.GetMonsterModelByID(poolKey);
                }
            }

            if (string.IsNullOrEmpty(poolKey))
            {
                ClearRoot();
                roleNames[i] = poolKey;
            }
            else
            {
                if (roleNames[i] != poolKey)
                {
                    GameHelper.Instance.SetAsyncObj(poolKey, poolPath, (go) =>
                    {
                        roleNames[idx] = poolKey;
                        ClearRoot();
                        go.SetActive(true);
                        go.transform.SetParent(uIViewHall.roleTranss[idx]);
                        go.transform.localPosition = Vector3.zero;
                        go.transform.localScale = Vector3.one;
                        go.transform.localScale = new Vector3(idx <= 2 ? -1 : 1, 1, 1);
                        SpineUtil.Instance.PlayLoop(go.transform.Find("spine"), GameConst.ModelAnim.stand.ToString(), timeScale);

                        if (idx == 0)
                            hallRole.SetSpineTransform(go.transform.Find("spine"));
                    }).Forget();
                }
            }

            void ClearRoot()
            {
                if (GameHelper.Instance != null && !isNullGo)
                    GameHelper.Instance.ClearObjToPool(uIViewHall.roleTranss[idx]);
            }
        }
    }

    [Subscriber(ThreadMode.Main)]
    //伙伴槽位信息刷新
    private void OnRefresh_PartnerSlot(GameEvent.Refresh_PartnerSlot _e)
    {
        if (isNullGo) return;

        SetRoleSpine();
        if (curRoleIndex != 0)
            ShowEquipUI(curRoleIndex);
    }
    #endregion

    #region pet
    private string petPool = string.Empty;
    private Coroutine petTalkCoroutine;
    private void SetPetSpine()
    {
        var myPet = PetM.Instance.GetOnePetByIdx(0);
        if (myPet.CfgId == 0)
        {
            SetPetTalkActive(false);
            petPool = string.Empty;
            GameHelper.Instance.ClearObjToPool(uIViewHall.petTrans);
        }
        else
        {
            string pooName = "Pet_" + myPet.CfgId;
            if (pooName == petPool) return;
            _ = GameHelper.Instance.SetAsyncObj(pooName, GameConst.GetMonsterModelByID(pooName), (petGo) =>
            {
                petPool = pooName;
                GameHelper.Instance.ClearObjToPool(uIViewHall.petTrans);
                petGo.transform.SetParent(uIViewHall.petTrans);
                petGo.transform.localScale = Vector3.one;
                petGo.transform.localPosition = new Vector3(0, -50, 0);
                petGo.SetActive(true);

                SetPetTalkActive(hallBattleState == HallBattleState.Rest);
            });
        }
    }

    //宠物刷新
    [Subscriber(ThreadMode.Main)]
    private void Refresh_PetsInfo(GameEvent.Refresh_PetsInfo _e)
    {
        if (isNullGo) return;
        SetPetSpine();
    }

    public void SetPetTalkActive(bool active)
    {
        if (string.IsNullOrEmpty(petPool)) return;

        uIViewHall.petTrans.gameObject.SetActive(active);
        uIViewHall.petTalkTrans.gameObject.SetActive(active);

        ClearPetTalkCoroutine();
        if (active)
        {
            // 使用 Coroutine
            petTalkCoroutine = CoroutineMgr.Instance.StartCoroutine(PetTalkCoroutine());
        }
    }

    private IEnumerator PetTalkCoroutine()
    {
        while (true)
        {
            ShowOnePetTalk();
            yield return new WaitForSeconds(5f);
            HidePetTalk();
            yield return new WaitForSeconds(8f);
        }
    }

    private void ShowOnePetTalk()
    {
        if (string.IsNullOrEmpty(petPool)) return;

        var item = CFG_PetSpeak.GetOne();
        if (item != null)
        {
            uIViewHall.petTalkTrans.gameObject.SetActive(true);
            uIViewHall.petTalkImg.gameObject.SetActive(item.type == 2);
            if (item.type == 1)
            {
                uIViewHall.petTalkText.text = item.text;
                //var width = uIViewHall.petTalkText.preferredWidth / 3f + 50;
                //uIViewHall.petTalkRectTrans.sizeDelta = new Vector2(width, 87);
                var rows = Math.Ceiling(uIViewHall.petTalkText.preferredWidth / 164f);
                var height = (float)rows * 22 + 26;
                if (height < 16) height = 164;
                uIViewHall.petTalkRectTrans.sizeDelta = new Vector2(164, height);
            }
            else
            {
                uIViewHall.petTalkText.text = string.Empty;
                uIViewHall.petTalkRectTrans.sizeDelta = new Vector2(164, 87);
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ChatAtlas, item.text, uIViewHall.petTalkImg);
            }
        }
        else
            HidePetTalk();
    }

    private void HidePetTalk()
    {
        //Debug.LogError("HidePetTalk");
        uIViewHall.petTalkTrans.gameObject.SetActive(false);
    }

    private void ClearPetTalkCoroutine()
    {
        if (petTalkCoroutine != null && CoroutineMgr.Instance != null && !isNullGo)
        {
            CoroutineMgr.Instance.StopCoroutine(petTalkCoroutine);
            petTalkCoroutine = null;
        }
    }
    #endregion

    #region level
    private GameObject passChallengeGo;
    private Text passChallengeText;
    private bool showPassChallenge;
    private float passChallengeInterval;
    private bool showNexPassChallengeText;

    private GameObject challengeGo;
    private bool showChallenge;
    private float challengeInterval;

    private void InitLevel()
    {
        // 休息
        ButtonAddListener(uIViewHall.levelRestItem.GetComponent<Button>(), BattleBackToRest);

        for (int i = 0; i < uIViewHall.levelItems.Length; i++)
        {
            var index = i;
            var levelItem = uIViewHall.levelItems[i];
            levelItem.btn.onClick.RemoveAllListeners();
            levelItem.btn.onClick.AddListener(() =>
            {
                OnClickLevelBtnAtIndex(index);
            });
        }

        // boss 页面
        uIViewHall.battleBossRewardTrans.gameObject.SetActive(false);
        ButtonAddListener(uIViewHall.battleBossRewardTrans.GetComponent<Button>(), () =>
        {
            MoveToNext();
            ClearBossRewardTimer();
            uIViewHall.battleBossRewardTrans.gameObject.SetActive(false);
        });

        uIViewHall.levelSliderImg.fillAmount = 0;
        RefreshLevelUI(RoleM.Instance.MainBattleId);
    }

    private void BattleBackToRest()
    {
        if (hallBattleState == HallBattleState.Battle)
        {
            hallBattleState = HallBattleState.Battle2Rest;

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                ChangeState(HallBattleState.Rest);
            }, 0.75f);

            uIViewHall.BattleToRestLoadingTrans.SetAsLastSibling();
            uIViewHall.BattleToRestLoadingTrans.gameObject.SetActive(true);
            SpineUtil.Instance.PlayOnce(uIViewHall.BattleToRestLoadingTrans, "animation", () =>
            {
                uIViewHall.BattleToRestLoadingTrans.gameObject.SetActive(false);
            });
        }
    }

    public void SetMainLevelUIActive()
    {
        uIViewHall.levelNodeTrans.gameObject.SetActive(curMainBattleCfg.flag != (int)MainBattleFlayType.Boss || hallBattleState != HallBattleState.Battle);
        if(passChallengeGo != null && hallBattleState != HallBattleState.Rest)
            passChallengeGo.SetActive(false);
    }

    // 直接攻打x关
    private void OnClickLevelBtnAtIndex(int index)
    {
        if (index < 0 || index > mainBattleCfgs.Length) return;
        var cfg = mainBattleCfgs[index];

        var curId = RoleM.Instance.MainBattleId;
        if (cfg.flag == (int)MainBattleFlayType.Boss)
        {
            if (hallBattleState == HallBattleState.Battle)
            {
                if (cfg.ID < curId)
                    UITool.ShowTips($"在休息室才可以进行连战模式。请先回休息室再尝试挑战！");
                else
                    Debug.Log("正在战斗中，请专心战斗！");
                return;
            }

            if (cfg.ID <= curId || cfg.ID == curId + 1)
            {
                isMopUp = cfg.ID <= curId;
                curMainBattleCfg = cfg;
                ChangeState(HallBattleState.Battle);
            }
            else
                UITool.ShowTips($"当前关卡尚未解锁!");
        }
        else
            Debug.Log($"当前第{cfg.ID}关卡不可战斗，目前可征战第{curId}关");
    }

    private void RefreshLevelUI(int curId)
    {
        for (int i = 0; i < uIViewHall.levelItems.Length; i++)
            RefreshOneLevelUI(i, curId);
    }

    private void RefreshCurLevelUI()
    {
        for (int i = 0; i < uIViewHall.levelItems.Length; i++)
        {
            var cfg = mainBattleCfgs[i];
            if (cfg.ID == curMainBattleCfg.ID)
            {
                RefreshOneLevelUI(i, curMainBattleCfg.ID);
                break;
            }
        }
    }

    private void RefreshOneLevelUI(int index, int curId)
    {
        var passId = RoleM.Instance.MainBattleId;
        var levelItem = uIViewHall.levelItems[index];
        var cfg = mainBattleCfgs[index];
        //Debug.LogError($"RefreshOneLevelUI, RoleM.Instance.MainBattleId: {passId}, cfg.ID: {cfg.ID}, cfg.flag:{cfg.flag}, curId:{curId}");

        levelItem.finishGo.SetActive(cfg.ID <= passId);
        // 小怪
        var isFlag = (cfg.flag == (int)MainBattleFlayType.None);
        levelItem.monsterNorIconGo.SetActive(isFlag && (hallBattleState != HallBattleState.Battle || cfg.ID != curId));

        // 精英
        isFlag = (cfg.flag == (int)MainBattleFlayType.Elite);
        levelItem.monsterEliteIconGo.SetActive(isFlag && (hallBattleState != HallBattleState.Battle || cfg.ID != curId));

        // 攻击剑特效
        levelItem.attackGo.SetActive(cfg.ID == curId && hallBattleState == HallBattleState.Battle);

        isFlag = (cfg.flag == (int)MainBattleFlayType.Boss);
        // 已经打过的boss
        levelItem.monsterPassBossIconGo.SetActive(isFlag && cfg.ID <= passId && (!passAllLevel || (passAllLevel && cfg.ID != passId)) && (hallBattleState != HallBattleState.Battle || cfg.ID != curId));

        // 未解锁的
        levelItem.monsterBossIconGo2.SetActive(isFlag && cfg.ID > passId + 1);

        // 火苗
        levelItem.monsterBossIconGo.SetActive(isFlag && (cfg.ID == passId + 1 || (passAllLevel && cfg.ID == passId)) && (hallBattleState != HallBattleState.Battle || cfg.ID != curId));

        // 挑衅
        if (levelItem.challengeGo != null)
        {
            challengeInterval = 0;
            passChallengeInterval = 0;
            showNexPassChallengeText = false;
            if (index > 0)
            {
                challengeGo = levelItem.challengeGo;
                showChallenge = levelItem.monsterBossIconGo.activeSelf;
                challengeGo.SetActive(showChallenge);
            }else
            {
                passChallengeGo = levelItem.challengeGo;
                passChallengeText = passChallengeGo.transform.Find("Text").GetComponent<Text>();
                showPassChallenge = levelItem.monsterPassBossIconGo.activeSelf;
                passChallengeText.text = (showNexPassChallengeText ? "扫荡" : "打门票");
                passChallengeGo.SetActive(showPassChallenge && hallBattleState == HallBattleState.Rest);
            }
        }

        if (cfg.ID == curId)
        {
            uIViewHall.levelTitleText.text = cfg.name;
            //NewPlayerGuideWidget.Instance.TryShowGuide();
            FuncOpenShowWidget.Instance.TryAddShowData();
            FuncOpenShowWidget.Instance.TryShow();
        }

        if (cfg.ID == passId)
            uIViewHall.levelSliderImg.fillAmount = (index + 1) / 4f;
    }

    private void ShowChallenge(float dTime)
    {
        if (showChallenge && challengeGo != null)
        {
            if (challengeInterval > 6)
                challengeInterval = 0;
            else if (challengeInterval > 3 && challengeGo.activeSelf)
            {
                challengeGo.SetActive(false);
                challengeInterval += dTime;
            }
            else if (challengeInterval == 0)
            {
                challengeGo.SetActive(true);
                challengeGo.transform.localScale = Vector3.one * 0.01f;
                challengeGo.transform.DOScale(1, 0.2f);
                challengeInterval += dTime;
            }
            else
                challengeInterval += dTime;
        }

        if (hallBattleState != HallBattleState.Rest)
            return;

        if (showPassChallenge && passChallengeGo != null)
        {
            if (passChallengeInterval > 8)
                passChallengeInterval = 0;
            else if (passChallengeInterval > 3 && passChallengeGo.activeSelf)
            {
                passChallengeGo.SetActive(false);
                passChallengeInterval += dTime;
            }
            else if (passChallengeInterval == 0)
            {
                passChallengeText.text = (showNexPassChallengeText ? "扫荡" : "打门票");
                passChallengeGo.transform.localScale = Vector3.one * 0.01f;
                passChallengeGo.transform.DOScale(1, 0.2f);
                passChallengeGo.SetActive(true);
                passChallengeInterval += dTime;

                showNexPassChallengeText = !showNexPassChallengeText;
            }
            else
                passChallengeInterval += dTime;
        }
    }
    #endregion

    #region 战斗血量变化
    /// <summary>
    ///  boss 相关
    /// </summary>
    private int bossHpbarSliderIndex; // slider轮换索引
    private int bossHpbarNum; // 血条数量
    private float bossSingleBarHp; // 每个血条的实际血量
    private float bossSingleBarCurHp; // 当前血条的血量
    private const int bossSingleBarMaxHp = 20000; // 每个血条的最大血量
    private int colorIndex; // 颜色索引
    private string[] colorNames = new string[] { "#2a76c4", "#5cc42a", "#d6c730", "#c4692a", "#c42a2a" };
    private Slider[] bossHpbarSliders;
    private Image[] bossHpbarSliderBgImgs;

    private void InitBossHp()
    {
        bossHpbarSliders = new Slider[] { uIViewHall.bossMonsterHPSlider2, uIViewHall.bossMonsterHPSlider1 };
        bossHpbarSliderBgImgs = new Image[] { uIViewHall.bossMonsterHPImg2, uIViewHall.bossMonsterHPImg1 };
        uIViewHall.boosMonsterHPHudGo.SetActive(false);

        uIViewHall.bossMonsterHPBtn.onClick.RemoveAllListeners();
        uIViewHall.bossMonsterHPBtn.onClick.AddListener(() =>
        {
            void BackToRest()
            {
                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    ChangeState(HallBattleState.Rest);
                }, 0.6f / timeScale);

                ShowZhuanLoading();
            }

            if(GuidanceM.Instance.isGuiding)
            {
                BackToRest();
                return;
            }

            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = "是否中止挑战，结束当前战斗",
                NoFun = () => { },
                YesFun = () =>
                {
                    BackToRest();
                }
            });
        });
    }

    public void SetMainBossHpUIActive()
    {
        uIViewHall.boosMonsterHPHudGo.SetActive(curMainBattleCfg.flag == (int)MainBattleFlayType.Boss && hallBattleState == HallBattleState.Battle);
    }

    // boss血条，每条血量 20000，最多99个
    public void ResetBossHpUIData(long totalHp)
    {
        if (curMainBattleCfg.flag != (int)MainBattleFlayType.Boss)
            return;

        colorIndex = colorNames.Length - 1;
        bossHpbarSliderIndex = 0;
        if (totalHp <= bossSingleBarMaxHp)
        {
            bossHpbarNum = 1;
            bossSingleBarHp = totalHp;
        }
        else
        {
            var num = totalHp / bossSingleBarMaxHp;
            bossHpbarNum = (int)num + (totalHp % bossSingleBarMaxHp == 0 ? 0 : 1);
            bossSingleBarHp = totalHp / bossHpbarNum;

            var cNum = (bossHpbarNum >= colorNames.Length ? (bossHpbarNum % colorNames.Length) : bossHpbarNum);
            colorIndex = (cNum == 0 ? 0 : (colorNames.Length - cNum % colorNames.Length));
        }
        uIViewHall.boosMonsterHPHudGo.SetActive(true);
        bossSingleBarCurHp = bossSingleBarHp;

        //Debug.LogError($"totalHp: {totalHp}, colorIndex:{colorIndex}, bossSingleBarCurHp:{bossSingleBarCurHp}, bossSingleBarHp:{bossSingleBarHp}, bossHpbarNum:{bossHpbarNum}");
        //Debug.LogError($"bossHpbarSliderIndex: {bossHpbarSliderIndex}, colorIndex:{colorIndex}");
        bossHpbarSliders[bossHpbarSliderIndex].transform.SetSiblingIndex(1);
        bossHpbarSliders[bossHpbarSliderIndex].value = 1;
        bossHpbarSliderBgImgs[bossHpbarSliderIndex].color = GetSliderColor(colorNames[colorIndex]);

        bossHpbarSliders[bossHpbarSliderIndex + 1].transform.SetSiblingIndex(1);
        bossHpbarSliders[bossHpbarSliderIndex + 1].value = 1;

        if (bossHpbarNum == 1)
            bossHpbarSliders[bossHpbarSliderIndex + 1].value = 0;
            
        if (colorIndex == colorNames.Length - 1)
            colorIndex = 0;
        else
            colorIndex++;

        bossHpbarSliderBgImgs[bossHpbarSliderIndex + 1].color = GetSliderColor(colorNames[colorIndex]);

        uIViewHall.bossMonsterTitleText.text = curMainBattleCfg.name;
        RefreshBossHpUI(0);
    }

    private Color GetSliderColor(string hexColor)
    {
        if (ColorUtility.TryParseHtmlString(hexColor, out Color color))
            return color;

        hexColor = colorNames[0];

        if (ColorUtility.TryParseHtmlString(hexColor, out Color color1))
            return color1;

        return Color.red;
    }

    /// <summary>
    /// 刷新Boss血条
    /// </summary>
    /// <param name="value">当前受击掉血</param>
    private void RefreshBossHpUI(long value)
    {
        if (value >= bossSingleBarCurHp)
        {
            while(value >= bossSingleBarCurHp)
            {
                bossHpbarNum--;
                if (bossHpbarNum == 0)
                {
                    if (!GuidanceM.Instance.isGuiding)
                        uIViewHall.boosMonsterHPHudGo.SetActive(false);

                    bossHpbarSliders[0].value = 0;
                    bossHpbarSliders[1].value = 0;
                    uIViewHall.bossMonsterScheduleText.text = "0%";
                    uIViewHall.bossMonsterNumText.text = string.Empty;
                    return;
                }

                if (bossHpbarNum > 1)
                {
                    bossHpbarSliders[bossHpbarSliderIndex].transform.SetSiblingIndex(1);
                    bossHpbarSliders[bossHpbarSliderIndex].value = 1;

                    if (colorIndex == colorNames.Length - 1)
                        colorIndex = 0;
                    else
                        colorIndex++;

                    bossHpbarSliderBgImgs[bossHpbarSliderIndex].color = GetSliderColor(colorNames[colorIndex]);
                }
                else
                    bossHpbarSliders[bossHpbarSliderIndex].value = 0;

                if (bossHpbarSliderIndex == 1)
                    bossHpbarSliderIndex = 0;
                else
                    bossHpbarSliderIndex++;

                value -= (long)bossSingleBarCurHp;
                bossSingleBarCurHp = bossSingleBarHp;
            }
        }

        bossSingleBarCurHp -= value;
        if (bossSingleBarCurHp < 0) bossSingleBarCurHp = 0;
        var per = bossSingleBarCurHp / bossSingleBarHp;
        bossHpbarSliders[bossHpbarSliderIndex].value = per;

        //Debug.LogError($"{value}, {bossSingleBarCurHp}, {bossSingleBarHp}, {bossHpbarSliderIndex}, per:{per}");
        uIViewHall.bossMonsterScheduleText.text = (per * 100).ToString("F1") + "%";
        uIViewHall.bossMonsterNumText.text = "x" + bossHpbarNum;
    }

    private void HallBattleBossBeHitEnt(HallBattleBossBeHit args)
    {
        RefreshBossHpUI((long)args.attackPower);
    }
    #endregion

    #region 怪物大招
    public void ShowBossMasterStroke()
    {
        uIViewHall.bossMasterStrokeTrans.SetParent(uIViewHall.roleTranss[0]);
        uIViewHall.bossMasterStrokeTrans.localPosition = new Vector3(1680, 1188, 0);
        uIViewHall.bossMasterStrokeTrans.localEulerAngles = new Vector3(0, 0, 40);
        uIViewHall.bossMasterStrokeTrans.gameObject.SetActive(true);
        SpineUtil.Instance.PlayOnce(uIViewHall.bossMasterStrokeTrans, "animation", () =>
        {
            uIViewHall.bossMasterStrokeTrans.gameObject.SetActive(false);
        }, timeScale);

        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            ResetBossMasterStroke();
        }, 0.5f / timeScale);
    }

    private void ResetBossMasterStroke()
    {
        uIViewHall.bossMasterStrokeTrans.SetParent(uIViewHall.transform);
        uIViewHall.bossMasterStrokeTrans.localScale = Vector3.one;
    }
    #endregion

    #region 打完一个小怪，处理数据
    public void RoundOver()
    {
        ResetBossMasterStroke();
    }
    #endregion

    #region 掉落武器
    private bool normalMoveEnable;

    private void HallBattleDropItemEnt(HallBattleDropItem args)
    {
        normalMoveEnable = true;
        var dropItems = HallBattleM.Instance.GetDropItem(curMainBattleCfg.monster.Length);

        if (dropItems == null || dropItems.Length <= 0)
            return;

        var maxX = uIViewHall.normalRewardRect.rect.width * 0.5f - 80;
        var limitX = maxX * (0.5f + curMainBattleCfg.flag * 0.25f);
        var rightLimitX = Math.Min(limitX, Screen.width * 0.5f - args.position.x - 80);
        limitX += args.position.x;

        for (int i = 0;i < dropItems.Length; i++)
        {
            var endPos = Vector3.zero;
            endPos.x = UnityEngine.Random.Range(-limitX + 150, rightLimitX);
            endPos.y = UnityEngine.Random.Range(-20, 200);
            DropItem(dropItems[i].ID, dropItems[i].Quality, args.position, endPos, true).Forget();
        }
    }

    public void ResetNormalReward()
    {
        var pos = uIViewHall.normalRewardTrans.localPosition;
        pos.x = 0;
        uIViewHall.normalRewardTrans.localPosition = pos;
    }

    // 普通关卡通关奖励随着地面一起移动
    public void NormalRewardMoveWithBg(float value)
    {
        if (uIViewHall.normalRewardTrans.childCount <= 1 || !normalMoveEnable)
            return;

        var pos = uIViewHall.normalRewardTrans.localPosition;
        pos.x -= value;
        uIViewHall.normalRewardTrans.localPosition = pos;
    }

    private async UniTask DropItem(int cfgID, int quality, Vector3 pos, Vector3 endPos, bool endDie)
    {
        //Debug.LogError($"HallBattleDropItemEnt: {cfgID}, {quality}");

        // 经验特殊处理
        if (RoleM.Instance.MainBattleId <= 5 && cfgID == 2)
            cfgID = 30;

        var itemCfg = CFG_Item.GetById(cfgID);
        if (itemCfg == null)
            return;

        var poolKey = "UIHall_Drop_Item";
        if (!PoolHelper.Instance.CheckPool(poolKey))
        {
            var go = GameObject.Instantiate(norRewardItemGo);
            PoolHelper.Instance.AddPool(poolKey, go);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
        var itemGo = PoolHelper.Instance.GetPoolObj(poolKey);

        // 设置资源
        var itemTrans = itemGo.transform;
        itemTrans.SetParent(uIViewHall.normalRewardTrans);
        var rewardScale = 1f;
        if (itemCfg.Type == 5)
            rewardScale = 1.5f;
        var lightTrans = itemTrans.Find("Item/LightGo");
        var iconTrans = itemTrans.Find("Item/IconImg");
        iconTrans.localScale = Vector3.one * rewardScale;
        lightTrans.localScale = Vector3.one * (quality > 1 ? 0.45f : 0.4f);

        GameUtil.Instance.SetItemIcon(cfgID, JUIExtensions.img(iconTrans));
        //Debug.LogError($"HallBattleDropItemEnt2: {cfgID}, {quality}, {itemCfg.Name}, {norRewardSpineNames[quality - 1]}");
        SpineUtil.Instance.PlayLoop(lightTrans, norRewardSpineNames[quality - 1], 1.8f);

        // 设置位置
        var startPos = -uIViewHall.normalRewardTrans.localPosition;
        startPos.x += pos.x;
        startPos.y = 0;
        itemTrans.localPosition = startPos;
        itemTrans.localScale = Vector3.one;

        var dropItem = itemTrans.Find("Item").GetComponent<HallDropItem>();
        if(dropItem == null)
            dropItem = itemTrans.Find("Item").gameObject.AddComponent<HallDropItem>();
        dropItem.Reset();
        dropItem.Set(Vector3.zero, endPos, () =>
        {
            if(endDie)
            {
                if(!isNullGo && itemGo != null && itemGo != null)
                {
                    DOTween.To(() => 0, (value) => { }, 1, UnityEngine.Random.Range(0f, 0.5f)).SetEase(Ease.Linear).OnComplete(() =>
                    {
                        if (isNullGo)
                            return;

                        var endPos = hallRole.transform.position;
                        var duTime = 0.35f + UnityEngine.Random.Range(-0.15f, 0.15f);
                        itemTrans.DOMove(endPos, duTime / timeScale).SetEase(Ease.OutCubic).OnComplete(() =>
                        {
                            if (!isNullGo && itemGo != null)
                            {
                                itemGo.SetActive(false);
                                PoolHelper.Instance.RecycObj(itemGo);
                            }
                        });

                        itemTrans.DOScale(Vector3.one * 0.3f, duTime / timeScale);
                    });
                }
            }else
            {
                var endPos = GetNormalRewardFlyEndPos(cfgID);
                endPos.x -= uIViewHall.normalRewardTrans.localPosition.x;
                dropItem.transform.DOLocalMove(endPos, 0.5f / timeScale).SetEase(Ease.OutCubic).OnComplete(() =>
                {
                    if (!isNullGo && itemGo != null)
                    {
                        itemGo.SetActive(false);
                        PoolHelper.Instance.RecycObj(itemGo);
                    }
                });

                dropItem.transform.DOScale(Vector3.one * 0.6f, 0.5f / timeScale);
            }
        });

        itemGo.SetActive(true);
    }
    #endregion

    #region 关卡结算
    private Action<bool> settlementCallBack;
    private GameObject norRewardItemGo;
    private TimeItem bossRewardTimeItem;
    private string[] norRewardSpineNames = new string[] {"1_hui", "2_lv", "3_lan", "4_zi", "5_fen", "6_cheng", "7_hong" };
    private bool curIsBossLevel = false;
    private TimeItem autoClothDelay;

    private void InitSettlement()
    {
        var rewardTrans = uIViewHall.normalRewardTrans;
        if (norRewardItemGo == null)
        {
            norRewardItemGo = rewardTrans.GetChild(0).gameObject;
            norRewardItemGo.SetActive(false);
        }
    }

    public void Settlement(Action<bool> callBack)
    {
        normalMoveEnable = false;
        settlementCallBack = callBack;
        TcpC.C2S_MainBattle(curMainBattleCfg.ID);
    }

    private void Sever_MainBattleResponse(MainBattleResponse resq)
    {
        //Debug.LogError($"打完第{resq.MainBattleId}关，{resq.Code}, {resq.Error}, {resq.Items}");
        if (resq.Code != 0)
        {
            UITool.ShowTips(resq.Error);
            return;
        }

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "首页战斗玩法",
            position = "",
            event_click = $"首页战斗通关",
        });

        if (GuidanceM.Instance.isGuiding)
            UniGameEvent.Send(new GameEvent.Refresh_Guide { completeType = GuidanceCompleteType.MainBattle, value = resq.MainBattleId });

        curIsBossLevel = false;
        // 重新刷关，可刷的只有前一个boss和通关前一章节
        var isBrush = (RoleM.Instance.MainBattleId >= resq.MainBattleId) && (curMainBattleCfg.flag == (int)MainBattleFlayType.Boss || CFG_MainBattle.IsPassAll(RoleM.Instance.MainBattleId));
        if (isBrush)
        {
            if (curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
            {
                curIsBossLevel = true;
                ShowBossRewards(resq.Items);
            }
            else
                ShowNormalReward(resq.Items);
        }
        else if (curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
        {
            curIsBossLevel = true;
            // boss 关卡，弹奖励，
            ShowBossRewards(resq.Items);
            setMainBattleLevel(resq.MainBattleId);
            RoleM.Instance.MainBattleId = resq.MainBattleId;
            passAllLevel = CFG_MainBattle.IsPassAll(resq.MainBattleId);

            CheckRed();
            CheckAutoSpecialSkill();

            // 最后一关为5-4BOSS关，前置一关为5-3，打完后自动回城，
            if (passAllLevel)
            {
                ChangeState(HallBattleState.Rest);
                return;
            }
            else
            {
                // 刷新level，新进度条
                RefreshConfigData();

                RefreshCurConfigData(resq.MainBattleId + 1);
                RefreshLevelUI(resq.MainBattleId + 1);
            }
        }else
        {
            // 普通关卡
            ShowNormalReward(resq.Items);

            setMainBattleLevel(resq.MainBattleId);

            RoleM.Instance.MainBattleId = resq.MainBattleId;

            CheckRed();
            CheckAutoSpecialSkill();

            RefreshConfigData();

            // 刷新新一关的关卡数据
            RefreshCurConfigData(resq.MainBattleId + 1);
            RefreshLevelUI(resq.MainBattleId + 1);
        }

        // 之前开宝箱UI
        for(int i = 0; i < resq.Items.Count; i++)
        {
            var itemCfg = CFG_Item.GetById(resq.Items[i].CfgId);
            if (itemCfg.Type == 5)
                RoleM.Instance.OpenBoxItems.Add(resq.Items[i]);
        }

        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            if (isNullGo) return;

            OnRefreshHomeBox(default);

            if(autoHallBattle)
            {
                if (uIViewHall.boxClothingBtn.gameObject.activeSelf)
                    OnClickBoxClothingButton();
                else
                    OnClickBoxDelButton();
            }
            else
            {
                //local server: advance to next stage even if auto-battle is off
                MoveToNext();
            }
        }, 1f);
    }

    private void setMainBattleLevel(int mainLevel)
    {
        //int curMainBattleID = RoleM.Instance.MainBattleId;
        //mainBattle变更时SendMessage
        bool mainBattleIsChange = RoleM.Instance.MainBattleId != 0 && RoleM.Instance.MainBattleId < mainLevel;

        RoleM.Instance.MainBattleId = mainLevel;

        //mainBattle变更时SendMessage
        if (mainBattleIsChange)
        {
            UniGameEvent.Send(new GameEvent.Refresh_MainBattleChanged());
        }
    }

    private void MoveToNext()
    {
        // 体力不足
        if (RoleM.Instance.Energy < curMainBattleCfg.cost)
        {
            EnergyDeficiencyBackToRest();
            return;
        }
        else if (curMainBattleCfg.flag == (int)MainBattleFlayType.Boss) // 下一关是boss
        {
            EnterBoss(() => {
                settlementCallBack?.Invoke(true);
                settlementCallBack = null;
            });
        }
        else if (curIsBossLevel) // 刚打完boss，开始下一关
            ExitBoss(true);
        else
        {
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (isNullGo) return;
                settlementCallBack?.Invoke(true);
                settlementCallBack = null;
            }, 0.6f / timeScale);
        }
    }

    // 体力不足，因体力不足回到休息室，给予飘字提示：“体力不足，战斗暂停！”此时跑动动画停止，切场动画---回到休息室
    private void EnergyDeficiencyBackToRest()
    {
        UITool.ShowTips($"体力不足，战斗暂停！");
        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            if (isNullGo) return;
            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (isNullGo) return;
                ChangeState(HallBattleState.Rest);
            }, 0.6f / timeScale);

            ShowZhuanLoading();
        }, 0.8f / timeScale);
    }

    // 进入boss关卡,过场--->特效（玩家需要跑动过程）--->BOSS出现   （可以先本地优化一版看效果）
    private void EnterBoss(Action callBack)
    {
        ShowZhuanLoading(() =>
        {
            EnterBossWithoutLoading(callBack);
        });
    }

    // boss来袭
    public void EnterBossWithoutLoading(Action callBack)
    {
        if (isNullGo) return;

        uIViewHall.bossCommingTrans.gameObject.SetActive(true);
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_BossComing);
        SpineUtil.Instance.PlayOnce(uIViewHall.bossCommingTrans, "animation", () =>
        {
            if (isNullGo) return;
            uIViewHall.bossCommingTrans.gameObject.SetActive(false);
        }, timeScale);

        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            if (isNullGo) return;

            callBack?.Invoke();
        }, 1 / timeScale);
    }

    private void ExitBoss(bool isSettlement)
    {
        GameHelper.Instance.SetDeayTimeDo(() =>
        {
            settlementCallBack?.Invoke(true);
            settlementCallBack = null;
        }, 0.6f / timeScale);

        ShowZhuanLoading();
    }

    // 普通关卡通关奖励
    private void ShowNormalReward(List<ItemInfo> items)
    {
        normalMoveEnable = true;
        ResetNormalReward();
        uIViewHall.normalRewardTrans.gameObject.SetActive(true);

        // 金币 / 10000；元宝跟随 flag + 1，其它一个 
        var limitX = uIViewHall.normalRewardRect.rect.width * 0.5f - 80;
        var index = 0;
        for (int i = 0; i < items.Count; i++)
        {
            // 不同物品飞的数量不同
            var itemInfo = items[i];
            var cfgID = itemInfo.CfgId;

            var itemCfg = CFG_Item.GetById(cfgID);
            if (itemCfg == null)
                continue;

            // 经验特殊处理
            if (cfgID == 2)
            {
                if (itemInfo.Num < 50000)
                    cfgID = 30;
                else if (itemInfo.Num < 100000)
                    cfgID = 31;
                else if (itemInfo.Num < 500000)
                    cfgID = 32;
                else if (itemInfo.Num < 1000000)
                    cfgID = 33;
                else
                    cfgID = 34;
            }

            var quality = itemCfg.Quality;
            var count = 1;
            if(itemCfg.Type == 5) // 装备
            {
                var dropItem = HallBattleM.Instance.GetDropItemByID(cfgID);
                if (dropItem != null)
                    quality = dropItem.Quality;
            }
            
            if (itemInfo.CfgId == 1)
                count = (int)itemInfo.Num / 10000; // 金币
            else if (itemInfo.CfgId == 5)
                count = UnityEngine.Random.Range(1, 3);

            //Debug.LogError($"{cfgID}, {count}");
            for (int j = 0; j < count; j++)
            {
                Vector3 targetPos = Vector3.zero;
                targetPos.x = UnityEngine.Random.Range(-limitX, limitX);
                targetPos.y = UnityEngine.Random.Range(-30, 180);
                DropItem(itemInfo.CfgId, quality, Vector3.zero, targetPos, false).Forget();
            }
            index++;
        }
    }

    private Vector3 GetNormalRewardFlyEndPos(int configId)
    {
        Transform trans;
        Vector3 pos;
        if (configId == 1)
            trans = topUI.GetView().txt_gold.transform;
        else if (configId == 5)
            trans = topUI.GetView().txt_coupon.transform;
        else
        {
            var bagTrans = downUI.GetView().Btns[4].trans;
            pos = bagTrans.localPosition + bagTrans.parent.localPosition;
            return pos;
        }

        pos = trans.localPosition + trans.parent.localPosition + trans.parent.parent.localPosition;
        return pos;
    }

    // boss通关奖励
    private void ShowBossRewards(List<ItemInfo> items)
    {
        var itemNode = uIViewHall.battleBossRewardTrans.Find("Items");
        for (int i = 0; i < itemNode.childCount; i++)
        {
            var rewardItem = itemNode.GetChild(i);
            rewardItem.gameObject.SetActive(items.Count > i);

            if (items.Count > i)
            {
                var cfgID = items[i].CfgId;
                GameUtil.Instance.SetItemIcon(cfgID, JUIExtensions.img(rewardItem.Find("icon")), JUIExtensions.img(rewardItem));
                JUIExtensions.txt(rewardItem, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(items[i].Num);

                JUIExtensions.btn(itemNode).onClick.RemoveAllListeners();
                JUIExtensions.btn(itemNode).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgID });
                });
            }
        }

        uIViewHall.battleBossRewardTrans.gameObject.SetActive(true);

        ClearBossRewardTimer();
        var timeText = uIViewHall.battleBossRewardTrans.Find("txt_time").GetComponent<Text>(); ;
        var count = 4f;
        bossRewardTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(count, 1, () =>
        {
            timeText.text = (int)(count - 1) + "秒后将自动退出";
            count -= timeScale;
            if (count <= 1)
            {
                ClearBossRewardTimer();
                uIViewHall.battleBossRewardTrans.gameObject.SetActive(false);
            }
        }));
    }

    private void ClearBossRewardTimer()
    {
        if (bossRewardTimeItem != null && GameHelper.Instance != null)
        {
            GameHelper.Instance.RemoveInvoke(bossRewardTimeItem);
            bossRewardTimeItem = null;
        }
    }
    #endregion

    #region 变强
    public void ShowStrengthenUI()
    {
        UITool.Open<UIStrengthen>();
    }
    #endregion
}