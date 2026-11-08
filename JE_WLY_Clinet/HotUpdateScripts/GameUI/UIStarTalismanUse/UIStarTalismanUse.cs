using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;

public class UIStarTalismanUse : UIBase<UIModelStarTalismanUse, UIViewStarTalismanUse>
{
    private UIViewStarTalismanUse viewItemDesignatedUse;

    private long selItemId;

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
        UITool.Close<UIStarTalismanUse>();
    }

    private void OnClickOkBtn()
    {
        if (selItemId == 0)
        {
            UITool.ShowTips($"尚未选择!");
            return;
        }

        TcpC.C2S_ItemUse(selItemId, 1);
        OnClickBackBtn();

    }

    public override void OnEnter(params object[] args)
    {
        selItemId = 0;

        for (int i = 0; i < viewItemDesignatedUse.designatItems.Count; i++)
        {
            viewItemDesignatedUse.GetDesignatItem(i).SetActive(false);
        }
        viewItemDesignatedUse.okBtn.gameObject.SetActive(false);
        viewItemDesignatedUse.tipGo.SetActive(true);

        var index = 0;
        for (int i = 25; i < 29; i++)
        {
            var itemInfo = BagM.Instance.GetItemByCFGID(i);

            var designatItem = viewItemDesignatedUse.GetDesignatItem(index);
            var itemCfg = CFG_Item.GetById(i);
            GameUtil.Instance.SetItemIcon(i, designatItem.iconImg, designatItem.bgImg);
            designatItem.nameText.text = GameUtil.Instance.GetFightStrByQuality((int)itemCfg.Quality, itemCfg.Name);
            designatItem.nameOutline.effectColor = GameUtil.Instance.GetQualityOutlineColor((int)itemCfg.Quality);
            designatItem.numText.text = "" + (itemInfo == null || itemInfo.Num == 0 ? 0 : itemInfo.Num);
            designatItem.selGo.SetActive(false);

            designatItem.itemBtn.onClick.RemoveAllListeners();
            designatItem.itemBtn.onClick.AddListener(() =>
            {
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = itemCfg.ID });
            });

            designatItem.SetActive(true);

            var clickIndex = index;
            designatItem.btn.onClick.RemoveAllListeners();
            designatItem.btn.onClick.AddListener(() =>
            {
                if (itemInfo == null || itemInfo.Num == 0)
                {
                    UITool.ShowTips($"所选择道具数量不足，请前往商城购买！");
                    viewItemDesignatedUse.okBtn.gameObject.SetActive(selItemId > 0);
                    viewItemDesignatedUse.tipGo.SetActive(selItemId == 0);
                }
                else
                {
                    viewItemDesignatedUse.okBtn.gameObject.SetActive(true);
                    viewItemDesignatedUse.tipGo.SetActive(false);
                    selItemId = itemInfo.Uid;
                    OnClickItem(clickIndex);
                }
            });

            index++;
        }

        if(index == 0)
        {
            UITool.ShowTips($"暂无升星符可用!");
            OnClickBackBtn();
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