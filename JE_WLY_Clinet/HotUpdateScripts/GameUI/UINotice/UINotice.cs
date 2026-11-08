using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;

public class UINotice : UIBase<UIModelNotice, UIViewNotice>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UINotice>();
        });
    }

    string noticeStr = "<b><size=24>《这岛有剧本》公告</size></b>\r\n\r\n亲爱的岛主们：​\r\n     《这岛有剧本》为了感谢大家对于本次公测的支持，公测期间充值的玩家，正式服可获双倍返还，规则如下： \r\n充值返还：\r\n1、返利仅限正式服单个游戏ID绑定，发放至该游戏ID； \r\n2、返还形式为游戏内通用货币（元宝，具体以游戏内为准）。\r\n此次公测期间凡充值玩家均可添加客服微信，微信号：zdyjb2025，期待在正式服与大家再次相遇，一同开启全新的游戏旅程！";
    string TiShenStr = "<b><size=24>《这岛有剧本》公告</size></b>\r\n\r\n亲爱的岛主们：​\r\n     感谢大家对于本次公测的支持";
    public override void OnEnter(params object[] args)
    {
        //GameHelper.Instance.SetDeayTimeDo(() => { 
        //    GetView().txt_Show.text = GameM.Instance.IsAudit ? TiShenStr: noticeStr;
        //},0.2f);

        GetView().txt_Show.text = GameM.Instance.IsAudit ? TiShenStr : noticeStr;

        GameHelper.Instance.SetDeayFrame(() => {
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_Show.GetComponent<RectTransform>());
            LayoutRebuilder.ForceRebuildLayoutImmediate(GetView().txt_Show.transform.parent.GetComponent<RectTransform>());
        });
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}