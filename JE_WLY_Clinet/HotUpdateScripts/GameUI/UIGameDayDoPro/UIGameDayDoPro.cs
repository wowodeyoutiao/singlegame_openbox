using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Helper;
using HotUpdateScripts.Game;
using JEngine.Event;

public class UIGameDayDoPro : UIBase<UIModelGameDayDoPro, UIViewGameDayDoPro>
{
    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }
    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);
        GetView().btn_Show.onClick.AddListener(delegate
        {
            OnUIShow(true);
        });

        GetView().btn_Close.onClick.AddListener(delegate
        {
            OnUIShow(false);
        });

        GetView().obj_Pro.SetActive(!GameRoomM.Instance.GetCanNoAttackGame());

        OnInitInfoShow();
        OnTimeShow();

        OnUniEventAddOrMove(true);
    }

    /// <summary>
    /// 升星石副本怪物IDs
    /// </summary>
    private int[] SXSMonsterIDs = { 602, 603, 604 };

    /// <summary>
    /// 魂石副本怪物IDs
    /// </summary>
    private int[] HSMonsterIDs = { 605, 606, 607 };

    public void OnTimeShow()
    {
        GetView().txt_Time.text = $"倒计时:{GameUtil.Instance.GetTimeSToHMS(GameRoomM.Instance.GameFightTime - GameUtil.Instance.GetCurrentTimeSeconds())}";
        OnCollectRefresh();
    }

    private CFG_XMLConfig.XML_Config.XML_DailyMap XMLConfig;
    private void OnInitInfoShow()
    {
        XMLConfig = CFG_XMLConfig.xmlConfig.DayRoomDic[(CFG_Room.RoomType)GameRoomM.Instance.CFGRoom.Roomtype];
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Gold || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.EXP)
        {
            GetView().showTips[0].txt_Des.text = "伤害评分:";
            GetView().showTips[0].txt_Num.text = XMLConfig.ScoreDic[0].key;
            GetView().showTips[1].txt_Des.text = "累计伤害:";
            GetView().showTips[1].txt_Num.text = "0";

            GetView().txt_desPro.text = "进度";
            GetView().obj_pf.SetActive(true);
            GetView().txt_pf1.text = XMLConfig.ScoreDic[0].key;
            GetView().txt_pf2.text = XMLConfig.ScoreDic[1].key;
            GetView().txt_Pro.text = "0/";
            GetView().slider_Pro.value = 0;
        }
        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.QHS || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SY)
        {
            GetView().showTips[0].txt_Des.text = "当前波次:";
            GetView().showTips[0].txt_Num.text = "0";
            GetView().showTips[1].txt_Des.text = "剩余怪物数量:";
            GetView().showTips[1].txt_Num.text = MonsterManager.Instance.MonesterCount.ToString();

            GetView().txt_desPro.text = "杀怪进度";
            GetView().obj_pf.SetActive(false);
            GetView().txt_Pro.text = "";
            GetView().slider_Pro.value = 1;
        }
        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SXS || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.HS)
        {
            var Ids = GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SXS ? SXSMonsterIDs : HSMonsterIDs;
            for (int i = 0; i < Ids.Length; i++)
            {
                GetView().showTips[i].txt_Des.text = GameUtil.Instance.GetStrByQuality(CFG_Monster.GetById(Ids[i]).Quality,CFG_Monster.GetById(Ids[i]).Name) ;
                GetView().showTips[i].txt_Num.text = MonsterManager.Instance.GetNumByMonsterId(Ids[i]).ToString();
            }
        }
    }

    [Subscriber(ThreadMode.Main)]
    private void OnRefresh(GameEvent.DailyMap_Score _e )
    {
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Gold || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.EXP)
        {
            GetView().showTips[0].txt_Num.text = _e.Scoring;
            GetView().showTips[1].txt_Num.text = GameUtil.Instance.FormatNumWithUnit(_e.TotalDamage);

            GetView().txt_pf1.text = _e.Scoring;

            int idx = XMLConfig.ScoreDic.Length;
            for (int i = 0; i < XMLConfig.ScoreDic.Length; i++)
            {
                if (XMLConfig.ScoreDic[i].key == _e.Scoring)
                {
                    idx = i;
                    break;
                }

            }
            GetView().txt_pf2.text = idx == (XMLConfig.ScoreDic.Length - 1) ? XMLConfig.ScoreDic[idx ].key : XMLConfig.ScoreDic[idx + 1].key;

            var maxHp = MonsterManager.Instance.getMonsterObjs()[0]._monser.maxHp;

            GetView().txt_Pro.text = GameUtil.Instance.FormatNumWithUnit(_e.TotalDamage)+ "/" + GameUtil.Instance.FormatNumWithUnit(maxHp);

            GetView().slider_Pro.value =(float) ((double)_e.TotalDamage / maxHp);
        }
        else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.QHS || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SY)
        {
            GetView().showTips[0].txt_Num.text = _e.Round.ToString();
        }

        OnMonsterRefresh(default);
    }

    private void OnUniEventAddOrMove(bool isAdd)
    {
        if (isAdd)
        {
            UniGameEvent.AddEventListener<GameEvent.Fight_Die>(OnMonsterRefresh);
            UniGameEvent.AddEventListener<GameEvent.Fight_HPChange>(OnHPChangeEff);
        }
        else
        {
            UniGameEvent.RemoveEventListener<GameEvent.Fight_HPChange>(OnHPChangeEff);
            UniGameEvent.RemoveEventListener<GameEvent.Fight_Die>(OnMonsterRefresh);
        }
    }
    private void OnMonsterRefresh(GameEvent.Fight_Die _e)
    {

        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.QHS || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SY)
        {
            GetView().showTips[1].txt_Num.text = MonsterManager.Instance.MonesterCount.ToString();

            GetView().txt_Pro.text = (GameRoomM.Instance.MonsterAllCount - MonsterManager.Instance.MonesterCount) + "/" + GameRoomM.Instance.MonsterAllCount;
            GetView().slider_Pro.value = (float)(GameRoomM.Instance.MonsterAllCount - MonsterManager.Instance.MonesterCount) / GameRoomM.Instance.MonsterAllCount;
        }
        OnCollectRefresh();
    }

    /// <summary>
    /// 收集状态刷新
    /// </summary>
    private void OnCollectRefresh()
    {
        if (GameRoomM.Instance.GetCanNoAttackGame())
        {
            var Ids = GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.SXS ? SXSMonsterIDs : HSMonsterIDs;
            for (int i = 0; i < Ids.Length; i++)
            {
                GetView().showTips[i].txt_Num.text = MonsterManager.Instance.GetNumByMonsterId(Ids[i]).ToString();
            }
        }
    }


    private void OnHPChangeEff(GameEvent.Fight_HPChange _e)
    {
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.Gold || GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.EXP)
        {
            var _monser = MonsterManager.Instance.getMonsterObjs()[0]._monser;
            var curAttack = _monser.maxHp - _monser.hp;
            GetView().showTips[1].txt_Num.text = GameUtil.Instance.FormatNumWithUnit(curAttack);
            GetView().txt_Pro.text = GameUtil.Instance.FormatNumWithUnit(curAttack) + "/" + GameUtil.Instance.FormatNumWithUnit(_monser.maxHp);

            GetView().slider_Pro.value = (float)(curAttack) / _monser.maxHp;
        }
    }


    private void OnUIShow(bool isShow)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        GetView().obj_UI.SetActive(isShow);
        GetView().btn_Show.gameObject.SetActive(!isShow);
    }

    public override void OnExit()
    {
        OnUniEventAddOrMove(false);
        JEvent.defaultEvent.Unregister(this);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}