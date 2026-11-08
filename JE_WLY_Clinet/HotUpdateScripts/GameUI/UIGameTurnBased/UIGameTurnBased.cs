using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;

using JEngine.Event;

using MetaFramework.UI;

using static GameEvent;

public class UIGameTurnBased : UIBase<UIModelGameTurnBased, UIViewGameTurnBased>
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
            ExitFighting();
        });

        GetView().btn_JumpGame.onClick.AddListener(() =>
        {
            TurnBasedFightCtrl.Instance.JumpGame();
            GetView().btn_JumpGame.gameObject.SetActive(false);
        });
    }

    private UICom_TopRes TopUI = null;
    private int timeSeconds = 10;

    public override void OnLayerChange()
    {
        if (TopUI == null)
        {
            UITool.LoadUII<UICom_TopRes>(GetView().Root, (ui) =>
            {
                TopUI = ui as UICom_TopRes;
                TopUI.OnSetLayer(_canvasOrder);
            });
        }
    }

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        GetView().txt_Rount.text = "";
        GetView().txt_Time.text = "";

        GetView().btn_JumpGame.gameObject.SetActive(false);
        GetView().btn_JumpGameDisabled.gameObject.SetActive(false);
        GetView().root_JumpTime.SetActive(false);
        Event_ShowUserName(default);
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_ShowUserName(PVP_TurnBasedShowUserName _e)
    {
        if (isNullGo) return;

        if (_e.isMy)
        {
            if (TurnBasedRoleManager.Instance.GetBlueProtagonist(out var blueRole))
                GetView().txt_MyName.text = blueRole.roleName;
            else
                GetView().txt_MyName.text = "";
        }
        else
        {
            if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.WuXing)
            {
                GetView().txt_OtherName.text = $"试炼{TurnBasedM.Instance.currentChallengeLayer}层";
            }
            else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More)
            {
                if (TurnBasedRoleManager.Instance.GetRedProtagonist(out var redRole))
                    GetView().txt_OtherName.text = redRole.roleName;
                else
                    GetView().txt_OtherName.text = "";
            }
            else
                GetView().txt_OtherName.text = "";
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void Event_ShowRound(PVP_TurnBasedRoundShow _e)
    {
        if (isNullGo)
            return;

        ShowRound(_e.currentRound, _e.totalRound);
    }

    private void ShowRound(int currentRound, int totalRound)
    {
        GetView().txt_Rount.text = $"回合数：{currentRound}/{CFG_XMLConfig.xmlConfig.GVGInfo.MaxRound}";
    }

    private void ExitFighting()
    {
        string showStr = "确认退出游戏？";

        UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
        {
            ShowStr = showStr,
            YesFun = () =>
            {
                TurnBasedFightCtrl.Instance.JumpGame();
                ShuShuClickSend();
            },
            BtnYes = "直接退出",
            BtnNo = "继续游戏",
        });
    }

    private void ShuShuClickSend()
    {
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.WuXing)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "五行试炼",
                component = "五行试炼战斗界面",
                position = "",
                event_click = $"退出战斗",
            });
        }
        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "竞技",
                position = "回合制挑战战斗界面",
                event_click = $"退出战斗",
            });
        }
    }

    // 倒计时结束后可跳过战斗
    TimeItem jumpGameTimeItem;

    [Subscriber(ThreadMode.Main)]
    private void ShowJumpGameTime(GameEvent.PVP_JumpGameCountdown _e)
    {
        if (isNullGo) return;

        GetView().btn_JumpGameDisabled.gameObject.SetActive(true);
        GetView().root_JumpTime.SetActive(true);
        GetView().txt_JumpTime.text = $"{timeSeconds}秒后可跳过战斗";

        jumpGameTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(timeSeconds + 1, 1, () =>
        {
            if (isNullGo) return;
            GetView().txt_JumpTime.text = $"{timeSeconds}秒后可跳过战斗";
            if (timeSeconds == 0)
            {
                GetView().txt_JumpTime.text = "";
                GetView().root_JumpTime.SetActive(false);
                GetView().btn_JumpGameDisabled.gameObject.SetActive(false);
                GetView().btn_JumpGame.gameObject.SetActive(true);
            }
            timeSeconds--;
        }));
    }

    public void ClearTime()
    {
        if(GameHelper.Instance != null && jumpGameTimeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(jumpGameTimeItem);
            jumpGameTimeItem = null;
        }
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        ClearTime();
    }

    public override void OnDestroy()
    {
        ClearTime();
    }
}