using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using SuperScrollView;

using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;

using UnityEngine;

using static CFG_RankingAward;

public class UIWorldBoss : UIBase<UIModelWorldBoss, UIViewWorldBoss>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIWorldBoss>();
        });

        GetView().btnHelp.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 15);
        });

        GetView().btn_Play.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss))
                return;

            if (!GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.WorldBoss.startTime, CFG_XMLConfig.xmlConfig.WorldBoss.endTime))
            {
                UITool.ShowTips($"每天 {CFG_XMLConfig.xmlConfig.WorldBoss.startTime}~{CFG_XMLConfig.xmlConfig.WorldBoss.endTime} 开启");
                //UITool.ShowTips($"每{GameUtil.Instance.GetWeekDayName(CFG_XMLConfig.xmlConfig.WorldBoss.weekSet)} {CFG_XMLConfig.xmlConfig.WorldBoss.startTime}~{CFG_XMLConfig.xmlConfig.WorldBoss.endTime} 开启");
                return;
            }

            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss))
                return;

            if (RoomCFG.Needenerge > RoleM.Instance.Energy)
            {

                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = "当前体力不足，是否前往获取？",
                    YesFun = () =>
                    {
                        OnEnegerAddBtnClick();
                    }
                });
                return;
            }
            TcpC.GetGameJoin(RoomCFG.ID, 0);
            //UITool.Close<UIWorldBoss>();
        });

        GetView().btn_RegionalServerRank.onClick.AddListener(() =>
        {
            isPersonInRank = false;
            rankClass = RankClass.WorldBossAreaServer;
            SetRankButton();
        });

        GetView().btn_PersonInjuryRank.onClick.AddListener(() =>
        {
            isPersonInRank = true;
            rankClass = RankClass.WorldBoss;
            SetRankButton();
        });


        GetView().infiniteScroll.InitListView(0, OnFillItem);
        OnBaseUI();
    }

    private void OnEnegerAddBtnClick()
    {
        GameUtil.Instance.OnEnegerAdd(null);
    }

    LoopListViewItem2 OnFillItem(LoopListView2 listView, int index)
    {
        var item = listView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        JUIExtensions.img(go.transform.Find("bg")).enabled = index % 2 == 1;
        HTTP_RankUserInfo userInfo = RankListInfo.rank_all.Length > index ? RankListInfo.rank_all[index] : default;
        OnRankItemShow(userInfo, index + 1, go, index >= RankListInfo.rank_all.Length);
        return item;
    }

    private CFGItem_Room RoomCFG;
    private bool isPersonInRank;

    private void SetRankButton()
    {
        GetView().btn_RegionalServerRankSelect.SetActive(!isPersonInRank);
        GetView().btn_PersonInjuryRankSelect.SetActive(isPersonInRank);
        GetView().txt_RankTypeDesc.text = isPersonInRank ?
            "排名              玩家                     累计伤害              奖励"
            : "排名              区服                     累计伤害              奖励";
        OnRankShow();
    }

    private void OnBaseUI()
    {
        GetView().txt_Skill.text = CFG_XMLConfig.xmlConfig.WorldBoss.bossDescribe;

        RoomCFG = CFG_Room.GetByRoomType(CFG_Room.RoomType.WorldBoss, RoleM.Instance.Level);

        for (int i = 0; i < GetView().root_Award.childCount; i++)
        {
            var go = GetView().root_Award.GetChild(i);
            go.gameObject.SetActive(RoomCFG.Droppreview.Length > i);

            if (RoomCFG.Droppreview.Length > i)
            {
                int cfgId = RoomCFG.Droppreview[i];
                GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(go, "icon"), JUIExtensions.img(go), JUIExtensions.img(go, "icon/nature"));
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }

        }
    }

    private TimeItem jAction;

    private void OnAutoGetInfo()
    {
        jAction = GameHelper.Instance.AddInvoke(new TimeItem(-1, 2, () =>
        {
            if (isNullGo)
                return;

            TcpC.C2S_CheckBoss(0);
        }, true));
    }

    private void ClearJation()
    {
        if (GameHelper.Instance != null && jAction != null)
        {
            GameHelper.Instance.RemoveInvoke(jAction);
            jAction = null;
        }
    }

    private int BossID = 0;

    //[Subscriber(ThreadMode.Main)]
    private void OnCheck_WorldBoss(GameEvent.Check_WorldBoss _e)
    {
        if (isNullGo)
            return;

        //Debug.Log($"世界boss信息： {_e.HP}  {_e.UsedCount}");
        var bossCFG = CFG_Monster.GetById(_e.BossID);

        if (BossID != _e.BossID)
        {
            BossID = _e.BossID;
            GetView().txt_BossName.text = bossCFG.Name;
            ResouseHelper.Instance.LoadSpriteByAtlas(bossCFG.Photo, bossCFG.Photoname, GetView().img_Moster);
        }

        int allHP = 0;
        for (int i = 0; i < bossCFG.BassProperty.Length; i++)
        {
            if (bossCFG.BassProperty[i].ElementAt(0).Key == ((int)pb.Attribute.MAX_HEALTH).ToString())
            {
                allHP = bossCFG.BassProperty[i].ElementAt(0).Value;
                break;
            }
        }
        GetView().slider_Hp.value = (float)_e.HP / allHP;
        GetView().txt_hp.text = $"{GameUtil.Instance.FormatNumWithUnit(_e.HP)}/{GameUtil.Instance.FormatNumWithUnit(allHP)}";

        //bool isOpen = GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.WorldBoss.startTime, CFG_XMLConfig.xmlConfig.WorldBoss.endTime, CFG_XMLConfig.xmlConfig.WorldBoss.weekSet);
        bool isOpen = GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.WorldBoss.startTime, CFG_XMLConfig.xmlConfig.WorldBoss.endTime);
        bool isShowBtn = isOpen && _e.HP > 0;

        // test code
        // bool isOpen = true;
        // bool isShowBtn = isOpen && _e.HP > 0;

        GetView().obj_Kill.gameObject.SetActive(_e.HP == 0);
        GetView().btn_Play.gameObject.SetActive(isShowBtn);
        GetView().txt_tip.gameObject.SetActive(!isOpen || _e.HP == 0);
        GetView().txt_Challenge.gameObject.SetActive(isShowBtn);

        if (isOpen)
        {
            GetView().txt_tip.text = "BOSS已被击杀";
            if (isShowBtn)
            {
                int remainingChallenges = CFG_XMLConfig.xmlConfig.WorldBoss.enterTime - _e.UsedCount;
                string remainingStr = remainingChallenges > 0 ? $"<color=#319726>{remainingChallenges}</color>" : $"<color=#E64E28>{remainingChallenges}</color>";
                GetView().txt_Challenge.text = $"今日剩余挑战次数:{remainingStr}/<color=#319726>{CFG_XMLConfig.xmlConfig.WorldBoss.enterTime}</color>";
            }
        }
        else
        {
            GetView().txt_tip.text = $"每天{CFG_XMLConfig.xmlConfig.WorldBoss.startTime}~{CFG_XMLConfig.xmlConfig.WorldBoss.endTime}可挑战";
            //GetView().txt_tip.text = $"每{GameUtil.Instance.GetWeekDayName(CFG_XMLConfig.xmlConfig.WorldBoss.weekSet)}{CFG_XMLConfig.xmlConfig.WorldBoss.startTime}~{CFG_XMLConfig.xmlConfig.WorldBoss.endTime}可挑战";
        }
    }

    private HTTP_RankDetail RankListInfo;
    private Dictionary<int, CFGItem_RankingAward> DicCfgRank;
    RankClass rankClass = RankClass.WorldBoss;

    private void OnRankShow()
    {
        HttpC.GetRankList(rankClass, (_r) =>
        {
            Debug.Log($"世界boss排行榜： {_r.data.rank_all.Length}个排行");
            if (isNullGo)
                return;

            RankListInfo = _r.data;
            DicCfgRank = GetRankDicByRankClass(rankClass);

            OnRankItemShow(RankListInfo.rank_own, RankListInfo.rank_own.num, GetView().root_MyItem.gameObject);

            GetView().infiniteScroll.SetListItemCount(DicCfgRank.Count);
            GetView().infiniteScroll.RefreshAllShownItem();
            GetView().infiniteScroll.MovePanelToItemIndex(0, 0);
        });
    }

    private void OnRankItemShow(HTTP_RankUserInfo useInfo, int Rankidx, GameObject go, bool isNull = false)
    {
        Transform award = go.transform.Find("Awards");
        for (int i = 0; i < 3; i++)
        {
            bool isShow = !isNull && DicCfgRank.ContainsKey(Rankidx) && DicCfgRank[Rankidx].award.Length > i;
            award.Find("Item_" + i).gameObject.SetActive(isShow);
            if (isShow)
            {
                int cfgId = int.Parse(DicCfgRank[Rankidx].award[i].ElementAt(0).Key);
                GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(award.Find("Item_" + i + "/icon")), JUIExtensions.img(award.Find("Item_" + i)));
                JUIExtensions.txt(award.Find("Item_" + i + "/txt_num")).text = GameUtil.Instance.FormatNumWithUnit(DicCfgRank[Rankidx].award[i].ElementAt(0).Value);
                JUIExtensions.btn(award.Find("Item_" + i)).onClick.RemoveAllListeners();
                JUIExtensions.btn(award.Find("Item_" + i)).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }
        }

        bool isTop3 = Rankidx >= 1 && Rankidx <= 3;
        JUIExtensions.img(go.transform, "img_Rank").enabled = isTop3;
        JUIExtensions.txt(go.transform.Find("txt_Rank")).enabled = !isTop3;
        if (Rankidx <= 0)
        {
            JUIExtensions.txt(go.transform.Find("txt_Rank")).text = "未上榜";
        }
        else if (isTop3)
        {
            ResouseHelper.Instance.LoadSpriteByAtlas("Common/Atlas/WorldBoss", "icon_" + Rankidx, JUIExtensions.img(go.transform, "img_Rank"));
        }
        else
        {
            JUIExtensions.txt(go.transform.Find("txt_Rank")).text = Rankidx.ToString();
        }

        if (isNull)
        {
            JUIExtensions.txt(go.transform, "txt_Name").text = "虚位以待";
            JUIExtensions.txt(go.transform, "txt_hurt").text = "";
        }
        else
        {
            string name = Rankidx <= 0 ? RoleM.Instance.Nickname : useInfo.name;
            string serverArea = "";

            if (rankClass == RankClass.WorldBoss)
            {
                serverArea = Regex.Match(useInfo.server_name ?? "", "\\d+").ToString();
                if (Rankidx > 0)
                {
                    if (string.IsNullOrEmpty(serverArea))
                        serverArea = $"[{useInfo.server_name}]";
                    else
                        serverArea = $"[{serverArea}服]";
                }
            }

            JUIExtensions.txt(go.transform, "txt_Name").text = rankClass == RankClass.WorldBoss ? $"{serverArea}{name}" : name;
            JUIExtensions.txt(go.transform, "txt_hurt").text = Rankidx <= 0 ? "" : GameUtil.Instance.FormatNumWithUnit(useInfo.value);
        }
    }

    public override void OnEnter(params object[] args)
    {
        isPersonInRank = false;
        rankClass = RankClass.WorldBossAreaServer;
        OnAutoGetInfo();
        SetRankButton();

        UniGameEvent.AddEventListener<GameEvent.Check_WorldBoss>(OnCheck_WorldBoss);
    }

    public override void OnExit()
    {
        UniGameEvent.RemoveEventListener<GameEvent.Check_WorldBoss>(OnCheck_WorldBoss);
        ClearJation();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}