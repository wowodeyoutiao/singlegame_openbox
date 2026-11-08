using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;
using JEngine.MyNet;
using pb;
using JEngine.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using HotUpdateScripts;

public class UIPVPOneReady : UIBase<UIModelPVPOneReady, UIViewPVPOneReady>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);

        GetView().btn_Cancel.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_MatchingRequest((int)matchType,true);
        });
    }

    private TimeItem TimeJa;
    private MatchType matchType;
    public override void OnEnter(params object[] args)
    {

        matchType = (MatchType)args[0];

        var baseInfo = matchType == MatchType.SINGLE_COMBAT ? CFG_XMLConfig.xmlConfig.PVPInfo : CFG_XMLConfig.xmlConfig.GVGInfo;

        int TimeCount = baseInfo.matchTime  + 2 ;
        ClearTimeJa();
        TimeJa = GameHelper.Instance.AddInvoke(new TimeItem(TimeCount, 1, () => {
            TimeCount--;
            GetView().txt_Time.text = $"匹配中(<color=#4EC823>{TimeCount}</color>S)";

            if (TimeCount <= 0)
                TcpC.C2S_MatchingRequest((int)matchType, true);
        }));
            
        GetView().btn_Cancel.gameObject.SetActive(true);
        GetView().obj_Over.SetActive(false);

        JUIExtensions.txt(GetView().root_Playes[0].Find("img/Txt_Name")).text = RoleM.Instance.Nickname;

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + RoleM.Instance.Job + RoleM.Instance.Sex, JUIExtensions.img(GetView().root_Playes[0].Find("bg_head/img_head")));

        for (int i = 1; i < GetView().root_Playes.Count; i++)
        {
            GetView().root_Playes[i].gameObject.SetActive(i <= baseInfo.playerNum);
            if (i <= baseInfo.playerNum)
            {
                JUIExtensions.txt(GetView().root_Playes[i].Find("img/Txt_Name")).text = "虚位以待";
                JUIExtensions.outline(GetView().root_Playes[i].Find("img/Txt_Name")).effectColor = JUIExtensions.ToColor("#696969");
                ResouseHelper.Instance.LoadSpriteByAtlas(Atlas, "bg_huise", JUIExtensions.img(GetView().root_Playes[i].Find("img")));
                GetView().root_Playes[i].Find("bg_head/img_head").gameObject.SetActive(false);
            }
        }
    }

    private string Atlas = "Common/Atlas/FightChoice";
    //[Subscriber(ThreadMode.Main)]
    public void OnMatchResult(GameEvent.PVP_MatchResult _e)
    {
        if (isNullGo)
            return;

        for (int i = 1; i < GetView().root_Playes.Count; i++)
        {
            bool isHave = _e.Enemys.Count > (i - 1);
            GetView().root_Playes[i].Find("bg_head/img_head").gameObject.SetActive(isHave);
            if (isHave)
            {
                JUIExtensions.txt(GetView().root_Playes[i].Find("img/Txt_Name")).text = _e.Enemys[i-1].Nickname;
                JUIExtensions.outline(GetView().root_Playes[i].Find("img/Txt_Name")).effectColor = JUIExtensions.ToColor("#7F3217");
                ResouseHelper.Instance.LoadSpriteByAtlas(Atlas, "bg_mingzi2", JUIExtensions.img(GetView().root_Playes[i].Find("img")));
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + _e.Enemys[i - 1].Job + _e.Enemys[i - 1].Gender, JUIExtensions.img(GetView().root_Playes[i].Find("bg_head/img_head")));
            }
            else
            {
                JUIExtensions.txt(GetView().root_Playes[i].Find("img/Txt_Name")).text = "虚位以待";
                JUIExtensions.outline(GetView().root_Playes[i].Find("img/Txt_Name")).effectColor = JUIExtensions.ToColor("#696969");
                ResouseHelper.Instance.LoadSpriteByAtlas(Atlas, "bg_huise", JUIExtensions.img(GetView().root_Playes[i].Find("img")));
            }
        }

        if (_e.IsBegin)
        {
            ClearTimeJa();
            GetView().txt_Time.text = "";
            GetView().btn_Cancel.gameObject.SetActive(false);
            GetView().obj_Over.SetActive(true);
            TimeJa = GameHelper.Instance.AddInvoke(new TimeItem(2, () =>
            {
                JEvent.defaultEvent.Post(new SenceChange
                {
                    scencePath = GameConst.NorFightPath,
                    CallBack = () =>
                    {
                        HotUpdateScripts.Game.FightController.Instance.OnInit();
                    }
                });
            }));
                
        }
    }

    private void ClearTimeJa()
    {
        if (TimeJa != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeJa);
            TimeJa = null;
        }
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        ClearTimeJa();
    }

    public override void OnDestroy()
    {
        ClearTimeJa();
        JEvent.defaultEvent.Unregister(this);
    }
}