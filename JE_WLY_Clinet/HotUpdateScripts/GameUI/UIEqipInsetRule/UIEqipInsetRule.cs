using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using pb;
using HotUpdateScripts.ProtoM;
using JEngine.UI;

public class UIEqipInsetRule : UIBase<UIModelEqipInsetRule, UIViewEqipInsetRule>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public enum UISelect
    {
        /// <summary>
        /// 魂石大师
        /// </summary>
        InsetMonster = 3,

        /// <summary>
        /// 自然大师
        /// </summary>
        NatureManster = 4,

    }

    private List<GameObject> ObjItemList = new List<GameObject>();

    private Dictionary<int, string> NatureStrDic = new Dictionary<int, string>
    {
        [1] = "风",
        [2] = "雨",
        [3] = "雷",
        [4] = "电",
    }; 

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIEqipInsetRule>();
        });
    }

    private void OnShowScrollView(int select)
    {
        var Master = CFG_EquipMaster.GetMastersByCatedory((CFG_EquipMaster.CategoryClass)select);

        string strTip = uISelect == UISelect.InsetMonster ? "全槽位镶嵌{0}颗{1}级以上魂石:({2}/{3})" : "4颗{0}系魂石可激活:({1}/{2})";

        for (int i = 0; i < Master.Count; i++)
        {
            GameObject go;
            if (ObjItemList.Count > i)
            {
                go = ObjItemList[i];
            }
            else
            {
                go = GameObject.Instantiate(GetView().obj_Item, GetView().root_Content);
                ObjItemList.Add(go);
            }
            go.SetActive(true);

            bool isCan = false;
            string titleStr = "";
            if (uISelect == UISelect.InsetMonster)
            {
                int count = GetJewelNum(Master[i].Limit);
                isCan = count >= CFG_XMLConfig.xmlConfig.Diamond_MasterDic[Master[i].UUID];
                titleStr = string.Format(strTip, CFG_XMLConfig.xmlConfig.Diamond_MasterDic[Master[i].UUID], Master[i].Limit, count, CFG_XMLConfig.xmlConfig.Diamond_MasterDic[Master[i].UUID]);
            }
            else
            {
                int count = GetJewelNatureNum(Master[i].Limit);
                isCan = count >= 4;
                titleStr = string.Format(strTip, NatureStrDic[Master[i].Limit] , count, 4);
            }
            JUIExtensions.txt(go.transform, "txt_title").text = isCan ? "<color=#76E344>·" + titleStr + "</color>" : titleStr;

            string arrStr = "";
            for (int j = 0; j < Master[i].Attribute.Length; j++)
            {
                arrStr += CFG_XMLConfig.GetAttributeStrByID(Master[i].Attribute[j].key, Master[i].Attribute[j].value) ;
                if (j < Master[i].Attribute.Length)
                {
                    arrStr +=  "\n";
                }
            }
            JUIExtensions.txt(go.transform, "txt_Des").text = isCan ? "<color=#76E344>" + arrStr + "</color>" : arrStr;
            go.transform.Find("txt_Des").GetComponent<ContentSizeFitter>().SetLayoutVertical();
        }

        if (Master.Count < ObjItemList.Count)
        {
            for (int i = Master.Count; i < ObjItemList.Count; i++)
            {
                ObjItemList[i].SetActive(false);
            }
        }

        GameHelper.Instance.SetDeayFrame(() =>
        {
            GetView().root_Content.GetComponent<ContentSizeFitter>().SetLayoutVertical();
        });
    }

    private int GetJewelNum(int Lv)
    {
        int Count = 0;

        if (TeamdIndex == 0)
        {
            for (int i = 1; i < (int)EquipSlot.MaxEquipSlot; i++)
            {
                var _SoltInfo = RoleM.Instance.GetSlotInfoBySolt((EquipSlot)i);
                for (int j = 0; j < _SoltInfo.Gemstone.Length; j++)
                {
                    if (_SoltInfo.Gemstone[j] != 0 && CFG_Item.GetById(_SoltInfo.Gemstone[j]).itemlevel >= Lv)
                    {
                        Count++;
                    }
                }
            }
        }
        else
        {
            var equips = PartnerM.Instance.GetOnePartnerSlotInfo(TeamdIndex).Equips;
            for (int i = 0; i < equips.Count; i++)
            {
                for (int j = 0; j < equips[i].Gemstone.Length; j++)
                {
                    if (equips[i].Gemstone[j] != 0 && CFG_Item.GetById(equips[i].Gemstone[j]).itemlevel >= Lv)
                    {
                        Count++;
                    }
                }
            }
        }
        return Count;
    }

    private int GetJewelNatureNum(int Nature)
    {
        int Count = 0;
        var _SoltInfo = TeamdIndex == 0 ? RoleM.Instance.GetSlotInfoBySolt((EquipSlot)EquipSolt) : PartnerM.Instance.GetOnePartnerOneEquipSolt(TeamdIndex, EquipSolt);
        for (int j = 0; j < _SoltInfo.Gemstone.Length; j++)
        {
            if (_SoltInfo.Gemstone[j] != 0 && CFG_Item.GetById(_SoltInfo.Gemstone[j]).NatureMark == Nature)
            {
                Count++;
            }
        }
        return Count;
    }

    private int EquipSolt;
    private UISelect uISelect;
    private int TeamdIndex;
    public override void OnEnter(params object[] args)
    {
        uISelect = (UISelect)args[0];
        EquipSolt = (int)args[1];
        TeamdIndex = (int)args[2];
        GetView().txt_Title.text = uISelect == UISelect.InsetMonster ? "魂石大师" : "自然大师";
        OnShowScrollView((int)uISelect);
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}