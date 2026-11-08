using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.LitJson;

using MetaFramework.UI;

using SuperScrollView;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

using static CFG_Room;

public class UIMail : UIBase<UIModelMail, UIViewMail>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {

        GetView().btn_Close.onClick.AddListener(CloseUI);
        GetView().btn_refresh.onClick.AddListener(OnBtnRefreshClick);
        GetView().btn_DelAll.onClick.AddListener(OnDelAllClick);
        GetView().btn_AllGet.onClick.AddListener(OnBtnAllGetClick);


        GetView().infiniteScroll.InitListView(0, OnFillItem) ;

        GetView().txt_MailTip.text = "邮件最多存储" + CFG_XMLConfig.xmlConfig.BaseValue.mailMaxNum + "封，超出后将无法收取，请及时清理！";
    }

    public void CloseUI()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        UITool.Close<UIMail>();
    }

    private void OnBtnRefreshClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        OnGetMail();

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "首页",
            component = "设置",
            position = "邮件",
            event_click = "刷新邮箱",
        });
    }

    private void OnBtnAllGetClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (_MailLsit.Count <= 0)
        {
            UITool.ShowTips( "暂无邮件");
            return;
        }
        HttpC.OnGetMailGetAll((_response) =>
        {
            if (_response.code == 0)
            {
                if (!isNullGo)
                    OnGetMail();

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "设置",
                    position = "邮件",
                    event_click = "全部领取并已读",
                });
            }
        });
       // OnGetMail();
    }

    //重新请求数据
    private void OnGetMail()
    {
        HttpC.OnGetMailInfo();
    }
    private void OnDelAllClick()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (_MailLsit.Count <= 0)
        {
            UITool.ShowTips("暂无邮件");
            return;
        }

        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack()
        {
            ShowStr = "确定要删除全部邮件吗？\n\n<color=#91F200>（附件未提取的邮件不会删除）</color>",
            YesFun = () => {
                HttpC.OnGetMailDelAll((_response) =>
                {
                    if (_response.code == 0 && !isNullGo)
                    {
                        OnGetMail();
                    }
                    else
                    {
                        UITool.ShowTips($"失败 code {_response.code}");
                    }
                });

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "设置",
                    position = "邮件",
                    event_click = "一键删除",
                });
            },
        });
    }
    private List<HTTP_MailInfo> _MailLsit = new List<HTTP_MailInfo>();
    private void SetMailRefrsh(bool isRefresh = false)
    {
        _MailLsit = MailM.Instance.GetMailList();
        int Count = _MailLsit.Count;
        GetView().tip_NoMail.SetActive(Count == 0);
        if (Count != 0)
        {
            if (!isRefresh)
            {
                GetView().infiniteScroll.SetListItemCount(Count);
            }
            GetView().infiniteScroll.RefreshAllShownItem();
        }
        else
        {
            GetView().infiniteScroll.SetListItemCount(0);
        }
    }

    //刷新邮件列表
    private LoopListViewItem2 OnFillItem(LoopListView2 listView, int idx)
    {
        if (idx < 0 || idx >= _MailLsit.Count)
        {
            return null;
        }
        LoopListViewItem2 item = listView.NewListViewItem("Item");
        var Info = _MailLsit[idx];
        GameObject go = item.gameObject;
        go.transform.Find("txt_Title").GetComponent<Text>().text = Info.title;
        go.transform.Find("txt_Des").GetComponent<Text>().text = Info.desc;
        go.transform.Find("txt_Gettime").GetComponent<Text>().text = GameUtil.Instance.GetTimeToYMD(Info.created_at);
        go.transform.Find("txt_EndTime").GetComponent<Text>().text = "剩余时间:"+ GameUtil.Instance.GetTimeToDHM(Info.expiry_at - GameUtil.Instance.GetCurrentTimeSeconds());
        go.transform.Find("tip_Add").gameObject.SetActive(Info.is_attachments != 0 && Info.is_extract == 0);
        go.transform.Find("tip_UnRead").gameObject.SetActive(Info.is_read == 0);
        go.transform.Find("tip_Read").gameObject.SetActive(Info.is_read != 0);
        go.GetComponent<Button>().onClick.RemoveAllListeners();
        go.GetComponent<Button>().onClick.AddListener(() => { OnMailItemClick(Info); });
        return item;
    }

    private void OnMailItemClick(HTTP_MailInfo info)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        HttpC.OnGetOneMail(info.mail_id, (_response) =>
        {
            if (_response.code == 0 && !isNullGo)
            {
                UITool.Open<UIMailDetail>(null, _response.data);
                MailM.Instance.SetOneMailInfo(_response.data);
                SetMailRefrsh(true);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "设置",
                    position = "邮件",
                    event_click = "点击单个邮件，查看详情",
                });
            }
            else
            {
                UITool.ShowTips($"失败 code {_response.code}");
            }
        });
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);

        _MailLsit = MailM.Instance.GetMailList();
        SetMailRefrsh();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);

    }
    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_Refresh(GameEvent.Refresh_MailUI _e)
    {
        SetMailRefrsh();
    }
}