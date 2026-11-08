using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Singleton;
using MetaFramework.UI;
using pb;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public enum eGuideStep
{ 
    None = 0,//没有可执行的引导

    GoToMainBattle  = 1,//首次打主线

    GoToBattle = 2,//打战役,与第3步连接
    InPlayBattleGuidance = 3,//在战役中引导

    AddNewRoleInFormation = 4,//引导玩家上阵

    ClocseItemBoxAndGoToMainFightAgain = 5,//关闭临时物品箱并引导再次战斗，非强制引导
}

public class NewPlayerGuideWidget : SingletonTemplate<NewPlayerGuideWidget>
{
    private int[] guidanceRecord = new int[0];
    public void TryShowGuide()
    {
        var curGuideStep = GetGuideStep();

        //Debug.LogError("=================="+curGuideStep);
        if (curGuideStep == eGuideStep.None) return;

        if (GameRoomM.Instance._GameState == GameRoomM.GameState.No)
        {
            UIManager.Instance.CloseAllExceptHall();
        }

        ShowGuide(curGuideStep);

        
    }

    public eGuideStep GetGuideStep()
    {
        if (!FinishGuidence(eGuideStep.GoToMainBattle)) //首个新手引导进游戏就触发
            return eGuideStep.GoToMainBattle;
        else if (!FinishGuidence(eGuideStep.InPlayBattleGuidance) && FuncOpenShowWidget.Instance.MainStageIsPass( CFG_XMLConfig.XML_Config.SwitchOpne.Play)) //完成第一步引导并解锁战役功能（由于UI在进战斗会被关闭，所以这里拆分了两步）
        {
            //Debug.LogError("======" + GameRoomM.Instance._GameState);
            if (GameRoomM.Instance._GameState == GameRoomM.GameState.No) //在大厅用这一步
                return eGuideStep.GoToBattle;
            else if (GameRoomM.Instance._GameState == GameRoomM.GameState.Play) //在战斗中用这一步
                return eGuideStep.InPlayBattleGuidance;
        }
        else if (FinishGuidence(eGuideStep.InPlayBattleGuidance) && PartnerSlotIsNull(1)) //引导上阵
            return eGuideStep.AddNewRoleInFormation;

        return eGuideStep.None;
    }

    public void TryShowGuideInPlayBattle()
    {
        if (FinishGuidence(eGuideStep.GoToMainBattle) && FuncOpenShowWidget.Instance.MainStageIsPass(CFG_XMLConfig.XML_Config.SwitchOpne.Play))
            TryShowGuide();
    }

    public bool CheckTaskIsDoing(int taskID)
    {
        var alltask = TaskM.Instance.GetTasksByClass(TaskClass.MAINLINE);
        if (alltask == null || alltask.Count == 0)
            return false;
        alltask.Sort((a, b) => b.CfgId.CompareTo(a.CfgId));

        if (alltask[0].CfgId != taskID) return false;

        var taskInfo = CFG_Task.GetById(alltask[0].CfgId);

        return alltask.Count > 0 && alltask[0].Status == TaskStatus.IN_PROGRESS;
    }

    public bool CheckMainStageIsPass(int stageID)
    {
        return RoleM.Instance.MainBattleId >= stageID;
    }

    //伙伴位置是否为空
    public bool PartnerSlotIsNull(int slot)
    {
        if (CFG_XMLConfig.CheckPartnerSoltIsOpenServer(slot))
        {
            UITool.ShowTips($"槽位尚未解锁，无法选择！");
            return false;
        }

        if (CFG_XMLConfig.GetOpenSoltLv(slot) > RoleM.Instance.Level)
        {
            UITool.ShowTips("槽位尚未解锁，无法选择！");
            return false;
        }

        var solt = PartnerM.Instance.GetOnePartnerSlotInfo(slot);
        if (solt.PartnerUid == 0)
            return true;

        return false;
    }

    public void ShowGuide(eGuideStep step)
    {
        UITool.Open<UINewPlayerGuide>(null, step);
    }

    public bool FinishGuidence(eGuideStep step)
    {
        if (guidanceRecord == null || guidanceRecord.Length == 0) return false;

        return guidanceRecord[0] >= (int)step;
    }

    public void CacheGuideStep(int[] guidance)
    {
        guidanceRecord = guidance;
    }

    #region
    private Dictionary<string, Transform> dicGuideGameObjec = new Dictionary<string, Transform>();

    public void AddGuideGameObject(string key, Transform target)
    {
        key = key.ToLower();
        if (dicGuideGameObjec.ContainsKey(key))
            dicGuideGameObjec[key] = target;
        else
            dicGuideGameObjec.Add(key, target);
    }

    public void SetPosition(Transform tf, string key)
    {
        Transform tfTarget = GetGuideObject<Transform>(key);
        if (tf == null || tfTarget == null) return;

        tf.position = tfTarget.position;
    }

    public (bool,Vector3) GetPosition(string key)
    {
        Transform tfTarget = GetGuideObject<Transform>(key);
        if (tfTarget == null) return (false, Vector3.zero);

        return (true, tfTarget.position);
    }

    public void InvokeButtonEvent(string key)
    {
        ExButton btn = GetGuideObject<ExButton>(key);
        if (btn != null)
        {
            btn.OnClick.Invoke();
            return;
        }

        Button btTarget = GetGuideObject<Button>(key);
        btTarget?.onClick.Invoke();
    }

    public T GetGuideObject<T>(string key)
    {
        key = key.ToLower();
        if (dicGuideGameObjec.ContainsKey(key) && dicGuideGameObjec[key] != null)
        {
            return dicGuideGameObjec[key].GetComponent<T>();
        }

        Debug.LogWarning($"新手warning：未找到key为{key}对应的组件信息");
        return default(T);
    }
    #endregion

    public void SaveGuideStepToSever(int curStep)
    {
        if (guidanceRecord == null || guidanceRecord.Length == 0)
            guidanceRecord = new int[1] { 0 }; 

        guidanceRecord[0] = curStep;

        Log.Print($"保存引导ID到服务器，passMaxId:{curStep}");
        TcpC.C2S_GuidanceUpdateRequest(guidanceRecord);
    }
}