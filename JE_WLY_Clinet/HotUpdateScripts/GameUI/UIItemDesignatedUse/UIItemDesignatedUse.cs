using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.MyNet;
using MetaFramework.UI;
using pb;
using System;
using UnityEngine;
using static UIViewItemDesignatedUse;

public class UIItemDesignatedUse : UIBase<UIModelItemDesignatedUse, UIViewItemDesignatedUse>
{
    private UIViewItemDesignatedUse viewItemDesignatedUse;

    private long itemUid;
    private int designatedItemCfgId;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        viewItemDesignatedUse = GetView();

        viewItemDesignatedUse.backBtn.onClick.AddListener(OnClickBackBtn);
        viewItemDesignatedUse.okBtn.onClick.AddListener(OnClickOkBtn);
    }
    
    private void OnClickBackBtn()
    {
        UITool.Close<UIItemDesignatedUse>();
    }

    private void OnClickOkBtn()
    {
        if(designatedItemCfgId == 0)
        {
            UITool.ShowTips($"尚未选择!");
            return;
        }

        UITool.ShowTips($"使用成功!");
        TcpC.C2S_ItemUse(itemUid, 1, designatedItemCfgId);
        OnClickBackBtn();
    }

    public override void OnEnter(params object[] args)
    {
        designatedItemCfgId = 0;
        var itemCfgId = (int)args[0];
        itemUid = (long)args[1];
        SetData(itemCfgId);
    }

    private void SetData(int cfgId)
    {
        var itemUseCfg = CFG_ItemUse.GetById(cfgId);
        if (itemUseCfg == null || itemUseCfg.Givegoods == null || itemUseCfg.Givegoods.Length == 0)
        {
            OnClickBackBtn();
            return;
        }

        viewItemDesignatedUse.titleText.text = itemUseCfg.Name;

        var count = Mathf.Max(viewItemDesignatedUse.designatItems.Count, itemUseCfg.Givegoods.Length);
        for (int i = 0; i < count; i++)
        {
            DesignatItem designatItem = viewItemDesignatedUse.GetDesignatItem(i);

            designatItem.SetActive(i < itemUseCfg.Givegoods.Length);
            if (i < itemUseCfg.Givegoods.Length)
            {
                var goods = itemUseCfg.Givegoods[i];
                var itemCfg = CFG_Item.GetById(goods.id);
                GameUtil.Instance.SetItemIcon(goods.id, designatItem.iconImg, designatItem.bgImg);
                designatItem.nameText.text = GameUtil.Instance.GetFightStrByQuality((int)itemCfg.Quality, itemCfg.Name);
                designatItem.nameOutline.effectColor = GameUtil.Instance.GetQualityOutlineColor((int)itemCfg.Quality);
                designatItem.selGo.SetActive(false);

                designatItem.itemBtn.onClick.RemoveAllListeners();
                designatItem.itemBtn.onClick.AddListener(() =>
                {
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = goods.id });
                });

                var index = i;
                designatItem.btn.onClick.RemoveAllListeners();
                designatItem.btn.onClick.AddListener(() =>
                {
                    designatedItemCfgId = goods.id;
                    OnClickItem(index);
                });
            }
        }
    }

    private void OnClickItem(int index)
    {
        for (int i = 0; i < viewItemDesignatedUse.designatItems.Count; i++)
        {
            viewItemDesignatedUse.designatItems[i].selGo.SetActive(i == index);
        }
    }

    public override void OnExit()
    {

    }
}