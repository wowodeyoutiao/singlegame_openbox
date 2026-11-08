using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.LitJson;
using JEngine.Event;
using JEngine.UI;
using SuperScrollView;

public class UIMailDetail : UIBase<UIModelMailDetail, UIViewMailDetail>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private HTTP_MailInfo _MailInfo;

    public override void OnInit()
    {
        GetView().infiniteScroll_Award.InitListView(0, OnFill) ;
        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIMailDetail>();
        });
        GetView().btn_DelAndGet.onClick.AddListener(OnDelAndGetBtnClick);
    }
    public override void OnEnter(params object[] args)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_EmailOpenAndClose);
        _MailInfo = (HTTP_MailInfo)args[0];
        //Debug.Log("邮件详情： is_extract：" + _MailInfo.is_extract + " / is_attachments: " + _MailInfo.is_attachments + " / attachments count = " + _MailInfo.attachments.Length);
        OnRefreshUI();
    }

    private void OnDelAndGetBtnClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        void OnDel()
        {
            HttpC.OnGetMailDelOne(_MailInfo.mail_id, (_response) =>
            {
                if (_response.code == 0)
                {
                    if (isNullGo) return;
                    MailM.Instance.SetMainDelByID(_MailInfo.mail_id);
                    JEvent.defaultEvent.Post(new GameEvent.Refresh_MailUI());
                    UITool.Close<UIMailDetail>();

                    HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                    {
                        page = "设置",
                        component = "邮件",
                        position = "邮件详情",
                        event_click = "提取并删除",
                    });
                }
                else
                {
                    UITool.ShowTips(_response.data);
                }
            });
        }
        if (_MailInfo.is_attachments == 0 || _MailInfo.is_extract == 1)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_EmailOpenAndClose);
            OnDel();
        }
        else
        {
            HttpC.OnGetMailGetOne(_MailInfo.mail_id,
            (_response) =>
            {
                if (_response.code == 0)
                {
                    if (isNullGo) return;
                    OnDel();
                }
                else
                {
                    UITool.ShowTips(_response.data);
                }
            });
        }
    }

    private LoopListViewItem2 OnFill(LoopListView2 loopListView,int Idx)
    {
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        var Info = _MailInfo.attachments[Idx];
        GameUtil.Instance.SetItemIcon(Info.user_props_id, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform));
        item.transform.Find("geted").gameObject.SetActive(_MailInfo.is_extract == 1);
        JUIExtensions.txt(go.transform.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(Info.num);
        JUIExtensions.btn(go.transform).onClick.RemoveAllListeners();
        JUIExtensions.btn(go.transform).onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = Info.user_props_id });
        });
        return item;
    }

    private void OnRefreshUI()
    {
        GetView().txt_des.text = _MailInfo.desc;
        GetView().txt_Title.text = _MailInfo.title;
        GetView().txt_NoAward.gameObject.SetActive(_MailInfo.is_attachments == 0);
        GetView().txt_BtnGetAndDel.text = (_MailInfo.is_attachments == 0 || _MailInfo.is_extract == 1) ? "删除邮件" : "提取并删除";

        // is_attachments 是否有附件 0 no, is_extract是否提取 0 no
        if (_MailInfo.is_attachments == 0 || _MailInfo.is_extract == 1)
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Set_Email.ToString(), MailM.Instance.GetDotNum());

        if (_MailInfo.attachments == null)
        {
            GetView().infiniteScroll_Award.SetListItemCount(0);
        }
        else
        {
            GetView().infiniteScroll_Award.SetListItemCount(_MailInfo.attachments.Length);
            GetView().infiniteScroll_Award.RefreshAllShownItem();
        }

    }

    public override void OnExit()
    {
        
    }
}