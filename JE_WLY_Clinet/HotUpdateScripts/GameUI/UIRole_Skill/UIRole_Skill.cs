using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using System.Linq;
using SuperScrollView;
using pb;
using JEngine.Event;

public class UIRole_Skill : UIBase<UIModelRole_Skill, UIViewRole_Skill>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => { 
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRole_Skill>();   
        });

        GetView().btn_Build.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
                return;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv))
                return;

            UITool.Close<UIRole_Skill>();
            UITool.Open<UIBuild>(null,UIBuild.SelectOne.Skill, UIBuild.SelectTwo.Skill_UpLv);

        });

        GetView().infiniteScroll_Skill.InitListView(0,OnFillSkillItem);
    }

    private LoopListViewItem2 OnFillSkillItem(LoopListView2 listView, int Idex)
    {
        if (Idex < 0 || Idex >= SkillShowList.Count)
            return null;

        var ViewItem = listView.NewListViewItem("Item");
        GameObject Item = ViewItem.gameObject;
        var _itemInfo = SkillShowList[Idex];
        JUIExtensions.txt(Item.transform.Find("txt_Name")).text = _itemInfo.Name;
        JUIExtensions.txt(Item.transform.Find("txt_Des")).text = _itemInfo.SkillDes;
        JUIExtensions.txt(Item.transform.Find("img_bg/img_2/txt_lv")).text = $"lv{_itemInfo.SkillLevel}";

        bool isEquip = SkillM.Instance.CheakSkillEquip(_itemInfo.UUID);
        int limetLv = TeamIdx == 0 ?(int) _itemInfo.RoleLevellimit : GetUnLockByParetner(_itemInfo.UUID);
        bool isLock = TeamIdx == 0 ? RoleM.Instance.Level < limetLv : PartnerInfo.Star < limetLv;
        Item.transform.Find("Lock").gameObject.SetActive(isLock && TeamIdx == 0);
        Item.transform.Find("Lock2").gameObject.SetActive(isLock && TeamIdx != 0);

        JUIExtensions.txt(Item.transform.Find("img_bg/img_3/txt_state")).text = isLock ? "未解锁" : (isEquip ?"已上阵" : "已解锁");

        JUIExtensions.txt(Item.transform.Find("img_bg/img_1/txt_type")).text = _itemInfo.ActiveorPassive == 1 ?"主" :"被";
        JUIExtensions.txt(Item.transform.Find("Lock/txt_lock")).text =  $"{limetLv}级解锁" ;
        JUIExtensions.txt(Item.transform.Find("Lock2/txt_lock")).text = limetLv + "星解锁";

        GameUtil.Instance.SetItemIcon((int)_itemInfo.SkillIcon, JUIExtensions.img(Item.transform.Find("img_bg/img_icon")));
        GameUtil.Instance.SetSkillBg((int)_itemInfo.UUID, JUIExtensions.img(Item.transform.Find("img_bg")));

        Item.transform.Find("up").gameObject.SetActive(isEquip);
        Item.transform.Find("unlock").gameObject.SetActive(!isLock && !isEquip);
        return ViewItem;
    }

    private int GetUnLockByParetner(int SkillId)
    {
        for (int i = 0; i < cfgOnePartner.Skills.Length; i++)
        {
            if (SkillId == cfgOnePartner.Skills[i].ElementAt(0).Value)
            {
                return int.Parse(cfgOnePartner.Skills[i].ElementAt(0).Key);
            }
        }
        return 0;
    }

    private List<CFGItem_Skill> SkillShowList = new List<CFGItem_Skill>();
    private CFGItem_Partner cfgOnePartner;
    private pb.PartnerInfo PartnerInfo;
    private void OnSkillPanel()
    {
        if (TeamIdx == 0)
            SkillShowList = SkillM.Instance.GetRolePanelSkillInfo();
        else
        {
            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx);
            PartnerInfo = PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid);
            cfgOnePartner = CFG_Partner.GetById(PartnerInfo.CfgId);
            SkillShowList.Clear();
            if (cfgOnePartner.Skills != null)
            {
                for (int i = 0; i < cfgOnePartner.Skills.Length; i++)
                {
                    var skillCFG = CFG_Skill.GetById(cfgOnePartner.Skills[i].ElementAt(0).Value);
                    for (int j = 0; j < PartnerInfo.Skills.Count; j++)
                    {
                        if (skillCFG.SkillID == CFG_Skill.GetById(PartnerInfo.Skills[j].CfgId).SkillID)
                        {
                            skillCFG = CFG_Skill.GetById(PartnerInfo.Skills[j].CfgId);
                            break;
                        }
                    }
                    SkillShowList.Add(skillCFG);
                }
            }
        }

        if (SkillShowList.Count != GetView().infiniteScroll_Skill.ItemTotalCount)
            GetView().infiniteScroll_Skill.SetListItemCount(SkillShowList.Count);

        GetView().infiniteScroll_Skill.RefreshAllShownItem();

            GetView().txt_title.text = TeamIdx == 0 ? "角色技能" : "伙伴技能";
        GetView().obj_NoSkill.SetActive(SkillShowList.Count == 0);

        //GetView().btn_Build.gameObject.SetActive(SkillShowList.Count != 0);
        GetView().btn_Build.gameObject.SetActive(false);
    }


    private int TeamIdx;
    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.AddEventListener<GameEvent.Refresh_LvUp>(OnRefresh_LvUp);

        TeamIdx = (int)args[0];
        OnSkillPanel();

        RedPointManager.Instance.RPS.Register(RedPointConst.Role_Skill_Build.ToString(), (node) =>
        {
            if (isNullGo) return;
            GetView().obj_BuildRed.SetActive(node.pointCount > 0);
        });

        BuildButtonLock();
    }

    public void BuildButtonLock()
    {
        if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv, false))
        {
            GetView().txt_BuildLock.text = $"{GameUtil.Instance.RemainingDaysUnlockOpen(CFG_XMLConfig.XML_Config.SwitchOpne.Build_SkillLv)}天后开启";
            GetView().txt_BuildLock.transform.parent.gameObject.SetActive(true);
        }
        else
        {
            GetView().txt_BuildLock.transform.parent.gameObject.SetActive(false);
        }
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_LvUp(GameEvent.Refresh_LvUp _e)
    {
        if (isNullGo)
            return;
        BuildButtonLock();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_CheckOpenTime(GameEvent.Refresh_ZeroTime _e)
    {
        if (isNullGo)
            return;
        BuildButtonLock();
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_ZeroTime>(OnSync_CheckOpenTime);
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_LvUp>(OnRefresh_LvUp);

    }

    public override void OnDestroy()
    {

    }
}