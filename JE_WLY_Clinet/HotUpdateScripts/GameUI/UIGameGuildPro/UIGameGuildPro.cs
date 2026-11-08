using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts;
using JEngine.UI;
using JEngine.Event;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;

public class UIGameGuildPro : UIBase<UIModelGameGuildPro, UIViewGameGuildPro>
{
    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }
    public override void OnInit()
    {
        GetView().btn_Show.onClick.AddListener(delegate
        {
            OnUIShow(true);
        });

        GetView().btn_Close.onClick.AddListener(delegate
        {
            OnUIShow(false);
        });

        for (int i = 0; i < 3; i++)
        {
            int idx = i;
            GetView().award1Items[i]._btn.onClick.AddListener(delegate
            {
                if (idx < GameRoomM.Instance.guildCopyReward.CurReward.Count)
                {
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = GameRoomM.Instance.guildCopyReward.CurReward[idx].K });
                }
            });

            GetView().award2Items[i]._btn.onClick.AddListener(delegate
            {
                if (idx < GameRoomM.Instance.guildCopyReward.NexReward.Count)
                {
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = GameRoomM.Instance.guildCopyReward.NexReward[idx].K });
                }
            });
        }

        JEvent.defaultEvent.Register(this);

        OnGuildProShow();
    }

    private void OnUIShow(bool isShow)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        GetView().obj_UI.SetActive(isShow);
        GetView().btn_Show.gameObject.SetActive(!isShow);
    }

    public void OnGuildProShow()
    {
        GetView().txt_NoAward.text = GameRoomM.Instance.guildCopyReward.CurReward.Count == 0 ? "无奖励" : "";

        for (int i = 0; i < GetView().award1Items.Length; i++)
        {
            GetView().award1Items[i]._btn.gameObject.SetActive(i < GameRoomM.Instance.guildCopyReward.CurReward.Count);
            if (i < GameRoomM.Instance.guildCopyReward.CurReward.Count)
            {
                GameUtil.Instance.SetItemIcon(GameRoomM.Instance.guildCopyReward.CurReward[i].K, GetView().award1Items[i]._imgIcon, GetView().award1Items[i]._imgBg);
                GetView().award1Items[i]._txtNum.text = GameUtil.Instance.FormatNumWithUnit(GameRoomM.Instance.guildCopyReward.CurReward[i].V);
            }
        }
        for (int i = 0; i < GetView().award2Items.Length; i++)
        {
            GetView().award2Items[i]._btn.gameObject.SetActive(i < GameRoomM.Instance.guildCopyReward.NexReward.Count);
            if (i < GameRoomM.Instance.guildCopyReward.NexReward.Count)
            {
                GameUtil.Instance.SetItemIcon(GameRoomM.Instance.guildCopyReward.NexReward[i].K, GetView().award2Items[i]._imgIcon, GetView().award2Items[i]._imgBg);
                GetView().award2Items[i]._txtNum.text = GameUtil.Instance.FormatNumWithUnit(GameRoomM.Instance.guildCopyReward.NexReward[i].V);
            }
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnRefsh(GameEvent.Refrsh_GuildCopyReward data)
    {
        OnGuildProShow();
    }
    public override void OnEnter(params object[] args)
    {
        
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
    }
}