using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Pkcs;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.Singleton;
using MetaFramework.UI;
using pb;
using System.Collections;
using System.Collections.Generic;
using System.Security.Policy;
using UnityEngine;
using UnityEngine.Windows.WebCam;
using static CFG_XMLConfig.XML_Config;
using static UIViewGameDayDoPro;

public enum EMainUIObject
{ 
    Devlop= 0,       //养成
    PartnerCall = 5,//召唤
    FightChoice = 10, //玩法
    //Guild = 15, //行会
    //Wuxing = 20, //五行
    //Title = 25,   //称号
    //Box = 30,   //宝箱

    LimitedTimeActivity = 35,//限时活动
    Competion = 40,            //赛事

    Zhufu = 45,        //祝福
}

public class FuncOpenShowWidget : SingletonTemplate<FuncOpenShowWidget>
{
    #region 功能开启数据操作
    //已经显示过的信息
    private HashSet<int> showedRecord = new HashSet<int>();

    //功能显示列表
    private List<int> showFuncStack = new List<int>();

    private long CurUnixDay = 0;

    /// <summary>
    /// 尝试显示UI，在主界面时调用
    /// </summary>
    public void TryShow()
    {
        if (showFuncStack.Count == 0) return;

        RecordFunOpenFlag(); //记录功能开启标记

        var curStep = NewPlayerGuideWidget.Instance.GetGuideStep();
        if (curStep != eGuideStep.None || UITool.IsOpen<UINewPlayerGuide>()) //有可执行的新手引导
        {
            (bool hasType, List<int> typeIDList) = hasExceptType();
            
            for (int i = 0; i < showFuncStack.Count; i++)
            {
                int typeID = typeIDList.Find((a) => a == showFuncStack[i]); //处理特殊情况
                if (typeID > 0)
                    continue;
                
                TcpC.C2S_RecordOpenID(showFuncStack[i]);
            }

            //if (hasType)
            //{
            //    UIManager.Instance.CloseAllExceptHall();
            //    UITool.Open<UIFuncOpenShow>(null, new List<int>(typeIDList));
            //}
            //NewPlayerGuideWidget.Instance.TryShowGuide();
            showFuncStack.Clear();

            return;
        }

        UIManager.Instance.CloseAllExceptHall();

        UITool.Open<UIFuncOpenShow>(null, new List<int>(showFuncStack));
        showFuncStack.Clear();
    }

    private (bool, List<int>) hasExceptType()
    {
        bool hasType = false;
        List<int> types = new List<int>();
        //int type = showFuncStack.Find((a) => a == (int)SwitchOpne.Play);
        //if (type > 0)
        //{
        //    hasType = true;
        //    types.Add(type);
        //} 

        return (hasType, types);
    }

    /// <summary>
    /// 尝试添加显示数据，在天数和关卡变化时调用
    /// </summary>
    public void TryAddShowData()
    {
        CurUnixDay = getUnixDayFromOpen();
        if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction == null) return;

        foreach (var item in CFG_XMLConfig.xmlConfig.SwitchOpenFunction)
        {
            if(!funcCanShow(item.Value.id)) continue;

            addShowData(item.Value.id);
        }
    }

    private void addShowData(int id)
    {
        var mid = showFuncStack.Find((a) => a == id);
        if (mid > 0) return;

        showFuncStack.Add(id);
    }

    /// <summary>
    /// 初始化显示记录，用于服务器返回信息时调用
    /// </summary>
    public void InitRecord(int[] data)
    {
        //每次服务器下发都做初始化操作。这样在游戏退出时就不需要做清空处理了
        showedRecord.Clear();

        if (data == null) return;

        for (int i = 0; i < data.Length; i++)
        {
            if (!showedRecord.Contains(data[i]))
                showedRecord.Add(data[i]);
        }
    }

    #endregion

    #region 功能开启UI逻辑
    private Dictionary<EMainUIObject, Transform> flyTargetDic = new Dictionary<EMainUIObject, Transform>();  //定义飞向的目标与Transform关系

    //定义具体模块与飞入目标对应表
    private Dictionary<CFG_XMLConfig.XML_Config.SwitchOpne, EMainUIObject> funcAndParentrelationShip = new() {
        //飞向养成
        { CFG_XMLConfig.XML_Config.SwitchOpne.OpenBox, EMainUIObject.Zhufu },                       //开箱玩法
        { CFG_XMLConfig.XML_Config.SwitchOpne.Build, EMainUIObject.Devlop },                            //装备强化
        { CFG_XMLConfig.XML_Config.SwitchOpne.Build_EquipInster, EMainUIObject.Devlop },            //魂石
        { CFG_XMLConfig.XML_Config.SwitchOpne.Title, EMainUIObject.Devlop },                               //称号
        { CFG_XMLConfig.XML_Config.SwitchOpne.Refined, EMainUIObject.Devlop },                         //槽位洗炼

         //飞向召唤
        { CFG_XMLConfig.XML_Config.SwitchOpne.CallRole, EMainUIObject.PartnerCall },                     //守卫海岛

        //飞向玩法
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play, EMainUIObject.FightChoice},                              //战役
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_Boss, EMainUIObject.FightChoice },                     //首领
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower, EMainUIObject.FightChoice },               //试炼塔
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_Relice, EMainUIObject.FightChoice },                  //遗迹
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_QHSAndGold, EMainUIObject.FightChoice},       //金币、强化副本
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_ExpAndSXS, EMainUIObject.FightChoice },          //经验副本+升星石副本
        { CFG_XMLConfig.XML_Config.SwitchOpne.DefendIsland, EMainUIObject.FightChoice },              //守卫海岛
        { CFG_XMLConfig.XML_Config.SwitchOpne.FightEnd, EMainUIObject.FightChoice },                     //一战到底
        { CFG_XMLConfig.XML_Config.SwitchOpne.WuXing, EMainUIObject.FightChoice },                       //五行试炼
        { CFG_XMLConfig.XML_Config.SwitchOpne.WarFight, EMainUIObject.FightChoice },                    //群雄战场

        //飞向限时活动
        { CFG_XMLConfig.XML_Config.SwitchOpne.WorldBoss, EMainUIObject.LimitedTimeActivity},        //即时PK

        //飞向赛事
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP, EMainUIObject.Competion },                     //即时PK
        { CFG_XMLConfig.XML_Config.SwitchOpne.Play_GVG, EMainUIObject.Competion },                     //回合PK
    };

    public void RegisterFlyToTarget(EMainUIObject mainUItag,Transform tfObj)
    {
        if (flyTargetDic.ContainsKey(mainUItag))
        {
            flyTargetDic[mainUItag] = tfObj;
            return;
        }

        flyTargetDic.Add(mainUItag, tfObj);
    }

    public Transform GetFlyToTarget(CFG_XMLConfig.XML_Config.SwitchOpne tag)
    {
        if (!funcAndParentrelationShip.ContainsKey(tag))
        {
            //Debug.Log("Cur Tag " + tag + "Fly To Devlop");
            //return null;
            return flyTargetDic[EMainUIObject.Devlop];
        }

        var mainTag = funcAndParentrelationShip[tag];
        if (flyTargetDic.ContainsKey(mainTag))
        {
            //Debug.Log("Cur Tag "+ tag + "Fly To " + flyTargetDic[mainTag]);
            return flyTargetDic[mainTag];
        }
            

        //return null;

        return flyTargetDic[EMainUIObject.Devlop];
    }
    #endregion


    #region

    /// <summary>
    /// 检查是否有新功能开启，有则在返回主界面时不弹二级界面
    /// </summary>
    /// <returns></returns>
    public bool CheckHasFuncOpen()
    {
        if(showFuncStack.Count > 0) return true;

        TryAddShowData();

        return showFuncStack.Count > 0;
    }

    private bool funcCanShow(int id)
    {
        int battleID = RoleM.Instance.MainBattleId;

        SwitchOpne curID = (SwitchOpne)id;

        XML_OpenFunction curFuncData = null; //当前展示功能数据
        if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(curID))
            curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[curID];

        if (curFuncData == null) return false;//没有读到配置就不显示

        if (string.IsNullOrEmpty(curFuncData.Icon)) return false;//如果没有配置图标就不显示

        if (showedRecord.Contains(id)) return false;//已经显示过了

        //第一判断条件，开服天数
        if (!dayIsEnough(curID)) return false;

        //第二判断条件，等级
        if(!levelIsEnough(curID)) return false;

        if (!MainStageIsPass(curID)) return false;
        return true;
    }

    //判断条件开服/注册天数
    private bool dayIsEnough(SwitchOpne curID)
    {
        return !GameUtil.Instance.CheckOpenServerIsLock(curID,false);
        //XML_OpenFunction curFuncData = null; //当前展示功能数据
        //if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(curID))
        //    curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[curID];

        //if (curFuncData == null) return false;//没有读到配置就不显示

        //int configDay = 1;
        //if (CFG_XMLConfig.xmlConfig.SwitchOpenDayDic.ContainsKey(curID))
        //    configDay = CFG_XMLConfig.xmlConfig.SwitchOpenDayDic[curID];

        //return configDay - 1 <= CurUnixDay;
    }

    private bool levelIsEnough(SwitchOpne curID)
    {
        XML_OpenFunction curFuncData = null; //当前展示功能数据
        if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(curID))
            curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[curID];

        if (curFuncData == null) return false;//没有读到配置就不显示

        int configLevel = 0;

        if (CFG_XMLConfig.xmlConfig.SwitchOpenLvDic.ContainsKey(curID))
            configLevel = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[curID];

        return configLevel <= RoleM.Instance.Level;
    }

    public bool MainStageIsPass(SwitchOpne curID)
    {
        XML_OpenFunction curFuncData = null; //当前展示功能数据
        if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(curID))
            curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[curID];

        if (curFuncData == null) return false;//没有读到配置就不显示

        return curFuncData.openMinBattel <= RoleM.Instance.MainBattleId;
    }

    /// <summary>
    /// 从服务器开启到现在的Unix时间（以天为单位）
    /// </summary>
    /// <returns></returns>
    private long getUnixDayFromOpen()
    {
        long openTime = GameUtil.Instance.GetServerOpenTimeMillisecond();
        long openZeroTime = GameUtil.Instance.ConvertToMidnightSeconds(openTime);
        long opday = (GameUtil.Instance.GetCurrentTimeSeconds() - openZeroTime) / 86400;
        return opday;
    }
    #endregion

    #region 功能new角标

    public bool FunctionIsLight(List<SwitchOpne> functionList)
    {
        if (functionList == null) return false;
        if (functionList.Count == 0) return false;
        if (funcFlagState.Count == 0) return false;

        foreach (var item in functionList)
        {
            int systemID = (int)item;
            if (!funcFlagState.ContainsKey(systemID)) continue;

            if (funcFlagState[systemID] == ModuleState.MS_REGISTER)
            {
                return true;
            }
        }

        return false;
    }

    private void RecordFunOpenFlag()
    {
        if (showFuncStack == null) return;

        for (int i = 0; i < showFuncStack.Count; i++)
        {
            if (funcFlagState.ContainsKey(showFuncStack[i]))
                continue;

            TcpC.C2S_ModuleNew(showFuncStack[i], pb.ModuleState.MS_REGISTER);
        }
    }

    public void HideFuncNewFlag(List<SwitchOpne> funcList)
    {
        if (funcList == null) return;

        for (int i = 0; i < funcList.Count; i++)
        {
            if (!funcFlagState.ContainsKey((int)funcList[i]))
                continue;

            if (funcFlagState[(int)funcList[i]] != ModuleState.MS_REGISTER)
                continue;

            TcpC.C2S_ModuleNew((int)funcList[i], pb.ModuleState.MS_USE);
        }
    }

    public Dictionary<int, ModuleState> funcFlagState = new();
    //上线推送一次
    public void InitFunOpenFlag(ModuleNewPush push)
    {
        funcFlagState.Clear();

        if (push == null || push.Infos == null) return;

        foreach (var item in push.Infos)
        {
            funcFlagState[item.CfgId] = item.State;
        }
    }

    public void OperateFunOpenFlag(ModuleInfo info)
    {
        if (info == null) return;

        funcFlagState[info.CfgId] = info.State;

        UniGameEvent.Send(new GameEvent.UIFunctionOpenOperate());
    }
    
    #endregion
}