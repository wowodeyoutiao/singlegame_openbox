using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using JEngine.UI;
using HotUpdateScripts.ProtoM;
using HotUpdateScripts;

public class UIRelicChallange : UIBase<UIModelRelicChallange, UIViewRelicChallange>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRelicChallange>();
        });
        GetView().btn_Play.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GameHelper.Instance.StarGame(BattleId);
        });

        GetView().btn_Sd.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (FightPlayM.Instance.GetInfo().relic_map.ChallengeTime == 0)
            {
                if (StoreM.Instance.GetCanBuyNum((int)CFG_Store.StroeID.BossChange) == 0)
                {
                    UITool.ShowTips("今日扫荡次数已不足") ;
                    return;
                }
                UITool.ShowTips("扫荡次数不足，请先去购买吧");
                return;
            }
            UITool.Open<UIGameSweepAndBuy>(null,UIGameSweepAndBuy.OpType.RelicSweep, BattleId);
            UITool.Close<UIRelicChallange>();
        } );
    }

    private int BattleId;

    public override void OnEnter(params object[] args)
    {
        BattleId = (int)args[0];
        int NextBattleId = (int)args[1];


        GetView().btn_Sd.gameObject.SetActive( FightPlayM.Instance.GetRelicLevelStar(BattleId) == 3);
        GetView().btn_Play.gameObject.SetActive(FightPlayM.Instance.GetRelicLevelStar(BattleId) != 3);

        var CFGBattle = CFG_Battle.GetById(BattleId);

        GetView().txt_title.text = $"第{CFGBattle.Chapter}-{CFGBattle.Layers}关";

        for (int i = 0; i < GetView().root_award.childCount; i++)
        {
            int idx = i;
            var go = GetView().root_award.GetChild(idx);
            go.gameObject.SetActive(idx < CFGBattle.Preview.Length);
            if (idx < CFGBattle.Preview.Length)
            {
                GameUtil.Instance.SetItemIcon(CFGBattle.Preview[idx], JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() => {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = CFGBattle.Preview[idx] });
                });
            }
        }

        for (int i = 0; i < GetView().root_des.childCount; i++)
        {
            bool isUnlock = FightPlayM.Instance.GetRelicLevelStar(BattleId) > i;
            GetView().root_des.GetChild(i).Find("on").gameObject.SetActive(isUnlock);
            GetView().root_des.GetChild(i).Find("off").gameObject.SetActive(!isUnlock);
            JUIExtensions.txt(GetView().root_des.GetChild(i)).text = $"<color=#{(!isUnlock ? "FFFFFF" : "94E771")}>{GameConst.RelicItemStrs[i]}</color>";
            JUIExtensions.outline(GetView().root_des.GetChild(i)).effectColor = !isUnlock ? new Color(188/255f, 188 / 255f, 188 / 255f, 1f):new Color(49/255f, 101/255f, 27 / 255f, 1f);
        }
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}