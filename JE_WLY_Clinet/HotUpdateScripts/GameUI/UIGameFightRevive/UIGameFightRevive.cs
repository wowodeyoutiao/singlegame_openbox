using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.MyNet;
using pb;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;

public class UIGameFightRevive : UIBase<UIModelGameFightRevive, UIViewGameFightRevive>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_Revive.onClick.AddListener(()=> {

            if (!GameUtil.Instance.ComCoinToCheck(CFG_Info.needMoneyID, CFG_Info.num, null))
            {
                return;
            }
            TcpC.C2S_Revive(1);

            UITool.Close<UIGameFightRevive>();
        });

        GetView().Btn_Back.onClick.AddListener(() => {
            if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace)
            {
                TcpC.C2S_Revive(0);
                UITool.Close<UIGameFightRevive>();
            }
            else
                TcpC.C2S_Exit();
        });
    }

    private CFG_XMLConfig.XML_Config.XML_Item_Revive CFG_Info;
    private TimeItem TimeJaction;

    private void ClearTime()
    {
        if (TimeJaction != null)
        {
            GameHelper.Instance.RemoveInvoke(TimeJaction);
            TimeJaction = null;
        }
    }
    public override void OnEnter(params object[] args)
    {
        var info = args[0] as ObjectDiePush;

        GetView().Btn_Revive.gameObject.SetActive(info.CanRevive);

        bool isWarFight = GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace;
        GetView().txt_Back.text = isWarFight ? "复活点复活" : "退出副本";

        int CDTime = info.Countdown;
        string ShowStr = "";
        if (info.CanRevive)
        {
            if (isWarFight)
            { 
                CFG_Info = info.ReviveTimes >= CFG_XMLConfig.xmlConfig.WarFight.revive.Length ? CFG_XMLConfig.xmlConfig.WarFight.revive[CFG_XMLConfig.xmlConfig.WarFight.revive.Length -1] : CFG_XMLConfig.xmlConfig.WarFight.revive[info.ReviveTimes]; //复活点复活

                ShowStr = "角色已死亡，是否消耗" + CFG_Info.num + GameUtil.Instance.GetItemName(CFG_Info.needMoneyID) + "立即原地复活 ?\n({0}S）";
            }
            else
            {
                CFG_Info = CFG_XMLConfig.xmlConfig.ReviveDic[1 + info.ReviveTimes];
                ShowStr = "角色已死亡，是否消耗" + CFG_Info.num + GameUtil.Instance.GetItemName(CFG_Info.needMoneyID) + "立即原地复活 ?\n {0}秒内无操作将自动退出副本";
            }
        }
        else
        {
            ShowStr = "当前复活次数已经用完,是否立即退出副本?{0}秒内无操作将自动退出";
        }

        GetView().Text_Tip.text = isWarFight ? "提示：倒计时结束无操作将会在复活点复活" : $"提示：当前还可复活{3 - info.ReviveTimes}次";
        ClearTime();
        GetView().Text_des.text = string.Format(ShowStr, CDTime);
        TimeJaction = GameHelper.Instance.AddInvoke(new TimeItem (CDTime,1, () =>
        {
            CDTime--;
            GetView().Text_des.text = string.Format(ShowStr, CDTime);
            if (CDTime == 0)
            {
                UITool.Close<UIGameFightRevive>();
            }
        }));
    }

    public override void OnExit()
    {
        ClearTime();
    }

    public override void OnDestroy()
    {
        ClearTime();
    }
}