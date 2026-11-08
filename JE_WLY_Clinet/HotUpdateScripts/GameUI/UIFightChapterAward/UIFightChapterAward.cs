using UnityEngine;
using MetaFramework.UI;
using JEngine.Helper;
using JEngine.UI;
using SuperScrollView;
using HotUpdateScripts.ProtoM;
using JEngine.Event;

public class UIFightChapterAward : UIBase<UIModelFightChapterAward, UIViewFightChapterAward>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btnClose.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightChapterAward>();
        });

        GetView().loopListView.InitListView(0, OnFillItem);
    }

    private LoopListViewItem2 OnFillItem(LoopListView2 listView, int index)
    {
        LoopListViewItem2 viewitem = listView.NewListViewItem("Item");
        Transform item = viewitem.transform;

        var cfgInfo = CFG_Battle.AllFightChapterAwardList[index];
        JUIExtensions.txt(item, "txt_Name").text = "通关" + cfgInfo.Chapter + "-" + cfgInfo.Layers + "关可领取";
        JUIExtensions.txt(item, "txt_Get").text = GameUtil.Instance.FormatNumWithUnit(cfgInfo.Powerlimit);

        var root_award = item.Find("awards");
        for (int i = 0; i < root_award.childCount; i++)
        {
            int idx = i;
            var go = root_award.GetChild(idx);
            go.gameObject.SetActive(idx < cfgInfo.ChapterAward.Length);
            if (idx < cfgInfo.ChapterAward.Length)
            {
                GameUtil.Instance.SetItemIcon(cfgInfo.ChapterAward[idx].id, JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                JUIExtensions.txt(go, "txt").text = GameUtil.Instance.FormatNumWithUnit(cfgInfo.ChapterAward[idx].num);
                JUIExtensions.btn(go).onClick.RemoveAllListeners();
                JUIExtensions.btn(go).onClick.AddListener(() => {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgInfo.ChapterAward[idx].id });
                });
            }
        }

        bool isPlay = FightPlayM.Instance.GetInfo().battle_map.RewardStage < cfgInfo.UID && (cfgInfo.UID <= FightPlayM.Instance.GetInfo().battle_map.BattleId);
        item.Find("btn_Play").gameObject.SetActive(isPlay);
        JUIExtensions.txt(item, "txt_Get").text = "";
        if (isPlay)
        {
            JUIExtensions.btn(item, "btn_Play").onClick.RemoveAllListeners();
            JUIExtensions.btn(item, "btn_Play").onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                TcpC.C2S_BattleReward(CFG_Room.RoomType.Fight, cfgInfo.Chapter);
            });
        }
        else
        {
            var txt = JUIExtensions.txt(item, "txt_Get");
            txt.text = cfgInfo.UID > FightPlayM.Instance.GetInfo().battle_map.BattleId ? "未通关":"已领取";
            txt.color = JUIExtensions.ToColor(cfgInfo.UID > FightPlayM.Instance.GetInfo().battle_map.BattleId ? "#E96C5C" : "#EDE6CD");
            JUIExtensions.outline(item, "txt_Get").effectColor = JUIExtensions.ToColor(cfgInfo.UID > FightPlayM.Instance.GetInfo().battle_map.BattleId ? "#9E2522" : "#4C4334");
        }

        return viewitem;
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        OnShowUI();
    }

    private void OnShowUI()
    {
        if (GetView().loopListView.ItemTotalCount != CFG_Battle.AllFightChapterAwardList.Count)
            GetView().loopListView.SetListItemCount(CFG_Battle.AllFightChapterAwardList.Count);
        GetView().loopListView.RefreshAllShownItem();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
    }

    [Subscriber(ThreadMode.Main)]
    private void Refresh_InstanceInfo(GameEvent.Refresh_InstanceInfo _e)
    {
        if (isNullGo) return;
            OnShowUI();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}