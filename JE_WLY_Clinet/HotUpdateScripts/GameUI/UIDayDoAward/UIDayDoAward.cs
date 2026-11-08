using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;

public class UIDayDoAward : UIBase<UIModelDayDoAward, UIViewDayDoAward>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_No.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            TcpC.C2S_DoubleReward(false);
        });

        GetView().Btn_Yes.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (!GameUtil.Instance.ComCoinToCheck(uesdID, uesdNum, null, false))
            {
                UITool.ShowTips("货币不足");
                return;
            }
            string timeStr = GameUtil.Instance.GetDataTime(GameUtil.Instance.GetCurrentTimeSeconds()).DayOfYear.ToString();

            if (PlayerPrefsExtend.GetString("UIDayDoAward") == timeStr)
            {
                TcpC.C2S_DoubleReward(true);
            }
            else
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = $"是否确定消耗{GameUtil.Instance.FormatNumWithUnit(uesdNum)}{GameUtil.Instance.GetItemName(uesdID)}进行双倍领取？",
                    isTogg = true,
                    TogNoFun = () => {
                        PlayerPrefsExtend.SetString("UIDayDoAward", "");
                    },
                    TogYesFun =() => {
                        PlayerPrefsExtend.SetString("UIDayDoAward", timeStr);
                    },
                    YesFun = delegate { TcpC.C2S_DoubleReward(true); },
                });
            }
        });
    }

    private int uesdID;
    private int uesdNum;
    public override void OnEnter(params object[] args)
    {
        var items = args[0] as List<pb.ItemInfo>;
        int RoomId = (int)args[1];
        for (int i = 0; i < GetView().root_award.childCount; i++)
        {
            var go = GetView().root_award.GetChild(i);
            go.gameObject.SetActive(i < items.Count);
            if (i < items.Count)
            {
                GameUtil.Instance.SetItemIcon(items[i].CfgId, JUIExtensions.img(go, "icon"), JUIExtensions.img(go));
                JUIExtensions.txt(go, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(items[i].Num);
            }
        }
        uesdID = CFG_XMLConfig.xmlConfig.DayRoomDic[CFG_Room.GetById(RoomId).Roomtype].doubleaccessID;
        uesdNum = CFG_XMLConfig.xmlConfig.DayRoomDic[CFG_Room.GetById(RoomId).Roomtype].doubleaccessNum;

        GetView().txt_num.text = GameUtil.Instance.GetCheakNumStr(uesdID, uesdNum, 
            "x"+ GameUtil.Instance.FormatNumWithUnit(uesdNum));

        GameRoomM.Instance.gameResultDailyMap = null;
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}