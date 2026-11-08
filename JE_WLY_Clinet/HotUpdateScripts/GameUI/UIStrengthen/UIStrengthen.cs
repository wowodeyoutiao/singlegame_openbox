using JEngine.Helper;
using MetaFramework.UI;
using System.Collections.Generic;
using UnityEngine.UI;

public class UIStrengthen : UIBase<UIModelStrengthen, UIViewStrengthen>
{
    private UIViewStrengthen strengthenView;

    private List<CFG_XMLConfig.XML_Config.SwitchOpne> switchOpnes = new List<CFG_XMLConfig.XML_Config.SwitchOpne>()
    {
        CFG_XMLConfig.XML_Config.SwitchOpne.Build,
        CFG_XMLConfig.XML_Config.SwitchOpne.Refined,
        CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar,
        CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster,
        CFG_XMLConfig.XML_Config.SwitchOpne.CallRole,
        CFG_XMLConfig.XML_Config.SwitchOpne.Pet
    };

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        strengthenView = GetView();

        strengthenView.closeBtn.onClick.AddListener(() =>
        {
            OnClickCloseButton();
        });
    }

    private void LoadButtonState(CFG_XMLConfig.XML_Config.SwitchOpne switchOpne, Image img, Text txt)
    {
        System.Text.StringBuilder tips = new("");
        bool isLockServer = GameUtil.Instance.CheckOpenServerIsLock(switchOpne, false);
        bool islockLv = GameUtil.Instance.CheakIsLock(switchOpne, false, tips);

        if (isLockServer)
            txt.text = $"{GameUtil.Instance.RemainingDaysUnlockOpen(switchOpne)}天后开启";
        else if (islockLv)
            txt.text = tips.ToString();
        else
            txt.text = string.Empty;

        img.material = (isLockServer || islockLv) ? strengthenView.grayMat : null;
    }

    private void OnClickButtonAtIndex(CFG_XMLConfig.XML_Config.SwitchOpne switchOpne)
    {
        if (GameUtil.Instance.CheckOpenServerIsLock(switchOpne))
            return;
        if (GameUtil.Instance.CheakIsLock(switchOpne))
            return;


        switch(switchOpne)
        {
            case CFG_XMLConfig.XML_Config.SwitchOpne.Build:// 装备强化
                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Reinforce);
                });
                break;
            case  CFG_XMLConfig.XML_Config.SwitchOpne.Refined://洗练
                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Refine);
                });
                break;
            case CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipStar:// 升星
                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Star);
                });
                break;
            case CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster:// 镶嵌
                UITool.Open<UIRole>(null, 0, () =>
                {
                    UITool.Open<UIBuild>(null, UIBuild.SelectOne.Equip, UIBuild.SelectTwo.Equip_Inset);
                });
                break;
            case  CFG_XMLConfig.XML_Config.SwitchOpne.CallRole:// 召唤
                UITool.Open<UIPartnerCall>();
                break;
            case  CFG_XMLConfig.XML_Config.SwitchOpne.Pet:// 宠物
                UITool.Open<UIPet>();
                break;
        }

        OnClickCloseButton();
    }

    private void OnClickCloseButton()
    {
        UITool.Close<UIStrengthen>();
    }

    public override void OnEnter(params object[] args)
    {
        for (int i = 0; i < strengthenView.btnTrans.childCount; i++)
        {
            var item = strengthenView.btnTrans.GetChild(i);

            // 按钮
            var switchOpne = switchOpnes[i];
            var btn = item.GetComponent<Button>();
            btn.onClick.RemoveAllListeners();
            btn.onClick.AddListener(() =>
            {
                OnClickButtonAtIndex(switchOpne);
            });

            var img = item.GetComponent<Image>();
            var txt = item.Find("Text").GetComponent<Text>();
            LoadButtonState(switchOpne, img, txt);
        }
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {

    }
}