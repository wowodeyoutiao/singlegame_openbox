using DG.Tweening;
using MetaFramework.UI;
using System;
using System.Diagnostics;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 当开启战役时的新手引导
/// </summary>
public class GoPlayBattle : NewPlayerGuideStepBase
{
    public Button goClick;
    public Button goStartChallengeClick;

    bool isFirst = false;
    public override void Enter(UIBase<UIModelNewPlayerGuide, UIViewNewPlayerGuide> parentController)
    {
        base.Enter(parentController);

        goClick.onClick.RemoveAllListeners();
        goStartChallengeClick.onClick.RemoveAllListeners();
        UniGameEvent.RemoveEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //移除UI打开监听

        goClick.onClick.AddListener(OnClickFight);
        goStartChallengeClick.onClick.AddListener(OnStartChallengeClickClick);
        //goPlayClick.onClick.AddListener(OnTryReleaseSkillClick);
        //设置首次点击区域位置
        NewPlayerGuideWidget.Instance.SetPosition(goClick.transform, "DownBtns_FightButton");
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
        goStartChallengeClick.onClick.RemoveAllListeners();

        UniGameEvent.RemoveEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //移除UI打开监听

        gameObject.SetActive(false);
    }


    //第一步骤，点击战斗按钮
    private void OnClickFight()
    {
        isFirst = false;
        if (CurrentStep > 0) return;

        Steps[CurrentStep]?.SetActive(false);

        UniGameEvent.AddEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //进入战斗后开始监听消息
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("DownBtns_FightButton");
    }

    private void WaitOpenUIFightChoice(GameEvent.UIPanel_Operate opened)
    {
        if (string.IsNullOrEmpty(opened.uiName)) return;

        if (opened.uiName.ToLower() == "UIFightChoice".ToLower() && opened.isOpen)
        {
            DoNextStep();
            UniGameEvent.RemoveEventListener<GameEvent.UIPanel_Operate>(WaitOpenUIFightChoice); //移除UI打开监听

            FlyButton(goStartChallengeClick.transform, goClick.transform.position, "FightChoice_FightButton", 1f);
        }
    }

    //第二步，点击开始挑战
    private void OnStartChallengeClickClick()
    {
        NewPlayerGuideWidget.Instance.InvokeButtonEvent("FightChoice_FightButton");
        DoNextStep();//Steps[CurrentStep]?.SetActive(false);//进入游戏后这个UI会被强制关闭。
    }

    //引导遥感
    private void WaitBattleStart(GameEvent.Refresh_NewPlayerGuide guide)
    {
        DoNextStep();
    }

    //int currentSkillIndex = 0;
    //private void OnShowReleaseTips(GameEvent.Refresh_NewPlayerGuide guide)
    //{
    //    UniGameEvent.RemoveEventListener<GameEvent.Refresh_NewPlayerGuide>(OnShowReleaseTips);//移除监听

    //    UniGameEvent.Send(new GameEvent.HallBallPauseOrPlay() { pause = true });//暂停游戏

    //    Data = guide.args;
    //    if (Data == null || Data.Length == 0)
    //    {
    //        UnityEngine.Debug.LogError("数据传输出错，请检查特殊技能信息传递参数！");
    //    }
    //    NewPlayerGuideWidget.Instance.SetPosition(goPlayClick.transform, $"Hall_SpecialSkill_{Data[currentSkillIndex]}");
    //    DoNextStep();
    //}

    //// 尝试发大招
    //private void OnTryReleaseSkillClick()
    //{
    //    NewPlayerGuideWidget.Instance.InvokeButtonEvent($"Hall_SpecialSkill_{Data[currentSkillIndex]}");

    //    currentSkillIndex++;

    //    //继续执行点击
    //    if (currentSkillIndex < Data.Length)
    //    {
    //        NewPlayerGuideWidget.Instance.SetPosition(goPlayClick.transform, $"Hall_SpecialSkill_{Data[currentSkillIndex]}");
    //        return;
    //    }

    //    DoNextStep();
    //}
}