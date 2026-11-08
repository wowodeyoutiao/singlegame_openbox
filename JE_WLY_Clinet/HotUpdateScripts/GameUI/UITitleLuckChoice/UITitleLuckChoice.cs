using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;
using System.Security.Cryptography;


/// <summary>
/// 称号幸运符选择
/// </summary>
public class UITitleLuckChoice : UIBase<UIModelTitleLuckChoice, UIViewTitleLuckChoice>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => { 
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UITitleLuckChoice>();
        });

        int idx = 0;
        foreach (int id in CFG_XMLConfig.xmlConfig.WarFight.luckTicketDic.Keys)
        {
            int cid = id;
            var touch = JUIExtensions.exbtn(GetView().root_Awards.GetChild(idx));

            touch.OnClick = () => {

                if (_SelectID == cid)
                    return;

                if (BagM.Instance.GetNumByCfgID(cid) <= 0)
                {
                    UITool.ShowTips("道具不足");
                    return;
                }

                _SelectID = cid;
                touch.transform.Find("on").gameObject.SetActive(true);
                if (_obj_On != null)
                    _obj_On.SetActive(false);
                _obj_On = touch.transform.Find("on").gameObject;

                GetView().txt_tip.gameObject.SetActive(false);
                GetView().btn_Select.gameObject.SetActive(true);
            };

            touch.OnPress = () => { 
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cid });
            };

            GameUtil.Instance.SetItemIcon(cid, JUIExtensions.img(GetView().root_Awards.GetChild(idx).Find("icon")), JUIExtensions.img(GetView().root_Awards.GetChild(idx)));
            idx++;
        }

        GetView().btn_Select.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Get<UIRoleTitile>().OnLuckItemSelect(_SelectID);
            UITool.Close<UITitleLuckChoice>();
        });
    }

    private int _SelectID;
    public GameObject _obj_On;

    public override void OnEnter(params object[] args)
    {
        int idx = 0;
        foreach (int id in CFG_XMLConfig.xmlConfig.WarFight.luckTicketDic.Keys)
        {
            int cid = id;
            JUIExtensions.txt(GetView().root_Awards.GetChild(idx).Find("txt_num")).text = BagM.Instance.GetNumByCfgID(cid).ToString();
            idx++;
        }

        _SelectID = 0;
        if (_obj_On != null)
            _obj_On.SetActive(false);
        _obj_On = null;
        GetView().txt_tip.gameObject.SetActive(true);
        GetView().btn_Select.gameObject.SetActive(false);
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}