using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using MetaFramework.UI;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 首次登录，引导主线战斗及释放大招
/// </summary>
public class GoMainBattle : NewPlayerGuideStepBase
{
    public Button goClick;

    public Button goPlayClick;

    bool isFirst = false;
    public override void Enter(UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> parentController)
    {
        base.Enter(parentController);

        goClick.onClick.RemoveAllListeners();
        goPlayClick.onClick.RemoveAllListeners();

        goClick.onClick.AddListener(OnClickFight);
        goPlayClick.onClick.AddListener(OnTryReleaseSkillClick);

        //设置首次点击区域位置
        NewPlayerGuideWidget.Instance.SetPosition(goClick.transform, "Hall_FightButton");
        DoNextStep();

        isFirst = true;
    }

    //此处强制关闭其他UI
    void Update()
    {
        if (isFirst)
        {
            var list = UIManager.Instance.GetAllOpenedUIByLayer(UILayer.Normal);
            if (list.Count > 0 && list[0].uiName != "UIHall")
            {
                UnityEngine.Debug.LogError($"close UI {list[0].uiName}");
                UIManager.Instance.CloseNormal(list[0].uiName);
            }
        }
    }

    /// <summary>
    /// 用于断线重连
    /// </summary>
    public override void ReEnter()
    {

    }

    public override void Exit()
    {
        goClick.onClick.RemoveAllListeners();
        goPlayClick.onClick.RemoveAllListeners();

        NewPlayerGuideWidget.Instance.SaveGuideStepToSever((int)eGuideStep.GoToMainBattle);

        gameObject.SetActive(false);

        UITool.Close<UINewPlayerGuide>();
    }


    //第一步骤，点击战斗按钮
    private void OnClickFight()
    {
        isFirst = false;
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("Hall_FightButton");

        Steps[CurrentStep]?.SetActive(false);
        UniGameEvent.AddEventListener<GameEvent.Refresh_NewPlayerGuide>(OnShowReleaseTips); //进入战斗后开始监听消息

        SendToShuShu(false, true);
    }

    int currentSkillIndex = 0;
    private void OnShowReleaseTips(GameEvent.Refresh_NewPlayerGuide guide)
    {
        UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(OnShowReleaseTips);//移除监听

        UniGameEvent.Send(new GameEvent.HallBallPauseOrPlay() { pause = true });//暂停游戏

        Data = guide.args;
        if (Data == null || Data.Length == 0)
        {
            UnityEngine.Debug.LogError("数据传输出错，请检查特殊技能信息传递参数！");
        }
        NewPlayerGuideWidget.Instance.SetPosition(goPlayClick.transform, $"Hall_SpecialSkill_Pos{Data[currentSkillIndex]}");
        DoNextStep();
    }

    // 尝试发大招
    private void OnTryReleaseSkillClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent($"Hall_SpecialSkill_{Data[currentSkillIndex]}");

        currentSkillIndex++;
        if(currentSkillIndex == 1)
            SendToShuShu(true, false);

        //继续执行点击
        if (currentSkillIndex < Data.Length)
        {
            NewPlayerGuideWidget.Instance.SetPosition(goPlayClick.transform, $"Hall_SpecialSkill_Pos{Data[currentSkillIndex]}");
            return;
        }

        DoNextStep();
    }

    private void SendToShuShu(bool isFirstCombo, bool is_Fistclick_fight)
    {
        Dictionary<string, object> properties = new Dictionary<string, object>();
        properties.Add("main_id", 1);
        properties.Add("is_first_combo", isFirstCombo);
        properties.Add("click_fail_combo", false);
        properties.Add("is_first_hand_combo", false);
        properties.Add("is_click_fight", is_Fistclick_fight);
        properties.Add("click_combo", 0);

        TDAnalyticsUtil.Track("combo", properties, true);
    }
}