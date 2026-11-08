using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using SuperScrollView;
using System.Collections.Generic;
using System.Text.RegularExpressions;
using UnityEngine;
using UnityEngine.UI;
using static CFG_XMLConfig.XML_Config;
using static HotUpdateScripts.ProtoM.ChatM;

public class UIChat : UIBase<UIModelChat, UIViewChat>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private int consumeId = CFG_XMLConfig.xmlConfig.BaseValue.chatConsume.itemid;
    private int consumeNum = CFG_XMLConfig.xmlConfig.BaseValue.chatConsume.num;

    private ChatShowType chatShowType = ChatShowType.None;
    List<ChatInfo> chatInfoList = new List<ChatInfo>();

    private int emojiCount = 66;

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIChat>();
        });

        GetView().btn_SendEmoji.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_emojiUI.gameObject.SetActive(true);
        });

        GetView().btn_emojiUI.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_emojiUI.gameObject.SetActive(false);
        });

        GetView().btn_Set.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIChatSet>();
        });

        GetView().btn_Send.onClick.AddListener(() =>
        {
            OnClickSendButton();
        });

        GetView().btn_SendPos.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            SendChat($"{GameRoomM.Instance.CFGRoom.Name}[{MainRole._MainRole.pointNow.X},{MainRole._MainRole.pointNow.Y}]<color=#00000000>{GameRoomM.Instance.CFGRoom.MapId}|{MainRole._MainRole.pointNow.X}|{MainRole._MainRole.pointNow.Y}</color>",false);
        });

        GetView().loopListView.InitListView(0, OnChatItemRefresh);
        GetView().loopGrid_Emoji.InitGridView(emojiCount, OnEmojiItemRefresh);

        InitTagButtons();
    }

    private void OnClickSendButton()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

        bool isLockDayLimt = GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.ChatSend, true);
        if (isLockDayLimt) return;

        var lockLv = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[SwitchOpne.ChatSend];
        if (RoleM.Instance.Level < lockLv)
        {
            UITool.ShowTips($"需要主角达到{lockLv}级才能发送聊天信息！");
            return;
        }

        if (!string.IsNullOrEmpty(GetView().input_chat.text))
        {
            SendChat(GetView().input_chat.text);
        }
    }

    private LoopGridViewItem OnEmojiItemRefresh(LoopGridView view, int index, int row, int column)
    {
        LoopGridViewItem item = view.NewListViewItem("Item");
        string SendStr = index == 65 ? "emoji_229" : "emoji_" + index;
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ChatAtlas, SendStr, JUIExtensions.img(item.transform));
        JUIExtensions.btn(item.transform).onClick.RemoveAllListeners();
        JUIExtensions.btn(item.transform).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            GetView().btn_emojiUI.gameObject.SetActive(false);
            SendChat(SendStr);
        });
        return item;
    }

    private LoopListViewItem2 OnChatItemRefresh(LoopListView2 list, int index)
    {
        LoopListViewItem2 item = list.NewListViewItem("Item");
        var oneChat = chatInfoList[index];
        bool isAll = chatShowType == ChatShowType.All;
        if (isAll)
            JUIExtensions.txt(item.transform).text = "            " + ChatM.Instance.GetChatByType(oneChat.ChatType, oneChat.Context, oneChat.Name);
        else
            JUIExtensions.txt(item.transform).text = ChatM.Instance.GetChatByType(oneChat.ChatType, oneChat.Context, oneChat.Name);

        bool isEmjoy = ChatM.Instance.CheckIsEmoji(oneChat.Context);
        item.transform.Find("emjoy").gameObject.SetActive(isEmjoy);
        RectTransform rt = item.GetComponent<RectTransform>();
        bool isZuoBiao = oneChat.Context.Contains("<color=#00000000>");
        item.transform.Find("dingwei").gameObject.SetActive(isZuoBiao);
        if (isEmjoy || isZuoBiao)
        {
            rt.SetSizeWithCurrentAnchors(RectTransform.Axis.Vertical, 65);
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.ChatAtlas, oneChat.Context, JUIExtensions.img(item.transform.Find("emjoy")));
        }
        else
        {
            rt.SetSizeWithCurrentAnchors(RectTransform.Axis.Vertical, 50);
        }

        Transform imgTip = item.transform.Find("img_Tip");
        imgTip.gameObject.SetActive(isAll);
        if (isAll)
        {
            (string atlas, string icon) = ChatM.Instance.GetChatTypeImg(oneChat.ChatType);
            ResouseHelper.Instance.LoadSpriteByAtlas(atlas, icon, JUIExtensions.img(item.transform.Find("img_Tip")));
        }

        var _btn = JUIExtensions.btn(item.transform);
        _btn.onClick.RemoveAllListeners();
        _btn.onClick.AddListener(() => {
            if (isZuoBiao)
            {

                string result = Regex.Match(oneChat.Context, @"<color=#00000000>([^<]+)</color>")
                                     .Groups[1].Value;
                var strs = result.Split('|');
                if (GameRoomM.Instance._GameState == GameRoomM.GameState.No || GameRoomM.Instance.CFGRoom.MapId != int.Parse(strs[0]))
                {
                    UITool.ShowTips("位置不可达！");
                    return;
                }
                MainRole._MainRole.targetCollectID = 0;
                MainRole._MainRole.targetID = 0;
                MainRole._MainRole._SendMovePos = new Point { X = int.Parse(strs[1]), Y = int.Parse(strs[2]) };
                UITool.Close<UIChat>();
            }
        });

        return item;
    }

    private void InitTagButtons()
    {

        for (int i = 0; i < GetView().ButtonTags.Length; i++)
        {
            int idx = i;
            string tagName = GetChatTypeName((ChatShowType)idx + 1);
            GetView().ButtonTags[i].transform.Find("Label").GetComponent<Text>().text = tagName;
            GetView().ButtonTags[i].transform.Find("Select/Label").GetComponent<Text>().text = tagName;
            GetView().ButtonTags[i].onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnShopToggleClick(idx);
            });

            if (i == 1) GetView().ButtonTags[i].gameObject.SetActive(false);

        }
    }

    private void OnShopToggleClick(int index = 0)
    {
        for (int i = 0; i < GetView().ButtonTags.Length; i++)
        {
            GetView().ButtonTags[i].transform.Find("Select").gameObject.SetActive(i == index);
        }
        if (index + 1 == (int)chatShowType)
            return;

        chatShowType = (ChatShowType)index + 1;
        OnChatUIShow();

        if (chatInfoList.Count > 0)
            GetView().loopListView.MovePanelToItemIndex(GetView().loopListView.ItemTotalCount - 1, 0);

        SetConstShow();
    }

    private string GetChatTypeName(ChatShowType chatShowType)
    {
        switch (chatShowType)
        {
            case ChatShowType.World:
                return "世界";
            case ChatShowType.All:
                return "综合";
            case ChatShowType.GUILD:
                return "行会";
            case ChatShowType.Zone:
                return "战区";
            default:
                return "";
        }
    }

    private void SetConstShow()
    {
        bool showConsume = chatShowType == ChatShowType.World;
        GetView().img_chatItemIcon.gameObject.SetActive(showConsume);
        if (showConsume)
        {
            GameUtil.Instance.SetItemIcon(consumeId, GetView().img_chatItemIcon);
            GetView().txt_chatItemCount.text = GameUtil.Instance.GetCheakNumStr(consumeId, consumeNum);
        }
    }

    private void OnChatUIShow()
    {
        chatInfoList = ChatM.Instance.GetChatsByType(chatShowType);

        GetView().btn_SendPos.gameObject.SetActive(chatShowType == ChatShowType.GUILD && GameRoomM.Instance._GameState != GameRoomM.GameState.No);
        GetView().loopListView.SetListItemCount(chatInfoList.Count);
        GetView().loopListView.RefreshAllShownItem();
    }

    private void SendChat(string chatStr,bool isLimit = true)
    {
        if (chatShowType == ChatShowType.GUILD && !RoleM.Instance.IsJoinedGuild())
        {
            UITool.ShowTips("请加入行会后再进行发言!");
            return;
        }

        if (chatStr.Length > 20 && isLimit)
        {
            UITool.ShowTips("最多可输入20个中文字符");
            return;
        }

        if (chatShowType == ChatShowType.World || chatShowType == ChatShowType.All)
        {
            if (!GameUtil.Instance.ComCoinToCheck(consumeId, consumeNum, null))
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = "金币不足，是否前往商店购买金币？",
                    YesFun = () =>
                    {
                        UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                    },
                });

                return;
            }
        }

        ChatType sendChatType = ChatType.WORLD;
        if (chatShowType == ChatShowType.GUILD)
        {
            sendChatType = ChatType.GUILD;
        }
        else if (chatShowType == ChatShowType.Zone)
        {
            sendChatType = ChatType.ZONE;
        }
        HttpC.SendChatMessage(sendChatType, chatStr);

        //TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Chat, MType.Request, TCPClientMgr.ProtoSerialize(new ChatRequest { Type = sendChatType, Context = chatStr }));


    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        OnShopToggleClick(args.Length > 0 ? (int)args[0] : 0);
        SetConstShow();
        GetView().input_chat.text = string.Empty;

        GetView().ButtonTags[(int)ChatShowType.Zone -1].gameObject.SetActive(GameRoomM.Instance._GameState != GameRoomM.GameState.No);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnSync_BagResouse);
        UniGameEvent.AddEventListener<GameEvent.Sync_TakeSource>(OnSync_TakeSource);
    }

    public override void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnSync_BagResouse);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_TakeSource>(OnSync_TakeSource);
    }

    [Subscriber(ThreadMode.Main)]

    private void OnChatCallBack(GameEvent.Sync_ChatSendBack _e)
    {
        if (isNullGo)
            return;

        GetView().input_chat.text = "";
    }

    [Subscriber(ThreadMode.Main)]

    private void OnRefreshChatUI(GameEvent.Sync_ChatPush _e)
    {
        if (isNullGo)
            return;

        if (_e.isOnlyHall)
            return;

        OnChatUIShow();

        if (chatInfoList.Count > 0)
            GetView().loopListView.MovePanelToItemIndex(GetView().loopListView.ItemTotalCount - 1, 0);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_BagResouse(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;

        SetConstShow();
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnSync_TakeSource(GameEvent.Sync_TakeSource _e)
    {
        if (isNullGo)
            return;

        SetConstShow();
    }

}