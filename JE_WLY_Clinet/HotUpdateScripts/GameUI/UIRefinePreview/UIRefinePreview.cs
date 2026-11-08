using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using System.Collections;
using System.Collections.Generic;
using System.Reflection;

using UnityEngine;
using UnityEngine.UI;

public class UIRefinePreview : UIBase<UIModelRefinePreview, UIViewRefinePreview>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    //攻击类部位可洗炼的属性种类
    int[] attackType = { 1, 2, 7, 8, 13, 14, 24, 25, 26, 36, 37, 38, 39, 53, 55 };
    //防御类部位可洗炼的属性种类
    int[] defensiveType = { 4, 5, 9, 12, 15, 16, 17, 30, 31, 32, 33, 34, 40, 54, 56 };

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIRefinePreview>();
        });

        for (int i = 0; i < GetView().buttons.Length; i++)
        {
            int index = i;
            GetView().buttons[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnClickAttribute(index);
            });
        }
    }

    EquipSlot equipSlot;

    public override void OnEnter(params object[] args)
    {
        equipSlot = (EquipSlot)args[0];
        SetUI();
    }

    public bool IsAttackEquipSlot()
    {
        switch (equipSlot)
        {
            case EquipSlot.Weapon:
            case EquipSlot.Necklace:
            case EquipSlot.LeftRing:
            case EquipSlot.RightRing:
            case EquipSlot.LeftBracelet:
            case EquipSlot.RightBracelet:
                return true;
            case EquipSlot.Clothes:
            case EquipSlot.Helmet:
            case EquipSlot.Belt:
            case EquipSlot.Boots:
            case EquipSlot.Greaves:
            case EquipSlot.Pendant:
                return false;
            default:
                return false;
        }
    }

    void SetUI()
    {
        GetView().txt_Title.text = GameM.Instance.EquipName[(int)equipSlot] + "可洗练属性预览";
        int count = IsAttackEquipSlot() ? attackType.Length : defensiveType.Length;
        for (int i = 0; i < GetView().buttons.Length; i++)
        {
            Button btn = GetView().buttons[i];
            btn.gameObject.SetActive(i < count);
            if (i < count)
            {
                int attrType = IsAttackEquipSlot() ? attackType[i] : defensiveType[i];
                string value = CFG_RandomDistribV2.GetAttributeName(attrType);
                JUIExtensions.txt(btn.transform, "txt_select").text = value;
                JUIExtensions.txt(btn.transform, "txt_unselect").text = value;
            }
        }

        OnClickAttribute(0);
    }

    private void SelectButton(int index)
    {
        for (int i = 0; i < GetView().buttons.Length; i++)
        {
            Button btn = GetView().buttons[i];
            bool selected = i == index;
            btn.transform.Find("select").gameObject.SetActive(selected);
            btn.transform.Find("txt_select").gameObject.SetActive(selected);
            btn.transform.Find("txt_unselect").gameObject.SetActive(!selected);
        }
    }

    void OnClickAttribute(int index)
    {
        int attrType = IsAttackEquipSlot() ? attackType[index] : defensiveType[index];
        var attrList = CFG_RandomDistribV2.DicRandom[attrType];
        int maxIndex = attrList.Count - 1;
        //string min = GameUtil.Instance.GetStrByQuality(attrList[0].QualityColor, CFG_XMLConfig.GetAttributeStrByID(attrList[0].AttributeType, attrList[0].ValueMin));
        string min = CFG_XMLConfig.GetAttributeStrByID(attrList[0].AttributeType, attrList[0].ValueMin);
        min = min.Insert(min.IndexOf(":") + 1, "\n");
        GetView().txt_MinValue.text = min;
        string max = CFG_XMLConfig.GetAttributeStrByID(attrList[maxIndex].AttributeType, attrList[maxIndex].ValueMax);
        max = max.Insert(max.IndexOf(":") + 1, "\n");
        GetView().txt_MaxValue.text = max;

        //JUIExtensions.outline(GetView().txt_MinValue.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(attrList[0].QualityColor);
        //JUIExtensions.outline(GetView().txt_MaxValue.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor(attrList[maxIndex].QualityColor);
        SelectButton(index);
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {

    }
}