using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using Spine;
using Spine.Unity;

using SuperScrollView;

using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;

using UnityEngine;
using UnityEngine.UI;

public class UIStarDetails : UIBase<UIModelStarDetails, UIViewStarDetails>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().looplist_Items.InitListView(0, OnItemRefresh);
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIStarDetails>();
        });
    }

    CFGItem_Partner cfgPartner;

    public override void OnEnter(params object[] args)
    {
        cfgPartner = (CFGItem_Partner)args[0];
        RefreshUI();
    }

    private LoopListViewItem2 OnItemRefresh(LoopListView2 list, int index)
    {
        var item = list.NewListViewItem("Item");
        Transform t = item.transform;
        JUIExtensions.txt(t.Find("txt_num")).text = cfgPartner.Skills[index].ElementAt(0).Key;
        CFGItem_Skill cfgSkill = CFG_Skill.GetById(cfgPartner.Skills[index].ElementAt(0).Value);
        CFGItem_Item itemCfg = CFG_Item.GetById((int)cfgSkill.SkillIcon);
        JUIExtensions.txt(t.Find("txt_desc")).text = itemCfg.Desc;
        return item;
    }

    private void RefreshUI()
    {
        GetView().txt_Name.text = GameUtil.Instance.GetItemName(cfgPartner.ItemID, cfgPartner.Quality);
        JUIExtensions.outline(GetView().txt_Name.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(cfgPartner.Quality);
        string poolKey = "Partner_" + cfgPartner.ID;
        SetRoleSpine(poolKey, poolKey == "" ? "" : GameConst.GetMonsterModelByID(poolKey));

        bool isNoSkill = cfgPartner.Skills == null || cfgPartner.Skills.Length == 0;
        GetView().txt_NoSkill.gameObject.SetActive(isNoSkill);
        GetView().looplist_Items.gameObject.SetActive(!isNoSkill);
        if (!isNoSkill)
        {
            GetView().looplist_Items.SetListItemCount(cfgPartner.Skills.Length);
            GetView().looplist_Items.RefreshAllShownItem();
            GetView().looplist_Items.MovePanelToItemIndex(0, 0);
        }
    }

    private void SetRoleSpine(string poolKey, string path)
    {

        if (GetView().root_Role.childCount > 0)
        {
            PoolHelper.Instance.RecycObj(GetView().root_Role.GetChild(0).gameObject);
        }

        if (!string.IsNullOrEmpty(poolKey))
        {
            _ = GameHelper.Instance.SetAsyncObj(poolKey, path, (_mainRole) =>
            {
                _mainRole.transform.SetParent(GetView().root_Role);
                _mainRole.transform.localScale = Vector3.one;
                _mainRole.transform.localPosition = Vector3.zero;
                _mainRole.SetActive(true);

                SpineUtil.Instance.PlayQueue(_mainRole.transform.Find("spine").GetComponent<SkeletonGraphic>(), new List<string> { "stand" }, null, true);
            });
        }
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}