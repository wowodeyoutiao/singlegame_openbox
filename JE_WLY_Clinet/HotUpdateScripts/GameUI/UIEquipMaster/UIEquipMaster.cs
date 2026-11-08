using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.UI;

public class UIEquipMaster : UIBase<UIModelEquipMaster, UIViewEquipMaster>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIEquipMaster>();
        });
    }

    private List<GameObject> ObjItemList = new List<GameObject>();

    private void SetDesUI(CFG_EquipMaster.CategoryClass cClass,int TeamIdx)
    {
        GetView().txt_Title.text = cClass == CFG_EquipMaster.CategoryClass.Reinforce ? "强化大师": "升星大师";
        int lv = 0;

        if (cClass == CFG_EquipMaster.CategoryClass.Reinforce)
        {
            if (TeamIdx == 0)
                lv = RoleM.Instance.GetEquipReinforce();
            else
            {
                lv = 9999;
                var equips = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).Equips;
                for (int i = 1; i < equips.Count; i++)
                {
                    lv = Mathf.Min(lv, equips[i].Enhance);
                }
            }
        }
        else
        {
            if (TeamIdx == 0)
                lv = RoleM.Instance.GetEquipStar();
            else
            {
                lv = 9999;
                var equips = PartnerM.Instance.GetOnePartnerSlotInfo(TeamIdx).Equips;
                for (int i = 1; i < equips.Count; i++)
                {
                    lv = Mathf.Min(lv, equips[i].Stars);
                }
            }
        }

        string strTip = cClass == CFG_EquipMaster.CategoryClass.Reinforce ? "全槽位强化至{0}级" : "全槽位升至{0}星";

        var Master = CFG_EquipMaster.GetMastersByCatedory(cClass);
        int Idx = 0;
        foreach (var item in Master)
        {
            GameObject go;
            if (Idx < ObjItemList.Count)
            {
                go = ObjItemList[Idx];
            }
            else
            {
                go = GameObject.Instantiate(GetView().Item.gameObject, GetView().root_Items);
                ObjItemList.Add(go);
            }
            go.SetActive(true);
            Idx++;

            JUIExtensions.txt(go.transform, "txt_title").text = $"<color=#{(item.Limit <= lv ? "76E344" : "D6D5D4")}>·{string.Format(strTip, item.Limit)}</color>";
           
            string str = "";
            int index = 1;
            foreach (var attr in item.Attribute)
            {
                str += CFG_XMLConfig.GetAttributeStrByID(attr.key, attr.value) ;
                if (index < item.Attribute.Length)
                {
                    str += "\n";
                }
                index++;
            }
            JUIExtensions.txt(go.transform, "txt_Des").text = $"<color=#{(item.Limit <= lv ? "76E344" : "A3A3A2")}>{str}</color>";
            go.transform.Find("txt_Des").GetComponent<ContentSizeFitter>().SetLayoutVertical();
        }
       
        if (Master.Count < ObjItemList.Count)
        {
            for (int i = Master.Count; i < ObjItemList.Count; i++)
            {
                ObjItemList[i].SetActive(false);
            }
        }

        //new JEngine.Core.JAction().DelayFrame(3).Do(() =>
        //{
            GetView().root_Items.GetComponent<ContentSizeFitter>().SetLayoutVertical();
        //}).Execute(true);
    }

    public override void OnEnter(params object[] args)
    {
        SetDesUI((CFG_EquipMaster.CategoryClass)args[0], (int)args[1]);
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}