using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Helper;
using pb;

public class UIPetShow : UIBase<UIModelPetShow, UIViewPetShow>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPetShow>();
        });

    }

    public override void OnEnter(params object[] args)
    {
        int CfgId = (int)args[0];
        int Level = (int)args[1];
        int Star = (int)args[2];

        bool isHave = args.Length > 3 ? (bool)args[3] : true;

        var petCfgInfo = CFG_Pet.GetById(CfgId);
        GetView().txt_name.text = GameUtil.Instance.GetStrByQuality(petCfgInfo.Quality, petCfgInfo.Name);

        GameUtil.Instance.SetItemIcon(petCfgInfo.ItemID, GetView().Img_icon, GetView().Img_Icon_bg, null, petCfgInfo.Quality);

        var cfgLv = CFG_Role.GetByJobAndLevel(21, Level);
        float addValue = petCfgInfo.Mark == 0 ? 1 :CFG_XMLConfig.xmlConfig.PetQualityDic[petCfgInfo.Quality];
        string Des = "";
        for (int i = 0; i < cfgLv.Property.Length; i++)
        {
            Des += (CFG_XMLConfig.GetAttributeStrByID(cfgLv.Property[i].Id, (int)(cfgLv.Property[i].Num * addValue)) + "\n");
        }
        GetView().txt_attribute.text = Des;

        var Arrt = petCfgInfo.GetStarAttributeByLv(Star);

        Des = "";
        for (int i = 0; i < Arrt.Length; i++)
        {
            if (Arrt[i].value > 0)
                Des += (CFG_XMLConfig.GetAttributeStrByID(Arrt[i].key, Arrt[i].value) + "\n");
        }
        GetView().txt_attribute1.text = Des;

        string EffStr = "无";
        if (petCfgInfo.EffectDes != null)
        {
            EffStr = "";
            for (int i = 0; i < petCfgInfo.EffectDes.Length; i++)
            {
                string str = petCfgInfo.EffectDes[i].key + "阶:";
                string[] strDess = petCfgInfo.EffectDes[i].values.Split('|');
                for (int j = 0; j < strDess.Length; j++)
                {
                    str += $"{strDess[j]}\n";
                }
                EffStr += $"<color=#{((petCfgInfo.EffectDes[i].key <= Star && isHave) ? "76E344" : "A3A3A2")}>{str}</color>\n";
            }
        }

        GetView().txt_Des.text = EffStr; // petCfgInfo.EffectDes != null ? petCfgInfo.EffectDes.Replace("|", "\n") : "无";

        GameHelper.Instance.SetDeayFrame(() => { 
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().rect_Center);
        });

    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}