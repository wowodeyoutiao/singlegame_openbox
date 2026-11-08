using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;

public class UIGuildDonation : UIBase<UIModelGuildDonation, UIViewGuildDonation>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuildDonation>();
        });

        GetView().Btn_Go_GuildDonation.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuildDonation>();
            // 跳转到工会信息
            JEvent.defaultEvent.Post(new UIGuild.SelectionSubPage { bottomButtonIndex = 0, subButtonIndex = 0 });
        });

        GetView().Btn_Go_SkillPray.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuildDonation>();
            // 跳转到工会技能
            JEvent.defaultEvent.Post(new UIGuild.SelectionSubPage { bottomButtonIndex = 0, subButtonIndex = 2 });
        });

        GetView().Btn_Go_ChallengeLeader.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGuildDonation>();
            // 跳转到工会地图
            JEvent.defaultEvent.Post(new UIGuild.SelectionSubPage { bottomButtonIndex = 0, subButtonIndex = 4 });
        });

    }

    public override void OnEnter(params object[] args)
    {

    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {
        
    }
}