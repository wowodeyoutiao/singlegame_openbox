using System.Collections.Generic;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using Spine;
using JEngine.Helper;
using JEngine.MyNet;
using pb;
using JEngine.Event;
using Spine.Unity;
using UnityEngine;
using HotUpdateScripts;

public class UIPartnerUpStar : UIBase<UIModelPartnerUpStar, UIViewPartnerUpStar>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => { 
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPartnerUpStar>(); 
        });

        GetView().btn_UpStar.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnUpStarBtnClick();
        });

        GetView().btn_SkillUI.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_SkillUI.gameObject.SetActive(false);
        });
        GetView().btn_Skill.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_SkillUI.gameObject.SetActive(true);
        });
    }

    private long Uid = 0;
    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);

        if (args.Length == 1)
        {
            var soltIdx = (int)args[0];
            Uid = PartnerM.Instance.GetOnePartnerSlotInfo(soltIdx).PartnerUid;
        }
        else
            Uid = (long)args[0];

        OnUIPartnerInfo();
    }

    private CFGItem_Partner.StarConsumeItem.Item[] userdData ;
    private void OnUIPartnerInfo()
    {
        var partner = PartnerM.Instance.GetOnePartnerInfoByUid(Uid);
        var cfgInfo = CFG_Partner.GetById(partner.CfgId);

        GetView().txt_name.text = cfgInfo.Name ;
        GetView().txt_Star.text = partner.Star.ToString();

        GetView().txt_nowStar.text =  partner.Star.ToString();
        string Des = "";
        var starInfo = cfgInfo.GetStarAttributeByLv(partner.Star);
        for (int i = 0; i < starInfo.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(starInfo[i].key, starInfo[i].value) + "\n");
        }
        GetView().txt_nowDes.text = Des == "" ?"无": Des;

        //下星级
        int nextStar = partner.Star == cfgInfo.StarConsume.Length ? partner.Star : partner.Star + 1;
        GetView().txt_nextStar.text = nextStar.ToString();
        Des = "";
        var nextstarInfo = cfgInfo.GetStarAttributeByLv(nextStar);
        for (int i = 0; i < nextstarInfo.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(nextstarInfo[i].key, nextstarInfo[i].value) + "\n");
        }
        GetView().txt_nextDes.text = Des;

        GetView().btn_UpStar.gameObject.SetActive(partner.Star < cfgInfo.StarAttribute.Length);

        //消耗
        if (partner.Star < cfgInfo.StarAttribute.Length)
        { 
            userdData = cfgInfo.GetStarConsumeByLv(nextStar);

            GameUtil.Instance.SetItemIcon(userdData[0].key, GetView().img_used1);
            GetView().txt_used1.text = GameUtil.Instance.GetCheakNumStr(userdData[0].key, userdData[0].value);

            GameUtil.Instance.SetItemIcon(userdData[1].key, GetView().img_used2);
            GetView().txt_used2.text = GameUtil.Instance.GetCheakNumStr(userdData[1].key, userdData[1].value);

            bool isRed = true;
            for (int i = 0; i < userdData.Length; i++)
            {
                if (!GameUtil.Instance.ComCoinToCheck(userdData[i].key, userdData[i].value,null,false))
                {
                    isRed = false; break;
                }
            }
            GetView().btn_UpStar.transform.Find("dot").gameObject.SetActive(isRed);
        }
        GetView().obj_used.SetActive(partner.Star < cfgInfo.StarAttribute.Length);

        int NextSkillId = cfgInfo.GetNextStarSkillIdByLv(nextStar);
        if (NextSkillId != 0)
        {
            var cfgSkill = CFG_Skill.GetById(NextSkillId);
            GetView().txt_nextSkill.text =  $"解锁{(cfgSkill.ActiveorPassive == 1 ?"主":"被")}动技能:";
            GetView().txt_nextSkillName.text = $"{cfgSkill.Name}";
            GetView().txt_ui_skillName.text = $"{cfgSkill.Name}";
            GetView().txt_ui_skilltype.text = $"{(cfgSkill.ActiveorPassive == 1 ? "主动" : "被动")}";
            GetView().txt_ui_skilldes.text = cfgSkill.SkillDes;
        }
        GetView().btn_Skill.gameObject.SetActive(NextSkillId != 0);
        GetView().btn_SkillUI.gameObject.SetActive(false);

        GetView().btn_UpStar.gameObject.SetActive(partner.Star != nextStar);
        string poolKey = "Partner_" + partner.CfgId;
        // string poolKey = "Partner_1" ;
        SetRoleSpine(poolKey, GameConst.GetMonsterModelByID(poolKey));
    }

    private void SetRoleSpine(string poolKey, string path)
    {

        if (GetView().root_Spine.childCount > 0)
        {
            PoolHelper.Instance.RecycObj(GetView().root_Spine.GetChild(0).gameObject);
        }

        if (!string.IsNullOrEmpty(poolKey))
        {
            _ = GameHelper.Instance.SetAsyncObj(poolKey, path, (_mainRole) =>
            {
                _mainRole.transform.SetParent(GetView().root_Spine);
                _mainRole.transform.localScale = Vector3.one;
                _mainRole.transform.localPosition = Vector3.zero;
                _mainRole.SetActive(true);

                SpineUtil.Instance.PlayLoop(_mainRole.transform.Find("spine"), GameConst.ModelAnim.stand.ToString());
            });
        }

    }

    private void OnUpStarBtnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        for (int i = 0; i < userdData.Length; i++)
        {
            if (!GameUtil.Instance.ComCoinToCheck(userdData[i].key, userdData[i].value, () => {
                UITool.Close<UIPartnerUpStar>();
                UITool.Close<UIRole>();
            }))
            {
                return;
            }
        }
        TcpC.C2S_PartnerStar(-1, Uid);

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "角色",
            component = "伙伴升星界面",
            position = "",
            event_click = "伙伴升星",
        });
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo && _e.cid == CmdId.C2S_PartnerStar && _e.code != 0) return;
        UITool.ShowTips( "升星成功");
        OnUIPartnerInfo();
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
    }

    public override void OnDestroy()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
    }
}